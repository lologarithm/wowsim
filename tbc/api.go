// The interface to the sim. All interactions with the sim should go through this file.
package tbc

/**
 * All the information needed to construct a sim and run it multiple times.
 */
type SimRequest struct {
	Options Options
	Gear EquipmentSpec
	Iterations int
	IncludeLogs bool
}

/**
 * The results / metrics from running a single sim one or more times.
 */
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

type CastMetric struct {
	Count int     `json:"count"`
	Dmg   float64 `json:"dmg"`
	Crits int     `json:"crits"`
}

func RunSimulation(request SimRequest) SimResult {
	return runSimulationImpl(request)
}

