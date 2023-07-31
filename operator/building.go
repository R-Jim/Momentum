package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type BuildingOperator struct {
	buildingStore *event.BuildingStore
	productStore  *event.ProductStore
}

func NewBuilding(buildingStore *event.BuildingStore, productStore *event.ProductStore) BuildingOperator {
	return BuildingOperator{
		buildingStore: buildingStore,
		productStore:  productStore,
	}
}

func (o BuildingOperator) Init(id uuid.UUID, buildingType event.BuildingType, pos math.Pos) error {
	event := o.buildingStore.NewBuildingInitEvent(id, buildingType, pos)

	if err := aggregator.NewBuildingAggregator(o.buildingStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.buildingStore, event); err != nil {
		return err
	}
	return nil
}

func (o BuildingOperator) EntityEnter(id, entityID uuid.UUID) error {
	event := o.buildingStore.NewBuildingEntityEnterEvent(id, entityID)

	if err := aggregator.NewBuildingAggregator(o.buildingStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.buildingStore, event); err != nil {
		return err
	}
	return nil
}

func (o BuildingOperator) EntityLeave(id, entityID uuid.UUID) error {
	event := o.buildingStore.NewBuildingEntityLeaveEvent(id, entityID)

	if err := aggregator.NewBuildingAggregator(o.buildingStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.buildingStore, event); err != nil {
		return err
	}

	return nil
}

type ActProductProgress struct {
	ProductID uuid.UUID
	Value     float64
}

func (o BuildingOperator) EntityAct(id, entityID uuid.UUID, productProgress ActProductProgress) error {
	actEvent := o.buildingStore.NewBuildingEntityActEvent(id, entityID)

	if err := aggregator.NewBuildingAggregator(o.buildingStore).Aggregate(actEvent); err != nil {
		return err
	}

	if productProgress.ProductID != uuid.Nil {
		productProgressEvent := o.productStore.NewProductProgressEvent(productProgress.ProductID, productProgress.Value)

		if err := aggregator.NewProductAggregator(o.productStore).Aggregate(productProgressEvent); err != nil {
			return err
		}
		if err := appendEvent(o.productStore, productProgressEvent); err != nil {
			return err
		}
	}
	if err := appendEvent(o.buildingStore, actEvent); err != nil {
		return err
	}

	return nil
}
