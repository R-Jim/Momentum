package aggregator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type MioState struct {
	ID         uuid.UUID
	Position   math.Pos
	StreetID   uuid.UUID
	BuildingID uuid.UUID
}

func NewMioAggregator(store *store.Store) Aggregator {
	return aggregateImpl{
		name:  "MIO",
		store: store,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{
			//"MIO_INIT"
			event.MioInitEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				if len(currentEvents) != 0 {
					return ErrAggregateFail
				}

				return nil
			},
			//"MIO_WALK"
			event.MioWalkEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}
				destinationPosition, err := event.ParseData[math.Pos](inputEvent)
				if err != nil {
					return err
				}

				_, _, distanceSqrt := math.GetDistances(state.Position, destinationPosition)
				if distanceSqrt > 1 {
					return ErrAggregateFail
				}

				return nil
			},
			//"MIO_RUN"
			event.MioRunEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}
				destinationPosition, err := event.ParseData[math.Pos](inputEvent)
				if err != nil {
					return err
				}

				_, _, distanceSqrt := math.GetDistances(state.Position, destinationPosition)
				if distanceSqrt < 2 || distanceSqrt > 5 {
					return ErrAggregateFail
				}

				return nil
			},
			//"MIO_IDLE"
			event.MioIdleEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				return nil
			},
			//"MIO_ENTER_STREET"
			event.MioEnterStreetEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
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
			//"MIO_ENTER_BUILDING"
			event.MioEnterBuildingEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}
				if state.StreetID.String() != uuid.Nil.String() {
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
			//"MIO_LEAVE_BUILDING"
			event.MioLeaveBuildingEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.StreetID.String() != uuid.Nil.String() {
					return ErrAggregateFail
				}
				if state.BuildingID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				buildingID, err := event.ParseData[uuid.UUID](inputEvent)
				if err != nil {
					return err
				}

				if state.BuildingID.String() != buildingID.String() {
					return ErrAggregateFail
				}

				return nil
			},
			//"MIO_ACT"
			event.MioActEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.BuildingID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				buildingID, err := event.ParseData[uuid.UUID](inputEvent)
				if err != nil {
					return err
				}

				if state.BuildingID.String() != buildingID.String() {
					return ErrAggregateFail
				}

				return nil
			},
		},
	}
}

func GetMioState(events []event.Event) (MioState, error) {
	return composeState(MioState{}, events, func(state MioState, e event.Event) (MioState, error) {
		switch e.Effect {
		case event.MioInitEffect:
			pos, err := event.ParseData[math.Pos](e)
			if err != nil {
				return state, err
			}
			state.ID = e.EntityID
			state.Position = pos
		case event.MioWalkEffect:
			pos, err := event.ParseData[math.Pos](e)
			if err != nil {
				return state, err
			}
			state.Position = pos
		case event.MioRunEffect:
			pos, err := event.ParseData[math.Pos](e)
			if err != nil {
				return state, err
			}
			state.Position = pos
		case event.MioIdleEffect:

		case event.MioEnterStreetEffect:
			streetID, err := event.ParseData[uuid.UUID](e)
			if err != nil {
				return state, err
			}
			state.StreetID = streetID
		case event.MioEnterBuildingEffect:
			buildingID, err := event.ParseData[uuid.UUID](e)
			if err != nil {
				return state, err
			}
			state.BuildingID = buildingID
		case event.MioLeaveBuildingEffect:
			state.BuildingID = uuid.Nil
		}
		return state, nil
	})
}
