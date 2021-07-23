package tbc

import (
	"math"
)

// AuraEffects will mutate a cast or simulation state.
type AuraEffect func(sim *Simulation, c *Cast)

type Aura struct {
	ID      int32
	Expires int // ticks aura will apply

	OnCast         AuraEffect
	OnCastComplete AuraEffect
	OnSpellHit     AuraEffect
	OnStruck       AuraEffect
	OnExpire       AuraEffect
	OnSpellMiss    AuraEffect
}

func AuraName(a int32) string {
	switch a {
	case MagicIDUnknown:
		return "Unknown"
	case MagicIDClBounce:
		return "Chain Lightning Bounce"
	case MagicIDLOTalent:
		return "Lightning Overload Talent"
	case MagicIDJoW:
		return "Judgement Of Wisdom Aura"
	case MagicIDEleFocus:
		return "Elemental Focus"
	case MagicIDEleMastery:
		return "Elemental Mastery"
	case MagicIDStormcaller:
		return "Stormcaller"
	case MagicIDBlessingSilverCrescent:
		return "Blessing of the Silver Crescent"
	case MagicIDQuagsEye:
		return "Quags Eye"
	case MagicIDFungalFrenzy:
		return "Fungal Frenzy"
	case MagicIDBloodlust:
		return "Bloodlust"
	case MagicIDSkycall:
		return "Skycall"
	case MagicIDEnergized:
		return "Energized"
	case MagicIDNAC:
		return "Nature Alignment Crystal"
	case MagicIDChaoticSkyfire:
		return "Chaotic Skyfire"
	case MagicIDInsightfulEarthstorm:
		return "Insightful Earthstorm"
	case MagicIDMysticSkyfire:
		return "Mystic Skyfire"
	case MagicIDMysticFocus:
		return "Mystic Focus"
	case MagicIDEmberSkyfire:
		return "Ember Skyfire"
	case MagicIDLB12:
		return "LB12"
	case MagicIDCL6:
		return "CL6"
	case MagicIDTLCLB:
		return "TLC-LB"
	case MagicIDISCTrink:
		return "Trink"
	case MagicIDNACTrink:
		return "NACTrink"
	case MagicIDPotion:
		return "Potion"
	case MagicIDRune:
		return "Rune"
	case MagicIDAllTrinket:
		return "AllTrinket"
	case MagicIDSpellPower:
		return "SpellPower"
	case MagicIDRubySerpent:
		return "RubySerpent"
	case MagicIDCallOfTheNexus:
		return "CallOfTheNexus"
	case MagicIDDCC:
		return "Darkmoon Card Crusade"
	case MagicIDDCCBonus:
		return "Aura of the Crusade"
	case MagicIDScryerTrink:
		return "Scryer Trinket"
	case MagicIDRubySerpentTrink:
		return "Ruby Serpent Trinket"
	case MagicIDXiriTrink:
		return "Xiri Trinket"
	case MagicIDDrums:
		return "Drums of Battle"
	case MagicIDDrum1:
		return "Drum #1"
	case MagicIDDrum2:
		return "Drum #2"
	case MagicIDDrum3:
		return "Drum #3"
	case MagicIDDrum4:
		return "Drum #4"
	case MagicIDNetherstrike:
		return "Netherstrike Set"
	case MagicIDTwinStars:
		return "Twin Stars Set"
	case MagicIDTidefury:
		return "Tidefury Set"
	case MagicIDSpellstrike:
		return "Spellstrike Set"
	case MagicIDSpellstrikeInfusion:
		return "Spellstrike Infusion"
	case MagicIDManaEtched:
		return "Mana-Etched Set"
	case MagicIDManaEtchedHit:
		return "Mana-EtchedHit"
	case MagicIDManaEtchedInsight:
		return "Mana-EtchedInsight"
	case MagicIDWindhawk:
		return "Windhawk Set Bonus"
	case MagicIDOrcBloodFury:
		return "Orc Blood Fury"
	case MagicIDTrollBerserking:
		return "Troll Berserking"
	case MagicIDEyeOfTheNight:
		return "EyeOfTheNight"
	case MagicIDChainTO:
		return "Chain of the Twilight Owl"
	case MagicIDCyclone2pc:
		return "Cyclone 2pc Bonus"
	case MagicIDCyclone4pc:
		return "Cyclone 4pc Bonus"
	case MagicIDCycloneMana:
		return "Cyclone Mana Cost Reduction"
	case MagicIDTLC:
		return "The Lightning Capacitor Aura"
	case MagicIDDestructionPotion:
		return "Destruction Potion"
	case MagicIDHexShunkHead:
		return "Hex Shunken Head"
	case MagicIDShiftingNaaru:
		return "Shifting Naaru Sliver"
	case MagicIDSkullGuldan:
		return "Skull of Guldan"
	case MagicIDNexusHorn:
		return "Nexus-Horn"
	case MagicIDSextant:
		return "Sextant of Unstable Currents"
	case MagicIDUnstableCurrents:
		return "Unstable Currents"
	case MagicIDEyeOfMag:
		return "Eye Of Mag"
	case MagicIDRecurringPower:
		return "Recurring Power"
	case MagicIDCataclysm4pc:
		return "Cataclysm 4pc Set Bonus"
	case MagicIDSkyshatter2pc:
		return "Skyshatter 2pc Set Bonus"
	case MagicIDSkyshatter4pc:
		return "Skyshatter 4pc Set Bonus"
	case MagicIDTotemOfPulsingEarth:
		return "Totem of Pulsing Earth"
	case MagicIDEssMartyrTrink:
		return "Essence of the Martyr Trinket"
	case MagicIDEssSappTrink:
		return "Restrained Essence of Sapphiron Trinket"

	}

	return "<<TODO: Add Aura name to switch!!>>"
}

