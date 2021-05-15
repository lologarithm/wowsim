package tbc

import (
	"fmt"
	"math/rand"
)

func debugFunc(sim *Simulation) func(string, ...interface{}) {
	return func(s string, vals ...interface{}) {
		fmt.Printf("[%0.1f] "+s, append([]interface{}{(float64(sim.CurrentTick) / float64(TicksPerSecond))}, vals...)...)
	}
}

type Simulation struct {
	CurrentMana float64

	SpellChooser func(*Simulation, bool) int // TODO: make more funtional. Return a cast instead of having function mutate sim itself.

	Stats       Stats
	Buffs       Stats     // temp increases
	Equip       Equipment // Current Gear
	activeEquip Equipment // cache of gear that can activate.

	bloodlustCasts    int
	destructionPotion bool
	Options           Options
	SpellRotation     []*Spell
	RotationIdx       int

	// ticks until cast is complete
	CastingSpell *Cast

	// timeToRegen := 0
	CDs   map[int32]int // Map of MagicID to ticks until CD is done. 'Advance' counts down these
	Auras []Aura        // this is array instaed of map to speed up browser perf.

	// Clears and regenerates on each Run call.
	metrics SimMetrics

	rando       *rand.Rand
	rseed       int64
	CurrentTick int
	endTick     int

	Debug func(string, ...interface{})
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
		Debug:         nil,
		SpellChooser:  ChooseSpell,
	}

	if options.UseAI {
		ai := NewAI(sim)
		sim.SpellChooser = ai.ChooseSpell
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
	return sim
}

// reset will set sim back and erase all current state.
// This is automatically called before every 'Run'
//  This includes resetting and reactivating always on trinkets, auras, set bonuses, etc
func (sim *Simulation) reset() {
	// sim.rseed++
	// sim.rando.Seed(sim.rseed)

	sim.bloodlustCasts = 0
	sim.CurrentTick = 0
	sim.CurrentMana = sim.Stats[StatMana]
	sim.CastingSpell = nil
	sim.Buffs = Stats{StatLen: 0}
	sim.CDs = map[int32]int{}
	sim.Auras = []Aura{}
	sim.metrics = SimMetrics{}

	if sim.Debug != nil {
		sim.Debug("SIM RESET\n")
		sim.Debug("----------------------\n")
	}

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

	sim.ActivateSets()

	if sim.Options.UseAI {
		// Reset a new AI
		// TODO: Can we take learnings from the last AI to modulate this AIs behavior?
		ai := NewAI(sim)
		sim.SpellChooser = ai.ChooseSpell
	}
}

