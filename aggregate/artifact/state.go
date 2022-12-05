package artifact

import "github.com/R-jim/Momentum/aggregate/common"

type Type string
type Status string

var (
	MockType Type = "mock"

	PlantedStatus      Status = "planted"
	TransportingStatus Status = "transporting"
	HarvestedStatus    Status = "harvested"
)

func (t Type) IsValid() bool {
	return t == MockType
}

func (s Status) IsValid() bool {
	return s == PlantedStatus ||
		s == TransportingStatus ||
		s == HarvestedStatus
}

type PositionState struct {
	X float64
	Y float64
}

type State struct {
	ID     string
	Status Status
	Type   Type
	Energy int

	SpikeIDs []string
}

func toState(events []Event) State {
	state := State{}

	for _, event := range events {
		switch event.Effect {
		case InitEffect:
			artifactType, _ := event.Data.(Type)

			state.ID = event.ID
			state.Type = artifactType
			state.Status = PlantedStatus
			state.Energy = 100

		case SpawnSpikeEffect:
			state.Energy -= 10
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
