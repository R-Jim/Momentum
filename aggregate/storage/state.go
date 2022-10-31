package storage

import "github.com/R-jim/Momentum/aggregate/common"

type Type string

var (
	FuelType Type = "fuel"
)

func (t Type) IsValid() bool {
	return t == FuelType
}

type State struct {
	ID       string
	Type     Type
	Quantity int
}

func toState(events []Event) State {
	state := State{}

	for _, event := range events {
		switch event.Effect {
		case InitEffect:
			storageType, _ := event.Data.(Type)

			state.ID = event.ID
			state.Type = storageType
			state.Quantity = 0

		case RefillEffect:
			quantity, _ := event.Data.(int)

			state.ID = event.ID
			state.Quantity += quantity

		case ConsumeEffect:
			quantity, _ := event.Data.(int)

			state.ID = event.ID
			state.Quantity -= quantity
		}
	}

	return state
}

func GetState(store Store, id string) (State, error) {
	events, err := store.getEventsByID(id)
	if err != nil {
		return State{}, err
	}

	if len(events) == 0 {
		return State{}, common.ErrEntityNotFound
	}

	return toState(events), nil
}
