package tbc

import (
	"fmt"
	"math/rand"
	"strings"
)

var IsDebug = false

func debug(s string, vals ...interface{}) {
	if IsDebug {
		fmt.Printf(s, vals...)
	}
}

type Simulation struct {
	CurrentMana float64

	Stats         Stats
	SpellRotation []string
	RotationIdx   int

	// ticks until cast is complete
	CastingSpell *Cast

	// timeToRegen := 0
	CDs   map[string]int
	Auras []Aura // this is array instaed of map to speed up browser perf.

	// Clears and regenerates on each Run call.
	metrics SimMetrics

	rando       *rand.Rand
	rseed       int64
	currentTick int
}

type SimMetrics struct {
	TotalDamage float64
	OOMAt       int
	Casts       []*Cast
}

func NewSim(stats Stats, spellOrder []string, rseed int64) *Simulation {
	rotIdx := 0
	if spellOrder[0] == "pri" {
		rotIdx = -1
		spellOrder = spellOrder[1:]
	}
	return &Simulation{
		RotationIdx:   rotIdx,
		Stats:         stats,
		SpellRotation: spellOrder,
		CDs:           map[string]int{},
		Auras:         []Aura{AuraLightningOverload()},
		rseed:         rseed,
		rando:         rand.New(rand.NewSource(rseed)),
	}
}

func (sim *Simulation) reset() {
	sim.rseed++
	sim.rando.Seed(sim.rseed)

	sim.CurrentMana = sim.Stats[StatMana]
	sim.CastingSpell = nil
	sim.CDs = map[string]int{}
	sim.Auras = []Aura{AuraLightningOverload()}
	sim.metrics = SimMetrics{}
}

func (sim *Simulation) Run(seconds int) SimMetrics {
	sim.reset()
	ticks := seconds * tickPerSecond
	for i := 0; i < ticks; i++ {
		sim.currentTick = i
		sim.Tick(i)
	}
	debug("(%0.0f/%0.0f mana)\n", sim.CurrentMana, sim.Stats[StatMana])
	return sim.metrics
}

func (sim *Simulation) Tick(i int) {
	if sim.CurrentMana < 0 {
		panic("you should never have negative mana.")
	}
	// MP5 regen
	sim.CurrentMana += (sim.Stats[StatMP5] / 5.0) / float64(tickPerSecond)

	// Unrelenting Storm 3/5 lol TODO: make this configurable input
	sim.CurrentMana += ((sim.Stats[StatInt] * 0.06) / 5.0) / float64(tickPerSecond)

	if sim.CurrentMana > sim.Stats[StatMana] {
		sim.CurrentMana = sim.Stats[StatMana]
	}

	if sim.Stats[StatMana]-sim.CurrentMana >= 1500 && sim.CDs["darkrune"] < 1 {
		// Restores 900 to 1500 mana. (2 Min Cooldown)
		sim.CurrentMana += float64(900 + sim.rando.Intn(1500-900))
		sim.CDs["darkrune"] = 120 * tickPerSecond
		debug("[%d] Used Mana Potion\n", i/tickPerSecond)
	}
	if sim.Stats[StatMana]-sim.CurrentMana >= 2250 && sim.CDs["potion"] < 1 {
		// Restores 1350 to 2250 mana. (2 Min Cooldown)
		sim.CurrentMana += float64(1350 + sim.rando.Intn(2250-1350))
		sim.CDs["potion"] = 120 * tickPerSecond
		debug("[%d] Used Mana Potion\n", i/tickPerSecond)
	}

	if sim.CastingSpell == nil && sim.metrics.OOMAt == 0 {
		if sim.CurrentMana < 500 {
			sim.metrics.OOMAt = i / tickPerSecond
		}
		// debug("(%0.0f/%0.0f mana)\n", sim.CurrentMana, sim.Stats[StatMana])
	}

	if sim.CastingSpell != nil {
		sim.CastingSpell.TicksUntilCast-- // advance state of current cast.
		if sim.CastingSpell.TicksUntilCast == 0 {
			sim.Cast(sim.CastingSpell)
		}
	}
	if sim.CastingSpell == nil {
		// Choose next spell
		sim.ChooseSpell()
		if sim.CastingSpell != nil {
			debug("[%d] Casting %s ...", i/tickPerSecond, sim.CastingSpell.Spell.ID)
		}
	}
	// CDS
	for k := range sim.CDs {
		sim.CDs[k]--
		if sim.CDs[k] <= 0 {
			delete(sim.CDs, k)
		}
	}

	todel := []int{}
	for i := range sim.Auras {
		if sim.Auras[i].Expires <= i {
			todel = append(todel, i)
		}
	}
	for i := len(todel) - 1; i >= 0; i-- {
		sim.cleanAura(i)
	}
}

