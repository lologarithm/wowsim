package tbc

import (
	"math"
)

// AuraEffects will mutate a cast or simulation state.
type AuraEffect func(sim *Simulation, c *Cast)

type Aura struct {
	ID      int32
	Expires float64 // time at which aura will be removed

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
		OnCastComplete: func(sim *Simulation, c *Cast) {
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

func AuraElementalFocus(currentTime float64) Aura {
	count := 2
	return Aura{
		ID:      MagicIDEleFocus,
		Expires: currentTime + 15,
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

func TryActivateEleMastery(sim *Simulation) {
	if sim.isOnCD(MagicIDEleMastery) || !sim.Options.Talents.ElementalMastery {
		return
	}

	sim.addAura(Aura{
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
			sim.setCD(MagicIDEleMastery, 180)
			sim.removeAuraByID(MagicIDEleMastery)
		},
	})
}

func AuraStormcaller(currentTime float64) Aura {
	return Aura{
		ID:      MagicIDStormcaller,
		Expires: currentTime + 8,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			c.Spellpower += 50
		},
	}
}

func createHasteActivate(id int32, haste float64, durSeconds float64) ItemActivation {
	// Implemented haste activate as a buff so that the creation of a new cast gets the correct cast time
	return func(sim *Simulation) Aura {
		sim.Buffs[StatHaste] += haste
		return Aura{
			ID:      id,
			Expires: sim.CurrentTime + durSeconds,
			OnExpire: func(sim *Simulation, c *Cast) {
				sim.Buffs[StatHaste] -= haste
			},
		}
	}
}

// createSpellDmgActivate creates a function for trinket activation that uses +spellpower
//  This is so we don't need a separate function for every spell power trinket.
func createSpellDmgActivate(id int32, sp float64, durSeconds float64) ItemActivation {
	return func(sim *Simulation) Aura {
		return Aura{
			ID:      id,
			Expires: sim.CurrentTime + durSeconds,
			OnCastComplete: func(sim *Simulation, c *Cast) {
				c.Spellpower += sp
			},
		}
	}
}

func ActivateQuagsEye(sim *Simulation) Aura {
	lastActivation := -math.MaxFloat64
	const hasteBonus = 320.0
	internalCD := 45.0
	return Aura{
		ID:      MagicIDQuagsEye,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation + internalCD < sim.CurrentTime && sim.rando.Float64() < 0.1 {
				sim.Buffs[StatHaste] += hasteBonus
				sim.addAura(AuraStatRemoval(sim.CurrentTime, 6.0, hasteBonus, StatHaste, MagicIDFungalFrenzy))
				lastActivation = sim.CurrentTime
			}
		},
	}
}

func ActivateNexusHorn(sim *Simulation) Aura {
	lastActivation := -math.MaxFloat64
	internalCD := 45.0
	const spellBonus = 225.0
	return Aura{
		ID:      MagicIDNexusHorn,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if lastActivation + internalCD < sim.CurrentTime && c.DidCrit && sim.rando.Float64() < 0.2 {
				sim.Buffs[StatSpellDmg] += spellBonus
				sim.addAura(AuraStatRemoval(sim.CurrentTime, 10.0, spellBonus, StatSpellDmg, MagicIDCallOfTheNexus))
				lastActivation = sim.CurrentTime
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
				Expires: sim.CurrentTime + 10,
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
func AuraStatRemoval(currentTime float64, duration float64, amount float64, stat Stat, id int32) Aura {
	return Aura{
		ID:      id,
		Expires: currentTime + duration,
		OnExpire: func(sim *Simulation, c *Cast) {
			if sim.Debug != nil {
				sim.Debug(" -%0.0f %s from %s\n", amount, stat.StatName(), AuraName(id))
			}
			sim.Buffs[stat] -= amount
		},
	}
}

func TryActivateDrums(sim *Simulation) {
	if sim.Options.NumDrums == 0  || sim.isOnCD(MagicIDDrums) {
		return
	}

	sim.Buffs[StatHaste] += 80
	sim.setCD(MagicIDDrums, 30)
	sim.addAura(AuraStatRemoval(sim.CurrentTime, 30, 80, StatHaste, MagicIDDrums))
}

func TryActivateBloodlust(sim *Simulation) {
	if sim.Options.NumBloodlust <= sim.bloodlustCasts || sim.isOnCD(MagicIDBloodlust) {
		return
	}

	dur := 40.0 // assumes that multiple BLs are different shaman.
	sim.setCD(MagicIDBloodlust, dur)
	sim.bloodlustCasts++ // TODO: will this break anything?
	sim.addAura(Aura{
		ID:      MagicIDBloodlust,
		Expires: sim.CurrentTime + dur,
		OnCast: func(sim *Simulation, c *Cast) {
			c.CastTime /= 1.3 // 30% faster.
			if c.CastTime < sim.Options.GCDMin {
				c.CastTime = sim.Options.GCDMin // can't cast faster than GCD
			}
		},
	})
}

func TryActivateRacial(sim *Simulation) {
	switch v := sim.Options.Buffs.Race; v {
	case RaceBonusOrc:
		if sim.isOnCD(MagicIDOrcBloodFury) {
			return
		}

		const spBonus = 143
		const dur = 15
		const cd = 120

		sim.Buffs[StatSpellDmg] += spBonus
		sim.setCD(MagicIDOrcBloodFury, cd)
		sim.addAura(AuraStatRemoval(sim.CurrentTime, dur, spBonus, StatSpellDmg, MagicIDOrcBloodFury))

	case RaceBonusTroll10, RaceBonusTroll30:
		if sim.isOnCD(MagicIDTrollBerserking) {
			return
		}

		hasteBonus := 1.1 // 10% haste
		if v == RaceBonusTroll30 {
			hasteBonus = 1.3 // 30% haste
		}
		const dur = 10
		const cd = 180

		sim.setCD(MagicIDTrollBerserking, cd)
		sim.addAura(Aura{
			ID:      MagicIDTrollBerserking,
			Expires: sim.CurrentTime + dur,
			OnCast: func(sim *Simulation, c *Cast) {
				c.CastTime /= hasteBonus
				if c.CastTime < sim.Options.GCDMin {
					c.CastTime = sim.Options.GCDMin // can't cast faster than GCD
				}
			},
		})
	}
}

func ActivateSkycall(sim *Simulation) Aura {
	const hasteBonus = 101
	const dur = 10
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
					Expires: sim.CurrentTime + dur,
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
		Expires: sim.CurrentTime + 20,
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
	lastActivation := -math.MaxFloat64
	const icd = 15.0
	return Aura{
		ID:      MagicIDInsightfulEarthstorm,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation + icd < sim.CurrentTime && sim.rando.Float64() < 0.04 {
				lastActivation = sim.CurrentTime
				if sim.Debug != nil {
					sim.Debug(" *Insightful Earthstorm Mana Restore - 300\n")
				}
				sim.CurrentMana += 300
			}
		},
	}
}

func ActivateMSD(sim *Simulation) Aura {
	lastActivation := -math.MaxFloat64
	const hasteBonus = 320.0
	const icd = 35.0
	return Aura{
		ID:      MagicIDMysticSkyfire,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation + icd < sim.CurrentTime && sim.rando.Float64() < 0.15 {
				sim.Buffs[StatHaste] += hasteBonus
				sim.addAura(AuraStatRemoval(sim.CurrentTime, 4.0, hasteBonus, StatHaste, MagicIDMysticFocus))
				lastActivation = sim.CurrentTime
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
					Expires: sim.CurrentTime + duration,
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
					Expires: sim.CurrentTime + duration,
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
	const icd = 2.5

	charges := 0
	lastActivation := 0.0
	return Aura{
		ID:      MagicIDTLC,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if lastActivation + icd >= sim.CurrentTime {
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
				lastActivation = sim.CurrentTime

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
		Expires: sim.CurrentTime + 30 * 60,
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

func TryActivateDestructionPotion(sim *Simulation) {
	if !sim.Options.Consumes.DestructionPotion || sim.isOnCD(MagicIDPotion) {
		return
	}

	// Only use dest potion if not using mana or if we haven't used it once.
	// If we are using mana, only use destruction potion on the pull.
	if sim.destructionPotion && sim.Options.Consumes.SuperManaPotion {
		return
	}

	const spBonus = 120
	const critBonus = 44.16
	const dur = 15

	sim.destructionPotion = true
	sim.setCD(MagicIDPotion, 120)
	sim.Buffs[StatSpellDmg] += spBonus
	sim.Buffs[StatSpellCrit] += critBonus

	sim.addAura(Aura{
		ID:      MagicIDDestructionPotion,
		Expires: sim.CurrentTime + dur,
		OnExpire: func(sim *Simulation, c *Cast) {
			sim.Buffs[StatSpellDmg] -= spBonus
			sim.Buffs[StatSpellCrit] -= critBonus
		},
	})
}

// TODO: This function doesn't really belong in auras.go, find a better home for it.
func TryActivateDarkRune(sim *Simulation) bool {
	if !sim.Options.Consumes.DarkRune || sim.isOnCD(MagicIDRune) {
		return false
	}

	// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
	totalRegen := sim.manaRegen() * 5
	if sim.Stats[StatMana] - (sim.CurrentMana + totalRegen) < 1500 {
		return false
	}

	// Restores 900 to 1500 mana. (2 Min Cooldown)
	sim.CurrentMana += 900 + (sim.rando.Float64() * 600)
	sim.setCD(MagicIDRune, 120)
	if sim.Debug != nil {
		sim.Debug("Used Dark Rune\n")
	}
	return true
}

// TODO: This function doesn't really belong in auras.go, find a better home for it.
func TryActivateSuperManaPotion(sim *Simulation) bool {
	if !sim.Options.Consumes.SuperManaPotion || sim.isOnCD(MagicIDPotion) {
		return false
	}

	// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
	totalRegen := sim.manaRegen() * 5
	if sim.Stats[StatMana] - (sim.CurrentMana + totalRegen) < 3000 {
		return false
	}

	// Restores 1800 to 3000 mana. (2 Min Cooldown)
	sim.CurrentMana += 1800 + (sim.rando.Float64() * 1200)
	sim.setCD(MagicIDPotion, 120)
	if sim.Debug != nil {
		sim.Debug("Used Mana Potion\n")
	}
	return true
}

func ActivateSextant(sim *Simulation) Aura {
	lastActivation := -math.MaxFloat64
	internalCD := 45.0
	const spellBonus = 190.0
	const dur = 15.0
	return Aura{
		ID:      MagicIDSextant,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if lastActivation + internalCD < sim.CurrentTime && c.DidCrit && sim.rando.Float64() < 0.2 {
				sim.Buffs[StatSpellDmg] += spellBonus
				sim.addAura(AuraStatRemoval(sim.CurrentTime, dur, spellBonus, StatSpellDmg, MagicIDUnstableCurrents))
				lastActivation = sim.CurrentTime
			}
		},
	}
}

func ActivateEyeOfMag(sim *Simulation) Aura {
	const spellBonus = 170.0
	const dur = 10
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
				Expires: sim.CurrentTime + dur,
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
	lastActivation := -math.MaxFloat64
	internalCD := 50.0
	const spellBonus = 130.0
	const dur = 10.0
	const proc = 0.2
	return Aura{
		ID:      MagicIDElderScribe,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			// This code is starting to look a lot like other ICD buff items. Perhaps we could DRY this out.
			if lastActivation + internalCD < sim.CurrentTime && sim.rando.Float64() < proc {
				sim.Buffs[StatSpellDmg] += spellBonus
				sim.addAura(AuraStatRemoval(sim.CurrentTime, dur, spellBonus, StatSpellDmg, MagicIDElderScribeProc))
				lastActivation = sim.CurrentTime
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
