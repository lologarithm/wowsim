package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/lologarithm/wowsim/tbc"
)

//go:embed ui
var uifs embed.FS

func main() {
	var useFS = flag.Bool("usefs", false, "Use local file system and wasm. Set to true for dev")
	var host = flag.String("host", ":3333", "URL to host the interface on.")

	flag.Parse()

	var fs http.Handler
	if *useFS {
		log.Printf("Using local file system for development.")
		fs = http.FileServer(http.Dir("."))
	} else {
		log.Printf("Embedded file server running.")
		fs = http.FileServer(http.FS(uifs))
	}

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		} else if strings.HasSuffix(req.URL.Path, "/ui.js") {
			var uijs []byte
			var err error
			if *useFS {
				// read file straight off disk
				uijs, err = ioutil.ReadFile("./ui/ui.js")
			} else {
				uijs, err = uifs.ReadFile("ui/ui.js")
				// modify so that simworker is replaced with networker.
				uijs = bytes.Replace(uijs, []byte(`this.worker = new window.Worker('simworker.js');`), []byte(`this.worker = new window.Worker('networker.js');`), 1)
			}
			if err != nil {
				log.Printf("Failed to open file..., %s", err)
			}
			resp.Write(uijs)
			return
		}
		fs.ServeHTTP(resp, req)
	})

	http.HandleFunc("/api", handleAPI)

	go func() {
		log.Printf("Launching interface on http://localhost%s/ui", *host)
		log.Printf("Closing: %s", http.ListenAndServe(*host, nil))
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if len(text) == 0 {
			continue
		}
		switch strings.TrimSpace(text) {
		case "profile":
			go func() {
				filename := fmt.Sprintf("profile_%d.cpu", time.Now().Unix())
				f, err := os.Create(filename)
				if err != nil {
					log.Fatal("could not create CPU profile: ", err)
				}
				if err := pprof.StartCPUProfile(f); err != nil {
					log.Fatal("could not start CPU profile: ", err)
				}
				fmt.Printf("Running profiling for 15 seconds, output to %s\n", filename)
				time.Sleep(time.Second * 15)
				pprof.StopCPUProfile()
				f.Close()
				fmt.Printf("Profiling complete.")
			}()
		case "quit":
			os.Exit(1)
		default:
			fmt.Printf("Unknown command: '%s'", text)
		}
	}
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	st := time.Now()
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

	log.Printf("API request took %v", time.Since(st))
}
