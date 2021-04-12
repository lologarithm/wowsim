package tbc

type Cast struct {
	Spell *Spell
	// Caster ... // Needed for onstruck effects?
	IsLO bool // stupid hack

	// Pre-hit Mutatable State
	TicksUntilCast int
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

func NewCast(sim *Simulation, sp *Spell) *Cast {
	cast := &Cast{
		Spell:      sp,
		ManaCost:   float64(sp.Mana),
		Spellpower: 0, // TODO: type specific bonuses...
		CritBonus:  1.5,
	}

	castTime := sp.CastTime
	isLB := sp.ID == MagicIDLB12
	isCL := sp.ID == MagicIDCL6

	if isLB || isCL {
		// Talent to reduce cast time.
		castTime -= 0.5 // Talent Lightning Mastery
	}
	castTime /= (1 + ((sim.Stats[StatHaste] + sim.Buffs[StatHaste]) / 1576)) // 15.76 rating grants 1% spell haste
	if castTime < 1.0 {
		castTime = 1.0 // can't cast faster than 1/sec even with max haste.
	}
	cast.TicksUntilCast = int(castTime * float64(TicksPerSecond))

	if isLB || isCL {
		cast.ManaCost *= 1 - (0.2 * float64(sim.Options.Talents.Convection))
	}

	// Apply any on cast effects.
	for _, aur := range sim.Auras {
		if aur.OnCast != nil {
			aur.OnCast(sim, cast)
		}
	}

	return cast
}

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

type DamageType byte

const (
	DamageTypeUnknown DamageType = iota
	DamageTypeFire
	DamageTypeNature
	DamageTypeFrost

	// who cares
	DamageTypeShadow
	DamageTypeHoly
	DamageTypeArcane
)

// spells
// TODO: DRP == (spellrankavailbetobetrained+11)/70
var spells = []Spell{
	// {ID: MagicIDLB4, Name: "LB4", Coeff: 0.795, CastTime: 2.0, MinDmg: 88, MaxDmg: 100, Mana: 50, DamageType: DamageTypeNature},
	// {ID: MagicIDLB10, Name: "LB10", Coeff: 0.795, CastTime: 2.5, MinDmg: 428, MaxDmg: 477, Mana: 265, DamageType: DamageTypeNature},
	{ID: MagicIDLB12, Name: "LB12", Coeff: 0.795, CastTime: 2.5, MinDmg: 563, MaxDmg: 643, Mana: 300, DamageType: DamageTypeNature},
	// {ID: MagicIDCL4, Name: "CL4", Coeff: 0.643, CastTime: 2, Cooldown: 6, MinDmg: 505, MaxDmg: 564, Mana: 605, DamageType: DamageTypeNature},
	{ID: MagicIDCL6, Name: "CL6", Coeff: 0.643, CastTime: 2, Cooldown: 6, MinDmg: 734, MaxDmg: 838, Mana: 760, DamageType: DamageTypeNature},
	// {ID: MagicIDES8, Name: "ES8", Coeff: 0.3858, CastTime: 1.5, Cooldown: 6, MinDmg: 658, MaxDmg: 692, Mana: 535, DamageType: DamageTypeNature},
	// {ID: MagicIDFrS5, Name: "FrS5", Coeff: 0.3858, CastTime: 1.5, Cooldown: 6, MinDmg: 640, MaxDmg: 676, Mana: 525, DamageType: DamageTypeFrost},
	// {ID: MagicIDFlS7, Name: "FlS7", Coeff: 0.15, CastTime: 1.5, Cooldown: 6, MinDmg: 377, MaxDmg: 420, Mana: 500, DotDmg: 100, DotDur: 6, DamageType: DamageTypeFire},
}

var spellmap = map[int32]*Spell{}

func init() {
	for _, sp := range spells {
		sp2 := sp //wtf go?
		spp := &sp2
		if spp.Coeff == 0 {
			spp.Coeff = spp.CastTime / 3.5
		}
		spellmap[sp.ID] = spp
	}
}
