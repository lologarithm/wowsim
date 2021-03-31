package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"gitlab.com/lologarithm/wowsim/tbc"
)

func main() {
	c := make(chan struct{}, 0)

	simfunc := js.FuncOf(Simulate)
	gearfunc := js.FuncOf(GearStats)
	gearlistfunc := js.FuncOf(GearList)

	js.Global().Set("tbcsim", simfunc)
	js.Global().Set("gearstats", gearfunc)
	js.Global().Set("gearlist", gearlistfunc)
	js.Global().Call("popgear")
	<-c
}

// GearList reports all items of gear to the UI to display.
func GearList(this js.Value, args []js.Value) interface{} {
	slot := -1

	if len(args) == 1 {
		slot = args[0].Int()
	}
	gears := "["
	for _, v := range tbc.ItemLookup {
		if slot != -1 && v.Slot != slot {
			continue
		}
		if len(gears) != 1 {
			gears += ","
		}
		gears += `{"name":"` + v.Name + `", "slot": ` + strconv.Itoa(v.Slot) + `}`
	}
	gears += "]"
	return gears
}

// GearStats takes a gear list and returns their total stats.
// This could power a simple 'current stats of all gear' UI.
func GearStats(this js.Value, args []js.Value) interface{} {
	return getGear(args[0]).Stats().Print()
}

// getGear converts js string array to a list of equipment items.
func getGear(val js.Value) tbc.Equipment {
	numGear := val.Length()
	gearstr := make([]string, numGear)
	for i := range gearstr {
		gearstr[i] = val.Index(i).String()
	}
	return tbc.NewEquipmentSet(gearstr...)
}

// Simulate takes in number of iterations and a gear list.
func Simulate(this js.Value, args []js.Value) interface{} {
	st := time.Now()

	// TODO: Accept talents, buffs, and consumes as inputs.

	if len(args) != 2 {
		return `{"error": "invalid arguments supplied"}`
	}

	gear := getGear(args[1])
	gearStats := gear.Stats()

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

	simi := args[0].Int()
	if simi == 1 {
		tbc.IsDebug = true
	}
	results := runTBCSim(stats, 120, simi)
	// for _, res := range results {
	// 	fmt.Printf("\n%s\n", res)
	// }
	fmt.Printf("Sim Took: %v\n", time.Now().Sub(st).String())

	jsstr := `{"results": [`
	for i, rslt := range results {
		jsstr += "\"" + strings.ReplaceAll(rslt, "\n", "\\n") + "\""
		if i != len(results)-1 {
			jsstr += ","
		}
	}
	jsstr += `], "stats": ` + stats.Print() + "}"
	return jsstr
}

func runTBCSim(stats tbc.Stats, seconds int, numSims int) []string {
	print("\nSim Duration:", seconds)
	print("\nNum Simulations: ", numSims)
	print("\n")
	spellOrders := [][]string{
		{"CL6", "LB12", "LB12", "LB12"},
		{"CL6", "LB12", "LB12", "LB12", "LB12"},
		// {"pri", "CL6", "LB12"}, // cast CL whenever off CD, otherwise LB
		{"LB12"}, // only LB
	}

	results := []string{}

	for _, spells := range spellOrders {
		simDmgs := []float64{}
		simOOMs := []int{}

		rseed := time.Now().Unix()
		sim := tbc.NewSim(stats, spells, rseed)
		for ns := 0; ns < numSims; ns++ {
			metrics := sim.Run(seconds)
			simDmgs = append(simDmgs, metrics.TotalDamage)
			simOOMs = append(simOOMs, metrics.OOMAt)
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
		output += fmt.Sprintf("Spell Order: %v\n", spells)
		output += "DPS:\n"
		output += fmt.Sprintf("Mean: %0.1f\n", (mean / float64(seconds)))
		output += fmt.Sprintf("Max: %0.1f\n", (max / float64(seconds)))
		output += fmt.Sprintf("Std.Dev: %0.1f\n", stdev/float64(seconds))
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
		output += "Went OOM: " + strconv.Itoa(numOoms) + "/" + strconv.Itoa(numSims) + " sims\n"
		if numOoms > 0 {
			output += fmt.Sprintf("Avg OOM Time: %d seconds\n", avg)
		}
		results = append(results, output)
	}

	return results

	// fmt.Printf("Casts: \n")
	// for k, v := range castStats {
	// 	fmt.Printf("\t%s: %d\n", k, v.Num)
	// }
}
