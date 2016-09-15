package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

type apiRoute struct {
	re *regexp.Regexp
	fn func([]string) map[string]interface{}
}

var apiRoutes = []*apiRoute{
	&apiRoute{re: regexp.MustCompile(`^/$`), fn: index},
	&apiRoute{re: regexp.MustCompile(`^/release/(\d+)$`), fn: release},
}

func apiHandler(ctx *fasthttp.RequestCtx) {
	path := strings.TrimPrefix(string(ctx.Path()), "/api")

	ctx.Response.Header.Set("Content-Type", "application/json")

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

	b, err := json.Marshal(ret)
	checkErr(err)
	fmt.Fprintf(ctx, string(b))
}

func index(params []string) map[string]interface{} {
	stuff := map[string]interface{}{
		"releases": fetchAll(`SELECT * FROM releases`),
		"sections": fetchAll(`SELECT * FROM sections`),
	}
	return stuff
}
func release(params []string) map[string]interface{} {
	if len(params) != 1 {
		return nil
	}
	id, err := strconv.Atoi(params[0])
	if err != nil {
		return nil
	}
	stuff := fetchAll(fmt.Sprintf("SELECT * FROM releases WHERE id = %d", id))
	if len(stuff) > 0 {
		return map[string]interface{}(stuff[0])
	}
	return nil
}
