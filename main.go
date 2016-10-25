package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

var (
	nocache bool
)

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
	flag.BoolVar(
		&nocache,
		"nocache",
		false,
		"disable caching of rendered content e.g. blog posts",
	)
	flag.Parse()

	listenAddr := fmt.Sprintf("%s:%d", addr, port)
	log.Println("listening on", listenAddr)
	log.Fatalln(fasthttp.ListenAndServe(listenAddr, requestHandler))
}
