package tbc

import (
	"math"
	"strings"
)

type Aura struct {
	ID      string
	Expires int // ticks aura will apply

	OnCast         AuraEffect
	OnCastComplete AuraEffect
	OnSpellHit     AuraEffect
	OnStruck       AuraEffect
	OnExpire       AuraEffect
}

// AuraEffects will mutate a cast or simulation state.
type AuraEffect func(sim *Simulation, c *Cast)

func AuraJudgementOfWisdom() Aura {
	return Aura{
		ID:      "jow",
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			debug(" -Judgement Of Wisdom +74 mana- ")
			sim.CurrentMana += 74
		},
	}
}

func AuraLightningOverload(lvl int) Aura {
	return Aura{
		ID:      "lotalent",
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if !strings.HasPrefix(c.Spell.ID, "LB") && !strings.HasPrefix(c.Spell.ID, "CL") {
				return
			}
			if sim.rando.Float64() < 0.04*float64(lvl) {
				debug("\tLightning Overload...")
				dmg := c.DidDmg
				if c.DidCrit {
					dmg /= 2
				}
				clone := &Cast{
					Spell:      c.Spell,
					Hit:        c.Hit,
					Crit:       c.Crit,
					Spellpower: c.Spellpower,
					DidDmg:     dmg,
					Effects: []AuraEffect{
						func(sim *Simulation, c *Cast) {
							c.DidDmg /= 2
						},
					},
				}
				sim.Cast(clone)
			}
		},
	}
}

func AuraElementalFocus(tick int) Aura {
	count := 2
	return Aura{
		ID:      "elefocus",
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
				sim.cleanAuraName("elefocus")
			}
		},
	}
}

func AuraEleMastery() Aura {
	return Aura{
		ID:      "elemastery",
		Expires: math.MaxInt32,
		OnCast: func(sim *Simulation, c *Cast) {
			debug(" -ele mastery active- ")
			c.Crit = 1.01 // 101% chance of crit
			c.ManaCost = 0
			sim.CDs["elemastery"] = 180 * TicksPerSecond
		},
		OnCastComplete: func(sim *Simulation, c *Cast) {
			sim.cleanAuraName("elemastery")
		},
	}
}

func AuraStormcaller(tick int) Aura {
	return Aura{
		ID:      "stormcaller",
		Expires: tick + (8 * TicksPerSecond),
		OnCast: func(sim *Simulation, c *Cast) {
			debug(" -stormcaller- ")
			c.Spellpower += 50
		},
	}
}

func ActivateSilverCrescent(sim *Simulation) Aura {
	debug(" -silver crescent active- ")
	return Aura{
		ID:      "silvercrescent",
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
		ID:      "quageye",
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation+(45*TicksPerSecond) < sim.currentTick && sim.rando.Float64() < 0.1 {
				debug(" -quags eye- ")
				sim.Buffs[StatHaste] += hasteBonus
				sim.addAura(AuraHasteRemoval(sim.currentTick, 6.0, hasteBonus, "fungalfrenzy"))
				lastActivation = sim.currentTick
			}
		},
	}
}

func AuraHasteRemoval(tick int, seconds int, amount float64, id string) Aura {
	return Aura{
		ID:      id,
		Expires: tick + (seconds * TicksPerSecond),
		OnExpire: func(sim *Simulation, c *Cast) {
			sim.Buffs[StatHaste] -= amount
		},
	}
}

func ActivateBloodlust(sim *Simulation) Aura {
	debug(" -BL Activated- ")
	sim.Buffs[StatHaste] += 472.8
	sim.CDs["bloodlust"] = 40 * TicksPerSecond // assumes that multiple BLs are different shaman.
	return AuraHasteRemoval(sim.currentTick, 40, 472.8, "bloodlust")
}

func ActivateSkycall(sim *Simulation) Aura {
	const hasteBonus = 101
	return Aura{
		ID:      "skycall",
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if sim.rando.Float64() < 0.1 { // TODO: what is actual proc rate?
				debug(" -skycall energized- ")
				sim.Buffs[StatHaste] += hasteBonus
				sim.addAura(AuraHasteRemoval(sim.currentTick, 10, hasteBonus, "energized"))
			}
		},
	}
}

func ActivateNAC(sim *Simulation) Aura {
	debug(" -NAC active- ")
	return Aura{
		ID:      "nac",
		Expires: sim.currentTick + 300*TicksPerSecond,
		OnCast: func(sim *Simulation, c *Cast) {
			c.Spellpower += 250
			c.ManaCost *= 1.2
		},
	}
}
