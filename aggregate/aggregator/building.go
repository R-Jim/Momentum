package aggregator

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/google/uuid"
)

type BuildingState struct {
	ID uuid.UUID
}

type buildingAggregateImpl struct {
	store        *store.Store
	aggregateSet map[event.Effect]func([]event.Event, event.Event) error
}

func NewBuildingAggregator(store *store.Store) Aggregator {
	return mioAggregateImpl{
		store:        store,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{},
	}
}

func (i buildingAggregateImpl) GetStore() *store.Store {
	return i.store
}

func (i buildingAggregateImpl) Aggregate(event event.Event) error {
	if err := aggregate(i.store, i.aggregateSet, event); err != nil {
		return fmt.Errorf("[BUILDING_AGGREGATE][%v] %v", event.Effect, err)
	}
	fmt.Printf("[BUILDING_AGGREGATE][%v] aggregated.\n", event.Effect)
	return nil
}

func GetBuildingState(events []event.Event) (BuildingState, error) {
	state := BuildingState{}

	for _, e := range events {
		switch e.Effect {
		}
	}
	return state, nil
}
