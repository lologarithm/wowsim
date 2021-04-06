package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/lologarithm/wowsim/tbc"
)

// This was used to parse a CSV from an item spreadsheet.
// Each spreadsheet has its own structure so we probably need to write a custom one for each...
//  We can try to write an 'auto parser' but some sheets have such different structure this could be hard.

func main() {
	data, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		panic(err)
	}

	r := csv.NewReader(strings.NewReader(string(data)))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	slot := tbc.EquipHead
	// items := []tbc.Item{}
	colHeader := []string{}
	for i, v := range records {
		if i == 0 {
			colHeader = v
		}
		if v[0] == colHeader[0] && v[1] == colHeader[1] {
			// Another header row...
			continue
		}
		if v[0] != "" && v[2] == "" {
			// Change slot
			switch v[0] {
			case "Helm", "Head":
				slot = tbc.EquipHead
			case "Neck":
				slot = tbc.EquipNeck
			case "Shoulder":
				slot = tbc.EquipShoulder
			case "Back":
				slot = tbc.EquipBack
			case "Chest":
				slot = tbc.EquipChest
			case "Bracer", "Wrist":
				slot = tbc.EquipWrist
			case "Hands":
				slot = tbc.EquipHands
			case "Belt", "Waist":
				slot = tbc.EquipWaist
			case "Boots", "Feet":
				slot = tbc.EquipFeet
			case "Legs":
				slot = tbc.EquipLegs
			case "Ring 1":
				slot = tbc.EquipFinger
			case "Trinket 1":
				slot = tbc.EquipTrinket
			case "Totem":
				slot = tbc.EquipTotem
			case "MH", "2H":
				slot = tbc.EquipWeapon
			case "OH", "Shield", "OH / Shield":
				slot = tbc.EquipOffhand
			}
			continue
		}
		if v[1] != "" && v[1] != "Name" {
			stm, _ := strconv.ParseFloat(v[4], 64)
			intv, _ := strconv.ParseFloat(v[5], 64)
			sph, _ := strconv.ParseFloat(v[9], 64)
			spc, _ := strconv.ParseFloat(v[7], 64)
			spd, _ := strconv.ParseFloat(v[6], 64)
			mp5, _ := strconv.ParseFloat(v[10], 64)
			haste, _ := strconv.ParseFloat(v[8], 64)
			// spp, _ := strconv.ParseFloat(v[], 64)

			numMeta, _ := strconv.Atoi(v[12])
			numRed, _ := strconv.Atoi(v[13])
			numYellow, _ := strconv.Atoi(v[14])
			numBlue, _ := strconv.Atoi(v[15])

			if len(v[15]) > 0 {
				// parse out socket bonus
				tokens := strings.Split(v[15], " ")
				val := strconv.Atoi()
				if len(tokens) == 2 {

				} else if len(tokens) == 3 {

				}
			}
			it := tbc.Item{
				Name:       v[1],
				SourceZone: v[2],
				Slot:       slot,
				Stats: tbc.Stats{
					tbc.StatStm:       stm,
					tbc.StatInt:       intv,
					tbc.StatSpellCrit: spc,
					tbc.StatSpellHit:  sph,
					tbc.StatSpellDmg:  spd,
					tbc.StatHaste:     haste,
					tbc.StatMP5:       mp5,
				},
				GemSlots: make([]tbc.GemColor, numMeta+numRed+numYellow+numBlue),
			}
			if numMeta > 0 { // its always 1 or 0
				it.GemSlots = append(it.GemSlots, tbc.GemColorMeta)
			}
			for i := 0; i < numRed; i++ {
				it.GemSlots = append(it.GemSlots, tbc.GemColorRed)
			}
			for i := 0; i < numYellow; i++ {
				it.GemSlots = append(it.GemSlots, tbc.GemColorYellow)
			}
			for i := 0; i < numBlue; i++ {
				it.GemSlots = append(it.GemSlots, tbc.GemColorBlue)
			}

			fmt.Fprintf(os.Stdout, "%#v,\n", it)

		}
	}
}
