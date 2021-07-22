package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"syscall/js"
	"time"

	"github.com/lologarithm/wowsim/tbc"
)

func main() {
	c := make(chan struct{}, 0)

	simfunc := js.FuncOf(Simulate)
	statfunc := js.FuncOf(StatWeight)
	statComputefunc := js.FuncOf(ComputeStats)
	gearlistfunc := js.FuncOf(GearList)
	packOptfunc := js.FuncOf(PackOptions)

	js.Global().Set("simulate", simfunc)
	js.Global().Set("statweight", statfunc)
	js.Global().Set("computestats", statComputefunc)
	js.Global().Set("gearlist", gearlistfunc)
	js.Global().Set("packopts", packOptfunc)
	js.Global().Call("wasmready")
	<-c
}

func PackOptions(this js.Value, args []js.Value) interface{} {
	opt := parseOptions(args[0])
	packedOpts := opt.Pack()
	js.CopyBytesToJS(js.Global().Get("results").Get(args[1].String()), packedOpts)
	return len(packedOpts)
}

// GearList reports all items of gear to the UI to display.
func GearList(this js.Value, args []js.Value) interface{} {
	slot := byte(128)

	if len(args) == 1 {
		slot = byte(args[0].Int())
	}
	gear := struct {
		Items    []tbc.Item
		Gems     []tbc.Gem
		Enchants []tbc.Enchant
	}{
		Items: make([]tbc.Item, 0, len(tbc.ItemsByID)),
	}
	for _, v := range tbc.ItemsByID {
		if slot != 128 && v.Slot != slot {
			continue
		}
		gear.Items = append(gear.Items, v)
	}
	gear.Gems = tbc.Gems
	gear.Enchants = tbc.Enchants

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
		gearStats = gearStats.CalculatedTotal()
		out, err := json.Marshal(gearStats)
		if err != nil {
			fmt.Printf("Failed to format JSON output: %s\n", err)
		}
		return string(out)
	}
	opt := parseOptions(args[1])
	stats := tbc.CalculateTotalStats(opt, gear)
	opt.UseAI = true // stupid complaining sim...maybe I should just default AI on.
	fakesim := tbc.NewSim(stats, gear, opt)
	sets := fakesim.ActivateSets()

	finalStats := stats
	for i, v := range fakesim.Buffs {
		finalStats[i] += v
	}
	out, err := json.Marshal(struct {
		Stats []float64
		Sets  []string
	}{Stats: finalStats, Sets: sets})
	if err != nil {
		fmt.Printf("Failed to format JSON output: %s\n", err)
	}
	return string(out)
}

// Simulate takes in number of iterations, duration, a gear list, and simulation options.
// (iterations, duration, gearlist, options, <optional, custom rotation)
func StatWeight(this js.Value, args []js.Value) interface{} {
	numSims := args[0].Int()
	seconds := args[1].Int()
	numClTargets := args[2].Int()
	gear := getGear(args[3])
	opts := parseOptions(args[4])
	stat := args[5].Int()
	statModVal := args[6].Float()

	if len(opts.Buffs.Custom) == 0 {
		opts.Buffs.Custom = tbc.Stats{tbc.StatLen: 0}
	}
	opts.Buffs.Custom[tbc.Stat(stat)] += statModVal
	opts.UseAI = true // use AI optimal rotation.

	if numSims == 1 {
		opts.Debug = true
	}

	simdmg := 0.0
	simdmgsq := 0.0
	simmet := make([]tbc.SimMetrics, 0, numSims)

	opts.RSeed = time.Now().Unix()
	opts.NumClTargets = numClTargets

	oomcount := 0
	sim := tbc.NewSim(tbc.CalculateTotalStats(opts, gear), gear, opts)
	for ns := 0; ns < numSims; ns++ {
		metrics := sim.Run(seconds)
		dps := metrics.TotalDamage / float64(seconds)
		simdmg += dps
		simdmgsq += dps * dps
		simmet = append(simmet, metrics)
		if metrics.OOMAt > 0 && metrics.OOMAt < seconds-5 {
			oomcount++
		}
	}

	meanSq := simdmgsq / float64(numSims)
	mean := simdmg / float64(numSims)
	stdev := math.Sqrt(meanSq - mean*mean)
	fmt.Printf("(Mod: %s) Mean: %0.1f, Stddev: %0.1f\n", tbc.Stat(stat).StatName(), mean, stdev)

	conf90 := 1.645 * stdev / math.Sqrt(float64(numSims))

	if float64(oomcount)/float64(numSims) > 0.25 {
		fmt.Printf("Many oom, invalid stat weights may be seen.")
		return fmt.Sprintf("%0.2f,%0.2f,%0.2f", mean, stdev, conf90)
	}

	return fmt.Sprintf("%0.2f,%0.2f,%0.2f", mean, stdev, conf90)
}

