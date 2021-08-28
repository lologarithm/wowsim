package tbc

import (
	"testing"
)

// Use same seed to get same result on every run.
var RSeed = int64(1)

var basicOptions = Options{
  RSeed:        RSeed,
  NumBloodlust: 1,
  NumDrums:     0,
  Buffs: Buffs{
    ArcaneInt:                true,
    GiftOftheWild:            true,
    BlessingOfKings:          true,
    ImprovedBlessingOfWisdom: true,
    JudgementOfWisdom:        false,
    Moonkin:                  true,
    SpriestDPS:               0,
    WaterShield:              true,
    Race:                     RaceBonusTroll10,
  },
  Encounter: Encounter{
    Duration: 60,
    NumClTargets: 1,
  },
  Talents: Talents{
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
}

var fullOptions = Options{
  RSeed:        RSeed,
	NumBloodlust: 1,
	NumDrums:     1,
	Buffs: Buffs{
		ArcaneInt:                true,
		GiftOftheWild:            true,
		BlessingOfKings:          true,
		ImprovedBlessingOfWisdom: true,
		JudgementOfWisdom:        true,
		Moonkin:                  true,
		SpriestDPS:               500,
		WaterShield:              true,
		Race:                     RaceBonusOrc,
	},
  Encounter: Encounter{
    Duration: 300,
    NumClTargets: 2,
  },
	Consumes: Consumes{
		FlaskOfBlindingLight: true,
		BrilliantWizardOil:   true,
		MajorMageblood:       false,
		BlackendBasilisk:     true,
    DestructionPotion:    true,
		SuperManaPotion:      true,
		DarkRune:             true,
	},
	Talents: Talents{
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
	Totems: Totems{
		TotemOfWrath: 1,
		WrathOfAir:   true,
		ManaStream:   true,
	},
}

var p1NearBisGear = EquipmentSpec{
  {Name: "Cyclone Faceguard (Tier 4)"},
  {Name: "Adornment of Stolen Souls"},
  {Name: "Cyclone Shoulderguards (Tier 4)"},
  {Name: "Ruby Drape of the Mysticant"},
  {Name: "Netherstrike Breastplate"},
  {Name: "Netherstrike Bracers"},
  {Name: "Soul-Eater's Handwraps"},
  {Name: "Netherstrike Belt"},
  {Name: "Stormsong Kilt"},
  {Name: "Windshear Boots"},
  {Name: "Ring of Unrelenting Storms"},
  {Name: "Ring of Recurrence"},
  {Name: "The Lightning Capacitor"},
  {Name: "Icon of the Silver Crescent"},
  {Name: "Totem of the Void"},
  {Name: "Nathrezim Mindblade"},
  {Name: "Mazthoril Honor Shield"},
}

func TestSimulateP1Basic(t *testing.T) {
  doSimulateTest(t, basicOptions, p1NearBisGear, 1297)
}

func TestSimulateP1Full(t *testing.T) {
  doSimulateTest(t, fullOptions, p1NearBisGear, 1515)
}

func BenchmarkSimulate(b *testing.B) {
	RunSimulation(SimRequest{
		Options: fullOptions,
		Gear: p1NearBisGear,
		Iterations:  b.N,
		IncludeLogs: false,
	})
}

// Performs a basic end-to-end test of the simulator.
//   This is where we can add more sophisticated checks if we would like.
//   Any changes to the damage output of an item set
func doSimulateTest(t *testing.T, options Options, gear EquipmentSpec, expectedDps float64) {
	result := RunSimulation(SimRequest{
		Options: options,
		Gear: gear,
		Iterations:  1,
		IncludeLogs: false,
	})

  tolerance := 0.5
	if result.DpsAvg < expectedDps - tolerance || result.DpsAvg > expectedDps + tolerance {
		t.Fatalf("Expected %0f dps from sim but was: %0f", expectedDps, result.DpsAvg)
	}
}
