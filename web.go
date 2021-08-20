package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/lologarithm/wowsim/tbc"
)

func init() {
	fs := http.FileServer(http.Dir("."))
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		// if strings.HasSuffix(req.URL.Path, "ui.js") {
		// 	this.worker = new window.Worker('simworker.js');
		// }
		log.Printf("Serving: %s", req.URL.String())
		fs.ServeHTTP(resp, req)
	})
	http.HandleFunc("/api", handleAPI)
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	// Assumes input is a JSON object as a string
	request := tbc.ApiRequest{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		panic("Error parsing request: " + err.Error() + "\nRequest: " + string(body))
	}
	result := tbc.ApiCall(request)
	resultData, err := json.Marshal(result)
	if err != nil {
		panic("Error marshaling result: " + err.Error())
	}
	w.Write(resultData)
}
