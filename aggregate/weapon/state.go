package weapon

import "github.com/R-jim/Momentum/aggregate/common"

type Type string

var (
	GunType Type = "gun"
)

func (t Type) IsValid() bool {
	return t == GunType
}

type State struct {
	ID   string
	Type Type
}

func toState(events []Event) State {
	state := State{}

	for _, event := range events {
		switch event.Effect {
		case InitEffect:
			artifactType, _ := event.Data.(Type)

			state.ID = event.ID
			state.Type = artifactType
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