// List of all magic effects and spells and items and stuff that can go on CD or have an aura.
const (
	MagicIDUnknown int32 = iota
	//Spells
	MagicIDLB12
	MagicIDCL6
	MagicIDTLCLB

	// Auras
	MagicIDClBounce
	MagicIDLOTalent
	MagicIDJoW
	MagicIDEleFocus
	MagicIDEleMastery
	MagicIDStormcaller
	MagicIDBlessingSilverCrescent
	MagicIDQuagsEye
	MagicIDFungalFrenzy
	MagicIDBloodlust
	MagicIDSkycall
	MagicIDEnergized
	MagicIDNAC
	MagicIDChaoticSkyfire
	MagicIDInsightfulEarthstorm
	MagicIDMysticSkyfire
	MagicIDMysticFocus
	MagicIDEmberSkyfire
	MagicIDSpellPower
	MagicIDRubySerpent
	MagicIDCallOfTheNexus
	MagicIDDCC
	MagicIDDCCBonus
	MagicIDDrums // drums effect
	MagicIDNetherstrike
	MagicIDTwinStars
	MagicIDTidefury
	MagicIDSpellstrike
	MagicIDSpellstrikeInfusion
	MagicIDManaEtched
	MagicIDManaEtchedHit
	MagicIDManaEtchedInsight
	MagicIDMisery
	MagicIDEyeOfTheNight
	MagicIDChainTO
	MagicIDCyclone2pc
	MagicIDCyclone4pc
	MagicIDCycloneMana // proc from 4pc
	MagicIDWindhawk
	MagicIDOrcBloodFury    // orc racials
	MagicIDTrollBerserking // troll racial
	MagicIDTLC             // aura on equip of TLC, stores charges
	MagicIDDestructionPotion
	MagicIDHexShunkHead
	MagicIDShiftingNaaru
	MagicIDSkullGuldan
	MagicIDNexusHorn
	MagicIDSextant          // Trinket Aura
	MagicIDUnstableCurrents // Sextant Proc Aura
	MagicIDEyeOfMag         // trinket aura
	MagicIDRecurringPower   // eye of mag proc aura
	MagicIDCataclysm4pc     // cyclone 4pc aura
	MagicIDSkyshatter2pc    // skyshatter 2pc aura
	MagicIDSkyshatter4pc    // skyshatter 4pc aura
	MagicIDElderScribe      // elder scribe robe item aura
	MagicIDElderScribeProc  // elder scribe robe temp buff
	MagicIDTotemOfPulsingEarth

	//Items
	MagicIDISCTrink
	MagicIDNACTrink
	MagicIDPotion
	MagicIDRune
	MagicIDAllTrinket
	MagicIDScryerTrink
	MagicIDRubySerpentTrink
	MagicIDXiriTrink
	MagicIDDrum1 // Party drum item CDs
	MagicIDDrum2
	MagicIDDrum3
	MagicIDDrum4
	MagicIDEyeOfTheNightTrink
	MagicIDChainTOTrink
	MagicIDHexTrink
	MagicIDShiftingNaaruTrink
	MagicIDSkullGuldanTrink
	MagicIDEssMartyrTrink
	MagicIDEssSappTrink

	// Always at end so we know how many magic IDs there are.
	MagicIDLen
)

