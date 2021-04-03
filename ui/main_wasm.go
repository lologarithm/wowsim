package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"syscall/js"
	"time"

	"github.com/lologarithm/wowsim/tbc"
)

func main() {
	c := make(chan struct{}, 0)

	simfunc := js.FuncOf(Simulate)
	gearfunc := js.FuncOf(GearStats)
	gearlistfunc := js.FuncOf(GearList)

	js.Global().Set("simulate", simfunc)
	js.Global().Set("gearstats", gearfunc)
	js.Global().Set("gearlist", gearlistfunc)
	js.Global().Call("popgear")
	<-c
}

// GearList reports all items of gear to the UI to display.
func GearList(this js.Value, args []js.Value) interface{} {
	slot := -1

	if len(args) == 1 {
		slot = args[0].Int()
	}
	gears := "["
	for _, v := range tbc.ItemLookup {
		if slot != -1 && v.Slot != slot {
			continue
		}
		if len(gears) != 1 {
			gears += ","
		}
		gears += `{"name":"` + v.Name + `", "slot": ` + strconv.Itoa(v.Slot) + `}`
	}
	gears += "]"
	return gears
}

// GearStats takes a gear list and returns their total stats.
// This could power a simple 'current stats of all gear' UI.
func GearStats(this js.Value, args []js.Value) interface{} {
	return getGear(args[0]).Stats().Print()
}

// getGear converts js string array to a list of equipment items.
func getGear(val js.Value) tbc.Equipment {
	numGear := val.Length()
	gearstr := make([]string, numGear)
	for i := range gearstr {
		gearstr[i] = val.Index(i).String()
	}
	return tbc.NewEquipmentSet(gearstr...)
}

func parseOptions(val js.Value) tbc.Options {
	opt := tbc.Options{
		NumBloodlust: val.Get("buffbl").Int(),
		NumDrums:     val.Get("buffdrum").Int(),
		Buffs: tbc.Buffs{
			ArcaneInt:                val.Get("buffai").Truthy(),
			GiftOftheWild:            val.Get("buffgotw").Truthy(),
			BlessingOfKings:          val.Get("buffbk").Truthy(),
			ImprovedBlessingOfWisdom: val.Get("buffibow").Truthy(),
			JudgementOfWisdom:        val.Get("debuffjow").Truthy(),
			Moonkin:                  val.Get("buffmoon").Truthy(),
			SpriestDPS:               val.Get("buffspriest").Int(),
			WaterShield:              val.Get("sbufws").Truthy(),
		},
		Consumes: tbc.Consumes{
			FlaskOfBlindingLight:     val.Get("confbl").Truthy(),
			FlaskOfMightyRestoration: val.Get("confmr").Truthy(),
			BrilliantWizardOil:       val.Get("conbwo").Truthy(),
			MajorMageblood:           val.Get("conmm").Truthy(),
			BlackendBasilisk:         val.Get("conbb").Truthy(),
			SuperManaPotion:          val.Get("consmp").Truthy(),
			DarkRune:                 val.Get("condr").Truthy(),
		},
		Talents: tbc.Talents{
			LightninOverload:   5,
			ElementalPrecision: 3,
			NaturesGuidance:    3,
			TidalMastery:       5,
			ElementalMastery:   true,
			UnrelentingStorm:   3,
			CallOfThunder:      5,
		},
		Totems: tbc.Totems{
			TotemOfWrath: val.Get("totwr").Int(),
			WrathOfAir:   val.Get("totwoa").Truthy(),
			ManaStream:   val.Get("totms").Truthy(),
		},
	}

	return opt
}

func parseRotation(val js.Value) [][]string {

	out := [][]string{}

	for i := 0; i < val.Length(); i++ {
		rot := []string{}
		jsrot := val.Index(i)
		for j := 0; j < jsrot.Length(); j++ {
			rot = append(rot, jsrot.Index(j).String())
		}
		out = append(out, rot)
	}

	return out
}

// Simulate takes in number of iterations, duration, a gear list, and simulation options.
// (iterations, duration, gearlist, options, <optional, custom rotation)
func Simulate(this js.Value, args []js.Value) interface{} {
	// TODO: Accept talents, buffs, and consumes as inputs.

	if len(args) < 4 {
		print("Expected 4 min arguments:  (#iterations, duration, gearlist, options)")
		return `{"error": "invalid arguments supplied"}`
	}

	customRotation := [][]string{}
	customHaste := 0.0
	if len(args) == 6 {
		customRotation = parseRotation(args[4])
		customHaste = args[5].Float()
	}
	gear := getGear(args[2])
	fmt.Printf("Gear Stats: %s", gear.Stats().Print())
	opt := parseOptions(args[3])
	stats := opt.StatTotal(gear)
	if customHaste != 0 {
		stats[tbc.StatHaste] = customHaste
	}
	fmt.Printf("Total Stats: %s", stats.Print())

	simi := args[0].Int()
	if simi == 1 {
		tbc.IsDebug = true
	}
	dur := args[1].Int()
	results := runTBCSim(opt, stats, gear, dur, simi, customRotation)

	st := time.Now()
	output, err := json.Marshal(results)
	if err != nil {
		print("Failed to json marshal results: ", err.Error())
	}
	fmt.Printf("Took %s to json marshal response.\n", time.Now().Sub(st))
	// output := "["

	// for i, res := range results {
	// 	output += fmt.Sprintf(`{"Duration": "%s"}`, res.Duration.String())

	// 	if i != len(results)-1 {
	// 		output += ","
	// 	}
	// }

	// output += "]"
	return string(output)
}

