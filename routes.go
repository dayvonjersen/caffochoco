package main

import (
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

// static assets

const STATIC_DIR = "./public"

var staticHandler = fasthttp.FSHandler(STATIC_DIR, 0)

// errors

func errorHandler(ctx *fasthttp.RequestCtx, statusCode int, data interface{}) {
	ctx.Response.SetStatusCode(statusCode)
	ctx.Response.Header.Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprintf(ctx, renderTemplate("./templates/error_"+strconv.Itoa(statusCode)+".html", data))
}
func notfoundHandler(ctx *fasthttp.RequestCtx) {
	errorHandler(ctx, 404, nil)
}

var routes = map[string]fasthttp.RequestHandler{
	"/api":   apiHandler,
	"/image": imageHandler,
	"/blog":  blogHandler,
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
	defer func() {
		if x := recover(); x != nil {
			errorHandler(ctx, 500, struct {
				Message, Stack string
			}{x.(string), getStack(debug.Stack())})
		}
	}()

	path := string(ctx.Path())

	switch {
	case path == "/", path == "/index.html":
		fasthttp.ServeFile(ctx, "./templates/index.html")
	case router(path, ctx):
	case fileExists(STATIC_DIR + path):
		staticHandler(ctx)
	default:
		notfoundHandler(ctx)
	}
}
