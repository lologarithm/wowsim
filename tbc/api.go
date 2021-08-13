// The interface to the sim. All interactions with the sim should go through this file.
package tbc

/**
 * Returns all items, enchants, and gems recognized by the sim.
 */
type GearListRequest struct {
}
type GearListResult struct {
	Items    []Item
	Enchants []Enchant
	Gems     []Gem
}

func GetGearList(request GearListRequest) GearListResult {
	return getGearListImpl(request)
}

/**
 * Returns character stats taking into account gear / buffs / consumes / etc
 */
type ComputeStatsRequest struct {
	Options     Options
	Gear        EquipmentSpec
}
type ComputeStatsResult struct {
	// Only from gear (no base stats!)
	GearOnly Stats

	// Stats with everything(base, gear, buffs, consumes, sets)
	FinalStats Stats

	// The sets used with FinalStats
	Sets []string
}

func ComputeStats(request ComputeStatsRequest) ComputeStatsResult {
	return computeStatsImpl(request)
}

/**
 * Returns stat weights and EP values, with standard deviations, for all stats.
 */
type StatWeightsRequest struct {
	Options     Options
	Gear        EquipmentSpec
	Iterations  int
}
type StatWeightsResult struct {
	// Increase in dps from 1 point of each stat
	Weights       Stats

	// Standard deviations for Weights
	WeightsStDev  Stats

	// EP value for each stat, basically just the weight but normalized so that
	// spell power always has an EP of 1
	EpValues      Stats

	// Standard deviations for EpValues
	EpValuesStDev Stats
}

func StatWeights(request StatWeightsRequest) StatWeightsResult {
	return statWeightsImpl(request)
}


/**
 * Runs multiple iterations of the sim with a single set of options / gear.
 */
type SimRequest struct {
	Options     Options
	Gear        EquipmentSpec
	Iterations  int
	IncludeLogs bool
}
type SimResult struct {
	ExecutionDurationMs int64
	Logs                string

	DpsAvg              float64
	DpsStDev            float64
	DpsMax              float64
	DpsHist             map[int]int          // rounded DPS to count

	NumOom              int
	OomAtAvg            float64
	DpsAtOomAvg         float64

	Casts               map[int32]CastMetric
}

/**
 * Metrics related to a specific type of cast, e.g. LB or CL
 */
type CastMetric struct {
	Count int
	Dmg   float64
	Crits int
}

func RunSimulation(request SimRequest) SimResult {
	return runSimulationImpl(request)
}

/**
 * Runs separate SimRequests for each provided request, and returns results in the same order.
 */
type BatchSimRequest struct {
	Requests []SimRequest
}
type BatchSimResult struct {
	Results  []SimResult
}

func RunBatchSimulation(request BatchSimRequest) BatchSimResult {
	return runBatchSimulationImpl(request)
}

type PackOptionsRequest struct {
	Options Options
}
type PackOptionsResult struct {
	Data []byte
	Length int // length of Data array
}

func PackOptions(request PackOptionsRequest) PackOptionsResult {
	data := request.Options.Pack()
	return PackOptionsResult{
		Data: data,
		Length: len(data),
	}
}

type ApiRequest struct {
	RequestType  string
	GearList     GearListRequest
	ComputeStats ComputeStatsRequest
	StatWeights  StatWeightsRequest
	Sim          SimRequest
	BatchSim     BatchSimRequest
	PackOptions  PackOptionsRequest
}

type ApiResult struct {
	GearList     GearListResult
	ComputeStats ComputeStatsResult
	StatWeights  StatWeightsResult
	Sim          SimResult
	BatchSim     BatchSimResult
	PackOptions  PackOptionsResult
}

func ApiCall(request ApiRequest) ApiResult {
	if request.RequestType == "GearList" {
		return ApiResult{
			GearList: GetGearList(request.GearList),
		}
	} else if request.RequestType == "ComputeStats" {
		return ApiResult{
			ComputeStats: ComputeStats(request.ComputeStats),
		}
	} else if request.RequestType == "StatWeights" {
		return ApiResult{
			StatWeights: StatWeights(request.StatWeights),
		}
	} else if request.RequestType == "Sim" {
		return ApiResult{
			Sim: RunSimulation(request.Sim),
		}
	} else if request.RequestType == "BatchSim" {
		return ApiResult{
			BatchSim: RunBatchSimulation(request.BatchSim),
		}
	} else if request.RequestType == "PackOptions" {
		return ApiResult{
			PackOptions: PackOptions(request.PackOptions),
		}
	} else {
		panic("Invalid request type: " + request.RequestType)
	}
}
