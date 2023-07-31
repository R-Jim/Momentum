package aggregator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type BuildingState struct {
	ID        uuid.UUID
	EntityMap map[uuid.UUID]bool
	WorkerMap map[uuid.UUID]bool
	Pos       math.Pos

	Type string
}

func NewBuildingAggregator(store *event.BuildingStore) Aggregator {
	s := event.Store(*store)

	return aggregateImpl{
		name:  "BUILDING",
		store: &s,
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

				_, err = event.ParseData[event.BuildingInitEventData](inputEvent)
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
			//"BUILDING_ENTITY_ACT"
			event.BuildingEntityActEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
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

				if entityID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				if !state.EntityMap[entityID] {
					return ErrAggregateFail
				}

				return nil
			},
			//"BUILDING_WORKER_ACT"
			event.BuildingWorkerActEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetBuildingState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}
				workerID, err := event.ParseData[uuid.UUID](inputEvent)
				if err != nil {
					return err
				}

				if workerID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				if !state.WorkerMap[workerID] {
					return ErrAggregateFail
				}

				return nil
			},
			//"BUILDING_WORKER_ASSIGN"
			event.BuildingWorkerAssignEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
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
			//"BUILDING_WORKER_UNASSIGN"
			event.BuildingWorkerUnassignEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
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

				if !state.WorkerMap[entityID] {
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
			data, err := event.ParseData[event.BuildingInitEventData](e)
			if err != nil {
				return state, err
			}

			state.ID = e.EntityID
			state.Type = string(data.Type)
			state.Pos = data.Pos
			state.EntityMap = map[uuid.UUID]bool{}
			state.WorkerMap = map[uuid.UUID]bool{}

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

		case event.BuildingWorkerAssignEffect:
			workerID, err := event.ParseData[uuid.UUID](e)
			if err != nil {
				return state, err
			}

			state.WorkerMap[workerID] = true
		case event.BuildingWorkerUnassignEffect:
			workerID, err := event.ParseData[uuid.UUID](e)
			if err != nil {
				return state, err
			}
			state.WorkerMap[workerID] = false
		}
		return state, nil
	})
}
