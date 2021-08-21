package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/lologarithm/wowsim/tbc"
)

//go:embed ui
var uifs embed.FS

func main() {
	var useFS = flag.Bool("usefs", false, "Use embedded file system and server. Set to false for dev")
	flag.Parse()

	var fs http.Handler
	if *useFS {
		log.Printf("Using local file system for development.")
		fs = http.FileServer(http.Dir("."))
	} else {
		log.Printf("Embedded Server running.")
		fs = http.FileServer(http.FS(uifs))

		dir, err := uifs.ReadDir("ui/icons")
		if err != nil {

		}
		for _, v := range dir {
			log.Printf("%s (%s) (%v)", v.Name(), v.Type(), v.IsDir())
		}
	}

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		} else if strings.HasSuffix(req.URL.Path, "/ui.js") {
			var uijs []byte
			var err error
			if *useFS {
				uijs, err = ioutil.ReadFile("./ui/ui.js")
			} else {
				uijs, err = uifs.ReadFile("ui/ui.js")
				uijs = bytes.Replace(uijs, []byte(`this.worker = new window.Worker('simworker.js');`), []byte(`this.worker = new window.Worker('networker.js');`), 1)
			}
			if err != nil {
				log.Printf("Failed to open file..., %s", err)
				// log.Printf("FS: %s", uifs.ReadDir("."))
			}
			resp.Write(uijs)
			return
		}
		log.Printf("Serving: %s", req.URL.String())
		fs.ServeHTTP(resp, req)
	})

	http.HandleFunc("/api", handleAPI)

	log.Printf("Closing: %s", http.ListenAndServe(":3333", nil))
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