func ActivateChainLightningBounce(sim *Simulation) Aura {
	return Aura{
		ID:      MagicIDClBounce,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if c.Spell.ID != MagicIDCL6 || c.IsClBounce {
				return
			}

			dmgCoeff := 1.0
			if c.IsLO {
				dmgCoeff = 0.5
			}
			for i := 1; i < sim.Options.NumClTargets; i++ {
				if sim.Options.Tidefury2Pc {
					dmgCoeff *= 0.83
				} else {
					dmgCoeff *= 0.7
				}
				clone := &Cast{
					IsLO:       c.IsLO,
					IsClBounce: true,
					Spell:      c.Spell,
					Crit:       c.Crit,
					CritBonus:  c.CritBonus,
					Effects:    []AuraEffect{func(sim *Simulation, c *Cast) { c.DidDmg *= dmgCoeff }},
				}
				sim.Cast(clone)
			}
		},
	}
}

func AuraLightningOverload(lvl int) Aura {
	chance := 0.04 * float64(lvl)
	return Aura{
		ID:      MagicIDLOTalent,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if c.Spell.ID != MagicIDLB12 && c.Spell.ID != MagicIDCL6 {
				return
			}
			if c.IsLO {
				return // can't proc LO on LO
			}
			actualChance := chance
			if c.Spell.ID == MagicIDCL6 {
				actualChance /= 3 // 33% chance of regular for CL LO
			}
			if sim.rando.Float64() < actualChance {
				if sim.Debug != nil {
					sim.Debug(" +Lightning Overload\n")
				}
				clone := &Cast{
					IsLO: true,
					// Don't set IsClBounce even if this is a bounce, so that the clone does a normal CL and bounces
					Spell:     c.Spell,
					Crit:      c.Crit,
					CritBonus: c.CritBonus,
					Effects:   []AuraEffect{func(sim *Simulation, c *Cast) { c.DidDmg /= 2 }},
				}
				sim.Cast(clone)
			}
		},
	}
}

func AuraJudgementOfWisdom() Aura {
	const mana = 74 / 2 // 50% proc
	return Aura{
		ID:      MagicIDJoW,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if sim.Debug != nil {
				sim.Debug(" +Judgement Of Wisdom: 37 mana (74 @ 50%% proc)\n")
			}
			sim.CurrentMana += mana
		},
	}
}

func AuraElementalFocus(tick int) Aura {
	count := 2
	return Aura{
		ID:      MagicIDEleFocus,
		Expires: tick + (15 * TicksPerSecond),
		OnCast: func(sim *Simulation, c *Cast) {
			c.ManaCost *= .6 // reduced by 40%
		},
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if c.ManaCost <= 0 {
				return // Don't consume charges from free spells.
			}
			count--
			if count == 0 {
				sim.removeAuraByID(MagicIDEleFocus)
			}
		},
	}
}

