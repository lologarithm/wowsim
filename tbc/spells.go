package tbc

import "math"

// Spell represents a single castable spell. This is all the data needed to begin a cast.
type Spell struct {
	ID         int32
	Name       string
	CastTime   float64
	Cooldown   float64
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

type Cast struct {
	Spell *Spell
	// Caster ... // Needed for onstruck effects?
	IsLO       bool // stupid hack
	IsClBounce bool // stupider hack

	// Pre-hit Mutatable State
	CastTime float64 // time in seconds to cast the spell
	ManaCost float64

	Hit        float64 // Direct % bonus... 0.1 == 10%
	Crit       float64 // Direct % bonus... 0.1 == 10%
	CritBonus  float64 // Multiplier to critical dmg bonus.
	Spellpower float64 // Bonus Spellpower to add at end of cast.

	// Calculated Values
	DidHit  bool
	DidCrit bool
	DidDmg  float64
	CastAt  float64 // simulation time the spell cast

	Effects []AuraEffect // effects applied ONLY to this cast.
}

// NewCast constructs a Cast from the current simulation and selected spell.
//  OnCast mechanics are applied at this time (anything that modifies the cast before its cast, usually just mana cost stuff)
func NewCast(sim *Simulation, sp *Spell) *Cast {
	cast := &Cast{
		Spell:      sp,
		ManaCost:   float64(sp.Mana),
		Spellpower: 0,
		CritBonus:  1.5,
	}

	castTime := sp.CastTime
	itsElectric := sp.ID == MagicIDLB12 || sp.ID == MagicIDCL6

	if itsElectric {
		// TODO: Add LightningMaster to talent list (this will never not be selected for an elemental shaman)
		castTime -= 0.5 // Talent Lightning Mastery
	}
	castTime /= (1 + ((sim.Stats[StatHaste] + sim.Buffs[StatHaste]) / 1576)) // 15.76 rating grants 1% spell haste
	castTime = math.Max(castTime, sim.Options.GCDMin)                        // can't cast faster than GCD
	cast.CastTime = castTime

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