// Simulate takes in number of iterations, duration, a gear list, and simulation options.
// (iterations, duration, gearlist, options, <optional, custom rotation)
func Simulate(this js.Value, args []js.Value) interface{} {
	if len(args) < 5 {
		print("Expected 5 min arguments:  (#iterations, duration, numClTargets, gearlist, options)")
		return `{"error": "invalid arguments supplied"}`
	}

	customRotation := [][]string{}
	customHaste := 0.0
	if len(args) >= 7 {
		if args[5].Truthy() {
			customRotation = parseRotation(args[5])
		}
		if args[6].Truthy() {
			customHaste = args[6].Float()
		}
	}
	gear := getGear(args[3])
	opt := parseOptions(args[4])
	stats := tbc.CalculateTotalStats(opt, gear)
	if customHaste != 0 {
		stats[tbc.StatHaste] = customHaste
	}

	simi := args[0].Int()
	if simi == 1 { // if single iteration, dump all logs to console.
		opt.Debug = true
	}
	dur := args[1].Int()
	numClTargets := args[2].Int()
	opt.NumClTargets = numClTargets
	fullLogs := false
	if len(args) > 7 {
		fullLogs = args[7].Truthy()
		fmt.Printf("Building Full Log:%v\n", fullLogs)
	}

	results := runTBCSim(opt, stats, gear, dur, simi, customRotation, fullLogs)
	st := time.Now()
	output, err := json.Marshal(results)
	if err != nil {
		print("Failed to json marshal results: ", err.Error())
	}
	fmt.Printf("Took %s to json marshal response.\n", time.Now().Sub(st))
	return string(output)
}

// getGear converts js string array to a list of equipment items.
func getGear(val js.Value) tbc.Equipment {
	numGear := val.Length()
	gearSet := make([]tbc.Item, numGear)
	for i := range gearSet {
		v := val.Index(i)
		name := v.Get("Name")
		id := v.Get("ID")
		var ic tbc.Item
		if !name.IsUndefined() {
			ic = tbc.ItemsByName[name.String()]
		} else if !id.IsUndefined() {
			ic = tbc.ItemsByID[int32(id.Int())]
		}
		gems := v.Get("Gems")
		gemids := v.Get("g")
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
		} else if !(gemids.IsUndefined() || gemids.IsNull()) && gemids.Length() > 0 {
			ic.Gems = make([]tbc.Gem, len(ic.GemSlots))
			for i := range ic.Gems {
				jsgem := gemids.Index(i)
				if jsgem.IsNull() || jsgem.IsUndefined() {
					continue
				}
				gv, ok := tbc.GemsByID[int32(jsgem.Int())]
				if !ok {
					continue // wasn't a valid gem
				}
				ic.Gems[i] = gv
			}
		}
		if !v.Get("Enchant").IsNull() && !v.Get("Enchant").IsUndefined() {
			ic.Enchant = tbc.EnchantLookup[v.Get("Enchant").String()]
		} else if !v.Get("e").IsNull() && !v.Get("e").IsUndefined() {
			ic.Enchant = tbc.EnchantByID[int32(v.Get("e").Int())]
		}
		gearSet[i] = ic
	}
	return tbc.Equipment(gearSet)
}

func parseOptions(val js.Value) tbc.Options {
	var custom = val.Get("custom")
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
			ImpSealofCrusader:        val.Get("debuffisoc").Truthy(),
			Misery:                   val.Get("debuffmis").Truthy(),
			Moonkin:                  val.Get("buffmoon").Truthy(),
			MoonkinRavenGoddess:      val.Get("buffmoonrg").Truthy(),
			SpriestDPS:               uint16(val.Get("buffspriest").Int()),
			WaterShield:              val.Get("sbufws").Truthy(),
			EyeOfNight:               val.Get("buffeyenight").Truthy(),
			TwilightOwl:              val.Get("bufftwilightowl").Truthy(),
			Race:                     tbc.RaceBonusType(val.Get("sbufrace").Int()),
			Custom: tbc.Stats{
				tbc.StatInt:       custom.Get("custint").Float(),
				tbc.StatSpellCrit: custom.Get("custsc").Float(),
				tbc.StatSpellHit:  custom.Get("custsh").Float(),
				tbc.StatSpellDmg:  custom.Get("custsp").Float(),
				tbc.StatHaste:     custom.Get("custha").Float(),
				tbc.StatMP5:       custom.Get("custmp5").Float(),
				tbc.StatMana:      custom.Get("custmana").Float(),
			},
		},
		Consumes: tbc.Consumes{
			FlaskOfBlindingLight:     val.Get("confbl").Truthy(),
			FlaskOfMightyRestoration: val.Get("confmr").Truthy(),
			BrilliantWizardOil:       val.Get("conbwo").Truthy(),
			MajorMageblood:           val.Get("conmm").Truthy(),
			BlackendBasilisk:         val.Get("conbb").Truthy(),
			DestructionPotion:        val.Get("condp").Truthy(),
			SuperManaPotion:          val.Get("consmp").Truthy(),
			DarkRune:                 val.Get("condr").Truthy(),
		},
		Talents: tbc.Talents{
			LightningOverload:   5,
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
			TotemOfWrath: val.Get("totwr").Int(),
			WrathOfAir:   val.Get("totwoa").Truthy(),
			Cyclone2PC:   val.Get("totcycl2p").Truthy(),
			ManaStream:   val.Get("totms").Truthy(),
		},
		DPSReportTime: val.Get("dpsReportTime").Int(),
		GCD:           val.Get("gcd").Float(),
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

