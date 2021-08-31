package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/lologarithm/wowsim/tbc"
)

// Maps agentType flag values to actual enum types
var agentTypesMap = map[string]tbc.AgentType{
	"3LB1CL":        tbc.AGENT_TYPE_FIXED_3LB_1CL,
	"4LB1CL":        tbc.AGENT_TYPE_FIXED_4LB_1CL,
	"5LB1CL":        tbc.AGENT_TYPE_FIXED_5LB_1CL,
	"6LB1CL":        tbc.AGENT_TYPE_FIXED_6LB_1CL,
	"7LB1CL":        tbc.AGENT_TYPE_FIXED_7LB_1CL,
	"8LB1CL":        tbc.AGENT_TYPE_FIXED_8LB_1CL,
	"9LB1CL":        tbc.AGENT_TYPE_FIXED_9LB_1CL,
	"10LB1CL":       tbc.AGENT_TYPE_FIXED_10LB_1CL,
	"LB":            tbc.AGENT_TYPE_FIXED_LB_ONLY,
	"Adaptive":      tbc.AGENT_TYPE_ADAPTIVE,
	"CLOnClearcast": tbc.AGENT_TYPE_CL_ON_CLEARCAST,
}

var DEFAULT_EQUIPMENT = tbc.EquipmentSpec{
	tbc.ItemSpec{Name: "Tidefury Helm"},
	tbc.ItemSpec{Name: "Charlotte's Ivy"},
	tbc.ItemSpec{Name: "Pauldrons of Wild Magic"},
	tbc.ItemSpec{Name: "Ogre Slayer's Cover"},
	tbc.ItemSpec{Name: "Tidefury Chestpiece"},
	tbc.ItemSpec{Name: "World's End Bracers"},
	tbc.ItemSpec{Name: "Earth Mantle Handwraps"},
	tbc.ItemSpec{Name: "Netherstrike Belt"},
	tbc.ItemSpec{Name: "Stormsong Kilt"},
	tbc.ItemSpec{Name: "Magma Plume Boots"},
	tbc.ItemSpec{Name: "Cobalt Band of Tyrigosa"},
	tbc.ItemSpec{Name: "Sparking Arcanite Ring"},
	tbc.ItemSpec{Name: "Mazthoril Honor Shield"},
	tbc.ItemSpec{Name: "Gavel of Unearthed Secrets"},
	tbc.ItemSpec{Name: "Natural Alignment Crystal"},
	tbc.ItemSpec{Name: "Icon of the Silver Crescent"},
	tbc.ItemSpec{Name: "Totem of the Void"},
}

var DEFAULT_OPTIONS = tbc.Options{
	NumBloodlust: 0,
	NumDrums:     0,
	Buffs: tbc.Buffs{
		ArcaneInt:                false,
		GiftOftheWild:            false,
		BlessingOfKings:          false,
		ImprovedBlessingOfWisdom: false,
		JudgementOfWisdom:        false,
		Moonkin:                  false,
		SpriestDPS:               0,
		WaterShield:              true,
		// Race:                     tbc.RaceBonusOrc,
		Custom: tbc.Stats{
			tbc.StatInt:       290,
			tbc.StatSpellDmg:  598 + 55,
			tbc.StatSpellHit:  24,
			tbc.StatSpellCrit: 120,
		},
	},
	Consumes: tbc.Consumes{
		// FlaskOfBlindingLight: true,
		// BrilliantWizardOil:   false,
		// MajorMageblood:       false,
		// BlackendBasilisk:     true,
		SuperManaPotion: false,
		// DarkRune:             false,
	},
	Talents: tbc.Talents{
		LightningOverload:  5,
		ElementalPrecision: 3,
		NaturesGuidance:    3,
		TidalMastery:       5,
		ElementalMastery:   true,
		UnrelentingStorm:   3,
		CallOfThunder:      5,
		Concussion:         5,
		Convection:         5,
	},
	Totems: tbc.Totems{
		TotemOfWrath: 1,
		WrathOfAir:   false,
		ManaStream:   true,
	},
}

// /script print(GetSpellBonusDamage(4))

