package example

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/template/event"
)

type sampleOperator struct {
	sampleStore *event.Store
}

func NewSample(sampleStore *event.Store) sampleOperator {
	return sampleOperator{
		sampleStore: sampleStore,
	}
}

func (o sampleOperator) SampleInit(id uuid.UUID) error {
	sampleInitEvent := NewSampleInitEvent(*o.sampleStore, id)

	if err := NewAggregator().Aggregate(o.sampleStore, sampleInitEvent); err != nil {
		return err
	}

	if err := o.sampleStore.AppendEvent(sampleInitEvent); err != nil {
		return err
	}
	return nil
}
