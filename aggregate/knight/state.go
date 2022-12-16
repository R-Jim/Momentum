package knight

import "github.com/R-jim/Momentum/aggregate/common"

type PositionState struct {
	X float64
	Y float64
}

type Health struct {
	Max   int
	Value int
}

type TargetType string

const (
	WaitPointTargetType TargetType = "wait_point"
	SpikeTargetType     TargetType = "spike"
)

type Target struct {
	ID       string
	Type     TargetType
	Position PositionState
}

type State struct {
	ID                  string
	Health              Health
	WeaponID            string
	Target              Target
	HarvestedArtifactID string
}

func toState(events []Event) State {
	state := State{}

	for _, event := range events {
		switch event.Effect {
		case InitEffect:
			health, _ := event.Data.(Health)

			state.ID = event.ID
			state.Health = health
		case DamageEffect:
			damage, _ := event.Data.(int)

			state.ID = event.ID
			state.Health.Value -= damage

		case ChangeTargetEffect:
			target, _ := event.Data.(Target)
			state.Target = target

		case ChangeWeaponEffect:
		case GatherArtifactEffect:
			artifactID, _ := event.Data.(string)
			state.HarvestedArtifactID = artifactID
		case DropArtifactEFfect:
			state.HarvestedArtifactID = ""
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
			state = position
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
