package jet

import "github.com/R-jim/Momentum/aggregate/common"

type PositionState struct {
	X int
	Y int
}

type InventoryState struct {
	FuelTankID string
}

type status int

var (
	LandedStatus   status = 1
	FlyingStatus   status = 2
	EngagingStatus status = 3
)

type CombatState struct {
	ID string
	// IsAttacking bool
	// IsLanded    bool
	Status   status
	TargetID string
}

func toCombatState(events []Event) CombatState {
	state := CombatState{}

	for _, event := range events {
		switch event.Effect {
		case InitEffect:
			state.ID = event.ID
			state.Status = LandedStatus

		case AttackEffect:
			targetID, _ := event.Data.(string)
			state.ID = event.ID
			state.TargetID = targetID
		case CancelAttackEffect:
			state.ID = event.ID
			state.TargetID = ""

		case EngageEffect:
			state.ID = event.ID
			state.Status = EngagingStatus
		case DisengageEffect:
			state.ID = event.ID
			state.Status = FlyingStatus

		case LandingEffect:
			state.ID = event.ID
			state.Status = LandedStatus
		case TakeOffEffect:
			state.ID = event.ID
			state.Status = FlyingStatus
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
		case FlyEffect:
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

func toInventoryState(events []Event) InventoryState {
	state := InventoryState{}

	for _, event := range events {
		switch event.Effect {
		case FlyEffect:
			inventory, _ := event.Data.(InventoryState)
			state.FuelTankID = inventory.FuelTankID
		}
	}
	return state
}