type SimResult struct {
	Casts        [][]CastMetric
	TotalDmgs    []float64
	DmgAtOOMs    []float64
	OOMAt        []float64 // oom time totals
	Rotation     []string
	SimSeconds   int
	RealDuration time.Duration
}

type CastMetric struct {
	Spell string
	Hit   bool
	Crit  bool
	Dmg   float64
	Time  float64 // seconds it took to cast this spell
}

func runTBCSim(opts tbc.Options, stats tbc.Stats, equip tbc.Equipment, seconds int, numSims int, customRotation [][]string) []SimResult {
	print("\nSim Duration:", seconds)
	print("\nNum Simulations: ", numSims)
	print("\n")

	spellOrders := [][]string{
		{"LB12"},               // only LB
		{"pri", "CL6", "LB12"}, // cast CL whenever off CD, otherwise LB
	}
	doingCustom := false
	if len(customRotation) > 0 {
		doingCustom = true
		spellOrders = customRotation
	}
	results := []SimResult{}
	var simMetrics SimResult

	lboom := 0
	lboomCount := 0

	prioom := 0
	prioomCount := 0

	pm := func(metrics tbc.SimMetrics, rotation int) {
		casts := make([]CastMetric, 0, len(metrics.Casts))
		for _, v := range metrics.Casts {
			casts = append(casts, CastMetric{
				Spell: v.Spell.ID,
				Hit:   v.DidHit,
				Crit:  v.DidCrit,
				Dmg:   v.DidDmg,
				Time:  float64(v.TicksUntilCast) / float64(tbc.TicksPerSecond),
			})
		}
		if metrics.OOMAt > 0 {
			// DmgAtOOMs
			simMetrics.DmgAtOOMs = append(simMetrics.DmgAtOOMs, metrics.DamageAtOOM)
			// Min/Max tracking
			if rotation == 0 {
				lboom += metrics.OOMAt
				lboomCount++
			} else if rotation == 1 {
				prioom += metrics.OOMAt
				prioomCount++
			}
		}
		simMetrics.OOMAt = append(simMetrics.OOMAt, float64(metrics.OOMAt))
		simMetrics.Casts = append(simMetrics.Casts, casts)
		simMetrics.TotalDmgs = append(simMetrics.TotalDmgs, metrics.TotalDamage)
	}

	dosim := func(spells []string, rotIdx int, simsec int) {
		simMetrics = SimResult{Rotation: spells}
		st := time.Now()
		rseed := time.Now().Unix()
		optNow := opts
		optNow.SpellOrder = spells
		optNow.RSeed = rseed
		sim := tbc.NewSim(stats, equip, optNow)
		for ns := 0; ns < numSims; ns++ {
			metrics := sim.Run(simsec)
			pm(metrics, rotIdx)
		}
		simMetrics.SimSeconds = simsec
		simMetrics.RealDuration = time.Now().Sub(st)
		results = append(results, simMetrics)

	}

	if !doingCustom {
		for i, spells := range spellOrders {
			dosim(spells, i, 600)
		}

		if lboomCount == 0 {
			lboom = 1000
		} else {
			lboom /= lboomCount
		}
		if prioomCount == 0 {
			prioom = 1000
		} else {
			prioom /= prioomCount
		}

		fmt.Printf("Avg LB OOM: %ds, Avg Pri OOM: %ds, Input %ds\n", lboom, prioom, seconds)

		if seconds >= lboom {
			fmt.Printf("LB only is optimal.\n")
			// LB spam is optimal.
			// Probably need to downrank.
			dosim(spellOrders[0], 3, seconds)
		} else if seconds < prioom {
			fmt.Printf("CL always is optimal.\n")
			dosim(spellOrders[1], 3, seconds)
			// Priority spam is optimal
		} else {
			bestRotMetrics, rotation := tbc.OptimalRotation(stats, opts, equip, seconds, numSims)
			simMetrics = SimResult{Rotation: rotation}
			for _, rotmet := range bestRotMetrics {
				pm(rotmet, 2)
			}
			simMetrics.SimSeconds = seconds
			results = append(results, simMetrics)
		}
	} else {
		for i, spells := range spellOrders {
			dosim(spells, i+3, seconds)
		}
	}
	return results
}