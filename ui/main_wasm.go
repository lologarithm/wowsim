package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"

	"github.com/lologarithm/wowsim/tbc"
)

func main() {
	c := make(chan struct{}, 0)

	simfunc := js.FuncOf(Simulate)
	statfunc := js.FuncOf(StatWeight)
	statsfunc := js.FuncOf(ComputeStats)
	gearlistfunc := js.FuncOf(GearList)

	js.Global().Set("simulate", simfunc)
	js.Global().Set("statweight", statfunc)
	js.Global().Set("computestats", statsfunc)
	js.Global().Set("gearlist", gearlistfunc)
	js.Global().Call("wasmready")
	<-c
}

// GearList reports all items of gear to the UI to display.
func GearList(this js.Value, args []js.Value) interface{} {
	slot := byte(128)

	if len(args) == 1 {
		slot = byte(args[0].Int())
	}
	gear := struct {
		Items []tbc.Item
		Gems  []tbc.Gem
	}{
		Items: make([]tbc.Item, 0, len(tbc.ItemLookup)),
	}
	for _, v := range tbc.ItemLookup {
		if slot != 128 && v.Slot != slot {
			continue
		}
		gear.Items = append(gear.Items, v)
	}
	gear.Gems = tbc.Gems
	output, err := json.Marshal(gear)
	if err != nil {
		// fmt.Printf("Failed to marshal gear list: %s", err)
		output = []byte(`{"error": ` + err.Error() + `}`)
	}
	// fmt.Printf("Item Output: %s", string(output))
	return string(output)
}

// GearStats takes a gear list and returns their total stats.
// This could power a simple 'current stats of all gear' UI.
func ComputeStats(this js.Value, args []js.Value) interface{} {
	gear := getGear(args[0])
	if len(args) != 2 {
		return `{"error": "incorrect args. expected computestats(gear, options)}`
	}
	if args[1].IsNull() {
		gearStats := gear.Stats()
		gearStats[tbc.StatSpellCrit] += (gearStats[tbc.StatInt] / 80) * 22.08
		gearStats[tbc.StatMana] += gearStats[tbc.StatInt] * 15
		return gearStats.Print(false)
	}
	opt := parseOptions(args[1])
	stats := opt.StatTotal(gear)
	return stats.Print(false)
}

// getGear converts js string array to a list of equipment items.
func getGear(val js.Value) tbc.Equipment {
	numGear := val.Length()
	gearSet := make([]tbc.Item, numGear)
	for i := range gearSet {
		v := val.Index(i)
		ic := tbc.ItemLookup[v.Get("Name").String()]
		gems := v.Get("Gems")
		if !(gems.IsUndefined() || gems.IsNull()) && gems.Length() > 0 {
			ic.Gems = make([]tbc.Gem, len(ic.GemSlots))
			for i := range ic.Gems {
				jsgem := gems.Index(i)
				if jsgem.IsNull() {
					continue
				}
				gv, ok := tbc.GemLookup[jsgem.String()]
				if !ok {
					continue // wasn't a valid gem
				}
				ic.Gems[i] = gv
			}
		}
		gearSet[i] = ic
	}
	return tbc.Equipment(gearSet)
}

