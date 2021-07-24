package tbc

import (
	"encoding/binary"
	"math"
)

type Options struct {
	SpellOrder []string
	NumClTargets int
	UseAI      bool // when set true, the AI will modulate the rotations to maximize DPS and mana.
	RSeed      int64
	ExitOnOOM  bool
	GCD        float64 // sets the GCD

	NumBloodlust int
	NumDrums     int

	Buffs    Buffs
	Consumes Consumes
	Talents  Talents
	Totems   Totems

	DPSReportTime int // how many seconds to calculate DPS for.

	Debug bool // enables debug printing.
	// TODO: could change this to be a func/stream consumer could provide,
	// make it easier to integrate into different output systems.

	// Hack indicating whether tidefury 2 piece bonus (CL bounce damage) is active
	// This is only set from the aura, not from actual options
	Tidefury2Pc bool
}

// Pack is how to convert all options/buffs/consumes/etc to reproduce the UI state
// so that a simulation can be shared. I am using byte packing here because most options are bools
// and this makes it easy to pack it all together.
// I also have a JSON parser for these in the command line interface. This compressed format
// is used to allow for shorter URLs
func (o Options) Pack() []byte {
	// first byte is version
	bytes := []byte{0, byte(o.NumBloodlust), byte(o.NumDrums)}
	bytes = append(bytes, o.Buffs.Pack()...)
	bytes = append(bytes, o.Consumes.Pack()...)
	bytes = append(bytes, o.Talents.Pack()...)
	bytes = append(bytes, o.Totems.Pack()...)
	return bytes
}

type Totems struct {
	TotemOfWrath int
	WrathOfAir   bool
	ManaStream   bool
	Cyclone2PC   bool // Cyclone set 2pc bonus
}

func (tt Totems) Pack() []byte {
	var opts byte
	if tt.WrathOfAir {
		opts = opts | 1
	}
	if tt.ManaStream {
		opts = opts | 2
	}
	if tt.Cyclone2PC {
		opts = opts | 4
	}
	bytes := []byte{byte(tt.TotemOfWrath), opts}

	return bytes
}

func (tt Totems) AddStats(s Stats) Stats {
	s[StatSpellCrit] += 66.24 * float64(tt.TotemOfWrath)
	s[StatSpellHit] += 37.8 * float64(tt.TotemOfWrath)
	if tt.WrathOfAir {
		s[StatSpellDmg] += 101
		if tt.Cyclone2PC {
			s[StatSpellDmg] += 20
		}
	}
	if tt.ManaStream {
		s[StatMP5] += 50
	}
	return s
}

type Talents struct {
	LightningOverload   int
	ElementalPrecision int
	NaturesGuidance    int
	TidalMastery       int
	ElementalMastery   bool
	UnrelentingStorm   int
	CallOfThunder      int
	Convection         int
	Concussion         float64 // temp hack to speed up not converting this to a int on every spell cast
}

func (t Talents) Pack() []byte {
	var elemast byte
	if t.ElementalMastery {
		elemast = 1
	}
	bytes := []byte{
		byte(t.LightningOverload),
		byte(t.ElementalPrecision),
		byte(t.NaturesGuidance),
		byte(t.TidalMastery),
		elemast,
		byte(t.UnrelentingStorm),
		byte(t.CallOfThunder),
		byte(t.Convection),
		byte(t.Concussion),
	}

	return bytes
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
	ImprovedDivineSpirit     bool

	// Party Buffs
	Moonkin             bool
	MoonkinRavenGoddess bool   // adds 20 spell crit to moonkin aura
	SpriestDPS          uint16 // adds Mp5 ~ 25% (dps*5%*5sec = 25%)
	EyeOfNight          bool   // Eye of night bonus from party member (not you)
	TwilightOwl         bool   // from party member

	// Self Buffs
	WaterShield    bool
	WaterShieldPPM byte // how many procs per minute does watershield get? Every 3 requires a recast.
	Race           RaceBonusType

	// Target Debuff
	JudgementOfWisdom bool
	ImpSealofCrusader bool
	Misery            bool

	// Custom
	Custom Stats
}

