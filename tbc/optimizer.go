package tbc

import (
	"fmt"
	"math"
	"strings"
)

var rotations = [][]string{
	{"CL6", "LB12", "LB12", "LB12"},
	{"CL6", "LB12", "LB12", "LB12", "LB12"},
	{"LB12"}, // only LB
	// {"pri", "CL6", "LB12"}, // cast CL whenever off CD, otherwise LB
}

type OptimalGemsResult struct {
	Sims []EquipmentResult
}

type EquipmentResult struct {
	Results []SimMetrics
	Equip   Equipment
}

func StatWeights(opts Options, equip Equipment, seconds int, numSims int) []float64 {
	type res struct {
		s   Stat
		val float64
		m   []SimMetrics
		d   float64
	}
	results := make(chan res, 10)
	base := 0.0

	doStat := func(mod Stat, value float64) {
		myopts := opts
		optClone2 := Stats{StatLen: 0}
		copy(optClone2, myopts.Buffs.Custom) // clone existing buff
		myopts.Buffs.Custom = optClone2
		myopts.Buffs.Custom[mod] += value
		stats := CalculateTotalStats(myopts, equip)

		myOpts := myopts
		myOpts.AgentType = AGENT_TYPE_ADAPTIVE
		simdmg := 0.0
		sim := NewSim(stats, equip, myOpts)
		simmet := make([]SimMetrics, 0, numSims)
		for ns := 0; ns < numSims; ns++ {
			metrics := sim.Run(seconds)
			simdmg += metrics.TotalDamage
			simmet = append(simmet, metrics)
		}
		results <- res{s: mod, val: value, m: simmet, d: simdmg / float64(numSims) / float64(seconds)}
		if base == 0 {
			base = simdmg / float64(numSims) / float64(seconds)
		}
	}

	doStat(StatSpellDmg, 0)

	// order of these doesn't matter, we put them back in the right order in the next loop.
	statsToTest := []Stat{StatInt, StatSpellDmg, StatSpellCrit, StatSpellHit, StatHaste, StatMP5}
	for _, v := range statsToTest {
		go doStat(v, 50)
	}

	modded := make([]float64, StatLen)
	for i := 0; i < len(statsToTest)+1; i++ {
		res := <-results
		// fmt.Printf("\n---- %s / +%0.0f  -  Diff: %0.1f\n", res.s.StatName(), res.val, res.d-base)
		if res.s == StatSpellDmg && res.val == 0 {
			continue
		} else {
			modded[res.s] = res.d - base
		}
		// printResult(res.m, seconds)
	}

	output := make([]float64, StatMP5+1) // one more than MP5
	for _, v := range statsToTest {
		output[v] = modded[v] / modded[StatSpellDmg]
	}
	return output
}

