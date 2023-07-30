package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type WorkerOperator struct {
	WorkerAggregator   aggregator.Aggregator
	BuildingAggregator aggregator.Aggregator
}

func (o WorkerOperator) Init(id uuid.UUID, position math.Pos) error {
	store := o.WorkerAggregator.GetStore()
	event := event.NewWorkerInitEvent(id, position)

	if err := o.WorkerAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}
	return nil
}

func (o WorkerOperator) AssignBuilding(id uuid.UUID, buildingID uuid.UUID) error {
	workerStore := o.WorkerAggregator.GetStore()
	buildingStore := o.BuildingAggregator.GetStore()

	workerEvents, err := (*workerStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}
	buildingEvents, err := (*buildingStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	workerAssignBuildingEvent := event.NewWorkerAssignEvent(id, buildingID, len(workerEvents)+1)
	buildingWorkerAssignEvent := event.NewBuildingWorkerAssignEvent(buildingID, id, len(buildingEvents)+1)

	if err := o.WorkerAggregator.Aggregate(workerAssignBuildingEvent); err != nil {
		return err
	}
	if err := o.BuildingAggregator.Aggregate(buildingWorkerAssignEvent); err != nil {
		return err
	}

	if err := (*workerStore).AppendEvent(workerAssignBuildingEvent); err != nil {
		return err
	}
	if err := (*buildingStore).AppendEvent(buildingWorkerAssignEvent); err != nil {
		return err
	}
	return nil
}

func (o WorkerOperator) UnassignBuilding(id uuid.UUID, buildingID uuid.UUID) error {
	workerStore := o.WorkerAggregator.GetStore()
	buildingStore := o.BuildingAggregator.GetStore()

	workerEvents, err := (*workerStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}
	buildingEvents, err := (*buildingStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	workerAssignBuildingEvent := event.NewWorkerUnassignEvent(id, buildingID, len(workerEvents)+1)
	buildingWorkerAssignEvent := event.NewBuildingWorkerUnassignEvent(buildingID, id, len(buildingEvents)+1)

	if err := o.WorkerAggregator.Aggregate(workerAssignBuildingEvent); err != nil {
		return err
	}
	if err := o.BuildingAggregator.Aggregate(buildingWorkerAssignEvent); err != nil {
		return err
	}

	if err := (*workerStore).AppendEvent(workerAssignBuildingEvent); err != nil {
		return err
	}
	if err := (*buildingStore).AppendEvent(buildingWorkerAssignEvent); err != nil {
		return err
	}
	return nil
}

func (o WorkerOperator) Act(id uuid.UUID, buildingID uuid.UUID) error {
	workerStore := o.WorkerAggregator.GetStore()
	buildingStore := o.BuildingAggregator.GetStore()

	workerEvents, err := (*workerStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}
	buildingEvents, err := (*buildingStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	workerActEvent := event.NewWorkerActEvent(id, buildingID, len(workerEvents)+1)
	buildingWorkerActEvent := event.NewBuildingWorkerActEvent(buildingID, id, len(buildingEvents)+1)

	if err := o.WorkerAggregator.Aggregate(workerActEvent); err != nil {
		return err
	}
	if err := o.BuildingAggregator.Aggregate(buildingWorkerActEvent); err != nil {
		return err
	}

	if err := (*workerStore).AppendEvent(workerActEvent); err != nil {
		return err
	}
	if err := (*buildingStore).AppendEvent(buildingWorkerActEvent); err != nil {
		return err
	}
	return nil
}

func (o WorkerOperator) Move(id uuid.UUID, position math.Pos) error {
	store := o.WorkerAggregator.GetStore()
	events := store.GetEvents()[id]

	event := event.NewWorkerMoveEvent(id, len(events)+1, position)

	if err := o.WorkerAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}
	return nil
}
