package aggregator

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
)

type state interface{}

type aggregateImpl struct {
	name         string
	store        *store.Store
	aggregateSet map[event.Effect]func([]event.Event, event.Event) error
}

type Aggregator interface {
	GetStore() *store.Store
	Aggregate(event event.Event) error
}

func (i aggregateImpl) GetStore() *store.Store {
	return i.store
}

func (i aggregateImpl) Aggregate(event event.Event) error {
	if err := aggregate(i.store, i.aggregateSet, event); err != nil {
		return fmt.Errorf("[%s][%v][%v] %v.\n", i.name, event.EntityID, event.Effect, err)
	}
	fmt.Printf("[%s][%v][%v] aggregated. %v\n", i.name, event.EntityID, event.Effect, event.Data)
	return nil
}

func aggregate(store *store.Store, aggregateSet map[event.Effect]func([]event.Event, event.Event) error, event event.Event) error {
	events, err := (*store).GetEventsByEntityID(event.EntityID)
	if err != nil {
		return err
	}

	aggregateFunc, isExist := aggregateSet[event.Effect]
	if !isExist {
		return ErrEffectNotSupported
	}

	return aggregateFunc(events, event)
}

func composeState[T state](state T, events []event.Event, composer func(T, event.Event) (T, error)) (T, error) {
	var err error
	for _, event := range events {
		state, err = composer(state, event)
		if err != nil {
			return state, fmt.Errorf("[COMPOSE_STATE_FAILED][%s]: %v", event.Effect, err)
		}
	}
	return state, nil
}