type SimResult struct {
	Rotation     []string
	SimSeconds   int
	RealDuration float64
	Logs         string
	DPSAvg       float64              `json:"dps"`
	DPSDev       float64              `json:"dev"`
	MaxDPS       float64              `json:"max"`
	OOMAt        float64              `json:"oomat"`
	NumOOM       int                  `json:"numOOM"`
	DPSAtOOM     float64              `json:"dpsAtOOM"`
	Casts        map[int32]CastMetric `json:"casts"`
	DPSHist      map[int]int          `json:"dpsHist"` // rounded DPS to count
}

type CastMetric struct {
	Count int     `json:"count"`
	Dmg   float64 `json:"dmg"`
	Crits int     `json:"crits"`
}

func runTBCSim(opts tbc.Options, stats tbc.Stats, equip tbc.Equipment, seconds int, numSims int, customRotation [][]string, fullLogs bool) []SimResult {
	print("\nSim Duration:", seconds)
	print("\nNum Simulations: ", numSims)
	print("\n")

	spellOrders := [][]string{}
	doingCustom := false
	if len(customRotation) > 0 {
		doingCustom = true
		spellOrders = customRotation
	}
	results := []SimResult{}
	logsBuffer := &strings.Builder{}

	dosim := func(spells []string, simsec int) {
		simMetrics := SimResult{
			DPSHist:  map[int]int{},
			Casts:    map[int32]CastMetric{},
			Rotation: spells,
		}
		if opts.UseAI {
			simMetrics.Rotation = []string{"AI Optimized"}
		}
		st := time.Now()
		rseed := time.Now().Unix()
		optNow := opts
		optNow.SpellOrder = spells
		optNow.RSeed = rseed
		sim := tbc.NewSim(stats, equip, optNow)

		var totalSq float64
		for ns := 0; ns < numSims; ns++ {
			if fullLogs {
				sim.Debug = func(s string, vals ...interface{}) {
					logsBuffer.WriteString(fmt.Sprintf("[%0.1f] "+s, append([]interface{}{(float64(sim.CurrentTick) / float64(tbc.TicksPerSecond))}, vals...)...))
				}
			}
			metrics := sim.Run(simsec)
			dps := metrics.TotalDamage / float64(simsec)
			if opts.DPSReportTime > 0 {
				dps = metrics.ReportedDamage / float64(opts.DPSReportTime)
			}
			totalSq += dps * dps
			simMetrics.DPSAvg += dps
			dpsRounded := int(math.Round(dps/10) * 10)
			simMetrics.DPSHist[dpsRounded] += 1
			if dps > simMetrics.MaxDPS {
				simMetrics.MaxDPS = dps
			}
			if (metrics.OOMAt) > 0 {
				simMetrics.OOMAt += float64(metrics.OOMAt)
				simMetrics.DPSAtOOM += float64(metrics.DamageAtOOM) / float64(metrics.OOMAt)
				simMetrics.NumOOM++
			}
			for _, cast := range metrics.Casts {
				var id = cast.Spell.ID
				if cast.IsLO {
					id = 1000 - cast.Spell.ID
				}
				cm := simMetrics.Casts[id]
				cm.Count++
				cm.Dmg += cast.DidDmg
				if cast.DidCrit {
					cm.Crits++
				}
				simMetrics.Casts[id] = cm
			}

		}

		meanSq := totalSq / float64(numSims)
		mean := simMetrics.DPSAvg / float64(numSims)
		stdev := math.Sqrt(meanSq - mean*mean)

		simMetrics.DPSDev = stdev
		simMetrics.DPSAvg /= float64(numSims)
		if simMetrics.NumOOM > 0 {
			simMetrics.OOMAt /= float64(simMetrics.NumOOM)
			simMetrics.DPSAtOOM /= float64(simMetrics.NumOOM)
		}

		simMetrics.Logs = logsBuffer.String()
		simMetrics.SimSeconds = simsec
		simMetrics.RealDuration = time.Now().Sub(st).Seconds()
		results = append(results, simMetrics)
	}

	if !doingCustom && opts.UseAI {
		dosim([]string{"AI Optimized"}, seconds) // Let AI determine best possible DPS
	} else {
		for _, spells := range spellOrders {
			dosim(spells, seconds)
		}
	}
	return results
}
