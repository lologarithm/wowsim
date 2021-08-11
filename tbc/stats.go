package tbc

import (
	"strconv"
)

type Stats [StatLen]float64

type Stat byte

const (
	StatInt Stat = iota
	StatStm
	StatSpellCrit
	StatSpellHit
	StatSpellDmg
	StatHaste
	StatMP5
	StatMana
	StatSpellPen
	StatSpirit

	StatLen
)

func (s Stat) StatName() string {
	switch s {
	case StatInt:
		return "StatInt"
	case StatStm:
		return "StatStm"
	case StatSpellCrit:
		return "StatSpellCrit"
	case StatSpellHit:
		return "StatSpellHit"
	case StatSpellDmg:
		return "StatSpellDmg"
	case StatHaste:
		return "StatHaste"
	case StatMP5:
		return "StatMP5"
	case StatMana:
		return "StatMana"
	case StatSpellPen:
		return "StatSpellPen"
	case StatSpirit:
		return "StatSpirit"
	}

	return "none"
}

func (st Stats) Clone() Stats {
	ns := Stats{}
	for i, v := range st {
		ns[i] = v
	}
	return ns
}

// func (st Stats) MarshalJSON() ([]byte, error) {
// 	output := bytes.Buffer{}
// 	output.WriteByte('{')
// 	printed := false
// 	for k, v := range st {
// 		name := Stat(k).StatName()
// 		if name == "none" {
// 			continue
// 		}
// 		if printed {
// 			printed = false
// 			output.WriteByte(',')
// 		}
// 		if v < 50 {
// 			printed = true
// 			output.WriteString("\"" + name + "\": " + strconv.FormatFloat(v, 'f', 3, 64))
// 		} else {
// 			printed = true
// 			output.WriteString("\"" + name + "\": " + strconv.FormatFloat(v, 'f', 0, 64))
// 		}
// 	}
// 	output.WriteByte('}')
// 	return output.Bytes(), nil
// }

// Print is debug print function
func (st Stats) Print() string {
	output := "{ "
	printed := false
	for k, v := range st {
		name := Stat(k).StatName()
		if name == "none" {
			continue
		}
		if printed {
			printed = false
			output += ",\n"
		}
		output += "\t"
		if v < 50 {
			printed = true
			output += "\"" + name + "\": " + strconv.FormatFloat(v, 'f', 3, 64)
		} else {
			printed = true
			output += "\"" + name + "\": " + strconv.FormatFloat(v, 'f', 0, 64)
		}
	}
	output += " }"
	return output
}

// CalculatedTotal will add Mana and Crit from Int and return the new stats.
func (s Stats) CalculatedTotal() Stats {
	stats := s

	// Add crit/mana from int
	stats[StatSpellCrit] += (stats[StatInt] / 80) * 22.08
	stats[StatMana] += stats[StatInt] * 15
	return stats
}

// CalculateTotalStats will take a set of equipment and options and add all stats/buffs/etc together
func CalculateTotalStats(o Options, e Equipment) Stats {
	gearStats := e.Stats()
	stats := BaseStats(o.Buffs.Race)
	for i := range stats {
		stats[i] += gearStats[i]
	}

	stats = o.Talents.AddStats(o.Buffs.AddStats(o.Consumes.AddStats(o.Totems.AddStats(stats))))

	if o.Buffs.BlessingOfKings {
		stats[StatInt] *= 1.1 // blessing of kings
	}
	if o.Buffs.ImprovedDivineSpirit {
		stats[StatSpellDmg] += stats[StatSpirit] * 0.1
	}

	stats = stats.CalculatedTotal()

	// Add stat increases from talents
	stats[StatMP5] += stats[StatInt] * (0.02 * float64(o.Talents.UnrelentingStorm))

	return stats
}

func BaseStats(race RaceBonusType) Stats {
	stats := Stats{
		StatInt:       104,    // Base int for troll,
		StatMana:      2678,   // level 70 shaman
		StatSpirit:    135,    // lvl 70 shaman
		StatSpellCrit: 48.576, // base crit for 70 sham
	}
	// TODO: Find race int differences.
	switch race {
	case RaceBonusOrc:

	}
	return stats
}
