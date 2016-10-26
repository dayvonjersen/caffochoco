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
		addr     string
		port     int
		nossl    bool
		certFile string
		keyFile  string
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
		443,
		"",
	)
	flag.BoolVar(
		&nocache,
		"nocache",
		false,
		"disable caching of rendered content e.g. blog posts",
	)
	flag.BoolVar(
		&nossl,
		"nossl",
		false,
		"disable https",
	)
	flag.StringVar(
		&certFile,
		"cert",
		"cert.pem",
		"path to cert",
	)
	flag.StringVar(
		&keyFile,
		"key",
		"key.pem",
		"path to key",
	)
	flag.Parse()

	listenAddr := fmt.Sprintf("%s:%d", addr, port)
	if nossl {
		log.Println("listening on", listenAddr)
		log.Fatalln(fasthttp.ListenAndServe(listenAddr, requestHandler))
	}

	// redirect http traffic to https
	log.Println("listening on", addr+":80")
	go fasthttp.ListenAndServe(addr+":80", func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.SetStatusCode(fasthttp.StatusTemporaryRedirect)
		ctx.Response.Header.Set(
			"Location",
			fmt.Sprintf("https://%s:%d%s",
				string(ctx.Request.Host()),
				port,
				string(ctx.Path()),
			),
		)
	})

	log.Println("listening on", listenAddr)
	log.Fatalln(fasthttp.ListenAndServeTLS(listenAddr, certFile, keyFile, requestHandler))
}