func (sim *Simulation) cleanAuraName(name string) {
	for i := range sim.Auras {
		if sim.Auras[i].ID == name {
			sim.cleanAura(i)
			break
		}
	}
}
func (sim *Simulation) cleanAura(i int) {
	// clean up mem
	sim.Auras[i].OnCast = nil
	sim.Auras[i].OnStruck = nil
	sim.Auras[i].OnSpellHit = nil
	debug("removing aura: %s", sim.Auras[i].ID)
	sim.Auras = sim.Auras[:i+copy(sim.Auras[i:], sim.Auras[i+1:])]
}

func (sim *Simulation) addAura(a Aura) {
	for i := range sim.Auras {
		if sim.Auras[i].ID == a.ID {
			sim.Auras[i] = a // replace
			return
		}
	}
	sim.Auras = append(sim.Auras, a)
}

func (sim *Simulation) ChooseSpell() {
	// Ele Mastery up, pop!
	if sim.CDs["elemastery"] <= 0 {
		// Apply auras
		sim.addAura(AuraEleMastery())
		sim.CDs["elemastery"] = 180 * tickPerSecond
	}

	if sim.RotationIdx == -1 {
		for i := 0; i < len(sim.SpellRotation); i++ {
			so := sim.SpellRotation[i]
			sp := spellmap[so]
			cast := NewCast(sim, sp, sim.Stats[StatSpellDmg], sim.Stats[StatSpellHit], sim.Stats[StatSpellCrit])
			if sim.CDs[so] == 0 && (sim.CurrentMana >= cast.ManaCost) {
				sim.CastingSpell = cast
				break
			}
		}
	} else {
		so := sim.SpellRotation[sim.RotationIdx]
		sp := spellmap[so]
		cast := NewCast(sim, sp, sim.Stats[StatSpellDmg], sim.Stats[StatSpellHit], sim.Stats[StatSpellCrit])
		if sim.CDs[so] == 0 && (sim.CurrentMana >= cast.ManaCost) {
			sim.CastingSpell = cast
			sim.RotationIdx++
			if sim.RotationIdx == len(sim.SpellRotation) {
				sim.RotationIdx = 0
			}
		}
	}
}

func (sim *Simulation) Cast(cast *Cast) {

	// Apply any on cast effects.
	for _, aur := range sim.Auras {
		if aur.OnCast != nil {
			aur.OnCast(sim, cast)
		}
	}

	if sim.rando.Float64() < cast.Hit {
		dmg := (float64(sim.rando.Intn(int(cast.Spell.MaxDmg-cast.Spell.MinDmg))) + cast.Spell.MinDmg) + (sim.Stats[StatSpellDmg] * cast.Spell.Coeff)

		cast.DidHit = true
		if sim.rando.Float64() < cast.Crit {
			cast.DidCrit = true
			dmg *= 2
			debug("crit")

			sim.addAura(AuraElementalFocus(sim.currentTick))
		} else {
			debug("hit")
		}

		if strings.HasPrefix(cast.Spell.ID, "LB") || strings.HasPrefix(cast.Spell.ID, "CL") {
			// Talent Concussion
			dmg *= 1.05
		}

		// Average Resistance = (Target's Resistance / (Caster's Level * 5)) * 0.75 "AR"
		// P(x) = 50% - 250%*|x - AR| <- where X is chance of resist
		// For now hardcode the 25% chance resist at 2.5% (this assumes bosses have 0 nature resist)
		if sim.rando.Float64() < 0.025 { // chance of 25% resist
			dmg *= .75
			debug("(partial resist)")
		}
		debug(": %0.0f\n", dmg)

		cast.DidDmg = dmg
		// Apply any on spell hit effects.
		for _, aur := range sim.Auras {
			if aur.OnSpellHit != nil {
				aur.OnSpellHit(sim, cast)
			}
		}

		sim.metrics.TotalDamage += cast.DidDmg
		sim.metrics.Casts = append(sim.metrics.Casts, cast)
	} else {
		debug("miss.\n")
	}

	sim.CurrentMana -= cast.ManaCost
	sim.CastingSpell = nil
	if cast.Spell.Cooldown > 0 {
		sim.CDs[cast.Spell.ID] = cast.Spell.Cooldown * tickPerSecond
	}
}