// Run will run the simulation for number of seconds.
// Returns metrics for what was cast and how much damage was done.
func (sim *Simulation) Run(seconds int) SimMetrics {
	sim.endTick = seconds * TicksPerSecond
	sim.reset()

	for i := 0; i < sim.endTick; {
		if sim.CurrentMana < 0 {
			panic("you should never have negative mana.")
		}

		sim.CurrentTick = i
		advance := sim.Spellcasting(i)

		if sim.Options.ExitOnOOM && sim.metrics.OOMAt > 0 {
			return sim.metrics
		}

		sim.Advance(i, advance)
		i += advance
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
	for _, aur := range sim.Auras {
		if aur.OnCastComplete != nil {
			aur.OnCastComplete(sim, cast)
		}
	}
	hit := 0.83 + ((sim.Stats[StatSpellHit] + sim.Buffs[StatSpellHit]) / 1260.0) + cast.Hit // 12.6 hit == 1% hit
	if hit > 0.99 {
		hit = 0.99 // can't get away from the 1% miss
	}

	if sim.Debug != nil {
		sim.Debug("Completed Cast (%s)\n", cast.Spell.Name)
	}
	dbgCast := cast.Spell.Name
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
				sim.addAura(AuraElementalFocus(sim.CurrentTick))
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
		sim.metrics.TotalDamage += cast.DidDmg
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
	sim.metrics.Casts = append(sim.metrics.Casts, cast)
	if sim.Debug != nil {
		sim.Debug("%s: %0.0f\n", dbgCast, cast.DidDmg)
	}
	sim.CurrentMana -= cast.ManaCost
	sim.CastingSpell = nil
	if cast.Spell.Cooldown > 0 {
		sim.CDs[cast.Spell.ID] = cast.Spell.Cooldown * TicksPerSecond
	}
}

func (sim *Simulation) ActivateRacial() {
	switch v := sim.Options.Buffs.Race; v {
	case RaceBonusOrc:
		const spBonus = 143
		const dur = 15
		if sim.CDs[MagicIDOrcBloodFury] < 1 {
			sim.Buffs[StatSpellDmg] += spBonus
			sim.addAura(AuraStatRemoval(sim.CurrentTick, dur, spBonus, StatSpellDmg, MagicIDOrcBloodFury))
			sim.CDs[MagicIDOrcBloodFury] = 120 * TicksPerSecond
		}
	case RaceBonusTroll10, RaceBonusTroll30:
		hasteBonus := 1.1 // 10% haste
		const dur = 10
		if v == RaceBonusTroll30 {
			hasteBonus = 1.3 // 30% haste
		}
		if sim.CDs[MagicIDTrollBerserking] < 1 {
			sim.addAura(ActivateBerserking(sim, hasteBonus))
		}
	}
}

func (sim *Simulation) ActivateSets() {
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

// Spellcasting performs the core logic of the advancement of simulation state.
// It will call 'Cast' on a spell ready to cast.
// If not casting it will activate ablities/trinkets that are off CD and then choose a new spell to cast.
//  It will pop mana potions if needed.
// Returns the number of system ticks until the next action will be ready.
func (sim *Simulation) Spellcasting(tickID int) int {
	// technically we dont really need this check with the new advancer.
	if sim.CastingSpell != nil && sim.CastingSpell.TicksUntilCast == 0 {
		sim.Cast(sim.CastingSpell)
	}

	if sim.CastingSpell == nil {
		if sim.Options.NumDrums > 0 && sim.CDs[MagicIDDrums] < 1 {
			// We have drums in the sim, and the drums aura isn't turned on.
			// Iterate our drum
			for i, v := range []int32{MagicIDDrum1, MagicIDDrum2, MagicIDDrum3, MagicIDDrum4} {
				if i == sim.Options.NumDrums {
					break
				}
				if sim.CDs[v] < 1 {
					sim.CDs[v] = 120 * TicksPerSecond // item goes on CD for 120s
					sim.addAura(ActivateDrums(sim))
					break
				}
			}
		}
		// Activate any specials
		if sim.Options.NumBloodlust > sim.bloodlustCasts && sim.CDs[MagicIDBloodlust] < 1 {
			sim.addAura(ActivateBloodlust(sim))
			sim.bloodlustCasts++ // TODO: will this break anything?
		}

		if sim.Options.Talents.ElementalMastery && sim.CDs[MagicIDEleMastery] < 1 {
			// Apply auras
			sim.addAura(AuraEleMastery())
		}

		sim.ActivateRacial()

		if sim.Options.Consumes.DestructionPotion && sim.CDs[MagicIDPotion] < 1 {
			// Only use dest potion if not using mana or if we haven't used it once.
			// If we are using mana, only use destruction potion on the pull.
			if !sim.Options.Consumes.SuperManaPotion || !sim.destructionPotion {
				sim.addAura(ActivateDestructionPotion(sim))
			}
		}

		didPot := false
		totalRegen := (sim.Stats[StatMP5] + sim.Buffs[StatMP5])
		// Pop potion before next cast if we have less than the mana provided by the potion minues 1mp5 tick.
		if sim.Options.Consumes.DarkRune && sim.Stats[StatMana]-sim.CurrentMana+totalRegen >= 1500 && sim.CDs[MagicIDRune] < 1 {
			// Restores 900 to 1500 mana. (2 Min Cooldown)
			sim.CurrentMana += 900 + (sim.rando.Float64() * 600)
			sim.CDs[MagicIDRune] = 120 * TicksPerSecond
			didPot = true
			if sim.Debug != nil {
				sim.Debug("Used Dark Rune\n")
			}
		}
		if sim.Options.Consumes.SuperManaPotion && sim.Stats[StatMana]-sim.CurrentMana+totalRegen >= 3000 && sim.CDs[MagicIDPotion] < 1 {
			// Restores 1800 to 3000 mana. (2 Min Cooldown)
			sim.CurrentMana += 1800 + (sim.rando.Float64() * 1200)
			sim.CDs[MagicIDPotion] = 120 * TicksPerSecond
			didPot = true
			if sim.Debug != nil {
				sim.Debug("Used Mana Potion\n")
			}
		}

		// Pop any on-use trinkets
		for _, item := range sim.activeEquip {
			if item.Activate == nil || item.ActivateCD == -1 { // ignore non-activatable, and always active items.
				continue
			}
			if sim.CDs[item.CoolID] > 0 {
				continue
			}
			if item.Slot == EquipTrinket && sim.CDs[MagicIDAllTrinket] > 0 {
				continue
			}
			sim.addAura(item.Activate(sim))
			sim.CDs[item.CoolID] = item.ActivateCD * TicksPerSecond
			if item.Slot == EquipTrinket {
				sim.CDs[MagicIDAllTrinket] = 30 * TicksPerSecond
			}
		}

		// Choose next spell
		ticks := sim.SpellChooser(sim, didPot)
		if sim.CastingSpell != nil {
			if sim.Debug != nil {
				sim.Debug("Start Casting %s Cast Time: %0.1fs\n", sim.CastingSpell.Spell.Name, float64(sim.CastingSpell.TicksUntilCast)/float64(TicksPerSecond))
			}
		}
		return ticks
	}

	return 1
}

// Advance moves time forward counting down auras, CDs, mana regen, etc
func (sim *Simulation) Advance(tickID int, ticks int) {

	if sim.CastingSpell != nil {
		sim.CastingSpell.TicksUntilCast -= ticks
	}

	// MP5 regen
	sim.CurrentMana += sim.manaRegen() * float64(ticks)

	if sim.CurrentMana > sim.Stats[StatMana] {
		sim.CurrentMana = sim.Stats[StatMana]
	}

	// CDS
	for k := range sim.CDs {
		sim.CDs[k] -= ticks
		if sim.CDs[k] < 1 {
			delete(sim.CDs, k)
		}
	}

	todel := []int{}
	for i := range sim.Auras {
		if sim.Auras[i].Expires <= (tickID + ticks) {
			todel = append(todel, i)
		}
	}
	for i := len(todel) - 1; i >= 0; i-- {
		sim.cleanAura(todel[i])
	}
}

func (sim *Simulation) manaRegen() float64 {
	return ((sim.Stats[StatMP5] + sim.Buffs[StatMP5]) / 5.0) / float64(TicksPerSecond)
}
