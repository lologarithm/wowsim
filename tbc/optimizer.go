package tbc

import (
	"fmt"
	"math"
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
		myopts.Buffs.Custom = Stats{StatLen: 0}
		myopts.Buffs.Custom[mod] += value
		stats := myopts.StatTotal(equip)

		myOpts := myopts
		myOpts.SpellOrder = []string{""}
		myOpts.UseAI = true
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

	statsToTest := []Stat{StatSpellDmg, StatSpellCrit, StatInt, StatSpellHit, StatMP5, StatHaste}
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

	output := make([]float64, StatLen)
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
	chryo := GemLookup["Rune Covered Chrysoprase"] // green  crit/mp5
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

	fmt.Printf("Set1:\n%s\nSet2:\n%s\n", opts.StatTotal(set1).Print(false), opts.StatTotal(set2).Print(false))

	simdmg := 0.0
	sim := NewSim(opts.StatTotal(set1), equip, opts)
	simmet := make([]SimMetrics, 0, numSims)
	for ns := 0; ns < numSims; ns++ {
		metrics := sim.Run(seconds)
		simdmg += metrics.TotalDamage
		simmet = append(simmet, metrics)
	}
	fmt.Printf("Set1: %0.0f\n", simdmg/float64(numSims)/float64(seconds))
	output.Sims = append(output.Sims, EquipmentResult{Results: simmet, Equip: set1})

	simdmg = 0.0
	sim = NewSim(opts.StatTotal(set2), equip, opts)
	simmet = make([]SimMetrics, 0, numSims)
	for ns := 0; ns < numSims; ns++ {
		metrics := sim.Run(seconds)
		simdmg += metrics.TotalDamage
		simmet = append(simmet, metrics)
	}
	fmt.Printf("Set2: %0.0f\n", simdmg/float64(numSims)/float64(seconds))
	output.Sims = append(output.Sims, EquipmentResult{Results: simmet, Equip: set2})

	return output
}

// Finds the optimal rotation for given parameters.
func OptimalRotation(stats Stats, opts Options, equip Equipment, seconds int, numSims int) ([]SimMetrics, []string) {

	// fmt.Printf("Starting optimize...\n")

	topDmg := 0.0
	topRot := []string{}
	topMets := []SimMetrics{}

	numLB := 10

	maxLB := 40
	minLB := 4

	for {
		rotation := []string{"CL6"}
		for i := 0; i < numLB; i++ {
			rotation = append(rotation, "LB12")
		}

		oomat := 0
		numoom := 0
		simdmg := 0.0
		// fmt.Printf("Starting opt sim: %v\n", rotation)
		sopts := opts
		sopts.SpellOrder = rotation
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
			topMets = simmet
			topRot = rotation
		}

		avgOOM := float64(numoom) / float64(numSims)
		// fmt.Printf("(%d LB: 1 CL) %0.0f DPS  OOM: %0.0f percent\n", numLB, simdmg/float64(seconds)/float64(numSims), avgOOM*100)

		if numLB == minLB || numLB == maxLB {
			if simdmg >= topDmg {
				// printResult(simmet, seconds)
				return simmet, rotation
			}
			// printResult(topMets, seconds)
			return topMets, topRot
		}
		// avgOOMAt := int(float64(oomat) / float64(numoom))
		if avgOOM < 0.1 {
			newLB := (numLB + minLB) / 2
			if numLB-minLB <= 1 { // im lazy and this is easy to write...
				newLB = minLB
			}
			if newLB == numLB {
				// printResult(simmet, seconds)
				return simmet, rotation
			}
			maxLB = numLB
			numLB = newLB
			continue
		} else if avgOOM > 0.33 {
			newLB := (numLB + maxLB) / 2
			if maxLB-numLB <= 1 { // im lazy and this is easy to write...
				newLB = maxLB
			} else if avgOOM > 0.9 {
				newLB = (newLB + maxLB) / 2 // skip ahead
			}

			if newLB == numLB {
				// printResult(simmet, seconds)
				return simmet, rotation
			}
			minLB = numLB
			numLB = newLB
			continue
		}

		// Optimal Found
		// fmt.Printf("Optimal Found: %0.0f DPS (%d LB : 1 CL)\n", simdmg/float64(seconds)/float64(numSims), numLB)
		// printResult(simmet, seconds)
		return simmet, rotation
	}

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

type EleAI struct {
	DidPot    bool
	LastMana  float64
	LastCheck int

	NumCasts int
	LB       *Spell
	CL       *Spell
}

func NewAI(sim *Simulation) *EleAI {
	ai := &EleAI{
		LB:       spellmap[MagicIDLB12],
		CL:       spellmap[MagicIDCL6],
		LastMana: sim.CurrentMana,
		NumCasts: 3,
	}
	sim.debug("[AI] initialized\n")
	return ai
}

