package tbc

import (
	"math"
	"testing"
)

// TestSimulate is a basic end-to-end test of the simulator.
//   This is where we can add more sophisticated checks if we would like.
//   Any changes to the damage output of an item set
func TestSimulate(t *testing.T) {
	options := Options{
		RSeed:        1, // Use same seed to get same result on every run.
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
		Encounter: Encounter{Duration: 60},
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
	result := RunSimulation(SimRequest{
		Options: options,
		Gear: EquipmentSpec{
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
		},
		Iterations:  1,
		IncludeLogs: false,
	})

	// Rounding because floats are dumb
	if math.Round(result.DpsAvg) != 1303 {
		t.Fatalf("Incorrect total damage in simulator: %0f", result.DpsAvg)
	}
}