// OptimalGems returns DPS for each equipment/gem set.
//   1. All +sp power
//   2. Follow colors to get socket bonuses
func OptimalGems(opts Options, equip Equipment, seconds int, numSims int) OptimalGemsResult {
	output := OptimalGemsResult{}

	set1 := equip.Clone()
	set2 := equip.Clone()
	// set3 := equip.Clone()

	ruby := GemLookup["Runed Living Ruby"] // red +sp
	// dawnstone := GemLookup["Gleaming Dawnstone"] // yellow  +crit
	nt := GemLookup["Potent Noble Topaz"] // orange  +sp/+crit
	// fo := GemLookup["Infused Fire Opal"] // orange +sp/+int
	chryo := GemLookup["Glowing Nightseye"] // green  crit/mp5
	// tanz := GemLookup["Glowing Tanzanite"] // purple  sp/stm
	// tala := GemLookup["Dazzling Talasite"] // green  int/mp5

	for i := range set1 {
		set1[i].Gems = make([]Gem, len(set1[i].GemSlots))
		for gs, color := range set1[i].GemSlots {
			if color != GemColorMeta {
				set1[i].Gems[gs] = ruby
			} else {
				set1[i].Gems[gs] = Gems[0]
			}
		}
	}

	for i := range set2 {
		set2[i].Gems = make([]Gem, len(set2[i].GemSlots))
		for gs, color := range set2[i].GemSlots {
			switch color {
			case GemColorRed:
				set2[i].Gems[gs] = ruby
			case GemColorYellow:
				set2[i].Gems[gs] = nt
			case GemColorBlue:
				set2[i].Gems[gs] = chryo
			case GemColorMeta:
				set2[i].Gems[gs] = Gems[0]
			}
		}
	}

	simdmg := 0.0
	sim := NewSim(CalculateTotalStats(opts, set1), equip, opts)
	simmet := make([]SimMetrics, 0, numSims)
	for ns := 0; ns < numSims; ns++ {
		metrics := sim.Run(seconds)
		simdmg += metrics.TotalDamage
		simmet = append(simmet, metrics)
	}
	fmt.Printf("All Red Gems: %0.0f DPS\n", simdmg/float64(numSims)/float64(seconds))
	output.Sims = append(output.Sims, EquipmentResult{Results: simmet, Equip: set1})

	simdmg = 0.0
	sim = NewSim(CalculateTotalStats(opts, set2), equip, opts)
	simmet = make([]SimMetrics, 0, numSims)
	for ns := 0; ns < numSims; ns++ {
		metrics := sim.Run(seconds)
		simdmg += metrics.TotalDamage
		simmet = append(simmet, metrics)
	}
	fmt.Printf("Matched Sockets: %0.0f DPS\n", simdmg/float64(numSims)/float64(seconds))
	output.Sims = append(output.Sims, EquipmentResult{Results: simmet, Equip: set2})

	return output
}

// Finds the optimal rotation for given parameters.
// This might not be needed now that the AI basically does this but faster.
func OptimalRotation(stats Stats, opts Options, equip Equipment, seconds int, numSims int) ([]SimMetrics, string) {
	topDmg := 0.0
	bestAgent := AGENT_TYPE_FIXED_LB_ONLY
	topMets := []SimMetrics{}

	for _, agentType := range ALL_AGENT_TYPES {
		if !strings.Contains(string(agentType), "Fixed") {
			continue
		}

		oomat := 0
		numoom := 0
		simdmg := 0.0
		// fmt.Printf("Starting opt sim: %v\n", rotation)
		sopts := opts
		sopts.AgentType = agentType
		sim := NewSim(stats, equip, sopts)

		simmet := make([]SimMetrics, 0, numSims)
		for ns := 0; ns < numSims; ns++ {
			metrics := sim.Run(seconds)

			if metrics.OOMAt != 0 {
				oomat += metrics.OOMAt
				numoom++
			}
			simdmg += metrics.TotalDamage

			simmet = append(simmet, metrics)
		}

		if simdmg > topDmg {
			topDmg = simdmg
			bestAgent = agentType
			topMets = simmet
		}
	}

	// Optimal Found
	// fmt.Printf("Optimal Found: %0.0f DPS (%d LB : 1 CL)\n", simdmg/float64(seconds)/float64(numSims), numLB)
	// printResult(simmet, seconds)
	return topMets, string(bestAgent)

}

func PrintResult(metrics []SimMetrics, seconds int) {
	numSims := len(metrics)
	simDmgs := make([]float64, 0, numSims)
	for _, metric := range metrics {
		simDmgs = append(simDmgs, metric.TotalDamage)
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

	fmt.Printf("DPS:\n")
	fmt.Printf("\tMean: %0.1f +/- %0.1f\n", (mean / float64(seconds)), stdev/float64(seconds))
	fmt.Printf("\tMax: %0.1f\n", (max / float64(seconds)))
	fmt.Printf("Total Casts:\n")

	// for k, v := range casts {
	// 	fmt.Printf("\t%s: %d\n", AuraName(k), v/numSims)
	// }
}
