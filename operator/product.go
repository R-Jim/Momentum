package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/google/uuid"
)

type ProductOperator struct {
	ProductAggregator aggregator.Aggregator
}

func (o ProductOperator) Init(id uuid.UUID, productType string) error {
	store := o.ProductAggregator.GetStore()

	event := event.NewProductInitEvent(id, productType)

	if err := o.ProductAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}
	return nil
}

func (o ProductOperator) Progress(id uuid.UUID, value float64) error {
	store := o.ProductAggregator.GetStore()
	events, err := (*store).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	event := event.NewProductProgressEvent(id, len(events)+1, value)

	if err := o.ProductAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}
	return nil
}