func (ai *EleAI) ChooseSpell(sim *Simulation, didPot bool) int {
	if ai.LastMana == 0 {
		ai.LastMana = sim.CurrentMana
	}
	ai.NumCasts++
	if didPot {
		// Use Potion to reset the calculation... only early on in fight.
		ai.LastMana = sim.CurrentMana
		ai.LastCheck = sim.currentTick
		ai.NumCasts = 0
	}
	if sim.CDs[MagicIDCL6] < 1 && ai.NumCasts > 3 {
		manaDrained := ai.LastMana - sim.CurrentMana
		timePassed := sim.currentTick - ai.LastCheck
		if timePassed == 0 {
			timePassed = 1
		}
		rate := manaDrained / float64(timePassed)
		timeRemaining := sim.endTick - sim.currentTick
		totalManaDrain := rate * float64(timeRemaining)
		buffer := ai.CL.Mana // mana buffer of 1 extra CL

		sim.debug("[AI] CL Ready: Mana/Tick: %0.1f, Est Mana Drain: %0.1f, CurrentMana: %0.1f\n", rate, totalManaDrain, sim.CurrentMana)
		// If we have enough mana to burn and CL is on CD, use it.
		if totalManaDrain < sim.CurrentMana-buffer {
			cast := NewCast(sim, ai.CL)
			if sim.CurrentMana >= cast.ManaCost {
				sim.debug("[AI] Selected CL\n")
				sim.CastingSpell = cast
				return cast.TicksUntilCast
			}
		}
	}
	cast := NewCast(sim, ai.LB)

	if sim.CurrentMana >= cast.ManaCost {
		sim.debug("[AI] Selected LB\n")
		sim.CastingSpell = cast
		return cast.TicksUntilCast
	}

	sim.debug("[AI] OOM Current Mana %0.0f, Cast Cost: %0.0f\n", sim.CurrentMana, cast.ManaCost)
	if sim.metrics.OOMAt == 0 {
		sim.metrics.OOMAt = sim.currentTick / TicksPerSecond
		sim.metrics.DamageAtOOM = sim.metrics.TotalDamage
	}
	return int(math.Ceil((cast.ManaCost - sim.CurrentMana) / sim.manaRegen()))
}

// 	if didPot {
// 		ai.DidPot = true
// 	}
// 	if ai.rotations == 1 {
// 		// Re-evaluate our rotation.
// 		} else {
// 			// continue current rotation
// 		}

// 		// Reset checker
// 		ai.LastCheck = sim.currentTick
// 		ai.LastMana = sim.CurrentMana
// 		ai.rotations = 0
// 	}

// }

// 	res := make(chan optRes, 1)

// 	// Sim each rotation but duration - #BL*40s without BL
// 	simRotations(res, stats, Options{}, equip, seconds-opts.NumBloodlust*40, numSims)

// 	final := make(chan optRes, 1)
// 	for range rotations {
// 		p1Res := <-res
// 		numOOM := 0
// 		avgMana := 0
// 		dps := 0.0
// 		for _, v := range p1Res.M {
// 			if v.OOMAt != 0 {
// 				numOOM++
// 			}
// 			avgMana += v.ManaAtEnd
// 			dps += v.TotalDamage
// 		}
// 		//  Get mana remaining
// 		avgMana /= numSims
// 		dps /= float64(numSims * (seconds - opts.NumBloodlust*40))
// fmt.Printf("Rot: %s, Mana Left: %v\n", p1Res.R, avgMana)
// 		// if avgMana < 1000 {
// 		// 	// If you dont have 1000 mana to spare on average, there is little value in BL.
// 		// 	//   This is disregarding 'critical' phases of boss fights where you might conserve mana until that spot, and then BL.
// 		// 	continue
// 		// }

// 		// Sim BL dur with the mana pool left, find top DPS rotation
// 		go func(v optRes) {
// 			ns := stats.Clone()
// 			ns[StatMana] = float64(avgMana)
// 			ns[StatMP5] = 0 // remove regen for this, its incorporated in the avg mana remaining.
// 			secRes := make(chan optRes, 1)
// 			simRotations(secRes, stats, Options{NumBloodlust: opts.NumBloodlust}, equip, opts.NumBloodlust*40, numSims)
// 			for range rotations {
// 				nr := <-secRes
// 				nr.OR = v.R
// 				nr.oDPS = dps
// 				final <- nr
// 			}
// 		}(p1Res)
// 	}

// 	// We now have an optimal BL rotation for each 'main' rotation.
// 	// Sim full duration, include optimal BL rotation for each main rotation.
// 	for i := 0; i < len(rotations)*len(rotations); i++ {
// 		nr := <-final
// 		bltotal := 0.0
// 		for _, v := range nr.M {
// 			bltotal += v.TotalDamage
// 		}
// 		bltotal /= float64(numSims)
// 		bldps := bltotal / float64(opts.NumBloodlust*40)

// 		totaltotal := int((nr.oDPS*float64(seconds-opts.NumBloodlust*40))+bltotal) / seconds
// fmt.Printf("BASE(%s): %0.0f, BL(%s): %0.0f, TOTAL: %d\n", nr.OR, nr.oDPS, nr.R, bldps, totaltotal)
// 	}
// }

// type optRes struct {
// 	M    []SimMetrics
// 	R    []string
// 	OR   []string
// 	oDPS float64
// }

// func simRotations(res chan optRes, stats Stats, opts Options, equip Equipment, seconds int, numSims int) {
// 	rseed := time.Now().Unix()
// 	for _, rot := range rotations {
// 		go func(v []string) {
// 			allRes := make([]SimMetrics, 0, numSims)
// 			sim := NewSim(stats, equip, Options{SpellOrder: v, RSeed: rseed, NumBloodlust: opts.NumBloodlust})
// 			for ns := 0; ns < numSims; ns++ {
// 				allRes = append(allRes, sim.Run(seconds))
// 			}
// 			res <- optRes{M: allRes, R: v}
// 		}(rot)
// 	}
// }
