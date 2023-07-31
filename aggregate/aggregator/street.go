package aggregator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type StreetState struct {
	ID        uuid.UUID
	EntityMap map[uuid.UUID]bool
	HeadA     math.Pos
	HeadB     math.Pos
}

func NewStreetAggregator(store *event.StreetStore) Aggregator {
	s := event.Store(*store)

	return aggregateImpl{
		name:  "STREET",
		store: &s,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{
			//"STREET_INIT"
			event.StreetInitEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetStreetState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() != uuid.Nil.String() {
					return ErrAggregateFail
				}

				if inputEvent.EntityID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				heads, err := event.ParseData[[]math.Pos](inputEvent)
				if err != nil {
					return err
				}
				if len(heads) != 2 {
					return ErrAggregateFail
				}

				return nil
			},
			//"STREET_ENTITY_ENTER"
			event.StreetEntityEnterEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetStreetState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}
				_, err = event.ParseData[uuid.UUID](inputEvent)
				if err != nil {
					return err
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
				entityID, err := event.ParseData[uuid.UUID](inputEvent)
				if err != nil {
					return err
				}

				if !state.EntityMap[entityID] {
					return ErrAggregateFail
				}

				return nil
			},
		},
	}
}

func GetStreetState(events []event.Event) (StreetState, error) {
	return composeState(StreetState{}, events, func(state StreetState, e event.Event) (StreetState, error) {
		switch e.Effect {
		case event.StreetInitEffect:
			heads, err := event.ParseData[[]math.Pos](e)
			if err != nil {
				return state, err
			}

			state.ID = e.EntityID
			state.HeadA = heads[0]
			state.HeadB = heads[1]
			state.EntityMap = map[uuid.UUID]bool{}

		case event.StreetEntityEnterEffect:
			entityID, err := event.ParseData[uuid.UUID](e)
			if err != nil {
				return state, err
			}

			state.EntityMap[entityID] = true
		case event.StreetEntityLeaveEffect:
			entityID, err := event.ParseData[uuid.UUID](e)
			if err != nil {
				return state, err
			}
			state.EntityMap[entityID] = false
		}
		return state, nil
	})
}
