package aggregator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type BuildingState struct {
	ID        uuid.UUID
	EntityMap map[uuid.UUID]bool
	Pos       math.Pos
}

func NewBuildingAggregator(store *store.Store) Aggregator {
	return aggregateImpl{
		name:  "BUILDING",
		store: store,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{
			//"BUILDING_INIT"
			event.BuildingInitEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetBuildingState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() != uuid.Nil.String() {
					return ErrAggregateFail
				}

				if inputEvent.EntityID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				_, err = event.ParseData[math.Pos](inputEvent)
				if err != nil {
					return err
				}

				return nil
			},
			//"BUILDING_ENTITY_ENTER"
			event.BuildingEntityEnterEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetBuildingState(currentEvents)
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
			//"BUILDING_ENTITY_LEAVE"
			event.BuildingEntityLeaveEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetBuildingState(currentEvents)
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

func GetBuildingState(events []event.Event) (BuildingState, error) {
	return composeState(BuildingState{}, events, func(state BuildingState, e event.Event) (BuildingState, error) {
		switch e.Effect {
		case event.BuildingInitEffect:
			pos, err := event.ParseData[math.Pos](e)
			if err != nil {
				return state, err
			}

			state.ID = e.EntityID
			state.Pos = pos
			state.EntityMap = map[uuid.UUID]bool{}

		case event.BuildingEntityEnterEffect:
			entityID, err := event.ParseData[uuid.UUID](e)
			if err != nil {
				return state, err
			}

			state.EntityMap[entityID] = true
		case event.BuildingEntityLeaveEffect:
			entityID, err := event.ParseData[uuid.UUID](e)
			if err != nil {
				return state, err
			}
			state.EntityMap[entityID] = false
		}
		return state, nil
	})
}
