package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type StreetOperator struct {
	streetAggregator aggregator.Aggregator

	streetAnimator animator.Animator
}

func NewStreet(streetAggregator aggregator.Aggregator, streetAnimator animator.Animator) StreetOperator {
	return StreetOperator{
		streetAggregator: streetAggregator,
		streetAnimator:   streetAnimator,
	}
}

func (o StreetOperator) Init(id uuid.UUID, headA, headB math.Pos) error {
	store := o.streetAggregator.GetStore()

	event := event.NewStreetInitEvent(id, headA, headB)

	if err := o.streetAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.streetAnimator != nil {
		o.streetAnimator.Animator().AppendEvent(event)
	}
	return nil
}

func (o StreetOperator) EntityEnter(id, entityID uuid.UUID) error {
	store := o.streetAggregator.GetStore()
	events, err := (*store).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	event := event.NewStreetEntityEnterEvent(id, len(events)+1, entityID)

	if err := o.streetAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.streetAnimator != nil {
		o.streetAnimator.Animator().AppendEvent(event)
	}
	return nil
}

func (o StreetOperator) EntityLeave(id, entityID uuid.UUID) error {
	store := o.streetAggregator.GetStore()
	events, err := (*store).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	event := event.NewStreetEntityLeaveEvent(id, len(events)+1, entityID)

	if err := o.streetAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.streetAnimator != nil {
		o.streetAnimator.Animator().AppendEvent(event)
	}
	return nil
}
