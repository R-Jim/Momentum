package aggregator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/google/uuid"
)

type ProductState struct {
	ID   uuid.UUID
	Type string

	Progress float64
}

func NewProductAggregator(store *store.Store) Aggregator {
	return aggregateImpl{
		name:  "PRODUCT",
		store: store,
		aggregateSet: map[event.Effect]func([]event.Event, event.Event) error{
			//"PRODUCT_INIT"
			event.ProductInitEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				if len(currentEvents) != 0 {
					return ErrAggregateFail
				}

				if inputEvent.EntityID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				_, err := event.ParseData[string](inputEvent)
				if err != nil {
					return ErrAggregateFail
				}

				return nil
			},
			//"PRODUCT_PROGRESS"
			event.ProductProgressEffect: func(currentEvents []event.Event, inputEvent event.Event) error {
				state, err := GetProductState(currentEvents)
				if err != nil {
					return err
				}
				if state.ID.String() == uuid.Nil.String() {
					return ErrAggregateFail
				}

				_, err = event.ParseData[float64](inputEvent)
				if err != nil {
					return ErrAggregateFail
				}

				return nil
			},
		},
	}
}

func GetProductState(events []event.Event) (ProductState, error) {
	return composeState(ProductState{}, events, func(state ProductState, e event.Event) (ProductState, error) {
		switch e.Effect {
		case event.ProductInitEffect:
			state.ID = e.EntityID

			value, err := event.ParseData[string](e)
			if err != nil {
				return state, err
			}

			state.Type = value
		case event.ProductProgressEffect:
			value, err := event.ParseData[float64](e)
			if err != nil {
				return state, err
			}

			state.Progress += value
		}
		return state, nil
	})
}

func (s ProductState) IsFinish() bool {
	return s.Progress >= 100
}