func AuraEleMastery() Aura {
	return Aura{
		ID:      MagicIDEleMastery,
		Expires: math.MaxInt32,
		OnCast: func(sim *Simulation, c *Cast) {
			if c.Spell.ID != MagicIDTLCLB {
				c.ManaCost = 0
			}
		},
		OnCastComplete: func(sim *Simulation, c *Cast) {
			c.Crit += 1.01 // 101% chance of crit
			// Remove the buff and put skill on CD
			sim.CDs[MagicIDEleMastery] = 180 * TicksPerSecond
			sim.removeAuraByID(MagicIDEleMastery)
		},
	}
}

func AuraStormcaller(tick int) Aura {
	return Aura{
		ID:      MagicIDStormcaller,
		Expires: tick + (8 * TicksPerSecond),
		OnCastComplete: func(sim *Simulation, c *Cast) {
			c.Spellpower += 50
		},
	}
}

func createHasteActivate(id int32, haste float64, durSeconds int) ItemActivation {
	// Implemented haste activate as a buff so that the creation of a new cast gets the correct cast time
	return func(sim *Simulation) Aura {
		sim.Buffs[StatHaste] += haste
		return Aura{
			ID:      id,
			Expires: sim.CurrentTick + durSeconds*TicksPerSecond,
			OnExpire: func(sim *Simulation, c *Cast) {
				sim.Buffs[StatHaste] -= haste
			},
		}
	}
}

// createSpellDmgActivate creates a function for trinket activation that uses +spellpower
//  This is so we don't need a separate function for every spell power trinket.
func createSpellDmgActivate(id int32, sp float64, durSeconds int) ItemActivation {
	return func(sim *Simulation) Aura {
		return Aura{
			ID:      id,
			Expires: sim.CurrentTick + durSeconds*TicksPerSecond,
			OnCastComplete: func(sim *Simulation, c *Cast) {
				c.Spellpower += sp
			},
		}
	}
}

func ActivateQuagsEye(sim *Simulation) Aura {
	lastActivation := math.MinInt32
	const hasteBonus = 320.0
	internalCD := 45 * TicksPerSecond
	return Aura{
		ID:      MagicIDQuagsEye,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation+internalCD < sim.CurrentTick && sim.rando.Float64() < 0.1 {
				sim.Buffs[StatHaste] += hasteBonus
				sim.addAura(AuraStatRemoval(sim.CurrentTick, 6.0, hasteBonus, StatHaste, MagicIDFungalFrenzy))
				lastActivation = sim.CurrentTick
			}
		},
	}
}

func ActivateNexusHorn(sim *Simulation) Aura {
	lastActivation := math.MinInt32
	internalCD := 45 * TicksPerSecond
	const spellBonus = 225.0
	return Aura{
		ID:      MagicIDNexusHorn,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if lastActivation+internalCD < sim.CurrentTick && c.DidCrit && sim.rando.Float64() < 0.2 {
				sim.Buffs[StatSpellDmg] += spellBonus
				sim.addAura(AuraStatRemoval(sim.CurrentTick, 10.0, spellBonus, StatSpellDmg, MagicIDCallOfTheNexus))
				lastActivation = sim.CurrentTick
			}
		},
	}
}

func ActivateDCC(sim *Simulation) Aura {
	const spellBonus = 8.0
	stacks := 0
	return Aura{
		ID:      MagicIDDCC,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if stacks < 10 {
				stacks++
				sim.Buffs[StatSpellDmg] += spellBonus
			}
			// Removal aura will refresh with new total spellpower based on stacks.
			//  This will remove the old stack removal buff.
			sim.addAura(Aura{
				ID:      MagicIDDCCBonus,
				Expires: sim.CurrentTick + (10 * TicksPerSecond),
				OnExpire: func(sim *Simulation, c *Cast) {
					sim.Buffs[StatSpellDmg] -= spellBonus * float64(stacks)
					stacks = 0
				},
			})
		},
	}
}

