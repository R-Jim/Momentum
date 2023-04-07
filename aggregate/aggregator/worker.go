package aggregator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/google/uuid"
)

type WorkerState struct {
	ID         uuid.UUID
	BuildingID uuid.UUID
}

func NewWorkerAggregator(store *store.Store) Aggregator {
	return aggregateImpl{
		name:  "WORKER",
		store: store,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{
			//"WORKER_INIT"
			event.WorkerInitEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				if len(currentEvents) != 0 {
					return ErrAggregateFail
				}

				if inputEvent.EntityID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				return nil
			},
			//"WORKER_ASSIGN"
			event.WorkerAssignEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetWorkerState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}
				if state.BuildingID.String() != uuid.Nil.String() {
					return ErrAggregateFail
				}

				_, err = event.ParseData[uuid.UUID](inputEvent)
				if err != nil {
					return err
				}
				return nil
			},
			//"WORKER_UNASSIGN"
			event.WorkerUnassignEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetWorkerState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}
				if state.BuildingID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				buildingID, err := event.ParseData[uuid.UUID](inputEvent)
				if err != nil {
					return err
				}

				if state.BuildingID != buildingID {
					return ErrAggregateFail
				}

				return nil
			},
		},
	}
}

func GetWorkerState(events []event.Event) (WorkerState, error) {
	return composeState(WorkerState{}, events, func(state WorkerState, e event.Event) (WorkerState, error) {
		switch e.Effect {
		case event.WorkerInitEffect:
			state.ID = e.EntityID
		case event.WorkerAssignEffect:
			buildingID, err := event.ParseData[uuid.UUID](e)
			if err != nil {
				return state, err
			}
			state.BuildingID = buildingID
		case event.WorkerUnassignEffect:
			state.BuildingID = uuid.Nil
		}
		return state, nil
	})
}
