package aggregator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/google/uuid"
)

type WorkerState struct {
	ID uuid.UUID
}

func NewWorkerAggregator(store *store.Store) Aggregator {
	return aggregateImpl{
		name:         "WORKER",
		store:        store,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{},
	}
}

func GetWorkerState(events []event.Event) (WorkerState, error) {
	return composeState(WorkerState{}, events, func(state WorkerState, e event.Event) (WorkerState, error) {
		switch e.Effect {
		}
		return state, nil
	})
}
