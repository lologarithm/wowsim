package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.com/lologarithm/wowsim/elesim"
	"gitlab.com/lologarithm/wowsim/tbc"
)

// /script print(GetSpellBonusDamage(4))

// Buffs
//  - Additive
//  - Multiplicative
//
// 'Aura' Effects
//  - On Hitting
//  - On Cast
//	- On Being Hit
//  - Always
//
// Applies to
//  - Specific Spell
//  - All Spells
//
// Modifiers
//  - Dmg
//  - Cast Time
//  - Hit Chance
//  - Mana Cost

func main() {
	var isDebug = flag.Bool("debug", false, "Include --debug to spew the entire simulation log.")
	var runWebUI = flag.Bool("web", false, "Use to run sim in web interface instead of in terminal")
	var useTBC = flag.Bool("tbc", false, "Use to run sim using TBC stats / effects")
	flag.Parse()

	elesim.IsDebug = *isDebug
	tbc.IsDebug = *isDebug

	if *runWebUI {
		log.Printf("Closing: %s", http.ListenAndServe(":3333", nil))
	}

	if *useTBC {
		gear := tbc.NewEquipmentSet(
			"Shamanistic Helmet of Second Sight",
			"Brooch of Heightened Potential",
			"Pauldrons of Wild Magic",
			"Ogre Slayer's Cover",
			"Tidefury Chestpiece",
			"World's End Bracers",
			"Earth Mantle Handwraps",
			"Wave-Song Girdle",
			"Stormsong Kilt",
			"Magma Plume Boots",
			"Cobalt Band of Tyrigosa",
			"Scintillating Coral Band",
			"Totem of the Void",
			"Khadgar's Knapsack",
			"Bleeding Hollow Warhammer",
		)

		gearStats := gear.Stats()
		fmt.Printf("Gear Stats:\n")
		gearStats.Print()

		stats := tbc.Stats{
			tbc.StatInt:       104,         // Base
			tbc.StatSpellCrit: 48.62 + 243, // Base + Talents + Totem
			tbc.StatSpellHit:  151.2,       // Totem + Talents
			tbc.StatSpellDmg:  101,         // Totem
			tbc.StatMP5:       50 + 25,     // Water Shield + Mana Stream
			tbc.StatMana:      2958,        // level 70 shaman
			tbc.StatHaste:     0,
			tbc.StatSpellPen:  0,
		}

		buffs := tbc.Stats{
			tbc.StatInt:       40, //arcane int
			tbc.StatSpellCrit: 0,
			tbc.StatSpellHit:  0,
			tbc.StatSpellDmg:  42, // sup wiz oil
			tbc.StatMP5:       0,
			tbc.StatMana:      0,
			tbc.StatHaste:     0,
			tbc.StatSpellPen:  0,
		}

		for i, v := range buffs {
			stats[i] += v
			stats[i] += gearStats[i]
		}

		stats[tbc.StatInt] *= 1.1 // blessing of kings

		stats[tbc.StatSpellCrit] += (stats[tbc.StatInt] / 80) / 100 // 1% crit per 59.5 int
		stats[tbc.StatMana] += stats[tbc.StatInt] * 15

		fmt.Printf("Final Stats:\n")
		stats.Print()
		sims := 10000
		if *isDebug {
			sims = 1
		}
		results := runTBCSim(stats, 120, sims)
		for _, res := range results {
			fmt.Printf("\n%s\n", res)
		}
	} else {
		stats := elesim.Stats{
			elesim.StatInt:       86 + 191,           // base + gear
			elesim.StatSpellCrit: 0.022 + 0.04 + .11, // base crit + gear + talents
			elesim.StatSpellHit:  0.03 + 0.03,        // talents + gear
			elesim.StatSpellDmg:  474,                // gear
			elesim.StatMP5:       33,                 //
			elesim.StatMana:      1240,
			elesim.StatSpellPen:  0,
		}

		buffs := elesim.Stats{
			elesim.StatSpellCrit: 0.18, // world buffs DMT+Ony+Songflower
			elesim.StatInt:       15,   // songflower
		}
		buffs[elesim.StatInt] += 31 // arcane brill
		buffs[elesim.StatInt] += 12 // GOTW
		// buffs[elesim.StatInt] += 10 // runn tum tuber

		for i, v := range buffs {
			// ZG Buff
			if elesim.Stat(i) == elesim.StatInt {
				stats[i] = stats[i] * 1.15
			}

			// I believe ZG buff applies before other buffs
			stats[i] += v
		}
		stats[elesim.StatSpellCrit] += (stats[elesim.StatInt] / 59.5) / 100 // 1% crit per 59.5 int
		stats[elesim.StatMana] += stats[elesim.StatInt] * 15
		stats.Print()

		// hardcode 120s 200 sims
		results := runSim(stats, 120, 500)
		for _, res := range results {
			fmt.Printf("\n%s\n", res)
		}
	}

}

