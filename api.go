package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type apiRoute struct {
	re *regexp.Regexp
	fn func([]string) map[string]interface{}
}

var apiRoutes = []*apiRoute{
	&apiRoute{re: regexp.MustCompile(`^/blog`), fn: blog},
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var ret map[string]interface{}
	for _, route := range apiRoutes {
		if route.re.MatchString(path) {
			ret = route.fn(route.re.FindStringSubmatch(path)[1:])
			break
		}
	}
	if ret == nil {
		w.WriteHeader(404)
		ret = map[string]interface{}{"error": 404}
	}

	b, err := json.MarshalIndent(ret, "", "    ")
	checkErr(err)
	fmt.Fprintf(w, string(b))
}

func blog(params []string) map[string]interface{} {
	return map[string]interface{}{"blog": getTocData()}
}
