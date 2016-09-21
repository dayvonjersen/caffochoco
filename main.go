package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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

func requestHandler(ctx *fasthttp.RequestCtx) {
	defer func() {
		log.Println("->", string(ctx.Method()), string(ctx.Path()))
		s := ctx.Response.StatusCode()
		log.Println("<-", s, fasthttp.StatusMessage(s))
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