func runSim(stats elesim.Stats, seconds int, numSims int) []string {
	fmt.Printf("\nSim Duration: %d sec\nNum Simulations: %d\n", seconds, numSims)
	spellOrders := [][]string{
		// {"CL4", "LB10", "LB10", "LB10"},
		{"LB10"},
		// {"LB10", "LB4"},
	}

	statchan := make(chan string, 3)
	for _, spells := range spellOrders {
		go func(spo []string) {
			simDmgs := []float64{}
			simOOMs := []int{}

			for ns := 0; ns < numSims; ns++ {
				dmg, oomat := elesim.Sim(seconds, stats, spo)
				simDmgs = append(simDmgs, dmg)
				simOOMs = append(simOOMs, oomat)
			}

			totalDmg := 0.0
			tdSq := totalDmg
			max := 0.0
			for _, dmg := range simDmgs {
				totalDmg += dmg
				tdSq += dmg * dmg

				if dmg > max {
					max = dmg
				}
			}

			meanSq := tdSq / float64(numSims)
			mean := totalDmg / float64(numSims)
			stdev := math.Sqrt(meanSq - mean*mean)

			output := ""
			output += fmt.Sprintf("Spell Order: %v\n", spo)
			output += fmt.Sprintf("DPS:")
			output += fmt.Sprintf("\tMean: %0.1f\n", (mean / float64(seconds)))
			output += fmt.Sprintf("\tMax: %0.1f\n", (max / float64(seconds)))
			output += fmt.Sprintf("\tStd.Dev: %0.1f\n", stdev/float64(seconds))

			ooms := 0
			numOoms := 0
			for _, oa := range simOOMs {
				ooms += oa
				if oa > 0 {
					numOoms++
				}
			}

			avg := 0
			if numOoms > 0 {
				avg = ooms / numOoms
			}
			output += fmt.Sprintf("Went OOM: %d/%d sims\n", numOoms, numSims)
			if numOoms > 0 {
				output += fmt.Sprintf("Avg OOM Time: %d seconds\n", avg)
			}
			statchan <- output
		}(spells)
	}

	results := []string{}
	for i := 0; i < len(spellOrders); i++ {
		results = append(results, <-statchan)
	}

	return results

	// fmt.Printf("Casts: \n")
	// for k, v := range castStats {
	// 	fmt.Printf("\t%s: %d\n", k, v.Num)
	// }
}

func runTBCSim(stats tbc.Stats, seconds int, numSims int) []string {
	fmt.Printf("\nSim Duration: %d sec\nNum Simulations: %d\n", seconds, numSims)
	spellOrders := [][]string{
		{"CL6", "LB12", "LB12", "LB12"},
		{"CL6", "LB12", "LB12", "LB12", "LB12"},
		// {"pri", "CL6", "LB12"}, // cast CL whenever off CD, otherwise LB
		{"LB12"}, // only LB
	}

	statchan := make(chan string, 3)
	for _, spells := range spellOrders {
		go func(spo []string) {
			simDmgs := []float64{}
			simOOMs := []int{}
			histogram := map[int]int{}

			rseed := time.Now().Unix()
			sim := tbc.NewSim(stats, spo, rseed)
			for ns := 0; ns < numSims; ns++ {
				metrics := sim.Run(seconds)
				simDmgs = append(simDmgs, metrics.TotalDamage)
				simOOMs = append(simOOMs, metrics.OOMAt)

				rv := int(math.Round(math.Round(metrics.TotalDamage/float64(seconds))/10) * 10)
				histogram[rv] += 1
			}

			out := ""
			for k, v := range histogram {
				out += strconv.Itoa(k) + "," + strconv.Itoa(v) + "\n"
			}
			ioutil.WriteFile(strings.Join(spo, ""), []byte(out), 0666)

			totalDmg := 0.0
			tdSq := totalDmg
			max := 0.0
			for _, dmg := range simDmgs {
				totalDmg += dmg
				tdSq += dmg * dmg

				if dmg > max {
					max = dmg
				}
			}

			meanSq := tdSq / float64(numSims)
			mean := totalDmg / float64(numSims)
			stdev := math.Sqrt(meanSq - mean*mean)

			// mean - dev*conf, mean + dev*conf
			// Z
			// 80% 	1.282
			// 85% 	1.440
			// 90% 	1.645
			// 95% 	1.960
			// 99% 	2.576
			// 99.5% 	2.807
			// 99.9% 	3.291

			output := ""
			output += fmt.Sprintf("Spell Order: %v\n", spo)
			output += fmt.Sprintf("DPS:")
			output += fmt.Sprintf("\tMean: %0.1f\n", (mean / float64(seconds)))
			output += fmt.Sprintf("\tMax: %0.1f\n", (max / float64(seconds)))
			output += fmt.Sprintf("\tStd.Dev: %0.1f\n", stdev/float64(seconds))

			ooms := 0
			numOoms := 0
			for _, oa := range simOOMs {
				ooms += oa
				if oa > 0 {
					numOoms++
				}
			}

			avg := 0
			if numOoms > 0 {
				avg = ooms / numOoms
			}
			output += fmt.Sprintf("Went OOM: %d/%d sims\n", numOoms, numSims)
			if numOoms > 0 {
				output += fmt.Sprintf("Avg OOM Time: %d seconds\n", avg)
			}
			statchan <- output
		}(spells)
		if tbc.IsDebug {
			time.Sleep(time.Second)
		}
	}

	results := []string{}
	for i := 0; i < len(spellOrders); i++ {
		results = append(results, <-statchan)
	}

	return results

	// fmt.Printf("Casts: \n")
	// for k, v := range castStats {
	// 	fmt.Printf("\t%s: %d\n", k, v.Num)
	// }
}
