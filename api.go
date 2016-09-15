package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

func apiHandler(ctx *fasthttp.RequestCtx) {
	path := strings.TrimPrefix(string(ctx.Path()), "/api")

	ctx.Response.Header.Set("Content-Type", "application/json")
	switch path {
	case "/":
		stuff := map[string]interface{}{
			"releases": fetchAll(`SELECT * FROM releases`),
			"sections": fetchAll(`SELECT * FROM sections`),
		}

		b, err := json.Marshal(stuff)
		checkErr(err)
		fmt.Fprintf(ctx, string(b))
	default:
		ctx.Response.SetStatusCode(404)
		fmt.Fprintf(ctx, `{"error":"404"}`)
	}
}
