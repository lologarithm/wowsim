package tbc

import "fmt"

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
	topNLB := 0

	numLB := 8

	maxLB := 20
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
			topNLB = numLB
			topDmg = simdmg
		}

		avgOOM := float64(numoom) / float64(numSims)
		fmt.Printf("(%d LB: 1 CL) %0.0f DPS  OOM: %0.0f percent\n", numLB, simdmg/float64(seconds)/float64(numSims), avgOOM*100)

		// avgOOMAt := int(float64(oomat) / float64(numoom))
		if avgOOM < 0.1 {
			newLB := (numLB + minLB) / 2
			if newLB == numLB {
				// I guess we fail.
				fmt.Printf("TopNLB = %d, %0.0f DPS\n", topNLB, topDmg/float64(seconds)/float64(numSims))
				fmt.Printf("Cant go lower... Found: %0.0f DPS (%d LB : 1 CL)\n", simdmg/float64(seconds)/float64(numSims), numLB)
				return simmet, rotation
			}
			maxLB = numLB
			numLB = newLB
			continue
		} else if avgOOM > 0.33 {
			newLB := (numLB + maxLB) / 2
			if newLB == numLB {
				// I guess we fail.
				fmt.Printf("TopNLB = %d, %0.0f DPS\n", topNLB, topDmg/float64(numSims))
				fmt.Printf("Cant go higher... Found: %0.0f DPS (%d LB : 1 CL)\n", simdmg/float64(seconds)/float64(numSims), numLB)
				return simmet, rotation
			}
			minLB = numLB
			numLB = newLB
			continue
		}

		// Optimal Found
		fmt.Printf("Optimal Found: %0.0f DPS (%d LB : 1 CL)\n", simdmg/float64(numSims), numLB)
		return simmet, rotation
	}

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
