// Top-level implementations for the api.go functions.
package tbc

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func getGearListImpl(request GearListRequest) GearListResult {
	return GearListResult{
		Items:    items,
		Enchants: Enchants,
		Gems:     Gems,
	}
}

func computeStatsImpl(request ComputeStatsRequest) ComputeStatsResult {
	equipment := NewEquipmentSet(request.Gear)

	gearOnlyStats := equipment.Stats().CalculatedTotal()

	request.Options.AgentType = AGENT_TYPE_ADAPTIVE
	stats := CalculateTotalStats(request.Options, equipment)
	fakeSim := NewSim(stats, equipment, request.Options)
	sets := fakeSim.ActivateSets()

	fakeSim.reset() // this will activate any perm-effect items as well

	finalStats := stats
	for stat, statValue := range fakeSim.Buffs {
		finalStats[stat] += statValue
	}

	return ComputeStatsResult{
		GearOnly:   gearOnlyStats,
		FinalStats: finalStats,
		Sets:       sets,
	}
}

func statWeightsImpl(request StatWeightsRequest) StatWeightsResult {
	request.Options.AgentType = AGENT_TYPE_ADAPTIVE

	baselineSimRequest := SimRequest{
		Options:    request.Options,
		Gear:       request.Gear,
		Iterations: request.Iterations,
	}
	baselineResult := RunSimulation(baselineSimRequest)

	var waitGroup sync.WaitGroup
	result := StatWeightsResult{}
	dpsHists := [StatLen]map[int]int{}

	doStat := func(stat Stat, value float64) {
		defer waitGroup.Done()

		simRequest := baselineSimRequest
		simRequest.Options.Buffs.Custom[stat] += value

		simResult := RunSimulation(simRequest)
		result.Weights[stat] = (simResult.DpsAvg - baselineResult.DpsAvg) / value
		dpsHists[stat] = simResult.DpsHist
	}

	// Spell hit mod shouldn't go over hit cap.
	computeStatsResult := ComputeStats(ComputeStatsRequest{
		Options: request.Options,
		Gear:    request.Gear,
	})
	spellHitMod := math.Max(0, math.Min(10, 202-computeStatsResult.FinalStats[StatSpellHit]))

	statMods := Stats{
		StatInt:       50,
		StatSpellDmg:  50,
		StatSpellCrit: 50,
		StatSpellHit:  spellHitMod,
		StatHaste:     50,
		StatMP5:       50,
	}

	for stat, mod := range statMods {
		if mod == 0 {
			continue
		}

		waitGroup.Add(1)
		go doStat(Stat(stat), mod)
	}

	waitGroup.Wait()

	for stat, mod := range statMods {
		if mod == 0 {
			continue
		}

		result.EpValues[stat] = result.Weights[stat] / result.Weights[StatSpellDmg]

		result.WeightsStDev[stat] = computeStDevFromHists(request.Iterations, statMods[stat], dpsHists[stat], baselineResult.DpsHist, nil, statMods[StatSpellDmg])
		result.EpValuesStDev[stat] = computeStDevFromHists(request.Iterations, statMods[stat], dpsHists[stat], baselineResult.DpsHist, dpsHists[StatSpellDmg], statMods[StatSpellDmg])
	}
	return result
}

func computeStDevFromHists(iters int, modValue float64, moddedStatDpsHist map[int]int, baselineDpsHist map[int]int, spellDmgDpsHist map[int]int, spellDmgModValue float64) float64 {
	sum := 0.0
	sumSquared := 0.0
	n := iters * 10
	for i := 0; i < n; {
		denominator := 1.0
		if spellDmgDpsHist != nil {
			denominator = float64(sampleFromDpsHist(spellDmgDpsHist, iters)-sampleFromDpsHist(baselineDpsHist, iters)) / spellDmgModValue
		}

		if denominator != 0 {
			ep := (float64(sampleFromDpsHist(moddedStatDpsHist, iters)-sampleFromDpsHist(baselineDpsHist, iters)) / modValue) / denominator
			sum += ep
			sumSquared += ep * ep
			i++
		}
	}
	epAvg := sum / float64(n)
	epStDev := math.Sqrt((sumSquared / float64(n)) - (epAvg * epAvg))
	return epStDev
}

