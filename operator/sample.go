package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/google/uuid"
)

type sampleOperator struct {
	sampleAggregator aggregator.Aggregator
}

func NewSample(sampleAggregator aggregator.Aggregator) sampleOperator {
	return sampleOperator{
		sampleAggregator: sampleAggregator,
	}
}

func (o sampleOperator) SampleOperate(id uuid.UUID) error {
	event := event.NewSampleEvent(id)

	if err := o.sampleAggregator.Aggregate(event); err != nil {
		return err
	}

	store := o.sampleAggregator.GetStore()
	if err := (*store).AppendEvent(event); err != nil {
		return err
	}
	return nil
}
