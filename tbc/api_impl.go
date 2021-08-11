// Top-level implementations for the api.go functions.
package tbc

import (
	"fmt"
	"math"
	"strings"
	"time"
)

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

	simResult := aggregator.getResult()
	simResult.Logs = logsBuffer.String()
	return simResult
}

type MetricsAggregator struct {
	startTime       time.Time
	numSims         int

	dpsSum          float64
	dpsSumSquared   float64
	dpsMax          float64
	dpsHist         map[int]int          // rounded DPS to count

	numOom          int
	oomAtSum        float64
	dpsAtOomSum     float64

	casts           map[int32]CastMetric
}

func NewMetricsAggregator() *MetricsAggregator {
	return &MetricsAggregator{
		startTime: time.Now(),
		dpsHist: make(map[int]int),
		casts: make(map[int32]CastMetric),
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

	dpsRounded := int(math.Round(dps / 10) * 10)
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
	result.DpsAvg      = aggregator.dpsSum / numSims
	result.DpsStDev    = math.Sqrt((aggregator.dpsSumSquared / numSims) - (result.DpsAvg * result.DpsAvg))
	result.DpsMax      = aggregator.dpsMax
	result.DpsHist     = aggregator.dpsHist

	result.NumOom      = aggregator.numOom
	result.OomAtAvg    = aggregator.oomAtSum / numSims
	result.DpsAtOomAvg = aggregator.dpsAtOomSum / numSims

	result.Casts       = aggregator.casts

	return result
}
