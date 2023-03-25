package aggregator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/google/uuid"
)

type BuildingState struct {
	ID uuid.UUID
}

func NewBuildingAggregator(store *store.Store) Aggregator {
	return aggregateImpl{
		name:         "BUILDING",
		store:        store,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{},
	}
}

func GetBuildingState(events []event.Event) (BuildingState, error) {
	return composeState(BuildingState{}, events, func(state BuildingState, e event.Event) (BuildingState, error) {
		switch e.Effect {
		}
		return state, nil
	})
}
