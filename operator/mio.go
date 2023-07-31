package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type MioOperator struct {
	mioStore      *event.MioStore
	buildingStore *event.BuildingStore
}

func NewMio(mioStore *event.MioStore, buildingStore *event.BuildingStore) MioOperator {
	return MioOperator{
		mioStore:      mioStore,
		buildingStore: buildingStore,
	}
}

func (o MioOperator) Init(id uuid.UUID, position math.Pos) error {
	event := o.mioStore.NewMioInitEvent(id, position)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, event); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) Walk(id uuid.UUID, posEnd math.Pos) error {
	event := o.mioStore.NewMioWalkEvent(id, posEnd)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, event); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) Run(id uuid.UUID, posEnd math.Pos) error {
	event := o.mioStore.NewMioRunEvent(id, posEnd)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, event); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) Idle(id uuid.UUID) error {
	event := o.mioStore.NewMioIdleEvent(id)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, event); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) EnterStreet(id uuid.UUID, streetID uuid.UUID) error {
	event := o.mioStore.NewMioEnterStreetEvent(id, streetID)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, event); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) SelectBuilding(id uuid.UUID, buildingID uuid.UUID) error {
	mioEnterBuildingEvent := o.mioStore.NewMioSelectBuildingEvent(id, buildingID)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(mioEnterBuildingEvent); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, mioEnterBuildingEvent); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) UnselectBuilding(id uuid.UUID, buildingID uuid.UUID) error {
	mioLeaveBuildingEvent := o.mioStore.NewMioUnselectBuildingEvent(id, buildingID)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(mioLeaveBuildingEvent); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, mioLeaveBuildingEvent); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) EnterBuilding(id uuid.UUID, buildingID uuid.UUID) error {
	mioEnterBuildingEvent := o.mioStore.NewMioEnterBuildingEvent(id, buildingID)
	buildingMioEnterEvent := o.buildingStore.NewBuildingEntityEnterEvent(buildingID, id)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(mioEnterBuildingEvent); err != nil {
		return err
	}
	if err := aggregator.NewBuildingAggregator(o.buildingStore).Aggregate(buildingMioEnterEvent); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, mioEnterBuildingEvent); err != nil {
		return err
	}
	if err := appendEvent(o.buildingStore, buildingMioEnterEvent); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) LeaveBuilding(id uuid.UUID, buildingID uuid.UUID) error {
	mioLeaveBuildingEvent := o.mioStore.NewMioLeaveBuildingEvent(id, buildingID)
	buildingMioLeaveEvent := o.buildingStore.NewBuildingEntityLeaveEvent(buildingID, id)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(mioLeaveBuildingEvent); err != nil {
		return err
	}
	if err := aggregator.NewBuildingAggregator(o.buildingStore).Aggregate(buildingMioLeaveEvent); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, mioLeaveBuildingEvent); err != nil {
		return err
	}
	if err := appendEvent(o.buildingStore, buildingMioLeaveEvent); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) Act(id uuid.UUID, buildingID uuid.UUID) error {
	mioActEvent := o.mioStore.NewMioActEvent(id, buildingID)
	buildingMioActEvent := o.buildingStore.NewBuildingEntityActEvent(buildingID, id)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(mioActEvent); err != nil {
		return err
	}
	if err := aggregator.NewBuildingAggregator(o.buildingStore).Aggregate(buildingMioActEvent); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, mioActEvent); err != nil {
		return err
	}
	if err := appendEvent(o.buildingStore, buildingMioActEvent); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) Stream(id uuid.UUID, value int) error {
	event := o.mioStore.NewMioStreamEvent(id, value)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, event); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) Eat(id uuid.UUID, value int) error {
	event := o.mioStore.NewMioEatEvent(id, value)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, event); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) Starve(id uuid.UUID, value int) error {
	event := o.mioStore.NewMioStarveEvent(id, value)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, event); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) Drink(id uuid.UUID, value int) error {
	event := o.mioStore.NewMioDrinkEvent(id, value)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, event); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) Sweat(id uuid.UUID, value int) error {
	event := o.mioStore.NewMioSweatEvent(id, value)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, event); err != nil {
		return err
	}
	return nil
}

func (o MioOperator) ChangePlannedPoses(id uuid.UUID, value []math.Pos) error {
	event := o.mioStore.NewMioChangePlannedPoses(id, value)

	if err := aggregator.NewMioAggregator(o.mioStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.mioStore, event); err != nil {
		return err
	}
	return nil
}
