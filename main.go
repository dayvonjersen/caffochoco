package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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
		http.HandleFunc("/", requestHandler)
		log.Fatalln(http.ListenAndServe(listenAddr, nil))
	}

	// redirect http traffic to https
	log.Println("listening on", addr+":80")
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(
				"Location",
				fmt.Sprintf("https://%s:%d%s",
					hostAddr(r),
					port,
					r.URL,
				),
			)
			w.WriteHeader(http.StatusTemporaryRedirect)
			fmt.Println("hello??")
		})

		s := &http.Server{
			Addr:    ":80",
			Handler: mux,
		}
		s.ListenAndServe()
	}()

	log.Println("listening on", listenAddr)
	http.HandleFunc("/", requestHandler)
	log.Fatalln(http.ListenAndServeTLS(listenAddr, certFile, keyFile, nil))
}
