package store

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/google/uuid"
)

type Store struct {
	counter int

	eventsSet map[uuid.UUID][]event.Event
}

func NewStore() Store {
	return Store{
		eventsSet: make(map[uuid.UUID][]event.Event),
	}
}

func (i Store) GetEventsByEntityID(id uuid.UUID) ([]event.Event, error) {
	return i.eventsSet[id], nil
}

// WARNING: action strictly used by AGGREGATOR ONLY
func (i *Store) AppendEvent(e event.Event) error {
	events := i.eventsSet[e.EntityID]
	if events == nil {
		events = []event.Event{}
	}

	i.eventsSet[e.EntityID] = append(events, e)
	i.counter++
	return nil
}

func (i Store) GetEvents() map[uuid.UUID][]event.Event {
	events := i.eventsSet // Shallow clone
	return events
}

func (i Store) GetCounter() int {
	return i.counter
}
