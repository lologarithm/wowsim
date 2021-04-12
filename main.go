package main

import (
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
	var runWebUI = flag.Bool("web", false, "Use to run sim in web interface instead of in terminal")
	flag.Parse()

	// this is silly lol
	tbc.IsDebug = *isDebug

	if *runWebUI {
		log.Printf("Closing: %s", http.ListenAndServe(":3333", nil))
	}

	gear := tbc.NewEquipmentSet(
		"Spellstrike Hood",
		"Brooch of Heightened Potential",
		"Pauldrons of Wild Magic",
		"Ogre Slayer's Cover",
		"Tidefury Chestpiece",
		"World's End Bracers",
		"Earth Mantle Handwraps",
		"Wave-Song Girdle",
		"Spellstrike Pants",
		"Magma Plume Boots",
		"Cobalt Band of Tyrigosa",
		"Scintillating Coral Band",
		"Mazthoril Honor Shield",
		"Bleeding Hollow Warhammer",
		// "Quagmirran's Eye",
		"Icon of the Silver Crescent",
		"Totem of the Void",
	)
	// gear[tbc.EquipHead].Gems[0] = tbc.GemLookup["Chaotic Skyfire Diamond"]

	gearStats := gear.Stats()
	fmt.Printf("Gear Stats:\n%s", gearStats.Print(true))

	opt := tbc.Options{
		// NumBloodlust: 1,
		// NumDrums:     2,
		Buffs: tbc.Buffs{
			ArcaneInt:                true,
			GiftOftheWild:            true,
			BlessingOfKings:          true,
			ImprovedBlessingOfWisdom: true,
			JudgementOfWisdom:        true,
			Moonkin:                  false,
			SpriestDPS:               0,
			WaterShield:              true,
		},
		Consumes: tbc.Consumes{
			BrilliantWizardOil: true,
			MajorMageblood:     true,
			BlackendBasilisk:   false,
			SuperManaPotion:    true,
			DarkRune:           true,
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
			TotemOfWrath: 1,
			WrathOfAir:   true,
			ManaStream:   true,
		},
	}

	sims := 10000
	if *isDebug {
		sims = 1
	}
	rotArray := []string{}
	if rotation != nil && len(*rotation) > 0 {
		rotArray = strings.Split(*rotation, ",")
	}

	results := runTBCSim(gear, opt, 300, sims, rotArray, *noopt)
	for _, res := range results {
		fmt.Printf("\n%s\n", res)
	}
}

func runTBCSim(equip tbc.Equipment, opt tbc.Options, seconds int, numSims int, customRotation []string, noopt bool) []string {
	fmt.Printf("\nSim Duration: %d sec\nNum Simulations: %d\n", seconds, numSims)

	stats := opt.StatTotal(equip)

	spellOrders := [][]string{
		// {"CL6", "LB12", "LB12", "LB12"},
		// {"CL6", "LB12", "LB12", "LB12", "LB12"},
		{"pri", "CL6", "LB12"}, // cast CL whenever off CD, otherwise LB
		{"LB12"},               // only LB
	}
	if len(customRotation) > 0 {
		fmt.Printf("Using Custom Rotation: %v\n", customRotation)
		spellOrders = [][]string{customRotation}
	}

	fmt.Printf("\nFinal Stats: %s\n", stats.Print(true))
	statchan := make(chan string, 3)
	for spi, spells := range spellOrders {
		go doSimMetrics(spells, stats, equip, opt, seconds, numSims, statchan)
		if tbc.IsDebug && spi != len(spellOrders)-1 {
			time.Sleep(time.Second * 2)
		}
	}

	results := []string{}
	for i := 0; i < len(spellOrders); i++ {
		results = append(results, <-statchan)
	}

	if !noopt {
		fmt.Printf("\n------- OPTIMIZING -------\n")
		optResult, optimalRotation := tbc.OptimalRotation(stats, opt, equip, seconds, numSims)
		fmt.Printf("\n-------   DONE  -------\n")
		fmt.Printf("Ratio: 1CL : %dLB\n", len(optimalRotation)-1)
		tbc.PrintResult(optResult, seconds)

		opt.UseAI = true
		go doSimMetrics([]string{"AI"}, stats, equip, opt, seconds, numSims, statchan)
		results = append(results, <-statchan)

		// opt.SpellOrder = optimalRotation
		// tbc.OptimalGems(opt, equip, seconds, numSims)

		// Gem it up.... for now
		ruby := tbc.GemLookup["Runed Living Ruby"]
		for i := range equip {
			equip[i].Gems = make([]tbc.Gem, len(equip[i].GemSlots))
			for gs, color := range equip[i].GemSlots {
				if color != tbc.GemColorMeta {
					equip[i].Gems[gs] = ruby
				} else {
					equip[i].Gems[gs] = tbc.Gems[0]
				}
			}
		}

		tbc.StatWeights(opt, equip, seconds, numSims)
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
	ioutil.WriteFile(strings.Join(spo, ""), []byte(out), 0666)

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
