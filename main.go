package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	pp "github.com/maruel/panicparse/stack"
	"github.com/valyala/fasthttp"
)

const (
	STATIC_DIR = "./public"
	IMAGE_DIR  = "./image"
)

var (
	db            *sql.DB
	staticHandler = fasthttp.FSHandler(STATIC_DIR, 0)
)

func notfoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(404)
	fmt.Fprintf(ctx, "file not found")
}

var routes = map[string]fasthttp.RequestHandler{
	"/api/":   apiHandler,
	"/image/": imageHandler,
	"/blog/":  blogHandler,
	// audioHandler
}

func router(path string, ctx *fasthttp.RequestCtx) bool {
	for route, handler := range routes {
		if strings.HasPrefix(path, route) {
			handler(ctx)
			return true
		}
	}
	return false
}

func getStack(stack []byte) string {
	in := bytes.NewBuffer(stack)
	trace, err := pp.ParseDump(in, ioutil.Discard)
	checkErr(err)
	p := &pp.Palette{}
	buckets := pp.SortBuckets(pp.Bucketize(trace, pp.AnyValue))
	src, pkg := pp.CalcLengths(buckets, false)
	ret := ""
	for _, bucket := range buckets {
		ret += p.StackLines(&bucket.Signature, src, pkg, false)
	}
	return ret
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	defer func() {
		log.Println("->", string(ctx.Method()), string(ctx.Path()))
		s := ctx.Response.StatusCode()
		log.Println("<-", s, fasthttp.StatusMessage(s))
	}()
	defer func() {
		if x := recover(); x != nil {
			ctx.SetStatusCode(500)
			fmt.Fprintf(ctx, renderTemplate("./templates/error_500.html", struct {
				Message, Stack string
			}{x.(string), getStack(debug.Stack())}))
		}
	}()

	path := string(ctx.Path())

	switch {
	case path == "/", path == "/index.html":
		fallthrough
	default:
		fasthttp.ServeFile(ctx, "./templates/index.html")
	case router(path, ctx):
	case fileExists(STATIC_DIR + path):
		staticHandler(ctx)

	}
}

func main() {
	var (
		addr string
		port int
	)
	flag.StringVar(
		&addr,
		"addr",
		"",
		"leave blank for 0.0.0.0",
	)
	flag.IntVar(
		&port,
		"port",
		8080,
		"",
	)
	flag.Parse()

	var err error
	db, err = sql.Open("mysql", "root@/caffochoco")
	checkErr(err)
	defer db.Close()

	listenAddr := fmt.Sprintf("%s:%d", addr, port)
	log.Println("listening on", listenAddr)
	log.Fatalln(fasthttp.ListenAndServe(listenAddr, requestHandler))
}

func checkErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func fileExists(filename string) bool {
	f, err := os.Open(filename)
	f.Close()
	if os.IsNotExist(err) {
		return false
	}
	checkErr(err)
	return true
}

func renderTemplate(filename string, data interface{}) string {
	t, err := template.ParseFiles(filename)
	checkErr(err)
	buf := new(bytes.Buffer)
	checkErr(t.Execute(buf, data))
	return buf.String()
}
