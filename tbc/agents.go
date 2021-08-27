package tbc

import (
	"fmt"
	"math"
)

/**
 * Controls the spell rotation, can be thought of as the 'player'.
 *
 * TODO: Decide what should count as a spell, so we can have agents do other things (drop
 * totems? cast lust? melee attacks?). One idea: anything on the GCD counts as a spell.
 */
type Agent interface {
	// Returns the action this Agent would like to take next.
	ChooseAction(*Simulation) AgentAction

	// This will be invoked if the chosen action is actually executed, so the Agent can update its state.
	OnActionAccepted(*Simulation, AgentAction)

	// Returns this Agent to its initial state.
	Reset(*Simulation)
}

// A single action that an Agent can take.
type AgentAction struct {
	// Exactly one of these should be set.
	Wait float64 // Duration to wait
	Cast *Cast
}

func NewWaitAction(duration float64) AgentAction {
	return AgentAction{
		Wait: duration,
	}
}

func NewCastAction(sim *Simulation, sp *Spell) AgentAction {
	return AgentAction{
		Cast: NewCast(sim, sp),
	}
}

// ################################################################
//                          FIXED ROTATION
// ################################################################
type FixedRotationAgent struct {
	numLBsPerCL       int // If -1, uses LB only
	numLBsSinceLastCL int
}

// Returns if any temporary haste buff is currently active.
// TODO: Figure out a way to make this automatic
func (agent *FixedRotationAgent) temporaryHasteActive(sim *Simulation) bool {
	return sim.hasAura(MagicIDBloodlust) ||
		sim.hasAura(MagicIDDrums) ||
		sim.hasAura(MagicIDTrollBerserking) ||
		sim.hasAura(MagicIDSkullGuldan) ||
		sim.hasAura(MagicIDFungalFrenzy)
}

func (agent *FixedRotationAgent) ChooseAction(sim *Simulation) AgentAction {
	if agent.numLBsPerCL == -1 {
		return NewCastAction(sim, spellmap[MagicIDLB12])
	}

	if agent.numLBsSinceLastCL < agent.numLBsPerCL {
		return NewCastAction(sim, spellmap[MagicIDLB12])
	}

	if !sim.isOnCD(MagicIDCL6) {
		return NewCastAction(sim, spellmap[MagicIDCL6])
	}

	// If we have a temporary haste effect (like bloodlust or quags eye) then
	// we should add LB casts instead of waiting
	if agent.temporaryHasteActive(sim) {
		return NewCastAction(sim, spellmap[MagicIDLB12])
	}

	return NewWaitAction(sim.getRemainingCD(MagicIDCL6))
}

func (agent *FixedRotationAgent) OnActionAccepted(sim *Simulation, action AgentAction) {
	if action.Cast == nil {
		return
	}

	if action.Cast.Spell.ID == MagicIDLB12 {
		agent.numLBsSinceLastCL++
	} else if action.Cast.Spell.ID == MagicIDCL6 {
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
	// Circular array buffer for recent mana snapshots, within a time window
	manaSnapshots      [manaSnapshotsBufferSize]ManaSnapshot
	numSnapshots       int32
	firstSnapshotIndex int32
}

const manaSpendingWindow = 60.0

// 2 * (# of seconds) should be plenty of slots
const manaSnapshotsBufferSize = int32(manaSpendingWindow) * 2

type ManaSnapshot struct {
	time      float64 // time this snapshot was taken
	manaSpent float64 // total amount of mana spent up to this time
}

func (agent *AdaptiveAgent) getOldestSnapshot() ManaSnapshot {
	return agent.manaSnapshots[agent.firstSnapshotIndex]
}

func (agent *AdaptiveAgent) purgeExpiredSnapshots(sim *Simulation) {
	expirationCutoff := sim.CurrentTime - manaSpendingWindow

	curIndex := agent.firstSnapshotIndex
	for agent.numSnapshots > 0 && agent.manaSnapshots[curIndex].time < expirationCutoff {
		curIndex = (curIndex + 1) % manaSnapshotsBufferSize
		agent.numSnapshots--
	}
	agent.firstSnapshotIndex = curIndex
}

func (agent *AdaptiveAgent) takeSnapshot(sim *Simulation) {
	snapshot := ManaSnapshot{
		time:      sim.CurrentTime,
		manaSpent: sim.metrics.ManaSpent,
	}

	nextIndex := (agent.firstSnapshotIndex + agent.numSnapshots) % manaSnapshotsBufferSize
	agent.manaSnapshots[nextIndex] = snapshot
	agent.numSnapshots++
}

func (agent *AdaptiveAgent) ChooseAction(sim *Simulation) AgentAction {
	if sim.isOnCD(MagicIDCL6) {
		return NewCastAction(sim, spellmap[MagicIDLB12])
	}

	agent.purgeExpiredSnapshots(sim)
	oldestSnapshot := agent.getOldestSnapshot()

	manaSpent := sim.metrics.ManaSpent - oldestSnapshot.manaSpent
	timeDelta := sim.CurrentTime - oldestSnapshot.time
	manaSpendingRate := manaSpent / math.Max(1.0, timeDelta)

	timeRemaining := sim.Options.Encounter.Duration - sim.CurrentTime
	projectedManaCost := manaSpendingRate * timeRemaining
	buffer := spellmap[MagicIDCL6].Mana // mana buffer of 1 extra CL

	if sim.Debug != nil {
		sim.Debug("[AI] CL Ready: Mana/Tick: %0.1f, Est Mana Cost: %0.1f, CurrentMana: %0.1f\n", manaSpendingRate, projectedManaCost, sim.CurrentMana)
	}

	// If we have enough mana to burn and CL is off CD, use it.
	if projectedManaCost < sim.CurrentMana-buffer {
		return NewCastAction(sim, spellmap[MagicIDCL6])
	}

	return NewCastAction(sim, spellmap[MagicIDLB12])
}
func (agent *AdaptiveAgent) OnActionAccepted(sim *Simulation, action AgentAction) {
	agent.takeSnapshot(sim)
}

func (agent *AdaptiveAgent) Reset(sim *Simulation) {
	agent.manaSnapshots = [manaSnapshotsBufferSize]ManaSnapshot{}
	agent.firstSnapshotIndex = 0
	agent.numSnapshots = 0
}

func NewAdaptiveAgent(sim *Simulation) *AdaptiveAgent {
	return &AdaptiveAgent{}
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
