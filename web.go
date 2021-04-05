package main

import (
	"log"
	"net/http"
	"strings"
)

func init() {
	fs := http.FileServer(http.Dir("."))
	http.HandleFunc("/ui/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		log.Printf("Serving: %s", req.URL.String())
		fs.ServeHTTP(resp, req)
	})
	// http.HandleFunc("/simtbc", simTBCPage)
}
