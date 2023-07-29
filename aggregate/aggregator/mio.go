package aggregator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/system"
	"github.com/google/uuid"
)

const (
	MAX_RUN_DISTANT  = 2
	MIN_RUN_DISTANT  = 1
	MAX_WALK_DISTANT = 1
)

type MioState struct {
	ID         uuid.UUID
	Position   math.Pos
	BuildingID uuid.UUID

	SelectedBuildingID uuid.UUID
	PlannedPoses       []math.Pos
}

type MioActivityState struct {
	MaxMood        int
	MaxEnergy      int
	MaxDehydration int

	Mood        int
	Energy      int
	Dehydration int
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
				if distanceSqrt > MAX_WALK_DISTANT {
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
				if distanceSqrt < MIN_RUN_DISTANT || distanceSqrt > MAX_RUN_DISTANT {
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
			//"MIO_SELECT_BUILDING"
			event.MioSelectBuildingEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}
				// if state.StreetID.String() != uuid.Nil.String() {
				// 	return ErrAggregateFail
				// }

				_, err = event.ParseData[uuid.UUID](inputEvent)
				if err != nil {
					return err
				}
				return nil
			},
			//"MIO_UNSELECT_BUILDING"
			event.MioUnselectBuildingEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
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
			//"MIO_ENTER_BUILDING"
			event.MioEnterBuildingEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
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
			//"MIO_LEAVE_BUILDING"
			event.MioLeaveBuildingEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
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
			//"MIO_STREAM"
			event.MioStreamEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				_, err = event.ParseData[int](inputEvent)
				if err != nil {
					return err
				}

				return nil
			},
			//"MIO_EAT"
			event.MioEatEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				_, err = event.ParseData[int](inputEvent)
				if err != nil {
					return err
				}

				return nil
			},
			//"MIO_STARVE"
			event.MioStarveEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				_, err = event.ParseData[int](inputEvent)
				if err != nil {
					return err
				}

				return nil
			},
			//"MIO_DRINK"
			event.MioDrinkEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				_, err = event.ParseData[int](inputEvent)
				if err != nil {
					return err
				}

				return nil
			},
			//"MIO_SWEAT"
			event.MioSweatEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				_, err = event.ParseData[int](inputEvent)
				if err != nil {
					return err
				}

				return nil
			},

			//"MIO_CHANGE_PLANNED_POSITIONS"
			event.MioChangePlannedPoses: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetMioState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				_, err = event.ParseData[[]math.Pos](inputEvent)
				if err != nil {
					return err
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
		case event.MioSelectBuildingEffect:
			buildingID, err := event.ParseData[uuid.UUID](e)
			if err != nil {
				return state, err
			}
			state.SelectedBuildingID = buildingID
		case event.MioUnselectBuildingEffect:
			state.SelectedBuildingID = uuid.Nil
		case event.MioEnterBuildingEffect:
			buildingID, err := event.ParseData[uuid.UUID](e)
			if err != nil {
				return state, err
			}
			state.BuildingID = buildingID
		case event.MioLeaveBuildingEffect:
			state.BuildingID = uuid.Nil

		case event.MioChangePlannedPoses:
			plannedPoses, err := event.ParseData[[]math.Pos](e)
			if err != nil {
				return state, err
			}
			state.PlannedPoses = plannedPoses
		}
		return state, nil
	})
}

func GetMioActivityState(events []event.Event) (MioActivityState, error) {
	return composeState(MioActivityState{}, events, func(state MioActivityState, e event.Event) (MioActivityState, error) {
		switch e.Effect {
		case event.MioInitEffect:
			state.MaxMood = system.STAT_CAP
			state.MaxEnergy = system.STAT_CAP
			state.MaxDehydration = system.STAT_CAP

			state.Mood = 70
			state.Energy = 70
			state.Dehydration = 70
		case event.MioStreamEffect:
			value, err := event.ParseData[int](e)
			if err != nil {
				return state, err
			}

			if state.Mood < value {
				state.Mood = 0
			} else {
				state.Mood -= value
			}
		case event.MioEatEffect:
			value, err := event.ParseData[int](e)
			if err != nil {
				return state, err
			}

			state.Energy += value
			if state.Energy > state.MaxEnergy {
				state.Energy = state.MaxEnergy
			}

			state.Mood += value / 4
			if state.Mood > state.MaxMood {
				state.Mood = state.MaxMood
			}
		case event.MioStarveEffect:
			value, err := event.ParseData[int](e)
			if err != nil {
				return state, err
			}

			if state.Energy < value {
				state.Energy = 0
			} else {
				state.Energy -= value
			}
		case event.MioDrinkEffect:
			value, err := event.ParseData[int](e)
			if err != nil {
				return state, err
			}

			state.Dehydration += value
			if state.Dehydration > state.MaxDehydration {
				state.Dehydration = state.MaxDehydration
			}

			state.Mood += value / 4
			if state.Mood > state.MaxMood {
				state.Mood = state.MaxMood
			}
		case event.MioSweatEffect:
			value, err := event.ParseData[int](e)
			if err != nil {
				return state, err
			}

			if state.Dehydration < value {
				state.Dehydration = 0
			} else {
				state.Dehydration -= value
			}
		}
		return state, nil
	})
}
