package tbc

import (
	"fmt"
	"math"
	"time"
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
	Wait time.Duration // Duration to wait
	Cast *Cast
}

func NewWaitAction(duration time.Duration) AgentAction {
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
	lb                *Spell
	cl                *Spell
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
		return NewCastAction(sim, agent.lb)
	}

	if agent.numLBsSinceLastCL < agent.numLBsPerCL {
		return NewCastAction(sim, agent.lb)
	}

	if !sim.isOnCD(MagicIDCL6) {
		return NewCastAction(sim, agent.cl)
	}

	// If we have a temporary haste effect (like bloodlust or quags eye) then
	// we should add LB casts instead of waiting
	if agent.temporaryHasteActive(sim) {
		return NewCastAction(sim, agent.lb)
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
		lb:                spellmap[MagicIDLB12],
		cl:                spellmap[MagicIDCL6],
	}
}

// ################################################################
//                          CL ON CLEARCAST
// ################################################################
type CLOnClearcastAgent struct {
	// Whether the last 2 spells procced clearcasting, either directly
	// or via lightning overload.
	prevCastProccedCC     bool
	prevPrevCastProccedCC bool

	lb *Spell
	cl *Spell
}

func (agent *CLOnClearcastAgent) ChooseAction(sim *Simulation) AgentAction {
	if sim.isOnCD(MagicIDCL6) || !agent.prevPrevCastProccedCC {
		return NewCastAction(sim, agent.lb)
	}

	return NewCastAction(sim, agent.cl)
}

func (agent *CLOnClearcastAgent) OnActionAccepted(sim *Simulation, action AgentAction) {
	agent.prevPrevCastProccedCC = agent.prevCastProccedCC
	agent.prevCastProccedCC = false
}

func (agent *CLOnClearcastAgent) Reset(sim *Simulation) {
	// Checking whether the EleFocus aura is active isn't enough; we need to know
	// how many charges it currently has. This information isn't available so instead
	// we infer by checking whether each spell is a crit.
	sim.addAura(Aura{
		ID:      MagicIDClearcastAgent,
		Expires: neverExpires,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			agent.prevCastProccedCC = agent.prevCastProccedCC ||
				(c.DidCrit && c.Spell.ID != MagicIDTLCLB)
		},
	})

	agent.prevCastProccedCC = false
	agent.prevPrevCastProccedCC = false
}

func NewCLOnClearcastAgent(sim *Simulation) *CLOnClearcastAgent {
	return &CLOnClearcastAgent{
		lb: spellmap[MagicIDLB12],
		cl: spellmap[MagicIDCL6],
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
	lb                 *Spell
	cl                 *Spell
}

const manaSpendingWindowNumSeconds = 60
const manaSpendingWindow = time.Second * manaSpendingWindowNumSeconds

// 2 * (# of seconds) should be plenty of slots
const manaSnapshotsBufferSize = manaSpendingWindowNumSeconds * 2

type ManaSnapshot struct {
	time      time.Duration // time this snapshot was taken
	manaSpent float64       // total amount of mana spent up to this time
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
	if agent.numSnapshots >= manaSnapshotsBufferSize {
		panic("Agent snapshot buffer full")
	}

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
		return NewCastAction(sim, agent.lb)
	}

	agent.purgeExpiredSnapshots(sim)
	oldestSnapshot := agent.getOldestSnapshot()

	manaSpent := sim.metrics.ManaSpent - oldestSnapshot.manaSpent
	timeDelta := sim.CurrentTime - oldestSnapshot.time
	manaSpendingRate := manaSpent / math.Max(1.0, timeDelta.Seconds())

	timeRemaining := sim.Duration - sim.CurrentTime
	projectedManaCost := manaSpendingRate * timeRemaining.Seconds()

	if sim.Debug != nil {
		sim.Debug("[AI] CL Ready: Mana/Tick: %0.1f, Est Mana Cost: %0.1f, CurrentMana: %0.1f\n", manaSpendingRate, projectedManaCost, sim.CurrentMana)
	}

	// If we have enough mana to burn and CL is off CD, use it.
	if projectedManaCost < sim.CurrentMana {
		return NewCastAction(sim, agent.cl)
	}

	return NewCastAction(sim, agent.lb)
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
	return &AdaptiveAgent{
		lb: spellmap[MagicIDLB12],
		cl: spellmap[MagicIDCL6],
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
	AGENT_TYPE_CL_ON_CLEARCAST
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
	AGENT_TYPE_CL_ON_CLEARCAST,
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
	case AGENT_TYPE_CL_ON_CLEARCAST:
		return NewCLOnClearcastAgent(sim)
	default:
		fmt.Printf("[ERROR] No rotation given to sim.\n")
		return nil
	}
}
