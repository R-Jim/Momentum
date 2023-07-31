package aggregator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/google/uuid"
)

const (
	SampleEffect event.Effect = "SAMPLE_EFFECT"
)

type SampleState struct {
	ID uuid.UUID
}

func NewSampleAggregator(store *event.SampleStore) Aggregator {
	s := event.Store(*store)

	return aggregateImpl{
		name:  "SAMPLE",
		store: &s,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{
			//"SAMPLE_EFFECTED"
			SampleEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetSampleState(append(currentEvents, inputEvent))
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}
				return nil
			},
		},
	}
}

func GetSampleState(events []event.Event) (SampleState, error) {
	state := SampleState{}

	return composeState(state, events, func(ss SampleState, e event.Event) (SampleState, error) {
		switch e.Effect {
		case SampleEffect:
			state.ID = e.ID
		}
		return ss, nil
	})
}
