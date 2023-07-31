package aggregator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/google/uuid"
)

type StageDescriptionState struct {
	Name        string
	Description string
}

type StageState struct {
	ID uuid.UUID

	GoalProductIDs []uuid.UUID
}

func NewStageAggregator(store *event.StageStore) Aggregator {
	s := event.Store(*store)

	return aggregateImpl{
		name:  "STAGE",
		store: &s,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{
			//"STAGE_INIT"
			event.StageInitEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				if len(currentEvents) != 0 {
					return ErrAggregateFail
				}

				if inputEvent.EntityID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				_, err := event.ParseData[[]uuid.UUID](inputEvent)
				if err != nil {
					return ErrAggregateFail
				}

				return nil
			},
		},
	}
}

func GetStageState(events []event.Event) (StageState, error) {
	return composeState(StageState{}, events, func(state StageState, e event.Event) (StageState, error) {
		switch e.Effect {
		case event.StageInitEffect:
			state.ID = e.EntityID
			value, err := event.ParseData[[]uuid.UUID](e)
			if err != nil {
				return state, err
			}

			state.GoalProductIDs = value
		}
		return state, nil
	})
}
