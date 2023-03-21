package store

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/google/uuid"
)

type impl struct {
	eventsSet map[uuid.UUID][]event.Event
}
type Store interface {
	GetEventsByEntityID(id uuid.UUID) ([]event.Event, error)
	// WARNING: action strictly used by AGGREGATOR ONLY
	AppendEvent(event.Event) error
	GetEvents() map[uuid.UUID][]event.Event
}

func NewStore() Store {
	return impl{
		eventsSet: make(map[uuid.UUID][]event.Event),
	}
}

func (i impl) GetEventsByEntityID(id uuid.UUID) ([]event.Event, error) {
	return i.eventsSet[id], nil
}

func (i impl) AppendEvent(e event.Event) error {
	events := i.eventsSet[e.EntityID]
	if events == nil {
		events = []event.Event{}
	}

	i.eventsSet[e.EntityID] = append(events, e)

	return nil
}

func (i impl) GetEvents() map[uuid.UUID][]event.Event {
	events := i.eventsSet // Shallow clone
	return events
}
