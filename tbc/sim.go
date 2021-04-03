package tbc

import (
	"fmt"
	"math"
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

	Stats       Stats
	Buffs       Stats     // temp increases
	Equip       Equipment // Current Gear
	activeEquip Equipment // cache of gear that can activate.

	Options       Options
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
	DamageAtOOM float64
	OOMAt       int
	Casts       []*Cast
	ManaAtEnd   int
	Rotation    []string
}

type Options struct {
	SpellOrder []string
	RSeed      int64
	ExitOnOOM  bool

	NumBloodlust int
	NumDrums     int

	Buffs    Buffs
	Consumes Consumes
	Talents  Talents
	Totems   Totems
}

func (o Options) StatTotal(e Equipment) Stats {
	gearStats := e.Stats()
	stats := o.BaseStats()
	for i := range stats {
		stats[i] += gearStats[i]
	}

	stats = o.Talents.AddStats(o.Buffs.AddStats(o.Consumes.AddStats(o.Totems.AddStats(stats))))

	if o.Buffs.BlessingOfKings {
		stats[StatInt] *= 1.1 // blessing of kings
	}

	// Final calculations
	stats[StatSpellCrit] += (stats[StatInt] / 80) / 100
	stats[StatMana] += stats[StatInt] * 15
	fmt.Printf("\fFinal MP5: %f", (stats[StatMP5] + (stats[StatInt] * 0.06)))

	return stats
}

func (o Options) BaseStats() Stats {
	stats := Stats{
		StatInt:  104,  // Base
		StatMana: 2958, // level 70 shaman
		StatLen:  0,
	}
	return stats
}

type Totems struct {
	TotemOfWrath int
	WrathOfAir   bool
	ManaStream   bool
}

func (tt Totems) AddStats(s Stats) Stats {
	s[StatSpellCrit] += 66.24 * float64(tt.TotemOfWrath)
	s[StatSpellHit] += 37.8 * float64(tt.TotemOfWrath)
	if tt.WrathOfAir {
		s[StatSpellDmg] += 104
	}
	if tt.ManaStream {
		s[StatMP5] += 50
	}
	return s
}

type Talents struct {
	LightninOverload   int
	ElementalPrecision int
	NaturesGuidance    int
	TidalMastery       int
	ElementalMastery   bool
	UnrelentingStorm   int
	CallOfThunder      int
	Convection         int
	Concussion         int
}

func (t Talents) AddStats(s Stats) Stats {
	s[StatSpellHit] += 25.2 * float64(t.ElementalPrecision)
	s[StatSpellHit] += 12.6 * float64(t.NaturesGuidance)
	s[StatSpellCrit] += 22.08 * float64(t.TidalMastery)
	s[StatSpellCrit] += 22.08 * float64(t.CallOfThunder)

	return s
}

type Buffs struct {
	// Raid buffs
	ArcaneInt                bool
	GiftOftheWild            bool
	BlessingOfKings          bool
	ImprovedBlessingOfWisdom bool

	// Party Buffs
	Moonkin    bool
	SpriestDPS int // adds Mp5 ~ 25% (dps*5%*5sec = 25%)

	// Self Buffs
	WaterShield    bool
	WaterShieldPPM int // how many procs per minute does watershield get? Every 3 requires a recast.

	// Target Debuff
	JudgementOfWisdom bool
}

func (b Buffs) AddStats(s Stats) Stats {
	if b.ArcaneInt {
		s[StatInt] += 40
	}
	if b.GiftOftheWild {
		s[StatInt] += 18 // assumes improved gotw, rounded down to nearest int... not sure if that is accurate.
	}
	if b.ImprovedBlessingOfWisdom {
		s[StatMP5] += 42
	}
	if b.Moonkin {
		s[StatSpellCrit] += 110.4
	}
	if b.WaterShield {
		s[StatMP5] += 50
	}
	s[StatMP5] += float64(b.SpriestDPS) * 0.25

	return s
}

type Consumes struct {
	// Buffs
	BrilliantWizardOil       bool
	MajorMageblood           bool
	FlaskOfBlindingLight     bool
	FlaskOfMightyRestoration bool
	BlackendBasilisk         bool

	// Used in rotations
	SuperManaPotion bool
	DarkRune        bool
}

func (c Consumes) AddStats(s Stats) Stats {
	if c.BrilliantWizardOil {
		s[StatSpellCrit] += 14
		s[StatSpellDmg] += 36
	}
	if c.MajorMageblood {
		s[StatMP5] += 16.0
	}
	if c.FlaskOfBlindingLight {
		s[StatSpellDmg] += 80
	}
	if c.FlaskOfMightyRestoration {
		s[StatMP5] += 25
	}
	if c.BlackendBasilisk {
		s[StatSpellDmg] += 23
	}
	return s
}

