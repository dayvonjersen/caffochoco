package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
)

const (
	STATIC_DIR string = "./public"
)

var (
	db            *sql.DB
	staticHandler = fasthttp.FSHandler(STATIC_DIR, 0)
)

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
		fasthttp.ServeFile(ctx, "./index.html")
	case strings.HasPrefix(path, "/api/"):
		apiHandler(ctx)
	// imageHandler
	// audioHandler
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