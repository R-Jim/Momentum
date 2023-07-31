package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/google/uuid"
)

type sampleOperator struct {
	sampleStore *event.SampleStore
}

func NewSample(sampleStore *event.SampleStore) sampleOperator {
	return sampleOperator{
		sampleStore: sampleStore,
	}
}

func (o sampleOperator) SampleOperate(id uuid.UUID) error {
	sampleEvent := o.sampleStore.NewSampleEvent(id)

	sampleAggregator := aggregator.NewSampleAggregator(o.sampleStore)
	if err := sampleAggregator.Aggregate(sampleEvent); err != nil {
		return err
	}

	if err := appendEvent(o.sampleStore, sampleEvent); err != nil {
		return err
	}
	return nil
}
