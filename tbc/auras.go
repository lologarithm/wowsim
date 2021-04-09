package tbc

import (
	"math"
)

type Aura struct {
	ID      int32
	Expires int // ticks aura will apply

	OnCast         AuraEffect
	OnCastComplete AuraEffect
	OnSpellHit     AuraEffect
	OnStruck       AuraEffect
	OnExpire       AuraEffect
}

func AuraName(a int32) string {
	switch a {
	case MagicIDUnknown:
		return "Unknown"
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
	}

	return "unknown"
}

// AuraEffects will mutate a cast or simulation state.
type AuraEffect func(sim *Simulation, c *Cast)

// List of all magic effects and spells and items and stuff that can go on CD or have an aura.
const (
	MagicIDUnknown int32 = iota
	// Auras
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

	//Spells
	MagicIDLB12
	MagicIDCL6

	//Items
	MagicIDISCTrink
	MagicIDNACTrink
	MagicIDPotion
	MagicIDRune
	MagicIDAllTrinket
)

func AuraJudgementOfWisdom() Aura {
	return Aura{
		ID:      MagicIDJoW,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			sim.debug(" +Judgement Of Wisdom: 74 mana\n")
			sim.CurrentMana += 74
		},
	}
}

func AuraLightningOverload(lvl int) Aura {
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
			if sim.rando.Float64() < 0.04*float64(lvl) {
				sim.debug(" +Lightning Overload\n")
				dmg := c.DidDmg
				if c.DidCrit {
					dmg /= 2
				}
				clone := &Cast{
					IsLO:       true,
					Spell:      c.Spell,
					Hit:        c.Hit,
					Crit:       c.Crit,
					Spellpower: c.Spellpower,
					Effects:    []AuraEffect{func(sim *Simulation, c *Cast) { c.DidDmg /= 2 }},
				}
				sim.Cast(clone)
			}
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
			c.Crit = 1.01 // 101% chance of crit
			c.ManaCost = 0
			sim.CDs[MagicIDEleMastery] = 180 * TicksPerSecond
		},
		OnCastComplete: func(sim *Simulation, c *Cast) {
			sim.removeAuraByID(MagicIDEleMastery)
		},
	}
}

func AuraStormcaller(tick int) Aura {
	return Aura{
		ID:      MagicIDStormcaller,
		Expires: tick + (8 * TicksPerSecond),
		OnCast: func(sim *Simulation, c *Cast) {
			sim.debug(" +Stormcaller\n")
			c.Spellpower += 50
		},
	}
}

func ActivateSilverCrescent(sim *Simulation) Aura {
	return Aura{
		ID:      MagicIDBlessingSilverCrescent,
		Expires: sim.currentTick + 20*TicksPerSecond,
		OnCast: func(sim *Simulation, c *Cast) {
			c.Spellpower += 155
		},
	}
}

func ActivateQuagsEye(sim *Simulation) Aura {
	lastActivation := math.MinInt32
	const hasteBonus = 320.0
	return Aura{
		ID:      MagicIDQuagsEye,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation+(45*TicksPerSecond) < sim.currentTick && sim.rando.Float64() < 0.1 {
				sim.debug(" +Quags Fungal Frenzy\n")
				sim.Buffs[StatHaste] += hasteBonus
				sim.addAura(AuraHasteRemoval(sim.currentTick, 6.0, hasteBonus, MagicIDFungalFrenzy))
				lastActivation = sim.currentTick
			}
		},
	}
}

func AuraHasteRemoval(tick int, seconds int, amount float64, id int32) Aura {
	return Aura{
		ID:      id,
		Expires: tick + (seconds * TicksPerSecond),
		OnExpire: func(sim *Simulation, c *Cast) {
			sim.debug(" -haste %0.0f from %s\n", amount, AuraName(id))
			sim.Buffs[StatHaste] += amount
		},
	}
}

func ActivateBloodlust(sim *Simulation) Aura {
	sim.Buffs[StatHaste] += 472.8
	sim.CDs[MagicIDBloodlust] = 40 * TicksPerSecond // assumes that multiple BLs are different shaman.
	return AuraHasteRemoval(sim.currentTick, 40, 472.8, MagicIDBloodlust)
}

func ActivateSkycall(sim *Simulation) Aura {
	const hasteBonus = 101
	return Aura{
		ID:      MagicIDSkycall,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if sim.rando.Float64() < 0.1 { // TODO: what is actual proc rate?
				sim.debug(" +Skycall Energized- \n")
				sim.Buffs[StatHaste] += hasteBonus
				sim.addAura(AuraHasteRemoval(sim.currentTick, 10, hasteBonus, MagicIDEnergized))
			}
		},
	}
}

func ActivateNAC(sim *Simulation) Aura {
	return Aura{
		ID:      MagicIDNAC,
		Expires: sim.currentTick + 300*TicksPerSecond,
		OnCast: func(sim *Simulation, c *Cast) {
			c.Spellpower += 250
			c.ManaCost *= 1.2
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
	return Aura{
		ID:      MagicIDInsightfulEarthstorm,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation+(15*TicksPerSecond) < sim.currentTick && sim.rando.Float64() < 0.04 {
				lastActivation = sim.currentTick
				sim.debug(" *Insightful Earthstorm Mana Restore - 300\n")
				sim.CurrentMana += 300
			}
		},
	}
}

func ActivateMSD(sim *Simulation) Aura {
	lastActivation := math.MinInt32
	const hasteBonus = 320.0
	return Aura{
		ID:      MagicIDMysticSkyfire,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation+(45*TicksPerSecond) < sim.currentTick && sim.rando.Float64() < 0.1 {
				sim.Buffs[StatHaste] += hasteBonus
				sim.addAura(AuraHasteRemoval(sim.currentTick, 4.0, hasteBonus, MagicIDMysticFocus))
				lastActivation = sim.currentTick
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
