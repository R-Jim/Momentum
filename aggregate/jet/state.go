package jet

import (
	"math"

	"github.com/R-jim/Momentum/aggregate/common"
)

type PositionState struct {
	X          float64
	Y          float64
	HeadPivotX float64
	HeadPivotY float64
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

type Target struct {
	ID   string
	Type string
}

type CombatState struct {
	ID     string
	Status Status
	Target *Target
}

func toCombatState(events []Event) CombatState {
	state := CombatState{}

	for _, event := range events {
		switch event.Effect {
		case InitEffect:
			state.ID = event.ID
			state.Status = IdleStatus

		case AttackEffect:
			target, _ := event.Data.(Target)
			state.ID = event.ID
			state.Target = &target
		case CancelAttackEffect:
			state.ID = event.ID
			state.Target = nil

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

	var lastX, lastY float64
	for _, event := range events {
		switch event.Effect {
		case FlyEffect:
			position, _ := event.Data.(PositionState)
			state.X = position.X
			state.Y = position.Y

			state.HeadPivotX, state.HeadPivotY = getSteps(lastX, lastY, state.X, state.Y)

			lastX = state.X
			lastY = state.Y
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

func getDistances(startX, startY, desX, desY float64) (distanceX, distanceY, distanceSqrt float64) {
	distanceX = math.Abs(desX - startX)
	distanceY = math.Abs(desY - startY)

	if distanceX == 0 && distanceY == 0 {
		return 0, 0, 0
	} else if distanceX == 0 {
		distanceSqrt = distanceY
	} else if distanceY == 0 {
		distanceSqrt = distanceX
	} else {
		distanceSqrt = math.RoundToEven(math.Sqrt(math.Pow(distanceX, 2)+math.Pow(distanceX, 2))*100) / 100
	}

	return distanceX, distanceY, distanceSqrt
}

func getSteps(startX, startY, desX, desY float64) (stepX, stepY float64) {
	distanceX, distanceY, distanceSqrt := getDistances(startX, startY, desX, desY)
	if distanceX == 0 && distanceY == 0 && distanceSqrt == 0 {
		return 0.5, 0.5
	}

	stepX = math.RoundToEven(distanceX/distanceSqrt*100) / 100
	stepY = math.RoundToEven(distanceY/distanceSqrt*100) / 100
	return stepX, stepY
}
