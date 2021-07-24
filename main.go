package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lologarithm/wowsim/tbc"
)

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
	var rotation = flag.String("rotation", "", "Custom comma separated rotation to simulate.\n\tFor Example: --rotation=CL6,LB12")
	var duration = flag.Int("duration", 300, "Custom fight duration in seconds.")
	var iterations = flag.Int("iter", 10000, "Custom number of iterations for the sim to run.")
	var runWebUI = flag.Bool("web", false, "Use to run sim in web interface instead of in terminal")
	var configFile = flag.String("config", "", "Specify an input configuration.")

	flag.Parse()

	if *runWebUI {
		log.Printf("Closing: %s", http.ListenAndServe(":3333", nil))
	}

	// Just some default gear if not provided in the config file.
	gear := tbc.NewEquipmentSet(
		"Tidefury Helm",
		"Charlotte's Ivy",
		"Pauldrons of Wild Magic",
		"Ogre Slayer's Cover",
		"Tidefury Chestpiece",
		"World's End Bracers",
		"Earth Mantle Handwraps",
		"Netherstrike Belt",
		"Stormsong Kilt",
		"Magma Plume Boots",
		"Cobalt Band of Tyrigosa",
		"Sparking Arcanite Ring",
		"Mazthoril Honor Shield",
		"Gavel of Unearthed Secrets",
		"Natural Alignment Crystal",
		"Icon of the Silver Crescent",
		"Totem of the Void",
	)

	// Auto gem the default gear above.
	ruby := tbc.GemLookup["Runed Living Ruby"]
	for i := range gear {
		gear[i].Gems = make([]tbc.Gem, len(gear[i].GemSlots))
		for gs, color := range gear[i].GemSlots {
			if color != tbc.GemColorMeta {
				gear[i].Gems[gs] = ruby
			} else {
				gear[i].Gems[gs] = tbc.Gems[0] // CSD
			}
		}
	}

	opt := tbc.Options{
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
			TotemOfWrath: 1,
			WrathOfAir:   false,
			ManaStream:   true,
		},
	}

	if *configFile != "" {
		data, err := ioutil.ReadFile(*configFile)
		if err != nil {
			log.Fatalf("Failed to open config file(%s): %s", *configFile, err)
		}
		gear, opt = getGear(data)
	}

	if *isDebug {
		*iterations = 1
		opt.Debug = true
	}
	rotArray := []string{}
	if rotation != nil && len(*rotation) > 0 {
		rotArray = strings.Split(*rotation, ",")
	}

	results := runTBCSim(gear, opt, *duration, *iterations, rotArray, *noopt)
	for _, res := range results {
		fmt.Printf("\n%s\n", res)
	}
}

type input struct {
	Options tbc.Options
	Gear    []tbc.Item
}

func getGear(val []byte) (tbc.Equipment, tbc.Options) {
	in := &input{}
	err := json.Unmarshal(val, in)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}
	gearSet := make([]tbc.Item, len(in.Gear))
	for i, v := range in.Gear {
		itemTemplate := tbc.ItemsByName[v.Name]
		if v.Name == "" && v.ID > 0 {
			itemTemplate = tbc.ItemsByID[v.ID]
		}
		ic := itemTemplate
		if len(v.Gems) > 0 {
			ic.Gems = make([]tbc.Gem, len(ic.GemSlots))
			for _, gem := range ic.Gems {
				gv, ok := tbc.GemLookup[gem.Name]
				if !ok {
					continue // wasn't a valid gem
				}
				ic.Gems[i] = gv
			}
		}
		ic.Enchant = tbc.EnchantLookup[v.Enchant.Name]
		gearSet[i] = ic
	}
	return tbc.Equipment(gearSet), in.Options
}

