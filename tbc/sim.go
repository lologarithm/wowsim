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
	Buffs         Stats     // temp increases
	Equip         Equipment // Current Gear
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

func NewSim(stats Stats, equip Equipment, spellOrder []string, rseed int64) *Simulation {
	rotIdx := 0
	if spellOrder[0] == "pri" {
		rotIdx = -1
		spellOrder = spellOrder[1:]
	}
	sim := &Simulation{
		RotationIdx:   rotIdx,
		Stats:         stats,
		SpellRotation: spellOrder,
		CDs:           map[string]int{},
		Buffs:         Stats{StatLen: 0},
		Auras:         []Aura{AuraLightningOverload()},
		Equip:         equip,
		rseed:         rseed,
		rando:         rand.New(rand.NewSource(rseed)),
	}
	return sim
}

func (sim *Simulation) reset() {
	sim.rseed++
	sim.rando.Seed(sim.rseed)

	sim.currentTick = 0
	sim.CurrentMana = sim.Stats[StatMana]
	sim.CastingSpell = nil
	sim.Buffs = Stats{StatLen: 0}
	sim.CDs = map[string]int{}
	sim.Auras = []Aura{}
	sim.metrics = SimMetrics{}

	// Activate all talent - TODO: make this configurable input

	// Lightning Overload 5/5
	sim.addAura(AuraLightningOverload())
	// Unrelenting Storm 3/5
	sim.Buffs[StatMP5] += sim.Stats[StatInt] * 0.06

	debug("Effective MP5: %0.1f\n", sim.Stats[StatMP5]+sim.Buffs[StatMP5])
	// Activate all permanent item effects.
	for _, item := range sim.Equip {
		if item.Activate != nil && item.ActivateCD == -1 {
			sim.addAura(item.Activate(sim))
		}
	}
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

func (sim *Simulation) Tick(tickID int) {
	if sim.CurrentMana < 0 {
		panic("you should never have negative mana.")
	}

	secondID := tickID / tickPerSecond
	// MP5 regen
	sim.CurrentMana += ((sim.Stats[StatMP5] + sim.Buffs[StatMP5]) / 5.0) / float64(tickPerSecond)

	if sim.CurrentMana > sim.Stats[StatMana] {
		sim.CurrentMana = sim.Stats[StatMana]
	}

	if sim.CastingSpell == nil && sim.metrics.OOMAt == 0 {
		if sim.CurrentMana < 500 {
			sim.metrics.OOMAt = secondID
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
		// Pop potion before next cast.
		if sim.Stats[StatMana]-sim.CurrentMana >= 1500 && sim.CDs["darkrune"] < 1 {
			// Restores 900 to 1500 mana. (2 Min Cooldown)
			sim.CurrentMana += float64(900 + sim.rando.Intn(1500-900))
			sim.CDs["darkrune"] = 120 * tickPerSecond
			debug("[%d] Used Mana Potion\n", secondID)
		}
		if sim.Stats[StatMana]-sim.CurrentMana >= 3000 && sim.CDs["potion"] < 1 {
			// Restores 1800 to 3000 mana. (2 Min Cooldown)
			sim.CurrentMana += float64(1800 + sim.rando.Intn(3000-1800))
			sim.CDs["potion"] = 120 * tickPerSecond
			debug("[%d] Used Mana Potion\n", secondID)
		}
		// Pop any on-use trinkets

		for _, item := range sim.Equip {
			if item.Activate == nil || item.ActivateCD == -1 { // ignore non-activatable, and always active items.
				continue
			}
			if sim.CDs[item.CoolID] > 0 {
				continue
			}
			sim.addAura(item.Activate(sim))
			sim.CDs[item.CoolID] = item.ActivateCD * tickPerSecond
		}

		// Choose next spell
		sim.ChooseSpell()
		if sim.CastingSpell != nil {
			debug("[%d] Casting %s (%0.1f) ...", secondID, sim.CastingSpell.Spell.ID, float64(sim.CastingSpell.TicksUntilCast)/float64(tickPerSecond))
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
		if sim.Auras[i].Expires <= tickID {
			todel = append(todel, i)
		}
	}
	for i := len(todel) - 1; i >= 0; i-- {
		sim.cleanAura(todel[i])
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
	if sim.Auras[i].OnExpire != nil {
		sim.Auras[i].OnExpire(sim, nil)
	}
	// clean up mem
	sim.Auras[i].OnCast = nil
	sim.Auras[i].OnStruck = nil
	sim.Auras[i].OnSpellHit = nil
	sim.Auras[i].OnExpire = nil

	debug(" -removed: %s- ", sim.Auras[i].ID)
	sim.Auras = sim.Auras[:i+copy(sim.Auras[i:], sim.Auras[i+1:])]
}

func (sim *Simulation) addAura(a Aura) {
	for i := range sim.Auras {
		if sim.Auras[i].ID == a.ID {
			// TODO: some auras can stack X values. Figure out plan
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
			for _, aur := range sim.Auras {
				if aur.OnCast != nil {
					aur.OnCast(sim, cast)
				}
			}
			if sim.CDs[so] == 0 && (sim.CurrentMana >= cast.ManaCost) {
				sim.CastingSpell = cast
				break
			}
		}
	} else {
		so := sim.SpellRotation[sim.RotationIdx]
		sp := spellmap[so]
		cast := NewCast(sim, sp, sim.Stats[StatSpellDmg], sim.Stats[StatSpellHit], sim.Stats[StatSpellCrit])
		// Apply any on cast effects.
		for _, aur := range sim.Auras {
			if aur.OnCast != nil {
				aur.OnCast(sim, cast)
			}
		}
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
	for _, aur := range sim.Auras {
		if aur.OnCastComplete != nil {
			aur.OnCastComplete(sim, cast)
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
		cast.DidDmg = dmg
		// Apply any on spell hit effects.
		for _, aur := range sim.Auras {
			if aur.OnSpellHit != nil {
				aur.OnSpellHit(sim, cast)
			}
		}
		// Apply any effects specific to this cast.
		for _, eff := range cast.Effects {
			eff(sim, cast)
		}
		debug(": %0.0f\n", cast.DidDmg)
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
