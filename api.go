package main

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

func apiHandler(ctx *fasthttp.RequestCtx) {
	path := strings.TrimPrefix(string(ctx.Path()), "/api")

	ctx.Response.Header.Set("Content-Type", "application/json")
	switch path {
	case "/":
		fmt.Fprintf(ctx, `{"hello":"world"}`)
	default:
		ctx.Response.SetStatusCode(404)
		fmt.Fprintf(ctx, `{"error":"404"}`)
	}
}
