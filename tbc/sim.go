package tbc

import (
	"fmt"
	"math"
	"math/rand"
)

var IsDebug = false

func debugFunc(sim *Simulation) func(string, ...interface{}) {
	return func(s string, vals ...interface{}) {
		fmt.Printf("[%0.1f] "+s, append([]interface{}{(float64(sim.currentTick) / float64(TicksPerSecond))}, vals...)...)
	}
}

type Simulation struct {
	CurrentMana float64

	SpellChooser func(*Simulation, bool) int // TODO: make more funtional. Return a cast instead of having function mutate sim itself.

	Stats       Stats
	Buffs       Stats     // temp increases
	Equip       Equipment // Current Gear
	activeEquip Equipment // cache of gear that can activate.

	bloodlustCasts int
	Options        Options
	SpellRotation  []*Spell
	RotationIdx    int

	// ticks until cast is complete
	CastingSpell *Cast

	// timeToRegen := 0
	CDs   map[int32]int
	Auras []Aura // this is array instaed of map to speed up browser perf.

	// Clears and regenerates on each Run call.
	metrics SimMetrics

	rando       *rand.Rand
	rseed       int64
	currentTick int
	endTick     int

	debug func(string, ...interface{})
}

type SimMetrics struct {
	TotalDamage float64
	DamageAtOOM float64
	OOMAt       int
	Casts       []*Cast
	ManaAtEnd   int
	Rotation    []string
}

// New sim contructs a simulator with the given stats / equipment / options.
//   Technically we can calculate stats from equip/options but want the ability to override those stats
//   mostly for stat weight purposes.
func NewSim(stats Stats, equip Equipment, options Options) *Simulation {
	if len(options.SpellOrder) == 0 && !options.UseAI {
		fmt.Printf("[ERROR] No rotation given to sim.\n")
		return nil
	}
	rotIdx := 0
	var rot []*Spell
	if !options.UseAI {
		if options.SpellOrder[0] == "pri" {
			rotIdx = -1
			options.SpellOrder = options.SpellOrder[1:]
		}
		rot = make([]*Spell, len(options.SpellOrder))
		for i, v := range options.SpellOrder {
			for _, sp := range spells {
				if sp.Name == v {
					rot[i] = &sp
					break
				}
			}
		}
	}
	sim := &Simulation{
		RotationIdx:   rotIdx,
		Stats:         stats,
		SpellRotation: rot,
		Options:       options,
		CDs:           map[int32]int{},
		Buffs:         Stats{StatLen: 0},
		Auras:         []Aura{},
		Equip:         equip,
		rseed:         options.RSeed,
		rando:         rand.New(rand.NewSource(options.RSeed)),
		debug:         func(a string, v ...interface{}) {},
		SpellChooser:  ChooseSpell,
	}

	if options.UseAI {
		ai := NewAI(sim)
		sim.SpellChooser = ai.ChooseSpell
	}

	if IsDebug {
		sim.debug = debugFunc(sim)
	}

	for _, eq := range equip {
		if eq.Activate != nil {
			sim.activeEquip = append(sim.activeEquip, eq)
		}
		for _, g := range eq.Gems {
			if g.Activate != nil {
				sim.activeEquip = append(sim.activeEquip, eq)
			}
		}
	}
	return sim
}

