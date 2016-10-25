package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/valyala/fasthttp"
)

type apiRoute struct {
	re *regexp.Regexp
	fn func([]string) map[string]interface{}
}

var apiRoutes = []*apiRoute{
	&apiRoute{re: regexp.MustCompile(`^/blog`), fn: blog},
}

func apiHandler(ctx *fasthttp.RequestCtx) {
	path := strings.TrimPrefix(string(ctx.Path()), "/api")

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")

	var ret map[string]interface{}
	for _, route := range apiRoutes {
		if route.re.MatchString(path) {
			ret = route.fn(route.re.FindStringSubmatch(path)[1:])
			break
		}
	}
	if ret == nil {
		ctx.Response.SetStatusCode(404)
		ret = map[string]interface{}{"error": 404}
	}

	b, err := json.MarshalIndent(ret, "", "    ")
	checkErr(err)
	fmt.Fprintf(ctx, string(b))
}

func blog(params []string) map[string]interface{} {
	return map[string]interface{}{"blog": getTocData()}
}