func (b Buffs) Pack() []byte {
	var opt1 byte
	var opt2 byte
	if b.ArcaneInt {
		opt1 = opt1 | 1
	}
	if b.GiftOftheWild {
		opt1 = opt1 | 1<<1
	}
	if b.BlessingOfKings {
		opt1 = opt1 | 1<<2
	}
	if b.ImprovedBlessingOfWisdom {
		opt1 = opt1 | 1<<3
	}
	if b.ImprovedDivineSpirit {
		opt1 = opt1 | 1<<4
	}
	if b.Moonkin {
		opt1 = opt1 | 1<<5
	}
	if b.MoonkinRavenGoddess {
		opt1 = opt1 | 1<<6
	}
	if b.EyeOfNight {
		opt1 = opt1 | 1<<7
	}
	if b.TwilightOwl {
		opt2 = opt2 | 1
	}
	if b.WaterShield {
		opt2 = opt2 | 1<<1
	}
	if b.JudgementOfWisdom {
		opt2 = opt2 | 1<<2
	}
	if b.ImpSealofCrusader {
		opt2 = opt2 | 1<<3
	}
	if b.Misery {
		opt2 = opt2 | 1<<4
	}

	bytes := []byte{
		opt1, opt2, b.WaterShieldPPM,
		0, 0, // spriest dps
		byte(b.Race),
		0,
	}

	binary.LittleEndian.PutUint16(bytes[3:], b.SpriestDPS)

	var customBytes []byte
	for _, v := range b.Custom {
		if v != 0 {
			bytes[6] = byte(len(b.Custom))
			customBytes = make([]byte, len(b.Custom)*8)
			for i, rv := range b.Custom {
				binary.LittleEndian.PutUint64(customBytes[i*8:], math.Float64bits(rv))
			}
			break
		}
	}
	if len(customBytes) > 0 {
		bytes = append(bytes, customBytes...)
	}
	return bytes
}

type RaceBonusType byte

// These values are used directly in the dropdown in index.html
const (
	RaceBonusNone RaceBonusType = iota
	RaceBonusDraenei
	RaceBonusTroll10
	RaceBonusTroll30
	RaceBonusOrc
)

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
		if b.MoonkinRavenGoddess {
			s[StatSpellCrit] += 20
		}
	}
	if b.TwilightOwl {
		s[StatSpellCrit] += 44.16
	}
	if b.EyeOfNight {
		s[StatSpellDmg] += 34
	}
	if b.WaterShield {
		s[StatMP5] += 50
	}
	if b.Race == RaceBonusDraenei {
		s[StatSpellHit] += 15.76 // 1% hit
	}
	s[StatMP5] += float64(b.SpriestDPS) * 0.25

	if b.ImpSealofCrusader {
		s[StatSpellCrit] += 66.24 // 3% crit
	}

	for k, v := range b.Custom {
		s[k] += v
	}
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
	DestructionPotion bool
	SuperManaPotion   bool
	DarkRune          bool
}

func (c Consumes) Pack() []byte {
	var opt1 byte
	if c.BrilliantWizardOil {
		opt1 = opt1 | 1
	}
	if c.MajorMageblood {
		opt1 = opt1 | 1<<1
	}
	if c.FlaskOfBlindingLight {
		opt1 = opt1 | 1<<2
	}
	if c.FlaskOfMightyRestoration {
		opt1 = opt1 | 1<<3
	}
	if c.BlackendBasilisk {
		opt1 = opt1 | 1<<4
	}
	if c.DestructionPotion {
		opt1 = opt1 | 1<<5
	}
	if c.SuperManaPotion {
		opt1 = opt1 | 1<<6
	}
	if c.DarkRune {
		opt1 = opt1 | 1<<7
	}
	return []byte{opt1}
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
