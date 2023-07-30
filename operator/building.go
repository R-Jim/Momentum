package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type BuildingOperator struct {
	BuildingAggregator aggregator.Aggregator
	ProductAggregator  aggregator.Aggregator
}

func (o BuildingOperator) Init(id uuid.UUID, buildingType event.BuildingType, pos math.Pos) error {
	store := o.BuildingAggregator.GetStore()

	event := event.NewBuildingInitEvent(id, buildingType, pos)

	if err := o.BuildingAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}
	return nil
}

func (o BuildingOperator) EntityEnter(id, entityID uuid.UUID) error {
	store := o.BuildingAggregator.GetStore()
	events, err := (*store).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	event := event.NewBuildingEntityEnterEvent(id, len(events)+1, entityID)

	if err := o.BuildingAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}
	return nil
}

func (o BuildingOperator) EntityLeave(id, entityID uuid.UUID) error {
	store := o.BuildingAggregator.GetStore()
	events, err := (*store).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	event := event.NewBuildingEntityLeaveEvent(id, len(events)+1, entityID)

	if err := o.BuildingAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	return nil
}

type ActProductProgress struct {
	ProductID uuid.UUID
	Value     float64
}

func (o BuildingOperator) EntityAct(id, entityID uuid.UUID, productProgress ActProductProgress) error {
	buildingStore := o.BuildingAggregator.GetStore()
	buildingEvents, err := (*buildingStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	actEvent := event.NewBuildingEntityActEvent(id, entityID, len(buildingEvents)+1)

	if err := o.BuildingAggregator.Aggregate(actEvent); err != nil {
		return err
	}

	if productProgress.ProductID != uuid.Nil {
		productStore := o.ProductAggregator.GetStore()
		productEvents, err := (*productStore).GetEventsByEntityID(productProgress.ProductID)
		if err != nil {
			return err
		}

		productProgressEvent := event.NewProductProgressEvent(productProgress.ProductID, len(productEvents)+1, productProgress.Value)

		if err := o.ProductAggregator.Aggregate(productProgressEvent); err != nil {
			return err
		}
		if err := (*productStore).AppendEvent(productProgressEvent); err != nil {
			return err
		}
	}
	if err := (*buildingStore).AppendEvent(actEvent); err != nil {
		return err
	}

	return nil
}
