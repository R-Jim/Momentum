package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type WorkerOperator struct {
	workerStore   *event.WorkerStore
	buildingStore *event.BuildingStore
}

func NewWorker(workerStore *event.WorkerStore, buildingStore *event.BuildingStore) WorkerOperator {
	return WorkerOperator{
		workerStore:   workerStore,
		buildingStore: buildingStore,
	}
}

func (o WorkerOperator) Init(id uuid.UUID, position math.Pos) error {
	event := o.workerStore.NewWorkerInitEvent(id, position)

	if err := aggregator.NewWorkerAggregator(o.workerStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.workerStore, event); err != nil {
		return err
	}
	return nil
}

func (o WorkerOperator) AssignBuilding(id uuid.UUID, buildingID uuid.UUID) error {
	workerAssignBuildingEvent := o.workerStore.NewWorkerAssignEvent(id, buildingID)
	buildingWorkerAssignEvent := o.buildingStore.NewBuildingWorkerAssignEvent(buildingID, id)

	if err := aggregator.NewWorkerAggregator(o.workerStore).Aggregate(workerAssignBuildingEvent); err != nil {
		return err
	}
	if err := aggregator.NewBuildingAggregator(o.buildingStore).Aggregate(buildingWorkerAssignEvent); err != nil {
		return err
	}

	if err := appendEvent(o.workerStore, workerAssignBuildingEvent); err != nil {
		return err
	}
	if err := appendEvent(o.buildingStore, buildingWorkerAssignEvent); err != nil {
		return err
	}
	return nil
}

func (o WorkerOperator) UnassignBuilding(id uuid.UUID, buildingID uuid.UUID) error {
	buildingWorkerAssignEvent := o.buildingStore.NewBuildingWorkerUnassignEvent(buildingID, id)
	workerAssignBuildingEvent := o.workerStore.NewWorkerUnassignEvent(id, buildingID)

	if err := aggregator.NewWorkerAggregator(o.workerStore).Aggregate(workerAssignBuildingEvent); err != nil {
		return err
	}
	if err := aggregator.NewBuildingAggregator(o.buildingStore).Aggregate(buildingWorkerAssignEvent); err != nil {
		return err
	}

	if err := appendEvent(o.workerStore, workerAssignBuildingEvent); err != nil {
		return err
	}
	if err := appendEvent(o.buildingStore, buildingWorkerAssignEvent); err != nil {
		return err
	}
	return nil
}

func (o WorkerOperator) Act(id uuid.UUID, buildingID uuid.UUID) error {
	workerActEvent := o.workerStore.NewWorkerActEvent(id, buildingID)
	buildingWorkerActEvent := o.buildingStore.NewBuildingWorkerActEvent(buildingID, id)

	if err := aggregator.NewWorkerAggregator(o.workerStore).Aggregate(workerActEvent); err != nil {
		return err
	}
	if err := aggregator.NewBuildingAggregator(o.buildingStore).Aggregate(buildingWorkerActEvent); err != nil {
		return err
	}

	if err := appendEvent(o.workerStore, workerActEvent); err != nil {
		return err
	}
	if err := appendEvent(o.buildingStore, buildingWorkerActEvent); err != nil {
		return err
	}
	return nil
}

func (o WorkerOperator) Move(id uuid.UUID, position math.Pos) error {
	event := o.workerStore.NewWorkerMoveEvent(id, position)

	if err := aggregator.NewWorkerAggregator(o.workerStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.workerStore, event); err != nil {
		return err
	}
	return nil
}

func (o WorkerOperator) ChangePlannedPoses(id uuid.UUID, value []math.Pos) error {
	event := o.workerStore.NewWorkerChangePlannedPoses(id, value)

	if err := aggregator.NewWorkerAggregator(o.workerStore).Aggregate(event); err != nil {
		return err
	}

	if err := appendEvent(o.workerStore, event); err != nil {
		return err
	}
	return nil
}
