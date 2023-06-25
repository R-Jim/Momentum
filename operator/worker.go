package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/animator"
	"github.com/google/uuid"
)

type WorkerOperator struct {
	WorkerAggregator   aggregator.Aggregator
	BuildingAggregator aggregator.Aggregator

	WorkerAnimator   animator.Animator
	BuildingAnimator animator.Animator
}

func (o WorkerOperator) Init(id uuid.UUID) error {
	store := o.WorkerAggregator.GetStore()
	event := event.NewWorkerInitEvent(id)

	if err := o.WorkerAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.WorkerAnimator != nil {
		o.WorkerAnimator.Animator().AppendEvent(event)
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

	if o.WorkerAnimator != nil {
		o.WorkerAnimator.Animator().AppendEvent(workerAssignBuildingEvent)
	}
	if o.BuildingAnimator != nil {
		o.BuildingAnimator.Animator().AppendEvent(buildingWorkerAssignEvent)
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

	if o.WorkerAnimator != nil {
		o.WorkerAnimator.Animator().AppendEvent(workerAssignBuildingEvent)
	}
	if o.BuildingAnimator != nil {
		o.BuildingAnimator.Animator().AppendEvent(buildingWorkerAssignEvent)
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

	if o.WorkerAnimator != nil {
		o.WorkerAnimator.Animator().AppendEvent(workerActEvent)
	}
	if o.BuildingAnimator != nil {
		o.BuildingAnimator.Animator().AppendEvent(buildingWorkerActEvent)
	}
	return nil
}
