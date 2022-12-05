package artifact

import "github.com/R-jim/Momentum/aggregate/common"

type impl struct {
	eventsSet map[string][]Event
}
type Store interface {
	getEventsByID(id string) ([]Event, error)
	// WARNING: action strictly used by AGGREGATOR ONLY
	appendEvent(Event) error
	GetEvents() map[string][]Event
}

func NewStore() Store {
	return impl{
		eventsSet: make(map[string][]Event),
	}
}

func (i impl) getEventsByID(id string) ([]Event, error) {
	return i.eventsSet[id], nil
}

// WARNING: action strictly used by AGGREGATOR ONLY
func (i impl) appendEvent(e Event) error {
	if !e.Effect.IsValid() {
		return common.ErrEventNotValid
	}

	events := i.eventsSet[e.ID]
	if events == nil {
		events = []Event{}
	}

	i.eventsSet[e.ID] = append(events, e)

	return nil
}

func (i impl) GetEvents() map[string][]Event {
	events := i.eventsSet // Shallow clone
	return events
}
