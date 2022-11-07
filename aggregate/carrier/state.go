package carrier

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/common"
)

type PositionState struct {
	X float64
	Y float64
}

type status int

var (
	IdleStatus   status = 1
	MovingStatus status = 2
)

type JetIDs []string

type CombatState struct {
	ID     string
	Status status
	Jets   JetIDs
}

func toCombatState(events []Event) CombatState {
	state := CombatState{}

	for _, event := range events {
		switch event.Effect {
		case InitEffect:
			state.ID = event.ID
			state.Status = IdleStatus
			state.Jets = make(JetIDs, 0)

		case LaunchJetEffect:
			jetID, _ := event.Data.(string)
			state.ID = event.ID

			state.Jets = state.Jets.removeJetID(jetID)
		case HouseJetEffect:
			jetID, _ := event.Data.(string)
			state.ID = event.ID

			jetIDs, err := state.Jets.append(jetID)
			if err != nil {
				fmt.Printf("[ERROR][toCombatState] %s", err.Error())
			}
			state.Jets = jetIDs

		case MoveEffect:
			state.ID = event.ID
			state.Status = MovingStatus
		case IdleEffect:
			state.ID = event.ID
			state.Status = IdleStatus
		}
	}

	return state
}

func GetCombatState(store Store, id string) (CombatState, error) {
	events, err := store.getEventsByID(id)
	if err != nil {
		return CombatState{}, err
	}

	if len(events) == 0 {
		return CombatState{}, common.ErrEntityNotFound
	}

	return toCombatState(events), nil
}

func toPositionState(events []Event) PositionState {
	state := PositionState{}

	for _, event := range events {
		switch event.Effect {
		case MoveEffect:
			position, _ := event.Data.(PositionState)
			state.X = position.X
			state.Y = position.Y
		}
	}
	return state
}

func GetPositionState(store Store, id string) (PositionState, error) {
	events, err := store.getEventsByID(id)
	if err != nil {
		return PositionState{}, err
	}

	if len(events) == 0 {
		return PositionState{}, common.ErrEntityNotFound
	}

	return toPositionState(events), nil
}
