package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type StreetOperator struct {
	streetStore *event.StreetStore
}

func NewStreet(streetStore *event.StreetStore) StreetOperator {
	return StreetOperator{
		streetStore: streetStore,
	}
}

func (o StreetOperator) Init(id uuid.UUID, headA, headB math.Pos) error {
	initEvent := o.streetStore.NewStreetInitEvent(id, headA, headB)

	if err := aggregator.NewStreetAggregator(o.streetStore).Aggregate(initEvent); err != nil {
		return err
	}

	if err := appendEvent(o.streetStore, initEvent); err != nil {
		return err
	}
	return nil
}

func (o StreetOperator) EntityEnter(id, entityID uuid.UUID) error {
	entityEnterEvent := o.streetStore.NewStreetEntityEnterEvent(id, entityID)

	if err := aggregator.NewStreetAggregator(o.streetStore).Aggregate(entityEnterEvent); err != nil {
		return err
	}

	if err := appendEvent(o.streetStore, entityEnterEvent); err != nil {
		return err
	}
	return nil
}

func (o StreetOperator) EntityLeave(id, entityID uuid.UUID) error {
	entityLeaveEvent := o.streetStore.NewStreetEntityLeaveEvent(id, entityID)

	if err := aggregator.NewStreetAggregator(o.streetStore).Aggregate(entityLeaveEvent); err != nil {
		return err
	}

	if err := appendEvent(o.streetStore, entityLeaveEvent); err != nil {
		return err
	}
	return nil
}
