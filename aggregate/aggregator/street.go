package aggregator

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/google/uuid"
)

type StreetState struct {
	ID uuid.UUID
}

type streetAggregateImpl struct {
	store        *store.Store
	aggregateSet map[event.Effect]func([]event.Event, event.Event) error
}

func NewStreetAggregator(store *store.Store) Aggregator {
	return mioAggregateImpl{
		store:        store,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{},
	}
}

func (i streetAggregateImpl) GetStore() *store.Store {
	return i.store
}

func (i streetAggregateImpl) Aggregate(event event.Event) error {
	if err := aggregate(i.store, i.aggregateSet, event); err != nil {
		return fmt.Errorf("[STREET_AGGREGATE][%v] %v", event.Effect, err)
	}
	fmt.Printf("[STREET_AGGREGATE][%v] aggregated.\n", event.Effect)
	return nil
}

func GetStreetState(events []event.Event) (MioState, error) {
	state := MioState{}

	for _, e := range events {
		switch e.Effect {
		}
	}
	return state, nil
}