// New sim contructs a simulator with the given stats / equipment / options.
//   Technically we can calculate stats from equip/options but want the ability to override those stats
//   mostly for stat weight purposes.
func NewSim(stats Stats, equip Equipment, options Options) *Simulation {
	if len(options.SpellOrder) == 0 {
		fmt.Printf("[ERROR] No rotation given to sim.\n")
		return nil
	}
	rotIdx := 0
	if options.SpellOrder[0] == "pri" {
		rotIdx = -1
		options.SpellOrder = options.SpellOrder[1:]
	}
	sim := &Simulation{
		RotationIdx:   rotIdx,
		Stats:         stats,
		SpellRotation: options.SpellOrder,
		Options:       options,
		CDs:           map[string]int{},
		Buffs:         Stats{StatLen: 0},
		Auras:         []Aura{},
		Equip:         equip,
		rseed:         options.RSeed,
		rando:         rand.New(rand.NewSource(options.RSeed)),
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

	// Activate all talents
	if sim.Options.Talents.LightninOverload > 0 {
		sim.addAura(AuraLightningOverload(sim.Options.Talents.LightninOverload))
	}

	// Judgement of Wisdom
	if sim.Options.Buffs.JudgementOfWisdom {
		sim.addAura(AuraJudgementOfWisdom())
	}

	// Activate all permanent item effects.
	for _, item := range sim.Equip {
		if item.Activate != nil && item.ActivateCD == -1 {
			sim.addAura(item.Activate(sim))
		}
	}

	debug("\nRotation: %v\n", sim.SpellRotation)
	debug("Effective MP5: %0.1f\n", sim.Stats[StatMP5]+sim.Buffs[StatMP5])
	debug("----------------------\n")
}

func (sim *Simulation) Run(seconds int) SimMetrics {
	// For now use the new 'event' driven state advancement.
	return sim.Run2(seconds)
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

func (sim *Simulation) ChooseSpell() int {
	if sim.RotationIdx == -1 {
		lowestWait := math.MaxInt32
		wasMana := false
		for i := 0; i < len(sim.SpellRotation); i++ {
			so := sim.SpellRotation[i]
			sp := spellmap[so]
			cast := NewCast(sim, sp, sim.Stats[StatSpellDmg], sim.Stats[StatSpellHit], sim.Stats[StatSpellCrit])
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

	so := sim.SpellRotation[sim.RotationIdx]
	sp := spellmap[so]
	cast := NewCast(sim, sp, sim.Stats[StatSpellDmg], sim.Stats[StatSpellHit], sim.Stats[StatSpellCrit])
	if sim.CDs[so] < 1 {
		if sim.CurrentMana >= cast.ManaCost {
			sim.CastingSpell = cast
			sim.RotationIdx++
			if sim.RotationIdx == len(sim.SpellRotation) {
				sim.RotationIdx = 0
			}
			return cast.TicksUntilCast
		} else {
			debug("Current Mana %0.0f, Cast Cost: %0.0f\n", sim.CurrentMana, cast.ManaCost)
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
	if sim.rando.Float64() < cast.Hit {
		dmg := (float64(sim.rando.Intn(int(cast.Spell.MaxDmg-cast.Spell.MinDmg))) + cast.Spell.MinDmg) + (sim.Stats[StatSpellDmg] * cast.Spell.Coeff)
		if cast.DidDmg != 0 { // use the pre-set dmg
			dmg = cast.DidDmg
		}
		cast.DidHit = true
		dbgCast := "hit"
		if sim.rando.Float64() < cast.Crit {
			cast.DidCrit = true
			dmg *= 2
			sim.addAura(AuraElementalFocus(sim.currentTick))
			dbgCast = "crit"
		}

		if sim.Options.Talents.Concussion > 0 && (strings.HasPrefix(cast.Spell.ID, "LB") || strings.HasPrefix(cast.Spell.ID, "CL")) {
			// Talent Concussion
			dmg *= 1 + (0.01 * float64(sim.Options.Talents.Concussion))
		}

		// Average Resistance = (Target's Resistance / (Caster's Level * 5)) * 0.75 "AR"
		// P(x) = 50% - 250%*|x - AR| <- where X is chance of resist
		// For now hardcode the 25% chance resist at 2.5% (this assumes bosses have 0 nature resist)
		if sim.rando.Float64() < 0.025 { // chance of 25% resist
			dmg *= .75
			debug("(partial resist)")
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
		debug("%s: %0.0f\n", dbgCast, cast.DidDmg)
		sim.metrics.TotalDamage += cast.DidDmg
		sim.metrics.Casts = append(sim.metrics.Casts, cast)
	} else {
		debug("miss.\n")
	}

	sim.CurrentMana -= cast.ManaCost
	sim.CastingSpell = nil
	if cast.Spell.Cooldown > 0 {
		sim.CDs[cast.Spell.ID] = cast.Spell.Cooldown * TicksPerSecond
	}
}
