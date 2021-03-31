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

func AuraLightningOverload() Aura {
	return Aura{
		ID:      "lotalent",
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if !strings.HasPrefix(c.Spell.ID, "LB") && !strings.HasPrefix(c.Spell.ID, "CL") {
				return
			}
			if sim.rando.Float64() < 0.2 {
				debug("\tLightning Overload...")
				clone := &Cast{
					Spell:      c.Spell,
					Hit:        c.Hit,
					Crit:       c.Crit,
					Spellpower: c.Spellpower,
					Effects: []AuraEffect{
						func(sim *Simulation, c *Cast) { c.DidDmg /= 2 },
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
		Expires: tick + (15 * tickPerSecond),
		OnCast: func(sim *Simulation, c *Cast) {
			if c.ManaCost <= 0 {
				return // Don't consume charges from free spells.
			}
			c.ManaCost *= .6 // reduced by 40%
		},
		OnCastComplete: func(sim *Simulation, c *Cast) {
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
			debug("ele mastery...")
			c.Crit = 1.01 // 101% chance of crit
			c.ManaCost = 0
			sim.cleanAuraName("elemastery")
		},
		OnCastComplete: func(sim *Simulation, c *Cast) {
			sim.cleanAuraName("elemastery")
		},
	}
}

func AuraStormcaller(tick int) Aura {
	return Aura{
		ID:      "stormcaller",
		Expires: tick + (8 * tickPerSecond),
		OnCast: func(sim *Simulation, c *Cast) {
			debug("stormcaller...")
			c.Spellpower += 50
		},
	}
}

func AuraQuagsEye() Aura {
	lastActivation := math.MinInt32
	return Aura{
		ID:      "quageye",
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation+(45*tickPerSecond) < sim.currentTick && sim.rando.Float64() < 0.1 {
				debug("\n\tQuags Eye Procced...\n")
				sim.addAura(AuraFungalFrenzy(sim.currentTick))
				sim.Buffs[StatHaste] += 601
				lastActivation = sim.currentTick
			}
		},
	}
}

func AuraFungalFrenzy(tick int) Aura {
	return Aura{
		ID:      "fungalfrenzy",
		Expires: tick + (6 * tickPerSecond),
		OnExpire: func(sim *Simulation, c *Cast) {
			sim.Buffs[StatHaste] -= 601
		},
	}
}
