// The interface to the sim. All interactions with the sim should go through this file.
package tbc

import (
	"encoding/base64"
)

/**
 * Returns all items, enchants, and gems recognized by the sim.
 */
func GetGearList(request GearListRequest) GearListResult {
	return getGearListImpl(request)
}

type GearListRequest struct {
}
type GearListResult struct {
	Items    []Item
	Enchants []Enchant
	Gems     []Gem
}

/**
 * Returns character stats taking into account gear / buffs / consumes / etc
 */
func ComputeStats(request ComputeStatsRequest) ComputeStatsResult {
	return computeStatsImpl(request)
}

type ComputeStatsRequest struct {
	Options Options
	Gear    EquipmentSpec
}
type ComputeStatsResult struct {
	// Only from gear (no base stats!)
	GearOnly Stats

	// Stats with everything(base, gear, buffs, consumes, sets)
	FinalStats Stats

	// The sets used with FinalStats
	Sets []string
}

/**
 * Returns stat weights and EP values, with standard deviations, for all stats.
 */
func StatWeights(request StatWeightsRequest) StatWeightsResult {
	return statWeightsImpl(request)
}

type StatWeightsRequest struct {
	Options    Options
	Gear       EquipmentSpec
	Iterations int
}
type StatWeightsResult struct {
	// Increase in dps from 1 point of each stat
	Weights Stats

	// Standard deviations for Weights
	WeightsStDev Stats

	// EP value for each stat, basically just the weight but normalized so that
	// spell power always has an EP of 1
	EpValues Stats

	// Standard deviations for EpValues
	EpValuesStDev Stats
}

/**
 * Runs multiple iterations of the sim with a single set of options / gear.
 */
func RunSimulation(request SimRequest) SimResult {
	return runSimulationImpl(request)
}

type SimRequest struct {
	Options     Options
	Gear        EquipmentSpec
	Iterations  int
	IncludeLogs bool
}
type SimResult struct {
	ExecutionDurationMs int64
	Logs                string

	DpsAvg   float64
	DpsStDev float64
	DpsMax   float64
	DpsHist  map[int]int // rounded DPS to count

	NumOom      int
	OomAtAvg    float64
	DpsAtOomAvg float64

	Casts map[int32]CastMetric
}

/**
 * Metrics related to a specific type of cast, e.g. LB or CL
 */
type CastMetric struct {
	Count int
	Dmg   float64
	Crits int
}

/**
 * Runs separate SimRequests for each provided request, and returns results in the same order.
 */
func RunBatchSimulation(request BatchSimRequest) BatchSimResult {
	return runBatchSimulationImpl(request)
}

type BatchSimRequest struct {
	Requests []SimRequest
}
type BatchSimResult struct {
	Results []SimResult
}

func PackOptions(request PackOptionsRequest) PackOptionsResult {
	return PackOptionsResult{
		Data: base64.StdEncoding.EncodeToString(request.Options.Pack()),
	}
}

type PackOptionsRequest struct {
	Options Options
}
type PackOptionsResult struct {
	// base64-encoded binary data
	Data string
}

func ApiCall(request ApiRequest) ApiResult {
	if request.GearList != nil {
		result := GetGearList(*request.GearList)
		return ApiResult{GearList: &result}
	} else if request.ComputeStats != nil {
		result := ComputeStats(*request.ComputeStats)
		return ApiResult{ComputeStats: &result}
	} else if request.StatWeights != nil {
		result := StatWeights(*request.StatWeights)
		return ApiResult{StatWeights: &result}
	} else if request.Sim != nil {
		result := RunSimulation(*request.Sim)
		return ApiResult{Sim: &result}
	} else if request.BatchSim != nil {
		result := RunBatchSimulation(*request.BatchSim)
		return ApiResult{BatchSim: &result}
	} else if request.PackOptions != nil {
		result := PackOptions(*request.PackOptions)
		return ApiResult{PackOptions: &result}
	} else {
		panic("Empty API request!")
	}
}

type ApiRequest struct {
	GearList     *GearListRequest
	ComputeStats *ComputeStatsRequest
	StatWeights  *StatWeightsRequest
	Sim          *SimRequest
	BatchSim     *BatchSimRequest
	PackOptions  *PackOptionsRequest
}
type ApiResult struct {
	GearList     *GearListResult
	ComputeStats *ComputeStatsResult
	StatWeights  *StatWeightsResult
	Sim          *SimResult
	BatchSim     *BatchSimResult
	PackOptions  *PackOptionsResult
}
