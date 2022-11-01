package jet

import "github.com/R-jim/Momentum/aggregate/common"

type PositionState struct {
	X int
	Y int
}

type InventoryState struct {
	FuelTankID string
}

type Status int

var (
	IdleStatus     Status = 1
	LandedStatus   Status = 2
	FlyingStatus   Status = 3
	EngagingStatus Status = 4
)

type CombatState struct {
	ID       string
	Status   Status
	TargetID string
}

func toCombatState(events []Event) CombatState {
	state := CombatState{}

	for _, event := range events {
		switch event.Effect {
		case InitEffect:
			state.ID = event.ID
			state.Status = IdleStatus

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
		case InitEffect:
			inventory, _ := event.Data.(InventoryState)
			state.FuelTankID = inventory.FuelTankID
		case FuelTankChangedEffect:
			fuelTankID, _ := event.Data.(string)
			state.FuelTankID = fuelTankID
		}
	}
	return state
}

func GetInventoryState(store Store, id string) (InventoryState, error) {
	events, err := store.getEventsByID(id)
	if err != nil {
		return InventoryState{}, err
	}

	if len(events) == 0 {
		return InventoryState{}, common.ErrEntityNotFound
	}

	return toInventoryState(events), nil
}