// AuraStatRemoval creates a general aura for removing any buff stat on expiring.
// This is useful for activations / effects that give temp stats.
func AuraStatRemoval(tick int, seconds int, amount float64, stat Stat, id int32) Aura {
	return Aura{
		ID:      id,
		Expires: tick + (seconds * TicksPerSecond),
		OnExpire: func(sim *Simulation, c *Cast) {
			if sim.Debug != nil {
				sim.Debug(" -%0.0f %s from %s\n", amount, stat.StatName(), AuraName(id))
			}
			sim.Buffs[stat] -= amount
		},
	}
}

func ActivateDrums(sim *Simulation) Aura {
	sim.Buffs[StatHaste] += 80
	sim.CDs[MagicIDDrums] = 30 * TicksPerSecond
	return AuraStatRemoval(sim.CurrentTick, 30, 80, StatHaste, MagicIDDrums)
}

func ActivateBloodlust(sim *Simulation) Aura {
	const dur = 40 * TicksPerSecond
	sim.CDs[MagicIDBloodlust] = dur // assumes that multiple BLs are different shaman.
	return Aura{
		ID:      MagicIDBloodlust,
		Expires: sim.CurrentTick + dur,
		OnCast: func(sim *Simulation, c *Cast) {
			c.CastTime /= 1.3 // 30% faster.
			if c.CastTime < sim.Options.GCD {
				c.CastTime = sim.Options.GCD // can't cast faster than GCD
			}
			c.TicksUntilCast = int(c.CastTime*float64(TicksPerSecond)) + 1
		},
	}
}

func ActivateBerserking(sim *Simulation, hasteBonus float64) Aura {
	const dur = 10 * TicksPerSecond
	const cd = 180 * TicksPerSecond
	sim.CDs[MagicIDTrollBerserking] = cd
	return Aura{
		ID:      MagicIDTrollBerserking,
		Expires: sim.CurrentTick + dur,
		OnCast: func(sim *Simulation, c *Cast) {
			c.CastTime /= hasteBonus
			if c.CastTime < sim.Options.GCD {
				c.CastTime = sim.Options.GCD // can't cast faster than GCD
			}
			c.TicksUntilCast = int(c.CastTime*float64(TicksPerSecond)) + 1
		},
	}
}

func ActivateSkycall(sim *Simulation) Aura {
	const hasteBonus = 101
	const dur = 10 * TicksPerSecond
	active := false
	return Aura{
		ID:      MagicIDSkycall,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if c.Spell.ID == MagicIDLB12 && sim.rando.Float64() < 0.15 {
				if !active {
					sim.Buffs[StatHaste] += hasteBonus
					active = true
				}
				sim.addAura(Aura{
					ID:      MagicIDEnergized,
					Expires: sim.CurrentTick + dur,
					OnExpire: func(sim *Simulation, c *Cast) {
						sim.Buffs[StatHaste] -= hasteBonus
						active = false
					},
				})
			}
		},
	}
}

func ActivateNAC(sim *Simulation) Aura {
	return Aura{
		ID:      MagicIDNAC,
		Expires: sim.CurrentTick + 20*TicksPerSecond,
		OnCast: func(sim *Simulation, c *Cast) {
			c.ManaCost *= 1.2
		},
		OnCastComplete: func(sim *Simulation, c *Cast) {
			c.Spellpower += 250
		},
	}
}

func ActivateCSD(sim *Simulation) Aura {
	return Aura{
		ID:      MagicIDChaoticSkyfire,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			c.CritBonus *= 1.03
		},
	}
}

func ActivateIED(sim *Simulation) Aura {
	lastActivation := math.MinInt32
	const icd = 15 * TicksPerSecond
	return Aura{
		ID:      MagicIDInsightfulEarthstorm,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation+icd < sim.CurrentTick && sim.rando.Float64() < 0.04 {
				lastActivation = sim.CurrentTick
				if sim.Debug != nil {
					sim.Debug(" *Insightful Earthstorm Mana Restore - 300\n")
				}
				sim.CurrentMana += 300
			}
		},
	}
}

