package tbc

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

func debugFunc(sim *Simulation) func(string, ...interface{}) {
	return func(s string, vals ...interface{}) {
		fmt.Printf("[%0.1f] "+s, append([]interface{}{sim.CurrentTime}, vals...)...)
	}
}

type Simulation struct {
	CurrentMana float64

	Agent Agent

	Stats       Stats
	Buffs       Stats     // temp increases
	Equip       Equipment // Current Gear
	activeEquip Equipment // cache of gear that can activate.

	bloodlustCasts    int
	destructionPotion bool
	Options           Options

	// timeToRegen := 0
	_CDs   []float64  // Map of MagicID to ticks until CD is done. 'Advance' counts down these
	Auras []Aura // this is array instaed of map to speed up browser perf.

	// Clears and regenerates on each Run call.
	metrics SimMetrics

	rando       *rand.Rand
	rseed       int64
	CurrentTime float64
	endTime     float64

	Debug func(string, ...interface{})
}

type SimMetrics struct {
	TotalDamage    float64
	ReportedDamage float64 // used when DPSReportTime is set
	DamageAtOOM    float64
	OOMAt          float64
	Casts          []*Cast
	ManaAtEnd      int
	Rotation       []string
}

// New sim contructs a simulator with the given stats / equipment / options.
//   Technically we can calculate stats from equip/options but want the ability to override those stats
//   mostly for stat weight purposes.
func NewSim(stats Stats, equip Equipment, options Options) *Simulation {
	if options.GCDMin == 0 {
		options.GCDMin = 0.75 // default to 0.75s GCD
	}

	sim := &Simulation{
		Stats:   stats,
		Options: options,
		_CDs:     make([]float64, MagicIDLen),
		Buffs:   Stats{StatLen: 0},
		Auras:   []Aura{},
		Equip:   equip,
		rseed:   options.RSeed,
		rando:   rand.New(rand.NewSource(options.RSeed)),
		Debug:   nil,
	}

	if options.Debug {
		sim.Debug = debugFunc(sim)
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

	sim.Agent = NewAgent(sim, options.AgentType)
	if sim.Agent == nil {
		fmt.Printf("[ERROR] No rotation given to sim.\n")
		return nil
	}

	return sim
}

// reset will set sim back and erase all current state.
// This is automatically called before every 'Run'
//  This includes resetting and reactivating always on trinkets, auras, set bonuses, etc
func (sim *Simulation) reset() {
	// sim.rseed++
	// sim.rando.Seed(sim.rseed)

	sim.destructionPotion = false
	sim.bloodlustCasts = 0
	sim.CurrentTime = 0.0
	sim.CurrentMana = sim.Stats[StatMana]
	sim.Buffs = Stats{StatLen: 0}
	sim._CDs = make([]float64, MagicIDLen)
	sim.Auras = []Aura{}
	sim.metrics = SimMetrics{
		Casts: make([]*Cast, 0, 1000),
	}

	if sim.Debug != nil {
		sim.Debug("SIM RESET\n")
		sim.Debug("----------------------\n")
	}

	// Activate all talents
	if sim.Options.Talents.LightningOverload > 0 {
		sim.addAura(AuraLightningOverload(sim.Options.Talents.LightningOverload))
	}

	// Chain lightning bounces
	if sim.Options.NumClTargets > 1 {
		sim.addAura(ActivateChainLightningBounce(sim))
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

	sim.ActivateSets()
	sim.Agent.Reset(sim)
}

// Run will run the simulation for number of seconds.
// Returns metrics for what was cast and how much damage was done.
func (sim *Simulation) Run(durationSeconds float64) SimMetrics {
	sim.endTime = durationSeconds
	sim.reset()

	for sim.CurrentTime < sim.endTime {
		sim.Spellcasting()

		if sim.Options.ExitOnOOM && sim.metrics.OOMAt > 0 {
			return sim.metrics
		}
		if sim.CurrentMana < 0 {
			panic("you should never have negative mana.")
		}
	}
	sim.metrics.ManaAtEnd = int(sim.CurrentMana)

	return sim.metrics
}

// Remove an aura by its ID, searches through auras
// and calls 'cleanAura'
func (sim *Simulation) removeAuraByID(id int32) {
	for i := range sim.Auras {
		if sim.Auras[i].ID == id {
			sim.cleanAura(i)
			break
		}
	}
}

// cleanAura will remove the given aura from the sim and release all references
// to prevent memory leaking.
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

	if sim.Debug != nil {
		sim.Debug(" -%s\n", AuraName(sim.Auras[i].ID))
	}
	sim.Auras = sim.Auras[:i+copy(sim.Auras[i:], sim.Auras[i+1:])]
}

// addAura will add a new aura to the simulation. If there is a matching aura ID
// it will be replaced with the newer aura.
// Auras with duration of 0 will be logged as activating but never added to simulation auras.
func (sim *Simulation) addAura(a Aura) {
	if sim.Debug != nil {
		sim.Debug(" +%s\n", AuraName(a.ID))
	}
	if a.Expires == 0 {
		return // no need to waste time adding aura that doesn't last.
	}
	for i := range sim.Auras {
		if sim.Auras[i].ID == a.ID {
			sim.Auras[i] = a // replace
			return
		}
	}
	sim.Auras = append(sim.Auras, a)
}

// Cast will actually cast and treat all casts as having no 'flight time'.
// This will activate any auras around casting, calculate hit/crit and add to sim metrics.
func (sim *Simulation) Cast(cast *Cast) {
	if sim.Debug != nil {
		sim.Debug("Current Mana %0.0f, Cast Cost: %0.0f\n", sim.CurrentMana, cast.ManaCost)
	}
	sim.CurrentMana -= cast.ManaCost

	for _, aur := range sim.Auras {
		if aur.OnCastComplete != nil {
			aur.OnCastComplete(sim, cast)
		}
	}
	hit := 0.83 + ((sim.Stats[StatSpellHit] + sim.Buffs[StatSpellHit]) / 1260.0) + cast.Hit // 12.6 hit == 1% hit
	hit = math.Min(hit, 0.99)                                                               // can't get away from the 1% miss

	dbgCast := cast.Spell.Name
	if sim.Debug != nil {
		sim.Debug("Completed Cast (%s)\n", dbgCast)
	}
	if sim.rando.Float64() < hit {
		sp := sim.Stats[StatSpellDmg] + sim.Buffs[StatSpellDmg] + cast.Spellpower
		dmg := (sim.rando.Float64() * (cast.Spell.MaxDmg - cast.Spell.MinDmg)) + cast.Spell.MinDmg
		dmg += (sp * cast.Spell.Coeff)
		if cast.DidDmg != 0 { // use the pre-set dmg
			dmg = cast.DidDmg
		}
		cast.DidHit = true

		crit := ((sim.Stats[StatSpellCrit] + sim.Buffs[StatSpellCrit]) / 2208.0) + cast.Crit // 22.08 crit == 1% crit
		if sim.rando.Float64() < crit {
			cast.DidCrit = true
			critBonus := 1.5 // fall back crit damage
			if cast.CritBonus != 0 {
				critBonus = cast.CritBonus // This means we had pre-set the crit bonus when the spell was created. CSD will modify this.
			}
			if cast.Spell.ID == MagicIDCL6 || cast.Spell.ID == MagicIDLB12 {
				critBonus *= 2 // This handles the 'Elemental Fury' talent which increases the crit bonus.
				critBonus -= 1 // reduce to multiplier instead of percent.
			}
			dmg *= critBonus
			if cast.Spell.ID != MagicIDTLCLB {
				// TLC does not proc focus.
				sim.addAura(AuraElementalFocus(sim.CurrentTime))
			}
			if sim.Debug != nil {
				dbgCast += " crit"
			}
		} else if sim.Debug != nil {
			dbgCast += " hit"
		}

		if sim.Options.Talents.Concussion > 0 && (cast.Spell.ID == MagicIDLB12 || cast.Spell.ID == MagicIDCL6) {
			// Talent Concussion
			dmg *= 1 + (0.01 * sim.Options.Talents.Concussion)
		}
		if sim.Options.Buffs.Misery {
			dmg *= 1.05
		}

		// Average Resistance (AR) = (Target's Resistance / (Caster's Level * 5)) * 0.75
		// P(x) = 50% - 250%*|x - AR| <- where X is %resisted
		// Using these stats:
		//    13.6% chance of
		resVal := sim.rando.Float64()
		if resVal < 0.18 { // 13% chance for 25% resist, 4% for 50%, 1% for 75%
			if sim.Debug != nil {
				dbgCast += " (partial resist: "
			}
			if resVal < 0.01 {
				dmg *= .25
				if sim.Debug != nil {
					dbgCast += "75%)"
				}
			} else if resVal < 0.05 {
				dmg *= .5
				if sim.Debug != nil {
					dbgCast += "50%)"
				}
			} else {
				dmg *= .75
				if sim.Debug != nil {
					dbgCast += "25%)"
				}
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
	} else {
		if sim.Debug != nil {
			dbgCast += " miss"
		}
		cast.DidDmg = 0
		cast.DidCrit = false
		cast.DidHit = false
		for _, aur := range sim.Auras {
			if aur.OnSpellMiss != nil {
				aur.OnSpellMiss(sim, cast)
			}
		}
	}

	if cast.Spell.Cooldown > 0 {
		sim.setCD(cast.Spell.ID, cast.Spell.Cooldown)
	}

	if sim.Debug != nil {
		sim.Debug("%s: %0.0f\n", dbgCast, cast.DidDmg)
	}

	sim.metrics.Casts = append(sim.metrics.Casts, cast)

	sim.metrics.TotalDamage += cast.DidDmg
	if sim.Options.DPSReportTime > 0 && sim.CurrentTime <= sim.Options.DPSReportTime {
		sim.metrics.ReportedDamage += cast.DidDmg
	}
}

// Activates set bonuses, returning the list of active bonuses.
func (sim *Simulation) ActivateSets() []string {
	active := []string{}
	// Activate Set Bonuses
	for _, set := range sets {
		itemCount := 0
		for _, i := range sim.Equip {
			if set.Items[i.Name] {
				itemCount++
				if bonus, ok := set.Bonuses[itemCount]; ok {
					active = append(active, set.Name+" ("+strconv.Itoa(itemCount)+"pc)")
					sim.addAura(bonus(sim))
				}
			}
		}
	}
	return active
}

// Spellcasting performs the core logic of the advancement of simulation state.
// It will call 'Cast' on a spell ready to cast.
// If not casting it will activate ablities/trinkets that are off CD and then choose a new spell to cast.
// It will pop mana potions if needed.
func (sim *Simulation) Spellcasting() {
	TryActivateDrums(sim)
	TryActivateBloodlust(sim)
	TryActivateEleMastery(sim)
	TryActivateRacial(sim)
	TryActivateDestructionPotion(sim)
	sim.TryActivateEquipment()

	didPot := false
	didPot = didPot || TryActivateDarkRune(sim)
	didPot = didPot || TryActivateSuperManaPotion(sim)

	// Choose next spell
	castingSpell := sim.Agent.ChooseSpell(sim, didPot)
	if castingSpell == nil {
		panic("Agent returned nil casting spell")
	}

	if sim.CurrentMana >= castingSpell.ManaCost {
		if sim.Debug != nil {
			sim.Debug("Start Casting %s Cast Time: %0.1fs\n", cast.Spell.Name, cast.CastTime)
		}

		sim.Agent.OnSpellAccepted(sim, castingSpell)
		sim.Advance(castingSpell.CastTime)
		sim.Cast(castingSpell)
	} else {
		// Not enough mana, wait until there is enough mana to cast the desired spell
		if sim.metrics.OOMAt == 0 {
			sim.metrics.OOMAt = sim.CurrentTime
			sim.metrics.DamageAtOOM = sim.metrics.TotalDamage
		}
		timeUntilRegen := (castingSpell.ManaCost - sim.CurrentMana) / sim.manaRegen()
		sim.Advance(timeUntilRegen)
		// Don't actually cast; let the next iteration do the cast, so we recheck for pots/CDs/etc
	}
}

// Advance moves time forward counting down auras, CDs, mana regen, etc
func (sim *Simulation) Advance(elapsedTime float64) {
	// MP5 regen
	sim.CurrentMana = math.Min(
		sim.Stats[StatMana],
		sim.CurrentMana + sim.manaRegen() * elapsedTime)

	// CDS
	for k := range sim._CDs {
		sim._CDs[k] = math.Max(0, sim._CDs[k] - elapsedTime)
	}

	todel := []int{}
	for i := range sim.Auras {
		if sim.Auras[i].Expires <= (sim.CurrentTime + elapsedTime) {
			todel = append(todel, i)
		}
	}
	for i := len(todel) - 1; i >= 0; i-- {
		sim.cleanAura(todel[i])
	}
	sim.CurrentTime += elapsedTime
}

// Returns rate of mana regen, as mana / second
func (sim *Simulation) manaRegen() float64 {
	return ((sim.Stats[StatMP5] + sim.Buffs[StatMP5]) / 5.0)
}

func (sim *Simulation) isOnCD(magicID int32) bool {
	return sim._CDs[magicID] > 0
}

func (sim *Simulation) getRemainingCD(magicID int32) float64 {
	return sim._CDs[magicID]
}

func (sim *Simulation) setCD(magicID int32, newCD float64) {
	sim._CDs[magicID] = newCD
}

// Pops any on-use trinkets / gear
func (sim *Simulation) TryActivateEquipment() {
	for _, item := range sim.activeEquip {
		if item.Activate == nil || item.ActivateCD == -1 { // ignore non-activatable, and always active items.
			continue
		}
		if sim.isOnCD(item.CoolID) {
			continue
		}
		if item.Slot == EquipTrinket && sim.isOnCD(MagicIDAllTrinket) {
			continue
		}
		sim.addAura(item.Activate(sim))
		sim.setCD(item.CoolID, item.ActivateCD)
		if item.Slot == EquipTrinket {
			sim.setCD(MagicIDAllTrinket, 30)
		}
	}
}