func sampleFromDpsHist(hist map[int]int, histNumSamples int) int {
	r := rand.Float64()
	sampleIdx := int(math.Floor(float64(histNumSamples) * r))

	curSampleIdx := 0
	for roundedDps, count := range hist {
		curSampleIdx += count
		if curSampleIdx >= sampleIdx {
			return roundedDps
		}
	}

	panic("Invalid dps histogram")
}

func runSimulationImpl(request SimRequest) SimResult {
	equipment := NewEquipmentSet(request.Gear)
	stats := CalculateTotalStats(request.Options, equipment)

	logsBuffer := &strings.Builder{}
	sim := NewSim(stats, equipment, request.Options)
	aggregator := NewMetricsAggregator()

	for i := 0; i < request.Iterations; i++ {
		if request.IncludeLogs {
			sim.Debug = func(s string, vals ...interface{}) {
				logsBuffer.WriteString(fmt.Sprintf("[%0.1f] "+s, append([]interface{}{sim.CurrentTime}, vals...)...))
			}
		}

		aggregator.addMetrics(request.Options, sim.Run())
	}

	result := aggregator.getResult()
	result.Logs = logsBuffer.String()
	return result
}

func runBatchSimulationImpl(request BatchSimRequest) BatchSimResult {
	result := BatchSimResult{
		Results: make([]SimResult, len(request.Requests)),
	}

	var waitGroup sync.WaitGroup
	runSimOnce := func(idx int) {
		defer waitGroup.Done()
		result.Results[idx] = RunSimulation(request.Requests[idx])
	}

	for i := 0; i < len(request.Requests); i++ {
		waitGroup.Add(1)
		go runSimOnce(i)
	}

	waitGroup.Wait()
	return result
}

type MetricsAggregator struct {
	startTime time.Time
	numSims   int

	dpsSum        float64
	dpsSumSquared float64
	dpsMax        float64
	dpsHist       map[int]int // rounded DPS to count

	numOom      int
	oomAtSum    float64
	dpsAtOomSum float64

	casts map[int32]CastMetric
}

func NewMetricsAggregator() *MetricsAggregator {
	return &MetricsAggregator{
		startTime: time.Now(),
		dpsHist:   make(map[int]int),
		casts:     make(map[int32]CastMetric),
	}
}

func (aggregator *MetricsAggregator) addMetrics(options Options, metrics SimMetrics) {
	aggregator.numSims++

	dps := metrics.TotalDamage / options.Encounter.Duration
	if options.DPSReportTime > 0 {
		dps = metrics.ReportedDamage / float64(options.DPSReportTime)
	}

	aggregator.dpsSum += dps
	aggregator.dpsSumSquared += dps * dps
	aggregator.dpsMax = math.Max(aggregator.dpsMax, dps)

	dpsRounded := int(math.Round(dps/10) * 10)
	aggregator.dpsHist[dpsRounded]++

	if metrics.OOMAt > 0 {
		aggregator.numOom++
		aggregator.oomAtSum += float64(metrics.OOMAt)
		aggregator.dpsAtOomSum += float64(metrics.DamageAtOOM) / float64(metrics.OOMAt)
	}

	for _, cast := range metrics.Casts {
		var id = cast.Spell.ID
		if cast.IsLO {
			id = 1000 - cast.Spell.ID
		}

		cm := aggregator.casts[id]
		cm.Count++
		cm.Dmg += cast.DidDmg
		if cast.DidCrit {
			cm.Crits++
		}

		aggregator.casts[id] = cm
	}
}

func (aggregator *MetricsAggregator) getResult() SimResult {
	result := SimResult{}
	result.ExecutionDurationMs = time.Since(aggregator.startTime).Milliseconds()

	numSims := float64(aggregator.numSims)
	result.DpsAvg = aggregator.dpsSum / numSims
	result.DpsStDev = math.Sqrt((aggregator.dpsSumSquared / numSims) - (result.DpsAvg * result.DpsAvg))
	result.DpsMax = aggregator.dpsMax
	result.DpsHist = aggregator.dpsHist

	result.NumOom = aggregator.numOom
	result.OomAtAvg = aggregator.oomAtSum / float64(aggregator.numOom)
	result.DpsAtOomAvg = aggregator.dpsAtOomSum / numSims

	result.Casts = aggregator.casts

	return result
}