func ActivateMSD(sim *Simulation) Aura {
	lastActivation := math.MinInt32
	const hasteBonus = 320.0
	const icd = 35 * TicksPerSecond
	return Aura{
		ID:      MagicIDMysticSkyfire,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation+icd < sim.CurrentTick && sim.rando.Float64() < 0.15 {
				sim.Buffs[StatHaste] += hasteBonus
				sim.addAura(AuraStatRemoval(sim.CurrentTick, 4.0, hasteBonus, StatHaste, MagicIDMysticFocus))
				lastActivation = sim.CurrentTick
			}
		},
	}
}

func ActivateESD(sim *Simulation) Aura {
	sim.Buffs[StatInt] += (sim.Stats[StatInt] + sim.Buffs[StatInt]) * 0.02
	return Aura{
		ID:      MagicIDEmberSkyfire,
		Expires: math.MaxInt32,
	}
}

func ActivateSpellstrike(sim *Simulation) Aura {
	const spellBonus = 92.0
	const duration = 10.0
	return Aura{
		ID:      MagicIDSpellstrike,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if sim.rando.Float64() < 0.05 { // TODO: validate
				sim.addAura(Aura{
					ID:      MagicIDSpellstrikeInfusion,
					Expires: sim.CurrentTick + (duration * TicksPerSecond),
					OnCastComplete: func(sim *Simulation, c *Cast) {
						c.Spellpower += spellBonus
					},
				})
			}
		},
	}
}

func ActivateManaEtched(sim *Simulation) Aura {
	const spellBonus = 110.0
	const duration = 15.0
	return Aura{
		ID:      MagicIDManaEtched,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if sim.rando.Float64() < 0.02 { // TODO: validate
				sim.addAura(Aura{
					ID:      MagicIDManaEtchedInsight,
					Expires: sim.CurrentTick + (duration * TicksPerSecond),
					OnCastComplete: func(sim *Simulation, c *Cast) {
						c.Spellpower += spellBonus
					},
				})
			}
		},
	}
}

func ActivateTLC(sim *Simulation) Aura {
	const spellBonus = 110.0
	const duration = 15.0

	tlcspell := spellmap[MagicIDTLCLB]
	const icd = 2.5 * TicksPerSecond

	charges := 0
	lastActivation := 0
	return Aura{
		ID:      MagicIDTLC,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if lastActivation+icd >= sim.CurrentTick {
				return
			}
			if !c.DidCrit {
				return
			}
			charges++
			if sim.Debug != nil {
				sim.Debug(" Lightning Capacitor Charges: %d\n", charges)
			}
			if charges >= 3 {
				if sim.Debug != nil {
					sim.Debug(" Lightning Capacitor Triggered!\n")
				}
				lastActivation = sim.CurrentTick

				clone := &Cast{
					Spell:     tlcspell,
					CritBonus: 1.5, // TLC does not get elemental fury
					// TLC does not get hit talents bonus, subtract them here. (since we dont conditionally apply them)
					Hit:  (-0.02 * float64(sim.Options.Talents.ElementalPrecision)) + (-0.01 * float64(sim.Options.Talents.NaturesGuidance)),
					Crit: (-0.01 * float64(sim.Options.Talents.TidalMastery)) + (-0.01 * float64(sim.Options.Talents.CallOfThunder)),
				}
				sim.Cast(clone)
				charges = 0
			}
		},
	}
}

func ActivateChainTO(sim *Simulation) Aura {
	if sim.Options.Buffs.TwilightOwl {
		return Aura{ID: MagicIDChainTO, Expires: 0}
	}
	return Aura{
		ID:      MagicIDChainTO,
		Expires: sim.CurrentTick + 30*60*TicksPerSecond,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			c.Crit += 0.02
		},
	}
}

func ActivateCycloneManaReduce(sim *Simulation) Aura {
	return Aura{
		ID:      MagicIDCyclone4pc,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if c.DidCrit && sim.rando.Float64() < 0.11 {
				sim.addAura(Aura{
					ID: MagicIDCycloneMana,
					OnCast: func(sim *Simulation, c *Cast) {
						// TODO: how to make sure this goes in before clearcasting?
						c.ManaCost -= 270
					},
					OnCastComplete: func(sim *Simulation, c *Cast) {
						sim.removeAuraByID(MagicIDCycloneMana)
					},
				})
			}
		},
	}
}

