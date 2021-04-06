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

// Finds the optimal rotation for given parameters.
func OptimalRotation(stats Stats, opts Options, equip Equipment, seconds int, numSims int) ([]SimMetrics, []string) {

	fmt.Printf("Starting optimize...\n")

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
		fmt.Printf("Starting opt sim: %v\n", rotation)
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
		fmt.Printf("(%d LB: 1 CL) %0.0f DPS  OOM: %0.0f percent\n", numLB, simdmg/float64(seconds)/float64(numSims), avgOOM*100)

		if numLB == minLB || numLB == maxLB {
			if simdmg >= topDmg {
				printResult(simmet, seconds)
				return simmet, rotation
			}
			printResult(topMets, seconds)
			return topMets, topRot
		}
		// avgOOMAt := int(float64(oomat) / float64(numoom))
		if avgOOM < 0.1 {
			newLB := (numLB + minLB) / 2
			if numLB-minLB <= 1 { // im lazy and this is easy to write...
				newLB = minLB
			}
			if newLB == numLB {
				printResult(simmet, seconds)
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
				printResult(simmet, seconds)
				return simmet, rotation
			}
			minLB = numLB
			numLB = newLB
			continue
		}

		// Optimal Found
		fmt.Printf("Optimal Found: %0.0f DPS (%d LB : 1 CL)\n", simdmg/float64(seconds)/float64(numSims), numLB)
		printResult(simmet, seconds)
		return simmet, rotation
	}

}

func printResult(metrics []SimMetrics, seconds int) {
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
	fmt.Printf("\tMean: %0.1f\n", (mean / float64(seconds)))
	fmt.Printf("\tMax: %0.1f\n", (max / float64(seconds)))
	fmt.Printf("\tStd.Dev: %0.1f\n", stdev/float64(seconds))
	fmt.Printf("Total Casts:\n")

	// for k, v := range casts {
	// 	fmt.Printf("\t%s: %d\n", tbc.AuraName(k), v/numSims)
	// }

}

type EleAI struct {
	LastMana    float64
	DidPot      bool
	LastCheck   int
	Rotation    []int32
	rotationIdx int

	rotations int
}

func NewAI(sim *Simulation) *EleAI {
	ai := &EleAI{
		Rotation: sim.SpellRotation,
		LastMana: sim.CurrentMana,
	}
	if sim.RotationIdx == -1 {
		ai.Rotation = []int32{MagicIDCL6, MagicIDLB12, MagicIDLB12, MagicIDLB12, MagicIDLB12}
	}
	sim.debug("[AI] initialized rotation: %#v\n", sim.SpellRotation)
	return ai
}

func (ai *EleAI) ChooseSpell(sim *Simulation, didPot bool) int {
	if didPot {
		ai.DidPot = true
	}
	if ai.rotationIdx == len(ai.Rotation) {
		ai.rotationIdx = 0
		ai.rotations++

		if ai.DidPot { // reset our behavior because potion messed up mana drain.
			sim.debug("[AI] Resetting rotation selection because potion was used.\n")
			ai.LastCheck = sim.currentTick
			ai.rotations = 0
			ai.LastMana = sim.CurrentMana
			ai.DidPot = false
		}
	}
	if ai.rotations == 1 {
		// Re-evaluate our rotation.
		manaDrained := ai.LastMana - sim.CurrentMana
		timePassed := sim.currentTick - ai.LastCheck
		rate := manaDrained / float64(timePassed)
		timeRemaining := sim.endTick - sim.currentTick
		totalManaDrain := rate * float64(timeRemaining)
		buffer := 0.0 // mana buffer to not module rotations too fast...

		sim.debug("[AI] End of rotation, recalculating best rotation. Rate: %0.1f, Total Drain: %0.1f, CurrentMana: %0.1f\n", rate, totalManaDrain, sim.CurrentMana)
		if totalManaDrain > sim.CurrentMana+buffer {
			// too much mana drain, less chain lightning.
			ai.Rotation = append(ai.Rotation, MagicIDLB12)
			sim.debug("[AI] - adding LB to %#v\n", ai.Rotation)
		} else if totalManaDrain < sim.CurrentMana-buffer {
			// more chain lightning
			if len(ai.Rotation) > 5 { // dont drop below 4xLB, 1xCL
				ai.Rotation = ai.Rotation[:len(ai.Rotation)-1]
				sim.debug("[AI] - dropping LB to %#v\n", ai.Rotation)
			}
		} else {
			// continue current rotation
		}

		// Reset checker
		ai.LastCheck = sim.currentTick
		ai.LastMana = sim.CurrentMana
		ai.rotations = 0
	}

	so := ai.Rotation[ai.rotationIdx]
	sp := spellmap[so]
	cast := NewCast(sim, sp, sim.Stats[StatSpellDmg], sim.Stats[StatSpellHit], sim.Stats[StatSpellCrit])
	if sim.CDs[so] < 1 {
		if sim.CurrentMana >= cast.ManaCost {
			sim.CastingSpell = cast
			ai.rotationIdx++
			return cast.TicksUntilCast
		} else {
			sim.debug("[AI] OOM Current Mana %0.0f, Cast Cost: %0.0f\n", sim.CurrentMana, cast.ManaCost)
			if sim.metrics.OOMAt == 0 {
				sim.metrics.OOMAt = sim.currentTick / TicksPerSecond
				sim.metrics.DamageAtOOM = sim.metrics.TotalDamage
			}
			return int(math.Ceil((cast.ManaCost - sim.CurrentMana) / sim.manaRegen()))
		}
	}
	return sim.CDs[so]

}

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
// 		fmt.Printf("Rot: %s, Mana Left: %v\n", p1Res.R, avgMana)
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
// 		fmt.Printf("BASE(%s): %0.0f, BL(%s): %0.0f, TOTAL: %d\n", nr.OR, nr.oDPS, nr.R, bldps, totaltotal)
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
