package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type streetOperator struct {
	streetAggregator aggregator.Aggregator

	streetAnimator animator.Animator
}

func NewStreet(streetAggregator aggregator.Aggregator, streetAnimator animator.Animator) streetOperator {
	return streetOperator{
		streetAggregator: streetAggregator,
		streetAnimator:   streetAnimator,
	}
}

func (o streetOperator) Init(id uuid.UUID, headA, headB math.Pos) error {
	store := o.streetAggregator.GetStore()

	event := event.NewStreetInitEvent(id, headA, headB)

	if err := o.streetAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.streetAnimator != nil {
		if err := animator.Draw(o.streetAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
	}
	return nil
}
