package elesim

type Aura struct {
	ID         string
	Duration   int // ticks aura will apply
	OnCast     func(c *Cast)
	OnSpellHit func(c *Cast)
	OnStruck   func(c *Cast)
}

type Cast struct {
	Spell *Spell
	// Caster ... // Needed for onstruck effects?

	// Mutatable State
	ManaCost   float64
	Hit        float64
	Crit       float64
	Spellpower float64
}

type Spell struct {
	ID         string
	CastTime   float64
	Cooldown   int
	Mana       float64
	MinDmg     float64
	MaxDmg     float64
	DamageType DamageType
	Coeff      float64
}

type DamageType int

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

func clearcasting() Aura {
	return Aura{
		Duration: 15 * tickPerSecond,
		OnCast: func(c *Cast) {
			debug("clearcasting...")
			c.ManaCost = 0
		},
	}
}

func zhc() Aura {
	dmgbonus := 204.0

	return Aura{
		Duration: 20 * tickPerSecond,
		OnCast: func(c *Cast) {
			debug("zhc(%.0f)...", dmgbonus)
			c.Spellpower += dmgbonus
			dmgbonus -= 17
		},
	}
}

func elemastery() Aura {
	return Aura{
		Duration: 99999999999999999,
		OnCast: func(c *Cast) {
			debug("ele mastery...")
			c.Crit = 1.01 // 101% chance of crit
			c.ManaCost = 0
		},
	}
}

func stormcaller() Aura {
	return Aura{
		Duration: 8 * tickPerSecond,
		OnCast: func(c *Cast) {
			debug("stormcaller...")
			c.Spellpower += 50
		},
	}
}

// spells
var spells = []Spell{
	{ID: "LB4", CastTime: 2.0, MinDmg: 88, MaxDmg: 100, Mana: 50, DamageType: DamageTypeNature},
	{ID: "LB10", CastTime: 3.0, MinDmg: 428, MaxDmg: 477, Mana: 265, DamageType: DamageTypeNature},
	{ID: "CL4", CastTime: 2.5, Cooldown: 6, MinDmg: 505, MaxDmg: 564, Mana: 605, DamageType: DamageTypeNature},
}

var spellmap = map[string]*Spell{}

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
