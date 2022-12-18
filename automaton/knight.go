package automaton

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/knight"
	"github.com/R-jim/Momentum/aggregate/spike"
	"github.com/R-jim/Momentum/info"
	"github.com/R-jim/Momentum/operator"
	"github.com/R-jim/Momentum/util"
)

/*
TODO:
knight will follow carrier by default
when on patrol mode with no enemy, moving in circular motion start from next closest position to carrier
when on patrol mode with enemy, move in between enemy and carrier
auto attack first enemy come into weapon range
*/

type KnightAutomaton interface {
	Auto(id string) error
}

type knightImpl struct {
	knightStore knight.Store
	spikeStore  spike.Store

	operator operator.Operator
}

func NewKnightAutomaton(knightStore knight.Store, spikeStore spike.Store, operator operator.Operator) KnightAutomaton {
	return knightImpl{
		knightStore: knightStore,
		spikeStore:  spikeStore,

		operator: operator,
	}
}

func (i knightImpl) Auto(id string) error {
	err := i.autoStrike(id)
	if err != nil {
		return err
	}
	return i.autoPatrol(id)
}

func (i knightImpl) autoPatrol(id string) error {
	radius := float64(4)

	combatState, err := knight.GetState(i.knightStore, id)
	if err != nil {
		return err
	}

	if combatState.Health.Value <= 0 || combatState.HarvestedArtifactID != "" {
		return nil
	}

	positionState, err := knight.GetPositionState(i.knightStore, id)
	if err != nil {
		return err
	}

	// set self wait point if no target
	if combatState.Target.ID == "" && combatState.Target.Type == "" {
		fmt.Printf("[AUTO_PATROL] change target")
		return i.operator.Knight.ChangeTarget(id, knight.Target{
			Type: knight.WaitPointTargetType,
			Position: knight.PositionState{
				X: positionState.X,
				Y: positionState.Y,
			},
		})
	}

	x, y := getNextStepXY(positionState.X, positionState.Y, 0, combatState.Target.Position.X, combatState.Target.Position.Y, radius, 2, float64(180))

	fmt.Printf("[AUTO_PATROL] %v:%v\n", x, y)
	return i.operator.Knight.Move(id, knight.PositionState{
		X: x,
		Y: y,
	})
}

func (i knightImpl) autoStrike(id string) error {
	combatState, err := knight.GetState(i.knightStore, id)
	if err != nil {
		return err
	}
	if combatState.Health.Value <= 0 {
		return nil
	}

	positionState, err := knight.GetPositionState(i.knightStore, id)
	if err != nil {
		return err
	}

	targetSpikeIDs := []string{}

	for _, spikeID := range i.spikeStore.GetEntityIDs() {
		spikeCombatState, err := spike.GetState(i.spikeStore, spikeID)
		if err != nil {
			return err
		}
		if spikeCombatState.Health.Value <= 0 {
			continue
		}

		spikePositionState, err := spike.GetPositionState(i.spikeStore, spikeID)
		if err != nil {
			return err
		}
		_, _, distance := util.GetDistances(positionState.X, positionState.Y, spikePositionState.X, spikePositionState.Y)
		if distance <= info.MELEE_RANGE {
			targetSpikeIDs = append(targetSpikeIDs, spikeID)
		}
	}

	for _, targetSpikeID := range targetSpikeIDs {
		err = i.operator.Knight.StrikeSpike(id, targetSpikeID)
		if err != nil {
			return err
		}
	}
	return nil
}
