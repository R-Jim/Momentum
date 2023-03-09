package aggregator

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/google/uuid"
)

const (
	SampleEffect event.Effect = "SAMPLE_EFFECTED"
)

type SampleState struct {
	ID uuid.UUID
}

type sampleAggregateImpl struct {
	store        *store.Store
	aggregateSet map[event.Effect]func([]event.Event, event.Event) error
}

func NewSampleAggregator(store *store.Store) Aggregator {
	return sampleAggregateImpl{
		store: store,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{
			//"SAMPLE_EFFECTED"
			SampleEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state := GetSampleState(append(currentEvents, inputEvent))
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}
				return nil
			},
		},
	}
}

func (i sampleAggregateImpl) GetStore() *store.Store {
	return i.store
}

func (i sampleAggregateImpl) Aggregate(event event.Event) error {
	if err := aggregate(i.store, i.aggregateSet, event); err != nil {
		return fmt.Errorf("[SAMPLE_AGGREGATE][%v] %v", event.Effect, err)
	}
	fmt.Printf("[SAMPLE_AGGREGATE][%v] aggregated.\n", event.Effect)
	return nil
}

func GetSampleState(events []event.Event) SampleState {
	state := SampleState{}

	for _, event := range events {
		switch event.Effect {
		case SampleEffect:
			state.ID = event.ID
		}
	}
	return state
}