func ActivateCataclysmLBDiscount(sim *Simulation) Aura {
	return Aura{
		ID:      MagicIDCataclysm4pc,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if c.DidCrit && sim.rando.Float64() < 0.25 {
				sim.CurrentMana += 120
			}
		},
	}
}

func ActivateSkyshatterImpLB(sim *Simulation) Aura {
	return Aura{
		ID:      MagicIDSkyshatter4pc,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if c.Spell.ID == MagicIDLB12 {
				c.DidDmg *= 1.05
			}
		},
	}
}

func ActivateDestructionPotion(sim *Simulation) Aura {
	sim.destructionPotion = true
	sim.Buffs[StatSpellDmg] += 120
	sim.Buffs[StatSpellCrit] += 44.16
	sim.CDs[MagicIDPotion] = 120 * TicksPerSecond

	const dur = 15 * TicksPerSecond
	return Aura{
		ID:      MagicIDDestructionPotion,
		Expires: sim.CurrentTick + dur,
		OnExpire: func(sim *Simulation, c *Cast) {
			sim.Buffs[StatSpellDmg] -= 120
			sim.Buffs[StatSpellCrit] -= 44.16
		},
	}
}

func ActivateSextant(sim *Simulation) Aura {
	lastActivation := math.MinInt32
	internalCD := 45 * TicksPerSecond
	const spellBonus = 190.0
	const dur = 15.0
	return Aura{
		ID:      MagicIDSextant,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if lastActivation+internalCD < sim.CurrentTick && c.DidCrit && sim.rando.Float64() < 0.2 {
				sim.Buffs[StatSpellDmg] += spellBonus
				sim.addAura(AuraStatRemoval(sim.CurrentTick, dur, spellBonus, StatSpellDmg, MagicIDUnstableCurrents))
				lastActivation = sim.CurrentTick
			}
		},
	}
}

func ActivateEyeOfMag(sim *Simulation) Aura {
	const spellBonus = 170.0
	const dur = 10 * TicksPerSecond
	active := false
	return Aura{
		ID:      MagicIDEyeOfMag,
		Expires: math.MaxInt32,
		OnSpellMiss: func(sim *Simulation, c *Cast) {
			if !active {
				sim.Buffs[StatSpellDmg] += spellBonus
				active = true
			}
			sim.addAura(Aura{
				ID:      MagicIDRecurringPower,
				Expires: sim.CurrentTick + dur,
				OnExpire: func(sim *Simulation, c *Cast) {
					sim.Buffs[StatSpellDmg] -= spellBonus
					active = false
				},
			})
		},
	}
}

func ActivateElderScribes(sim *Simulation) Aura {
	// Gives a chance when your harmful spells land to increase the damage of your spells and effects by up to 130 for 10 sec. (Proc chance: 20%, 50s cooldown)
	lastActivation := math.MinInt32
	internalCD := 50 * TicksPerSecond
	const spellBonus = 130.0
	const dur = 10.0
	const proc = 0.2
	return Aura{
		ID:      MagicIDElderScribe,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			// This code is starting to look a lot like other ICD buff items. Perhaps we could DRY this out.
			if lastActivation+internalCD < sim.CurrentTick && sim.rando.Float64() < proc {
				sim.Buffs[StatSpellDmg] += spellBonus
				sim.addAura(AuraStatRemoval(sim.CurrentTick, dur, spellBonus, StatSpellDmg, MagicIDElderScribeProc))
				lastActivation = sim.CurrentTick
			}
		},
	}
}

func ActivateTotemOfPulsingEarth(sim *Simulation) Aura {
	return Aura{
		ID:      MagicIDTotemOfPulsingEarth,
		Expires: math.MaxInt32,
		OnCast: func(sim *Simulation, c *Cast) {
			if c.Spell.ID == MagicIDLB12 {
				// TODO: how to make sure this goes in before clearcasting?
				c.ManaCost = math.Max(c.ManaCost-27, 0)
			}
		},
	}
}