func (sim *Simulation) reset() {
	sim.rseed++
	sim.rando.Seed(sim.rseed)

	if sim.Options.UseAI {
		// Reset a new AI
		// TODO: Can we take learnings from the last AI to modulate this AIs behavior?
		ai := NewAI(sim)
		sim.SpellChooser = ai.ChooseSpell
	}

	sim.bloodlustCasts = 0
	sim.currentTick = 0
	sim.CurrentMana = sim.Stats[StatMana]
	sim.CastingSpell = nil
	sim.Buffs = Stats{StatLen: 0}
	sim.CDs = map[int32]int{}
	sim.Auras = []Aura{}
	sim.metrics = SimMetrics{}

	sim.debug("\nSIM RESET\nRotation: %v\n", sim.SpellRotation)
	sim.debug("Effective MP5: %0.1f\n", sim.Stats[StatMP5]+sim.Buffs[StatMP5])
	sim.debug("----------------------\n")

	// Activate all talents
	if sim.Options.Talents.LightninOverload > 0 {
		sim.addAura(AuraLightningOverload(sim.Options.Talents.LightninOverload))
	}

	// Judgement of Wisdom
	if sim.Options.Buffs.JudgementOfWisdom {
		sim.addAura(AuraJudgementOfWisdom())
	}

	// Activate all permanent item effects.
	for _, item := range sim.activeEquip {
		if item.Activate != nil && item.ActivateCD == -1 {
			sim.addAura(item.Activate(sim))
		}
		for _, g := range item.Gems {
			if g.Activate != nil {
				sim.addAura(g.Activate(sim))
			}
		}
	}

	// Activate Set Bonuses
	for _, set := range sets {
		itemCount := 0
		for _, i := range sim.Equip {
			if set.Items[i.Name] {
				itemCount++
				if bonus, ok := set.Bonuses[itemCount]; ok {
					sim.addAura(bonus(sim))
				}
			}
		}
	}
}

func (sim *Simulation) Run(seconds int) SimMetrics {
	sim.endTick = seconds * TicksPerSecond
	// For now use the new 'event' driven state advancement.
	return sim.Run2(seconds)
}

