package elesim

import "fmt"

var tickPerSecond = 30

type Stats []float64

type Stat int

const (
	StatInt Stat = iota
	StatSpellCrit
	StatSpellHit
	StatSpellDmg
	StatMP5
	StatMana
	StatSpellPen
)

func (s Stat) StatName() string {
	switch s {
	case StatInt:
		return "StatInt"
	case StatSpellCrit:
		return "StatSpellCrit"
	case StatSpellHit:
		return "StatSpellHit"
	case StatSpellDmg:
		return "StatSpellDmg"
	case StatMP5:
		return "StatMP5"
	case StatMana:
		return "StatMana"
	case StatSpellPen:
		return "StatSpellPen"
	}

	return "none"
}

func (st Stats) Print() {
	fmt.Printf("Stats:\n")

	for k, v := range st {
		if v < 50 {
			fmt.Printf("\t%s: %0.3f\n", Stat(k).StatName(), v)
		} else {
			fmt.Printf("\t%s: %0.0f\n", Stat(k).StatName(), v)
		}

	}
}
