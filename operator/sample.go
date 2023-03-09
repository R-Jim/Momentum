package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/animator"
	"github.com/google/uuid"
)

type sampleOperator struct {
	sampleAggregator aggregator.Aggregator

	sampleAnimator animator.Animator
}

func NewSample(sampleAggregator aggregator.Aggregator, sampleAnimator animator.Animator) sampleOperator {
	return sampleOperator{
		sampleAggregator: sampleAggregator,
		sampleAnimator:   sampleAnimator,
	}
}

func (o sampleOperator) SampleOperate(id uuid.UUID) error {
	event := event.NewSampleChangeEvent(id)

	if err := o.sampleAggregator.Aggregate(event); err != nil {
		return err
	}

	store := o.sampleAggregator.GetStore()
	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if err := animator.Draw(o.sampleAnimator.GetAnimateSet(), event); err != nil {
		return err
	}
	return nil
}