func parseOptions(val js.Value) tbc.Options {
	opt := tbc.Options{
		ExitOnOOM:    val.Get("exitoom").Truthy(),
		NumBloodlust: val.Get("buffbl").Int(),
		NumDrums:     val.Get("buffdrums").Int(),
		UseAI:        val.Get("useai").Truthy(),
		Buffs: tbc.Buffs{
			ArcaneInt:                val.Get("buffai").Truthy(),
			GiftOftheWild:            val.Get("buffgotw").Truthy(),
			BlessingOfKings:          val.Get("buffbk").Truthy(),
			ImprovedBlessingOfWisdom: val.Get("buffibow").Truthy(),
			ImprovedDivineSpirit:     val.Get("buffids").Truthy(),
			JudgementOfWisdom:        val.Get("debuffjow").Truthy(),
			Misery:                   val.Get("debuffmis").Truthy(),
			Moonkin:                  val.Get("buffmoon").Truthy(),
			MoonkinRavenGoddess:      val.Get("buffmoonrg").Truthy(),
			SpriestDPS:               val.Get("buffspriest").Int(),
			WaterShield:              val.Get("sbufws").Truthy(),
			EyeOfNight:               val.Get("buffeyenight").Truthy(),
			TwilightOwl:              val.Get("bufftwilightowl").Truthy(),
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
			Cyclone2PC:   val.Get("totcycl2p").Truthy(),
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
func StatWeight(this js.Value, args []js.Value) interface{} {
	numSims := args[0].Int()
	if numSims == 1 {
		tbc.IsDebug = true
	}
	seconds := args[1].Int()
	gear := getGear(args[2])
	opts := parseOptions(args[3])
	stat := args[4].Int()
	statModVal := args[5].Float()

	opts.Buffs.Custom = tbc.Stats{tbc.StatLen: 0}
	opts.Buffs.Custom[tbc.Stat(stat)] += statModVal
	opts.UseAI = true // use AI optimal rotation.

	simdmg := 0.0
	fmt.Printf("StatWeight(%d, %0.0f) Stats: %s\n", stat, statModVal, opts.StatTotal(gear).Print(false))
	simmet := make([]tbc.SimMetrics, 0, numSims)

	// allcasts := []*tbc.Cast{}
	for ns := 0; ns < numSims; ns++ {
		sim := tbc.NewSim(opts.StatTotal(gear), gear, opts)
		metrics := sim.Run(seconds)
		simdmg += metrics.TotalDamage
		simmet = append(simmet, metrics)
	}

	return simdmg / float64(numSims) / float64(seconds)
}

// Simulate takes in number of iterations, duration, a gear list, and simulation options.
// (iterations, duration, gearlist, options, <optional, custom rotation)
func Simulate(this js.Value, args []js.Value) interface{} {
	if len(args) < 4 {
		print("Expected 4 min arguments:  (#iterations, duration, gearlist, options)")
		return `{"error": "invalid arguments supplied"}`
	}

	customRotation := [][]string{}
	customHaste := 0.0
	if len(args) == 6 {
		if args[4].Truthy() {
			customRotation = parseRotation(args[4])
		}
		if args[5].Truthy() {
			customHaste = args[5].Float()
		}
	}
	gear := getGear(args[2])
	opt := parseOptions(args[3])
	doOptimize := args[3].Get("doopt").Truthy() || opt.UseAI
	stats := opt.StatTotal(gear)
	if customHaste != 0 {
		stats[tbc.StatHaste] = customHaste
	}

	simi := args[0].Int()
	if simi == 1 {
		tbc.IsDebug = true
	}
	dur := args[1].Int()

	results := runTBCSim(opt, stats, gear, dur, simi, customRotation, doOptimize)
	st := time.Now()
	output, err := json.Marshal(results)
	if err != nil {
		print("Failed to json marshal results: ", err.Error())
	}
	fmt.Printf("Took %s to json marshal response.\n", time.Now().Sub(st))
	return string(output)
}

type SimResult struct {
	Casts        [][]CastMetric
	TotalDmgs    []float64
	DmgAtOOMs    []float64
	OOMAt        []float64 // oom time totals
	Rotation     []string
	SimSeconds   int
	RealDuration float64
}

type CastMetric struct {
	ID   int32
	Hit  bool
	Crit bool
	IsLO bool
	Dmg  float64
	Time float64 // seconds it took to cast this spell
}

func runTBCSim(opts tbc.Options, stats tbc.Stats, equip tbc.Equipment, seconds int, numSims int, customRotation [][]string, doOptimize bool) []SimResult {
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

	pm := func(metrics tbc.SimMetrics) {
		casts := make([]CastMetric, 0, len(metrics.Casts))
		for _, v := range metrics.Casts {
			casts = append(casts, CastMetric{
				ID:   v.Spell.ID,
				IsLO: v.IsLO,
				Hit:  v.DidHit,
				Crit: v.DidCrit,
				Dmg:  v.DidDmg,
				Time: float64(v.TicksUntilCast) / float64(tbc.TicksPerSecond),
			})
		}
		if metrics.OOMAt > 0 {
			// DmgAtOOMs
			simMetrics.DmgAtOOMs = append(simMetrics.DmgAtOOMs, metrics.DamageAtOOM)
		}
		simMetrics.OOMAt = append(simMetrics.OOMAt, float64(metrics.OOMAt))
		simMetrics.Casts = append(simMetrics.Casts, casts)
		simMetrics.TotalDmgs = append(simMetrics.TotalDmgs, metrics.TotalDamage)
	}

	dosim := func(spells []string, simsec int) {
		simMetrics = SimResult{Rotation: spells}
		if opts.UseAI {
			simMetrics.Rotation = []string{"AI Optimized"}
		}
		st := time.Now()
		rseed := time.Now().Unix()
		optNow := opts
		optNow.SpellOrder = spells
		optNow.RSeed = rseed
		for ns := 0; ns < numSims; ns++ {
			sim := tbc.NewSim(stats, equip, optNow)
			metrics := sim.Run(simsec)
			pm(metrics)
		}
		simMetrics.SimSeconds = simsec
		simMetrics.RealDuration = time.Now().Sub(st).Seconds()
		results = append(results, simMetrics)
	}

	if !doingCustom && doOptimize {
		opts.UseAI = true
		dosim([]string{"AI Optimized"}, seconds) // Let AI determine best possible DPS
	} else {
		for _, spells := range spellOrders {
			dosim(spells, seconds)
		}
	}
	return results
}
