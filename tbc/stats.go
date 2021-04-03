package tbc

import (
	"strconv"
)

var TicksPerSecond = 30

type Stats []float64

type Stat int

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
	}

	return "none"
}

func (st Stats) Clone() Stats {
	ns := make(Stats, StatLen)
	for i, v := range st {
		ns[i] = v
	}
	return ns
}

func (st Stats) Print() string {
	output := "{ "
	for k, v := range st {
		name := Stat(k).StatName()
		if name == "none" {
			continue
		}
		if v < 50 {
			output += "\t\"" + name + "\": " + strconv.FormatFloat(v, 'f', 3, 64)
		} else {
			output += "\t\"" + name + "\": " + strconv.FormatFloat(v, 'f', 0, 64)
		}
		if k != len(st)-1 {
			output += ",\n"
		}
	}
	output += " }"
	return output
}
