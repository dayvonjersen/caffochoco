package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

// static assets

const STATIC_DIR = "./public"

var staticHandler = http.StripPrefix("/", http.FileServer(http.Dir(STATIC_DIR)))

// capture http response code sent because net/http is """good"""
type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// errors

func errorHandler(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, renderTemplate("./templates/error_"+strconv.Itoa(statusCode)+".html", data))
}
func notfoundHandler(w http.ResponseWriter, r *http.Request) {
	errorHandler(w, r, 404, nil)
}

var routes = map[string]func(w http.ResponseWriter, r *http.Request){
	"/api":   apiHandler,
	"/image": imageHandler,
	"/blog":  blogHandler,
}

func router(path string, w http.ResponseWriter, r *http.Request) bool {

	for route, handler := range routes {
		if strings.HasPrefix(path, route) {
			handler(w, r)
			return true
		}
	}
	return false
}

func requestHandler(ww http.ResponseWriter, r *http.Request) {
	w := &logResponseWriter{ww, http.StatusOK}
	defer func() {
		log.Println(req(r))
		s := w.statusCode
		log.Println("<-", s, http.StatusText(s))
	}()
	defer func() {
		if x := recover(); x != nil {
			errorHandler(w, r, 500, struct {
				Message, Stack string
			}{x.(string), getStack(debug.Stack())})
		}
	}()

	path := r.URL.Path

	switch {
	case path == "/", path == "/index.html":
		http.ServeFile(w, r, "./templates/index.html")
	case router(path, w, r):
	case fileExists(STATIC_DIR + path):
		staticHandler.ServeHTTP(w, r)
	default:
		notfoundHandler(w, r)
	}
}

func req(r *http.Request) string {
	return fmt.Sprint(
		hostAddr(r), " <-> ", remoteAddr(r), "\n",
		r.Header.Get("User-Agent"), "\n",
		" -> ", r.Method, " ", r.URL, "\n",
	)
}

func hostAddr(r *http.Request) string {
	host := r.Header.Get("Host")
	if host != "" {
		return host
	}
	return r.Host
}

func remoteAddr(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		if strings.Contains(ip, ",") {
			return ip[:strings.LastIndex(ip, ",")]
		}
		return ip
	}
	return r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
}
