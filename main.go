package main

import (
	_ "expvar"
	"log"
	"net/http"
	"net/http/httputil"
	_ "net/http/pprof"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/handlers"
)

func main() {
	err := http.ListenAndServe(":80", handlers.CombinedLoggingHandler(os.Stdout, newLimiter()))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type limiter struct {
	root  http.Handler
	index http.Handler
}

func newLimiter() *limiter {
	return &limiter{
		root:  httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "151.101.194.2"}),
		index: httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "f2.shared.global.fastly.net"}),
	}
}

var counter = 0

func (l limiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Host == "index.rubygems.org" && strings.Contains(r.RequestURI, "/api/v1/dependencies") {
		counter = counter + 1
		if counter > 5 {
			w.WriteHeader(429)
		} else {
			l.index.ServeHTTP(w, r)
		}
	} else if r.Host == "index.rubygems.org" {
		l.index.ServeHTTP(w, r)
	} else if r.Host == "rubygems.org" {
		l.root.ServeHTTP(w, r)
	} else {
		http.DefaultServeMux.ServeHTTP(w, r)
	}
}
