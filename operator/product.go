package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/google/uuid"
)

type ProductOperator struct {
	productStore *event.ProductStore
}

func NewProduct(productStore *event.ProductStore) ProductOperator {
	return ProductOperator{
		productStore: productStore,
	}
}

func (o ProductOperator) Init(id uuid.UUID, productType string) error {
	event := o.productStore.NewProductInitEvent(id, productType)

	if err := aggregator.NewProductAggregator(o.productStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.productStore, event); err != nil {
		return err
	}
	return nil
}

func (o ProductOperator) Progress(id uuid.UUID, value float64) error {
	event := o.productStore.NewProductProgressEvent(id, value)

	if err := aggregator.NewProductAggregator(o.productStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.productStore, event); err != nil {
		return err
	}
	return nil
}