func main() {
	// f, err := os.Create("profile2.cpu")
	// if err != nil {
	// 	log.Fatal("could not create CPU profile: ", err)
	// }
	// defer f.Close() // error handling omitted for example
	// if err := pprof.StartCPUProfile(f); err != nil {
	// 	log.Fatal("could not start CPU profile: ", err)
	// }
	// defer pprof.StopCPUProfile()

	var isDebug = flag.Bool("debug", false, "Include --debug to spew the entire simulation log.")
	var noopt = flag.Bool("noopt", false, "If included it will disable optimization.")
	var agentTypeStr = flag.String("agentType", "Adaptive", "Custom comma separated agent type to simulate.\n\tFor Example: --rotation=3LB1CL")
	var duration = flag.Float64("duration", 300, "Custom fight duration in seconds.")
	var iterations = flag.Int("iter", 10000, "Custom number of iterations for the sim to run.")
	var configFile = flag.String("config", "", "Specify an input configuration.")

	flag.Parse()

	simRequest := tbc.SimRequest{}
	if *configFile != "" {
		fileData, err := ioutil.ReadFile(*configFile)
		if err != nil {
			log.Fatalf("Failed to open config file(%s): %s", *configFile, err)
		}

		_ = json.Unmarshal([]byte(fileData), &simRequest)
	} else {
		simRequest.Options = DEFAULT_OPTIONS
		simRequest.Gear = DEFAULT_EQUIPMENT
	}

	if *isDebug {
		*iterations = 1
		simRequest.IncludeLogs = true
	}
	if agentTypeStr == nil {
		simRequest.Options.AgentType = tbc.AGENT_TYPE_ADAPTIVE
	} else if agentType, ok := agentTypesMap[*agentTypeStr]; ok {
		simRequest.Options.AgentType = agentType
	} else {
		panic(fmt.Sprintf("Invalid agent type: %s", *agentTypeStr))
	}
	simRequest.Options.Encounter.Duration = *duration
	simRequest.Options.RSeed = time.Now().Unix()
	simRequest.Iterations = *iterations

	runTBCSim(simRequest, *noopt)
}

func runTBCSim(simRequest tbc.SimRequest, noopt bool) {
	fmt.Printf(
		"\nSim Duration: %0.1f sec\nNum Simulations: %d\n",
		simRequest.Options.Encounter.Duration,
		simRequest.Iterations)

	equipment := tbc.NewEquipmentSet(simRequest.Gear)
	stats := tbc.CalculateTotalStats(simRequest.Options, equipment)
	fmt.Printf("\nFinal Stats: %s\n", stats.Print())

	if !noopt {
		statWeightsResult := tbc.StatWeights(tbc.StatWeightsRequest{
			Options:    simRequest.Options,
			Gear:       simRequest.Gear,
			Iterations: simRequest.Iterations,
		})
		fmt.Printf("Weights: [\n")
		for stat, weight := range statWeightsResult.Weights {
			if tbc.Stat(stat) == tbc.StatStm {
				continue
			}
			fmt.Printf("%s: %0.2f\t", tbc.Stat(stat).StatName(), weight)
		}
		fmt.Printf("\n]\n")
	}

	fmt.Printf("Starting main simulation with agent: %#v", simRequest.Options.AgentType)
	simResult := tbc.RunSimulation(simRequest)
	fmt.Printf("\nLogs:\n%s\n", simResult.Logs)
	fmt.Printf("\n%s\n", simResultsToString(simRequest, simResult))
}

func simResultsToString(request tbc.SimRequest, result tbc.SimResult) string {
	output := ""
	output += fmt.Sprintf("Agent Type: %v\n", string(request.Options.AgentType))
	output += fmt.Sprintf("DPS:")
	output += fmt.Sprintf("\tMean: %0.1f +/- %0.1f\n", result.DpsAvg, result.DpsStDev)
	output += fmt.Sprintf("\tMax: %0.1f\n", result.DpsMax)
	output += fmt.Sprintf("Total Casts:\n")

	for castId, cast := range result.Casts {
		if castId > tbc.MagicIDLen {
			name := tbc.AuraName(1000 - castId)
			output += fmt.Sprintf("\t%s (LO): %0.1f\n", name, float64(cast.Count)/float64(request.Iterations))
		} else {
			output += fmt.Sprintf("\t%s: %0.1f\n", tbc.AuraName(castId), float64(cast.Count)/float64(request.Iterations))
		}
	}

	output += fmt.Sprintf("Went OOM: %d/%d sims\n", result.NumOom, request.Iterations)
	if result.NumOom > 0 {
		output += fmt.Sprintf("Avg OOM Time: %0.1f seconds\n", result.OomAtAvg)
		output += fmt.Sprintf("Avg DPS At OOM: %0.0f\n", result.DpsAtOomAvg)
	}
	output += fmt.Sprintf("Sim execution took %s", time.Duration(result.ExecutionDurationMs)*time.Millisecond)
	return output
}
