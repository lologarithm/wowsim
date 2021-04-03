package tbc

// Old fully ticking style simulator.
func (sim *Simulation) Run1(seconds int) SimMetrics {
	sim.reset()

	ticks := seconds * TicksPerSecond
	for i := 0; i < ticks; i++ {
		sim.currentTick = i
		sim.Tick1(i)
	}
	debug("(%0.0f/%0.0f mana)\n", sim.CurrentMana, sim.Stats[StatMana])
	return sim.metrics
}

func (sim *Simulation) Tick1(tickID int) {
	if sim.CurrentMana < 0 {
		panic("you should never have negative mana.")
	}

	secondID := tickID / TicksPerSecond
	// MP5 regen
	sim.CurrentMana += ((sim.Stats[StatMP5] + sim.Buffs[StatMP5]) / 5.0) / float64(TicksPerSecond)

	if sim.CurrentMana > sim.Stats[StatMana] {
		sim.CurrentMana = sim.Stats[StatMana]
	}

	if sim.CastingSpell != nil {
		sim.CastingSpell.TicksUntilCast-- // advance state of current cast.
		if sim.CastingSpell.TicksUntilCast == 0 {
			sim.Cast(sim.CastingSpell)
		}
	}

	if sim.CastingSpell == nil {
		// Pop potion before next cast.
		if sim.Stats[StatMana]-sim.CurrentMana+sim.Stats[StatMP5] >= 1500 && sim.CDs["darkrune"] < 1 {
			// Restores 900 to 1500 mana. (2 Min Cooldown)
			sim.CurrentMana += float64(900 + sim.rando.Intn(1500-900))
			sim.CDs["darkrune"] = 120 * TicksPerSecond
			debug("[%d] Used Mana Potion\n", secondID)
		}
		if sim.Stats[StatMana]-sim.CurrentMana+sim.Stats[StatMP5] >= 3000 && sim.CDs["potion"] < 1 {
			// Restores 1800 to 3000 mana. (2 Min Cooldown)
			sim.CurrentMana += float64(1800 + sim.rando.Intn(3000-1800))
			sim.CDs["potion"] = 120 * TicksPerSecond
			debug("[%d] Used Mana Potion\n", secondID)
		}
		// Pop any on-use trinkets

		for _, item := range sim.Equip {
			if item.Activate == nil || item.ActivateCD == -1 { // ignore non-activatable, and always active items.
				continue
			}
			if sim.CDs[item.CoolID] > 0 {
				continue
			}
			sim.addAura(item.Activate(sim))
			sim.CDs[item.CoolID] = item.ActivateCD * TicksPerSecond
		}

		// Choose next spell
		sim.ChooseSpell()
		if sim.CastingSpell != nil {
			debug("[%d] Casting %s (%0.1f) ...", secondID, sim.CastingSpell.Spell.ID, float64(sim.CastingSpell.TicksUntilCast)/float64(TicksPerSecond))
		}
	}

	// CDS
	for k := range sim.CDs {
		sim.CDs[k]--
		if sim.CDs[k] <= 0 {
			delete(sim.CDs, k)
		}
	}

	todel := []int{}
	for i := range sim.Auras {
		if sim.Auras[i].Expires <= tickID {
			todel = append(todel, i)
		}
	}
	for i := len(todel) - 1; i >= 0; i-- {
		sim.cleanAura(todel[i])
	}
}
