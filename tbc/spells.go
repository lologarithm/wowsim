package tbc

import "math"

type Cast struct {
	Spell *Spell
	// Caster ... // Needed for onstruck effects?
	IsLO bool // stupid hack

	// Pre-hit Mutatable State
	TicksUntilCast int
	CastTime       float64 // time in seconds to cast the spell
	ManaCost       float64

	Hit        float64 // Direct % bonus... 0.1 == 10%
	Crit       float64 // Direct % bonus... 0.1 == 10%
	CritBonus  float64 // Multiplier to critical dmg bonus.
	Spellpower float64 // Bonus Spellpower to add at end of cast.

	// Calculated Values
	DidHit  bool
	DidCrit bool
	DidDmg  float64
	CastAt  int // simulation tick the spell cast

	Effects []AuraEffect // effects applied ONLY to this cast.
}

// NewCast constructs a Cast from the current simulation and selected spell.
//  OnCast mechanics are applied at this time (anything that modifies the cast before its cast, usually just mana cost stuff)
func NewCast(sim *Simulation, sp *Spell) *Cast {
	cast := &Cast{
		Spell:      sp,
		ManaCost:   float64(sp.Mana),
		Spellpower: 0, // TODO: type specific bonuses...
		CritBonus:  1.5,
	}

	castTime := sp.CastTime
	itsElectric := sp.ID == MagicIDLB12 || sp.ID == MagicIDCL6

	if itsElectric {
		// TODO: Add LightningMaster to talent list (this will never not be selected for an elemental shaman)
		castTime -= 0.5 // Talent Lightning Mastery
	}
	castTime /= (1 + ((sim.Stats[StatHaste] + sim.Buffs[StatHaste]) / 1576)) // 15.76 rating grants 1% spell haste
	if castTime < 1.0 {
		castTime = 1.0 // can't cast faster than 1/sec even with max haste.
	}
	cast.CastTime = castTime
	cast.TicksUntilCast = int(castTime*float64(TicksPerSecond)) + 1 // round up

	if itsElectric {
		cast.ManaCost *= 1 - (0.02 * float64(sim.Options.Talents.Convection))
	}

	// Apply any on cast effects.
	for _, aur := range sim.Auras {
		if aur.OnCast != nil {
			aur.OnCast(sim, cast)
		}
	}

	return cast
}

// ChooseSpell is a basic rotation spell selector. This is the default
// spell selection logic if not using the 'ai optimizer' for selecting spells.
func ChooseSpell(sim *Simulation, _ bool) int {
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
			sim.metrics.OOMAt = sim.CurrentTick / TicksPerSecond
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
			if sim.Debug != nil {
				sim.Debug("Current Mana %0.0f, Cast Cost: %0.0f\n", sim.CurrentMana, cast.ManaCost)
			}
			if sim.metrics.OOMAt == 0 {
				sim.metrics.OOMAt = sim.CurrentTick / TicksPerSecond
				sim.metrics.DamageAtOOM = sim.metrics.TotalDamage
			}
			return int(math.Ceil((cast.ManaCost - sim.CurrentMana) / sim.manaRegen()))
		}
	}
	return sim.CDs[so]
}

// Spell represents a single castable spell. This is all the data needed to begin a cast.
type Spell struct {
	ID         int32
	Name       string
	CastTime   float64
	Cooldown   int
	Mana       float64
	MinDmg     float64
	MaxDmg     float64
	DamageType DamageType
	Coeff      float64

	DotDmg float64
	DotDur float64
}

// DamageType is currently unused.
type DamageType byte

const (
	DamageTypeUnknown DamageType = iota
	DamageTypeFire
	DamageTypeNature
	DamageTypeFrost

	// Who cares about these fake damage types.
	DamageTypeShadow
	DamageTypeHoly
	DamageTypeArcane
)

// All Spells
// TODO: Downrank Penalty == (spellrankavailbetobetrained+11)/70
//    Might not even be worth calculating because I don't think there is much case for ever downranking.
var spells = []Spell{
	// {ID: MagicIDLB4, Name: "LB4", Coeff: 0.795, CastTime: 2.5, MinDmg: 88, MaxDmg: 100, Mana: 50, DamageType: DamageTypeNature},
	// {ID: MagicIDLB10, Name: "LB10", Coeff: 0.795, CastTime: 2.5, MinDmg: 428, MaxDmg: 477, Mana: 265, DamageType: DamageTypeNature},
	{ID: MagicIDLB12, Name: "LB12", Coeff: 0.795, CastTime: 2.5, MinDmg: 571, MaxDmg: 652, Mana: 300, DamageType: DamageTypeNature},
	// {ID: MagicIDCL4, Name: "CL4", Coeff: 0.643, CastTime: 2, Cooldown: 6, MinDmg: 505, MaxDmg: 564, Mana: 605, DamageType: DamageTypeNature},
	{ID: MagicIDCL6, Name: "CL6", Coeff: 0.643, CastTime: 2, Cooldown: 6, MinDmg: 734, MaxDmg: 838, Mana: 760, DamageType: DamageTypeNature},
	// {ID: MagicIDES8, Name: "ES8", Coeff: 0.3858, CastTime: 1.5, Cooldown: 6, MinDmg: 658, MaxDmg: 692, Mana: 535, DamageType: DamageTypeNature},
	// {ID: MagicIDFrS5, Name: "FrS5", Coeff: 0.3858, CastTime: 1.5, Cooldown: 6, MinDmg: 640, MaxDmg: 676, Mana: 525, DamageType: DamageTypeFrost},
	// {ID: MagicIDFlS7, Name: "FlS7", Coeff: 0.15, CastTime: 1.5, Cooldown: 6, MinDmg: 377, MaxDmg: 420, Mana: 500, DotDmg: 100, DotDur: 6, DamageType: DamageTypeFire},
	{ID: MagicIDTLCLB, Name: "TLCLB", Coeff: 0.0, CastTime: 0, MinDmg: 694, MaxDmg: 807, Mana: 0, DamageType: DamageTypeNature},
}

// Spell lookup map to make lookups faster.
var spellmap = map[int32]*Spell{}

func init() {
	for _, sp := range spells {
		// Turns out to increase efficiency go 'range' will actually only allocate a single struct and mutate.
		// If we want to create a pointer we need to clone the struct.
		sp2 := sp
		spp := &sp2
		spellmap[sp.ID] = spp
	}
}