func runTBCSim(equip tbc.Equipment, opt tbc.Options, seconds int, numSims int, customRotation []string, noopt bool) []string {
	fmt.Printf("\nSim Duration: %d sec\nNum Simulations: %d\n", seconds, numSims)

	stats := tbc.CalculateTotalStats(opt, equip)

	spellOrders := [][]string{
		// {"CL6", "LB12", "LB12", "LB12"},
		// {"CL6", "LB12", "LB12", "LB12", "LB12"},
		// {"CL6", "LB12", "LB12", "LB12", "LB12", "LB12"},
		// {"pri", "CL6", "LB12"}, // cast CL whenever off CD, otherwise LB
		// {"LB12"},               // only LB
	}
	if len(customRotation) > 0 {
		fmt.Printf("Using Custom Rotation: %v\n", customRotation)
		spellOrders = [][]string{customRotation}
	}

	fmt.Printf("\nFinal Stats: %s\n", stats.Print())
	statchan := make(chan string, 3)
	for spi, spells := range spellOrders {
		go doSimMetrics(spells, stats, equip, opt, seconds, numSims, statchan)
		if opt.Debug && spi != len(spellOrders)-1 {
			time.Sleep(time.Second * 2)
		}
	}

	results := []string{}
	for i := 0; i < len(spellOrders); i++ {
		results = append(results, <-statchan)
	}

	opt.UseAI = true
	go doSimMetrics([]string{"AI"}, stats, equip, opt, seconds, numSims, statchan)
	results = append(results, <-statchan)

	if !noopt {
		// fmt.Printf("\n------- OPTIMIZING -------\n")
		// optResult, optimalRotation := tbc.OptimalRotation(stats, opt, equip, seconds, numSims)
		// fmt.Printf("\n-------   DONE  -------\n")
		// fmt.Printf("Ratio: 1CL : %dLB\n", len(optimalRotation)-1)
		// tbc.PrintResult(optResult, seconds)

		tbc.OptimalGems(opt, equip, seconds, numSims)
		weights := tbc.StatWeights(opt, equip, seconds, numSims)
		// fmt.Printf("Weights: [ SP: %0.2f,  Int: %0.2f,  Crit: %0.2f,  Hit: %0.2f,  Haste: %0.2f,  MP5: %0.2f ]\n", weights[0], weights[1], weights[2], weights[3], weights[4], weights[5])
		fmt.Printf("Weights: [\n")
		for i, v := range weights {
			if tbc.Stat(i) == tbc.StatStm {
				continue
			}
			fmt.Printf("%s: %0.2f\t", tbc.Stat(i).StatName(), v)
		}
		fmt.Printf("\n]\n")
	}

	return results
}

func doSimMetrics(spo []string, stats tbc.Stats, equip tbc.Equipment, opt tbc.Options, seconds int, numSims int, statchan chan string) {
	simDmgs := []float64{}
	simOOMs := []int{}
	histogram := map[int]int{}
	casts := map[int32]int{}
	manaSpent := 0.0
	manaLeft := 0.0
	oomdps := 0.0
	ooms := 0
	numOoms := 0

	rseed := time.Now().Unix()
	opt.SpellOrder = spo
	opt.RSeed = rseed
	sim := tbc.NewSim(stats, equip, opt)
	for ns := 0; ns < numSims; ns++ {
		metrics := sim.Run(seconds)
		simDmgs = append(simDmgs, metrics.TotalDamage)
		simOOMs = append(simOOMs, metrics.OOMAt)
		manaLeft += float64(metrics.ManaAtEnd)
		oomdps += metrics.DamageAtOOM

		ooms += metrics.OOMAt
		if metrics.OOMAt > 0 {
			numOoms++
		}

		for _, cast := range metrics.Casts {
			casts[cast.Spell.ID] += 1
			manaSpent += cast.ManaCost
		}

		rv := int(math.Round(math.Round(metrics.TotalDamage/float64(seconds))/10) * 10)
		histogram[rv] += 1
	}

	oomdps /= float64(numOoms)

	// TODO: do this better... for now just dumping histograph data to disk lol.
	out := ""
	for k, v := range histogram {
		out += strconv.Itoa(k) + "," + strconv.Itoa(v) + "\n"
	}
	// ioutil.WriteFile(strings.Join(spo, ""), []byte(out), 0666)

	totalDmg := 0.0
	tdSq := totalDmg
	max := 0.0
	for _, dmg := range simDmgs {
		totalDmg += dmg
		tdSq += dmg * dmg

		if dmg > max {
			max = dmg
		}
	}

	meanSq := tdSq / float64(numSims)
	mean := totalDmg / float64(numSims)
	stdev := math.Sqrt(meanSq - mean*mean)

	output := ""
	output += fmt.Sprintf("Spell Order: %v\n", spo)
	output += fmt.Sprintf("DPS:")
	output += fmt.Sprintf("\tMean: %0.1f +/- %0.1f\n", (mean / float64(seconds)), stdev/float64(seconds))
	output += fmt.Sprintf("\tMax: %0.1f\n", (max / float64(seconds)))
	output += fmt.Sprintf("Total Casts:\n")

	for k, v := range casts {
		output += fmt.Sprintf("\t%s: %d\n", tbc.AuraName(k), v/numSims)
	}
	// output += fmt.Sprintf("Avg Mana Spent: %d\n", int(manaSpent)/numSims)
	// output += fmt.Sprintf("Avg Mana Left: %d\n", int(manaLeft)/numSims)

	// avgleft := (manaLeft) / float64(numSims)
	// extraCL := int(avgleft / 414) // 414 is cost of difference casting CL instead of LB
	// output += fmt.Sprintf("Add CL: %d\n", extraCL)

	avgoomSec := 0
	if numOoms > 0 {
		avgoomSec = ooms / numOoms
	}
	output += fmt.Sprintf("Went OOM: %d/%d sims\n", numOoms, numSims)
	if numOoms > 0 {
		output += fmt.Sprintf("Avg OOM Time: %d seconds\n", avgoomSec)
		output += fmt.Sprintf("Avg DPS At OOM: %0.0f\n", oomdps/float64(avgoomSec))
	}
	statchan <- output
}
