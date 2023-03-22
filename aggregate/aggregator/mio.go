package aggregator

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
)

type MioState struct {
	ID       uuid.UUID
	Position math.Pos
}

type mioAggregateImpl struct {
	store        *store.Store
	aggregateSet map[event.Effect]func([]event.Event, event.Event) error
}

func NewMioAggregator(store *store.Store) Aggregator {
	return mioAggregateImpl{
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
				destinationPosition, ok := inputEvent.Data.(math.Pos)
				if !ok {
					return pkgerrors.WithStack(fmt.Errorf("failed to parse data for effect: %s", event.MioWalkEffect))
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
				destinationPosition, ok := inputEvent.Data.(math.Pos)
				if !ok {
					return pkgerrors.WithStack(fmt.Errorf("failed to parse data for effect: %s", event.MioWalkEffect))
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
		},
	}
}

func (i mioAggregateImpl) GetStore() *store.Store {
	return i.store
}

func (i mioAggregateImpl) Aggregate(event event.Event) error {
	if err := aggregate(i.store, i.aggregateSet, event); err != nil {
		return fmt.Errorf("[MIO_AGGREGATE][%v] %v", event.Effect, err)
	}
	fmt.Printf("[MIO_AGGREGATE][%v] aggregated.\n", event.Effect)
	return nil
}

func GetMioState(events []event.Event) (MioState, error) {
	state := MioState{}

	for _, e := range events {
		switch e.Effect {
		case event.MioInitEffect:
			pos, ok := e.Data.(math.Pos)
			if !ok {
				return state, pkgerrors.WithStack(fmt.Errorf("failed to compose state for effect: %s", e.Effect))
			}
			state.ID = e.EntityID
			state.Position = pos
		case event.MioWalkEffect:
			pos, ok := e.Data.(math.Pos)
			if !ok {
				return state, pkgerrors.WithStack(fmt.Errorf("failed to compose state for effect: %s", e.Effect))
			}
			state.Position = pos
		case event.MioRunEffect:
			pos, ok := e.Data.(math.Pos)
			if !ok {
				return state, pkgerrors.WithStack(fmt.Errorf("failed to compose state for effect: %s", e.Effect))
			}
			state.Position = pos
		case event.MioIdleEffect:
		}
	}
	return state, nil
}
