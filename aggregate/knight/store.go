package knight

import "github.com/R-jim/Momentum/aggregate/common"

type impl struct {
	events map[string][]Event
}
type Store interface {
	getEventsByID(id string) ([]Event, error)
	// WARNING: action strictly used by AGGREGATOR ONLY
	appendEvent(Event) error
}

func NewStore() Store {
	return impl{
		events: make(map[string][]Event),
	}
}

func (i impl) getEventsByID(id string) ([]Event, error) {
	return i.events[id], nil
}

// WARNING: action strictly used by AGGREGATOR ONLY
func (i impl) appendEvent(e Event) error {
	if !e.Effect.IsValid() {
		return common.ErrEventNotValid
	}

	events := i.events[e.ID]
	if events == nil {
		events = []Event{}
	}

	i.events[e.ID] = append(events, e)

	return nil
}
