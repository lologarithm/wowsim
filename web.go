package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lologarithm/wowsim/tbc"
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
	http.HandleFunc("/simtbc", simTBCPage)
}

func simTBCPage(w http.ResponseWriter, r *http.Request) {
	fileData, err := ioutil.ReadFile("tbc/ui/index.html")
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	if r.ContentLength > 0 {
		// parse form.
		r.ParseForm()
		intv, _ := strconv.Atoi(r.FormValue("int"))
		sph, _ := strconv.ParseFloat(r.FormValue("spellhit"), 64)
		spc, _ := strconv.ParseFloat(r.FormValue("spellcrit"), 64)
		spd, _ := strconv.ParseFloat(r.FormValue("spelldmg"), 64)
		mp5, _ := strconv.ParseFloat(r.FormValue("mp5"), 64)
		haste, _ := strconv.ParseFloat(r.FormValue("haste"), 64)
		spp, _ := strconv.ParseFloat(r.FormValue("spellpen"), 64)

		stats := tbc.Stats{
			tbc.StatInt:       float64(intv) + 86, // Add base stats
			tbc.StatSpellCrit: spc + 151,          // Add base+talents to gear
			tbc.StatSpellHit:  sph/100 + 0.03,     // Add talent hit
			tbc.StatSpellDmg:  spd,                // gear
			tbc.StatMP5:       mp5,                // gear
			tbc.StatHaste:     haste,
			tbc.StatMana:      1240, // Base Mana L60 Troll Shaman
			tbc.StatSpellPen:  spp,  //
		}
		stats[tbc.StatSpellCrit] += (stats[tbc.StatInt] / 59.5) / 100 // 1% crit per 59.5 int
		stats[tbc.StatMana] += stats[tbc.StatInt] * 15

		stats.Print()

		results := runTBCSim(stats, 300, 500)
		fileData = append(fileData, "<pre>"...)
		for _, res := range results {
			fileData = append(fileData, res...)
			fileData = append(fileData, "\n"...)
		}
		fileData = append(fileData, "</pre>"...)
	}

	w.Write(fileData)
}
