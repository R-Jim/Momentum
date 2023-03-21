package aggregator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
)

type Aggregator interface {
	GetStore() *store.Store
	Aggregate(event event.Event) error
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
