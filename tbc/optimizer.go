package tbc

import (
	"sync"
)

/**
 * Returns weights for each offensive stat.
 *
 * Weights are NOT EP values. They represent the avg DPS increase from 1 point of that stat.
 */
func StatWeights(simRequest SimRequest) []float64 {
	simRequest.Options.AgentType = AGENT_TYPE_ADAPTIVE
	baselineResult := RunSimulation(simRequest)

	var waitGroup sync.WaitGroup
	statWeights := make([]float64, StatMP5 + 1) // MP5 is the tested stat with the highest idx

	doStat := func(stat Stat, value float64) {
		defer waitGroup.Done()

		request := simRequest
		request.Options.Buffs.Custom[stat] += value

		result := RunSimulation(request)
		statWeights[stat] = (result.DpsAvg - baselineResult.DpsAvg) / value
	}

	statsToTest := []Stat{StatInt, StatSpellDmg, StatSpellCrit, StatSpellHit, StatHaste, StatMP5}
	for _, v := range statsToTest {
		waitGroup.Add(1)
		go doStat(v, 50)
	}

	waitGroup.Wait()
	return statWeights
}
