package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type mioOperator struct {
	mioAggregator aggregator.Aggregator

	mioAnimator animator.Animator
}

func NewMio(mioAggregator aggregator.Aggregator, mioAnimator animator.Animator) mioOperator {
	return mioOperator{
		mioAggregator: mioAggregator,
		mioAnimator:   mioAnimator,
	}
}

func (o mioOperator) Init(id uuid.UUID, position math.Pos) error {
	store := o.mioAggregator.GetStore()
	event := event.NewMioInitEvent(id, position)

	if err := o.mioAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.mioAnimator != nil {
		if err := animator.Draw(o.mioAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
	}
	return nil
}

func (o mioOperator) Walk(id uuid.UUID, posEnd math.Pos) error {
	store := o.mioAggregator.GetStore()
	events, err := (*store).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	event := event.NewMioWalkEvent(id, len(events)+1, posEnd)

	if err := o.mioAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.mioAnimator != nil {
		if err := animator.Draw(o.mioAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
	}
	return nil
}

func (o mioOperator) Run(id uuid.UUID, posEnd math.Pos) error {
	store := o.mioAggregator.GetStore()
	events, err := (*store).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	event := event.NewMioRunEvent(id, len(events)+1, posEnd)

	if err := o.mioAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.mioAnimator != nil {
		if err := animator.Draw(o.mioAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
	}
	return nil
}

func (o mioOperator) Idle(id uuid.UUID) error {
	store := o.mioAggregator.GetStore()
	events, err := (*store).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	event := event.NewMioIdleEvent(id, len(events)+1)

	if err := o.mioAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.mioAnimator != nil {
		if err := animator.Draw(o.mioAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
	}
	return nil
}
