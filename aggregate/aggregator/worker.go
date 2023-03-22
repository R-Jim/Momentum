package aggregator

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/google/uuid"
)

type WorkerState struct {
	ID uuid.UUID
}

type workerAggregateImpl struct {
	store        *store.Store
	aggregateSet map[event.Effect]func([]event.Event, event.Event) error
}

func NewWorkerAggregator(store *store.Store) Aggregator {
	return mioAggregateImpl{
		store:        store,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{},
	}
}

func (i workerAggregateImpl) GetStore() *store.Store {
	return i.store
}

func (i workerAggregateImpl) Aggregate(event event.Event) error {
	if err := aggregate(i.store, i.aggregateSet, event); err != nil {
		return fmt.Errorf("[WORKER_AGGREGATE][%v] %v", event.Effect, err)
	}
	fmt.Printf("[WORKER_AGGREGATE][%v] aggregated.\n", event.Effect)
	return nil
}

func GetWorkerState(events []event.Event) (WorkerState, error) {
	state := WorkerState{}

	for _, e := range events {
		switch e.Effect {
		}
	}
	return state, nil
}
