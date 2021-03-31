package elesim

import (
	"fmt"
	"math/rand"
	"strings"
)

var IsDebug = false

func debug(s string, vals ...interface{}) {
	if IsDebug {
		fmt.Printf(s, vals...)
	}
}

func Sim(seconds int, stats Stats, spellOrder []string) (float64, int) {
	ticks := seconds * tickPerSecond

	// ticks until cast is complete
	currentMana := stats[StatMana]
	casting := 0
	// timeToRegen := 0
	castingSpell := ""

	spellIdx := 0
	cds := map[string]int{}
	auras := map[string]Aura{}

	totalDmg := 0.0
	castStats := map[string]struct {
		Num   int
		Total float64
	}{}
	oomAt := 0

	for i := 0; i < ticks; i++ {
		casting--

		// MP5 regen
		currentMana += (stats[StatMP5] / 5.0) / float64(tickPerSecond)
		if currentMana > stats[StatMana] {
			currentMana = stats[StatMana]
		}

		if stats[StatMana]-currentMana >= 2250 && cds["potion"] < 1 {
			// Restores 1350 to 2250 mana. (2 Min Cooldown)
			currentMana += float64(1350 + rand.Intn(2250-1350))
			cds["potion"] = 120 * tickPerSecond
			debug("[%d] Used Mana Potion\n", i/tickPerSecond)
		}
		if stats[StatMana]-currentMana >= 1500 && cds["darkrune"] < 1 {
			// Restores 900 to 1500 mana. (2 Min Cooldown)
			currentMana += float64(900 + rand.Intn(1500-900))
			cds["darkrune"] = 120 * tickPerSecond
			debug("[%d] Used Mana Potion\n", i/tickPerSecond)
		}

		if castingSpell == "" && oomAt == 0 {
			// debug("(%0.0f/%0.0f mana)\n", currentMana, stats[StatMana])
		}

		if casting <= 0 {
			if castingSpell != "" {
				sp := spellmap[castingSpell]
				cast := &Cast{
					Spell:      sp,
					ManaCost:   float64(sp.Mana),
					Spellpower: stats[StatSpellDmg], // TODO: type specific bonuses...
				}
				if strings.HasPrefix(sp.ID, "LB") || strings.HasPrefix(sp.ID, "CL") {
					// totem of the storm
					cast.Spellpower += 33
					// Talent Convection
					cast.ManaCost *= 0.9
				}

				cast.Hit = 0.83 + stats[StatSpellHit]
				cast.Crit = stats[StatSpellCrit]

				for _, aur := range auras {
					if aur.OnCast != nil {
						aur.OnCast(cast)
					}
				}

				// TODO: generalize aura removal.
				if _, ok := auras["elemastery"]; ok {
					delete(auras, "elemastery")
				} else if _, ok := auras["cc"]; ok {
					delete(auras, "cc")
				}

				if rand.Float64() < cast.Hit {
					dmg := (float64(rand.Intn(int(sp.MaxDmg-sp.MinDmg))) + sp.MinDmg) + (stats[StatSpellDmg] * sp.Coeff)

					if rand.Float64() < cast.Crit {
						dmg *= 2
						debug("crit")
					} else {
						debug("hit")
					}

					if strings.HasPrefix(sp.ID, "LB") || strings.HasPrefix(sp.ID, "CL") {
						// Talent Concussion
						dmg *= 1.05
					}

					// Average Resistance = (Target's Resistance / (Caster's Level * 5)) * 0.75 "AR"
					// P(x) = 50% - 250%*|x - AR| <- where X is chance of resist
					// For now hardcode the 25% chance resist at 2.5% (this assumes bosses have 0 nature resist)
					if rand.Float64() < 0.025 { // chance of 25% resist
						dmg *= .75
						debug("(partial resist)")
					}
					debug(": %0.0f\n", dmg)

					totalDmg += dmg
					stat := castStats[sp.ID]
					stat.Num += 1
					stat.Total += dmg
					castStats[sp.ID] = stat
				} else {
					debug("miss.\n")
				}

				currentMana -= cast.ManaCost

				castingSpell = ""
				if sp.Cooldown > 0 {
					cds[sp.ID] = sp.Cooldown * tickPerSecond
				}
				if rand.Float64() < 0.1 {
					// TODO: make Elemental Focus an aura that applies an aura.
					// talent for clearcast chance on cast
					auras["cc"] = clearcasting()
					debug("\tGained Clearcasting.\n")
				}
				if rand.Float64() < 0.2 {
					auras["stc"] = stormcaller()
					debug("\tGained Stormcaller.\n")
				}
				continue
			} else {
				// Choose next spell
				so := spellOrder[spellIdx]

				isclearcasting := auras["cc"].Duration > 0 || auras["elemastery"].Duration > 0

				// anytime ZHC AND Ele Mastery are up, pop!
				if cds["zhc"] <= 0 && cds["elemastery"] <= 0 {
					// Apply auras
					auras["zhc"] = zhc()
					auras["elemastery"] = elemastery()

					cds["zhc"] = 120 * tickPerSecond
					cds["elemastery"] = 180 * tickPerSecond
				}

				sp := spellmap[so]
				cost := sp.Mana
				if strings.HasPrefix(sp.ID, "LB") || strings.HasPrefix(sp.ID, "CL") {
					// Talent to reduce mana cost by 10%
					cost *= 0.9
				}
				if cds[so] == 0 && (currentMana >= sp.Mana || isclearcasting) {
					castTime := spellmap[so].CastTime
					if strings.HasPrefix(sp.ID, "LB") || strings.HasPrefix(sp.ID, "CL") {
						// Talent to reduce cast time.
						castTime -= 1
					}
					casting = int(castTime * float64(tickPerSecond))
					castingSpell = sp.ID
					spellIdx++
					if spellIdx == len(spellOrder) {
						spellIdx = 0
					}
					debug("[%d] Casting %s ...", i/tickPerSecond, sp.ID)
				} else if !isclearcasting && currentMana < 200 && oomAt == 0 {
					oomAt = i / tickPerSecond
				}

			}
		}
		// CDS
		for k := range cds {
			cds[k]--
			if cds[k] <= 0 {
				delete(cds, k)
			}
		}
		for k, v := range auras {
			nv := v
			nv.Duration--
			if nv.Duration > 0 {
				auras[k] = nv
			} else {
				delete(auras, k)
			}
		}
	}

	return totalDmg, oomAt
}