func (sim *Simulation) removeAuraByID(id int32) {
	for i := range sim.Auras {
		if sim.Auras[i].ID == id {
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
	sim.Auras[i].OnCastComplete = nil
	sim.Auras[i].OnStruck = nil
	sim.Auras[i].OnSpellHit = nil
	sim.Auras[i].OnExpire = nil

	sim.debug(" -%s\n", AuraName(sim.Auras[i].ID))
	sim.Auras = sim.Auras[:i+copy(sim.Auras[i:], sim.Auras[i+1:])]
}

func (sim *Simulation) addAura(a Aura) {
	sim.debug(" +%s\n", AuraName(a.ID))
	for i := range sim.Auras {
		if sim.Auras[i].ID == a.ID {
			sim.Auras[i] = a // replace
			return
		}
	}
	sim.Auras = append(sim.Auras, a)
}

func ChooseSpell(sim *Simulation, didPot bool) int {
	if sim.RotationIdx == -1 {
		lowestWait := math.MaxInt32
		wasMana := false
		for i := 0; i < len(sim.SpellRotation); i++ {
			sp := sim.SpellRotation[i]
			so := sp.ID
			cast := NewCast(sim, sp)
			if sim.CDs[so] > 0 { // if
				if sim.CDs[so] < lowestWait {
					lowestWait = sim.CDs[so]
				}
				continue
			}
			if sim.CurrentMana >= cast.ManaCost {
				sim.CastingSpell = cast
				return cast.TicksUntilCast
			}
			manaRegenTicks := int(math.Ceil((cast.ManaCost - sim.CurrentMana) / sim.manaRegen()))
			if manaRegenTicks < lowestWait {
				lowestWait = manaRegenTicks
				wasMana = true
			}
		}
		if wasMana && sim.metrics.OOMAt == 0 { // loop only completes if no spell was found.
			sim.metrics.OOMAt = sim.currentTick / TicksPerSecond
			sim.metrics.DamageAtOOM = sim.metrics.TotalDamage
		}
		return lowestWait
	}

	sp := sim.SpellRotation[sim.RotationIdx]
	so := sp.ID
	cast := NewCast(sim, sp)
	if sim.CDs[so] < 1 {
		if sim.CurrentMana >= cast.ManaCost {
			sim.CastingSpell = cast
			sim.RotationIdx++
			if sim.RotationIdx == len(sim.SpellRotation) {
				sim.RotationIdx = 0
			}
			return cast.TicksUntilCast
		} else {
			sim.debug("Current Mana %0.0f, Cast Cost: %0.0f\n", sim.CurrentMana, cast.ManaCost)
			if sim.metrics.OOMAt == 0 {
				sim.metrics.OOMAt = sim.currentTick / TicksPerSecond
				sim.metrics.DamageAtOOM = sim.metrics.TotalDamage
			}
			return int(math.Ceil((cast.ManaCost - sim.CurrentMana) / sim.manaRegen()))
		}
	}
	return sim.CDs[so]
}

func (sim *Simulation) Cast(cast *Cast) {
	for _, aur := range sim.Auras {
		if aur.OnCastComplete != nil {
			aur.OnCastComplete(sim, cast)
		}
	}
	hit := 0.83 + ((sim.Stats[StatSpellHit] + sim.Buffs[StatSpellHit]) / 1260.0) + cast.Hit // 12.6 hit == 1% hit
	if hit > 1.0 {
		hit = 0.99 // can't get away from the 1% miss
	}

	sim.debug("Completed Cast (%s)\n", cast.Spell.Name)
	dbgCast := cast.Spell.Name
	if sim.rando.Float64() < hit {
		sp := sim.Stats[StatSpellDmg] + sim.Buffs[StatSpellDmg] + cast.Spellpower
		dmg := (sim.rando.Float64() * (cast.Spell.MaxDmg - cast.Spell.MinDmg)) + cast.Spell.MinDmg + (sp * cast.Spell.Coeff)
		if cast.DidDmg != 0 { // use the pre-set dmg
			dmg = cast.DidDmg
		}
		cast.DidHit = true

		crit := ((sim.Stats[StatSpellCrit] + sim.Buffs[StatSpellCrit]) / 2208.0) + cast.Crit // 22.08 crit == 1% crit
		if sim.rando.Float64() < crit {
			cast.DidCrit = true
			// TODO: Make Elemental Fury talent check here.
			dmg *= (cast.CritBonus * 2) - 1 // if CSD equipped the cast crit bonus will be modified during 'onCastComplete.'
			sim.addAura(AuraElementalFocus(sim.currentTick))
			if IsDebug {
				dbgCast += " crit"
			}
		} else if IsDebug {
			dbgCast += " hit"
		}

		if sim.Options.Talents.Concussion > 0 && cast.Spell.ID == MagicIDLB12 || cast.Spell.ID == MagicIDCL6 {
			// Talent Concussion
			dmg *= 1 + (0.01 * sim.Options.Talents.Concussion)
		}

		// Average Resistance (AR) = (Target's Resistance / (Caster's Level * 5)) * 0.75
		// P(x) = 50% - 250%*|x - AR| <- where X is %resisted
		// For now hardcode the 25% chance resist at 2.5% (this assumes bosses have 0 nature resist)
		// .... Using level 70 shows a 0.35% chance of resist instead of 2.5%... not sure what is correct.
		if sim.rando.Float64() < 0.025 { // chance of 25% resist
			dmg *= .75
			if IsDebug {
				dbgCast += " (partial resist)"
			}
		}
		cast.DidDmg = dmg
		// Apply any effects specific to this cast.
		for _, eff := range cast.Effects {
			eff(sim, cast)
		}
		// Apply any on spell hit effects.
		for _, aur := range sim.Auras {
			if aur.OnSpellHit != nil {
				aur.OnSpellHit(sim, cast)
			}
		}
		sim.metrics.TotalDamage += cast.DidDmg
	} else {
		if IsDebug {
			dbgCast += " miss"
		}
	}
	sim.metrics.Casts = append(sim.metrics.Casts, cast)
	sim.debug("%s: %0.0f\n", dbgCast, cast.DidDmg)
	sim.CurrentMana -= cast.ManaCost
	sim.CastingSpell = nil
	if cast.Spell.Cooldown > 0 {
		sim.CDs[cast.Spell.ID] = cast.Spell.Cooldown * TicksPerSecond
	}
}
