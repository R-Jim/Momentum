package aggregator

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
)

type StreetState struct {
	ID        uuid.UUID
	EntityMap map[uuid.UUID]bool
}

type streetAggregateImpl struct {
	store        *store.Store
	aggregateSet map[event.Effect]func([]event.Event, event.Event) error
}

func NewStreetAggregator(store *store.Store) Aggregator {
	return mioAggregateImpl{
		store: store,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{
			//"STREET_ENTITY_ENTER"
			event.StreetEntityEnterEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetStreetState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}
				_, ok := inputEvent.Data.(uuid.UUID)
				if !ok {
					return pkgerrors.WithStack(fmt.Errorf("failed to parse data for effect: %s", event.StreetEntityEnterEffect))
				}

				return nil
			},
			//"STREET_ENTITY_LEAVE"
			event.StreetEntityLeaveEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetStreetState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}
				entityID, ok := inputEvent.Data.(uuid.UUID)
				if !ok {
					return pkgerrors.WithStack(fmt.Errorf("failed to parse data for effect: %s", event.StreetEntityLeaveEffect))
				}

				if !state.EntityMap[entityID] {
					return ErrAggregateFail
				}

				return nil
			},
		},
	}
}

func (i streetAggregateImpl) GetStore() *store.Store {
	return i.store
}

func (i streetAggregateImpl) Aggregate(event event.Event) error {
	if err := aggregate(i.store, i.aggregateSet, event); err != nil {
		return fmt.Errorf("[STREET_AGGREGATE][%v] %v", event.Effect, err)
	}
	fmt.Printf("[STREET_AGGREGATE][%v] aggregated.\n", event.Effect)
	return nil
}

func GetStreetState(events []event.Event) (StreetState, error) {
	state := StreetState{}

	for _, e := range events {
		switch e.Effect {
		case event.StreetEntityEnterEffect:
			entityID, ok := e.Data.(uuid.UUID)
			if !ok {
				return state, pkgerrors.WithStack(fmt.Errorf("failed to compose state for effect: %s", e.Effect))
			}
			state.EntityMap[entityID] = true
		case event.StreetEntityLeaveEffect:
			entityID, ok := e.Data.(uuid.UUID)
			if !ok {
				return state, pkgerrors.WithStack(fmt.Errorf("failed to compose state for effect: %s", e.Effect))
			}
			state.EntityMap[entityID] = false
		}
	}
	return state, nil
}
