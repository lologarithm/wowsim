package tbc

import (
	"fmt"
)

/**
 * Controls the spell rotation, can be thought of as the 'player'.
 *
 * TODO: Decide what should count as a spell, so we can have agents do other things (drop
 * totems? cast lust? melee attacks?). One idea: anything on the GCD counts as a spell.
 */
type Agent interface {
	// Returns the spell this Agent would like to cast next. Should never return nil.
	ChooseSpell(*Simulation, bool) *Cast

	// This will be invoked if the chosen spell is actually cast, so the Agent can update its state.
	OnSpellAccepted(*Simulation, *Cast)

	// Returns this Agent to its initial state.
	Reset(*Simulation)
}

// ################################################################
//                          FIXED ROTATION
// ################################################################
type FixedRotationAgent struct {
	numLBsPerCL       int // If -1, uses LB only
	numLBsSinceLastCL int
}

func (agent *FixedRotationAgent) ChooseSpell(sim *Simulation, didPot bool) *Cast {
	if agent.numLBsPerCL == -1 {
		return NewCast(sim, spellmap[MagicIDLB12])
	}

	if sim.CDs[MagicIDCL6] < 1 && agent.numLBsSinceLastCL >= agent.numLBsPerCL {
		return NewCast(sim, spellmap[MagicIDCL6])
	} else {
		return NewCast(sim, spellmap[MagicIDLB12])
	}
}

func (agent *FixedRotationAgent) OnSpellAccepted(sim *Simulation, cast *Cast) {
	if cast.Spell.ID == MagicIDLB12 {
		agent.numLBsSinceLastCL++
	} else if cast.Spell.ID == MagicIDCL6 {
		agent.numLBsSinceLastCL = 0
	}
}

func (agent *FixedRotationAgent) Reset(sim *Simulation) {
	agent.numLBsSinceLastCL = agent.numLBsPerCL
}

func NewFixedRotationAgent(sim *Simulation, numLBsPerCL int) *FixedRotationAgent {
	return &FixedRotationAgent{
		numLBsPerCL:       numLBsPerCL,
		numLBsSinceLastCL: numLBsPerCL, // This lets us cast CL first
	}
}

// ################################################################
//                             ADAPTIVE
// ################################################################
type AdaptiveAgent struct {
	LastMana  float64
	LastCheck int
	NumCasts  int
}

func (agent *AdaptiveAgent) ChooseSpell(sim *Simulation, didPot bool) *Cast {
	if agent.LastMana == 0 {
		agent.LastMana = sim.CurrentMana
	}
	agent.NumCasts++
	if didPot {
		// Use Potion to reset the calculation...
		agent.LastMana = sim.CurrentMana
		agent.LastCheck = sim.CurrentTick
		agent.NumCasts = 0
	}
	// Always give a couple casts before we figure out mana drain.
	if sim.CDs[MagicIDCL6] < 1 && agent.NumCasts > 3 {
		manaDrained := agent.LastMana - sim.CurrentMana
		timePassed := sim.CurrentTick - agent.LastCheck
		if timePassed == 0 {
			timePassed = 1
		}
		rate := manaDrained / float64(timePassed)
		timeRemaining := sim.endTick - sim.CurrentTick
		totalManaDrain := rate * float64(timeRemaining)
		buffer := spellmap[MagicIDCL6].Mana // mana buffer of 1 extra CL

		if sim.Debug != nil {
			sim.Debug("[AI] CL Ready: Mana/Tick: %0.1f, Est Mana Drain: %0.1f, CurrentMana: %0.1f\n", rate, totalManaDrain, sim.CurrentMana)
		}
		// If we have enough mana to burn and CL is off CD, use it.
		if totalManaDrain < sim.CurrentMana-buffer {
			cast := NewCast(sim, spellmap[MagicIDCL6])
			if sim.CurrentMana >= cast.ManaCost {
				return cast
			}
		}
	}
	return NewCast(sim, spellmap[MagicIDLB12])
}
func (agent *AdaptiveAgent) OnSpellAccepted(sim *Simulation, cast *Cast) {
}

func (agent *AdaptiveAgent) Reset(sim *Simulation) {
	agent.LastMana = sim.CurrentMana
	agent.LastCheck = 0
	agent.NumCasts = 3
}

func NewAdaptiveAgent(sim *Simulation) *AdaptiveAgent {
	return &AdaptiveAgent{
		LastMana: sim.CurrentMana,
		NumCasts: 3, // This lets us cast CL first.
	}
}

type AgentType int

// This must be kept in sync with the enum in ui.js
const (
	AGENT_TYPE_FIXED_3LB_1CL AgentType = iota
	AGENT_TYPE_FIXED_4LB_1CL
	AGENT_TYPE_FIXED_5LB_1CL
	AGENT_TYPE_FIXED_6LB_1CL
	AGENT_TYPE_FIXED_7LB_1CL
	AGENT_TYPE_FIXED_8LB_1CL
	AGENT_TYPE_FIXED_9LB_1CL
	AGENT_TYPE_FIXED_10LB_1CL
	AGENT_TYPE_FIXED_LB_ONLY
	AGENT_TYPE_ADAPTIVE
)

var ALL_AGENT_TYPES = []AgentType{
	AGENT_TYPE_FIXED_3LB_1CL,
	AGENT_TYPE_FIXED_4LB_1CL,
	AGENT_TYPE_FIXED_5LB_1CL,
	AGENT_TYPE_FIXED_6LB_1CL,
	AGENT_TYPE_FIXED_7LB_1CL,
	AGENT_TYPE_FIXED_8LB_1CL,
	AGENT_TYPE_FIXED_9LB_1CL,
	AGENT_TYPE_FIXED_10LB_1CL,
	AGENT_TYPE_FIXED_LB_ONLY,
	AGENT_TYPE_ADAPTIVE,
}

func NewAgent(sim *Simulation, agentType AgentType) Agent {
	switch agentType {
	case AGENT_TYPE_FIXED_3LB_1CL:
		return NewFixedRotationAgent(sim, 3)
	case AGENT_TYPE_FIXED_4LB_1CL:
		return NewFixedRotationAgent(sim, 4)
	case AGENT_TYPE_FIXED_5LB_1CL:
		return NewFixedRotationAgent(sim, 5)
	case AGENT_TYPE_FIXED_6LB_1CL:
		return NewFixedRotationAgent(sim, 6)
	case AGENT_TYPE_FIXED_7LB_1CL:
		return NewFixedRotationAgent(sim, 7)
	case AGENT_TYPE_FIXED_8LB_1CL:
		return NewFixedRotationAgent(sim, 8)
	case AGENT_TYPE_FIXED_9LB_1CL:
		return NewFixedRotationAgent(sim, 9)
	case AGENT_TYPE_FIXED_10LB_1CL:
		return NewFixedRotationAgent(sim, 10)
	case AGENT_TYPE_FIXED_LB_ONLY:
		return NewFixedRotationAgent(sim, -1)
	case AGENT_TYPE_ADAPTIVE:
		return NewAdaptiveAgent(sim)
	default:
		fmt.Printf("[ERROR] No rotation given to sim.\n")
		return nil
	}
}
