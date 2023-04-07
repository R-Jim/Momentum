package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type MioOperator struct {
	MioAggregator      aggregator.Aggregator
	BuildingAggregator aggregator.Aggregator

	MioAnimator      animator.Animator
	BuildingAnimator animator.Animator
}

func (o MioOperator) Init(id uuid.UUID, position math.Pos) error {
	store := o.MioAggregator.GetStore()
	event := event.NewMioInitEvent(id, position)

	if err := o.MioAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.MioAnimator != nil {
		if err := animator.Draw(o.MioAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
	}
	return nil
}

func (o MioOperator) Walk(id uuid.UUID, posEnd math.Pos) error {
	store := o.MioAggregator.GetStore()
	events, err := (*store).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	event := event.NewMioWalkEvent(id, len(events)+1, posEnd)

	if err := o.MioAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.MioAnimator != nil {
		if err := animator.Draw(o.MioAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
	}
	return nil
}

func (o MioOperator) Run(id uuid.UUID, posEnd math.Pos) error {
	store := o.MioAggregator.GetStore()
	events, err := (*store).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	event := event.NewMioRunEvent(id, len(events)+1, posEnd)

	if err := o.MioAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.MioAnimator != nil {
		if err := animator.Draw(o.MioAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
	}
	return nil
}

func (o MioOperator) Idle(id uuid.UUID) error {
	store := o.MioAggregator.GetStore()
	events, err := (*store).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	event := event.NewMioIdleEvent(id, len(events)+1)

	if err := o.MioAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.MioAnimator != nil {
		if err := animator.Draw(o.MioAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
	}
	return nil
}

func (o MioOperator) EnterStreet(id uuid.UUID, streetID uuid.UUID) error {
	store := o.MioAggregator.GetStore()
	events, err := (*store).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	event := event.NewMioEnterStreetEvent(id, streetID, len(events)+1)

	if err := o.MioAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.MioAnimator != nil {
		if err := animator.Draw(o.MioAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
	}
	return nil
}

func (o MioOperator) EnterBuilding(id uuid.UUID, buildingID uuid.UUID) error {
	mioStore := o.MioAggregator.GetStore()
	buildingStore := o.BuildingAggregator.GetStore()

	mioEvents, err := (*mioStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}
	buildingEvents, err := (*buildingStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	mioEnterBuildingEvent := event.NewMioEnterBuildingEvent(id, buildingID, len(mioEvents)+1)
	buildingMioEnterEvent := event.NewBuildingEntityEnterEvent(buildingID, len(buildingEvents)+1, id)

	if err := o.MioAggregator.Aggregate(mioEnterBuildingEvent); err != nil {
		return err
	}
	if err := o.BuildingAggregator.Aggregate(buildingMioEnterEvent); err != nil {
		return err
	}

	if err := (*mioStore).AppendEvent(mioEnterBuildingEvent); err != nil {
		return err
	}
	if err := (*buildingStore).AppendEvent(buildingMioEnterEvent); err != nil {
		return err
	}

	if o.MioAnimator != nil {
		if err := animator.Draw(o.MioAnimator.GetAnimateSet(), mioEnterBuildingEvent); err != nil {
			return err
		}
	}
	if o.BuildingAnimator != nil {
		if err := animator.Draw(o.BuildingAnimator.GetAnimateSet(), buildingMioEnterEvent); err != nil {
			return err
		}
	}
	return nil
}

func (o MioOperator) LeaveBuilding(id uuid.UUID, buildingID uuid.UUID) error {
	mioStore := o.MioAggregator.GetStore()
	buildingStore := o.BuildingAggregator.GetStore()

	mioEvents, err := (*mioStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}
	buildingEvents, err := (*buildingStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	mioLeaveBuildingEvent := event.NewMioLeaveBuildingEvent(id, buildingID, len(mioEvents)+1)
	buildingMioLeaveEvent := event.NewBuildingEntityLeaveEvent(buildingID, len(buildingEvents)+1, id)

	if err := o.MioAggregator.Aggregate(mioLeaveBuildingEvent); err != nil {
		return err
	}
	if err := o.BuildingAggregator.Aggregate(buildingMioLeaveEvent); err != nil {
		return err
	}

	if err := (*mioStore).AppendEvent(mioLeaveBuildingEvent); err != nil {
		return err
	}
	if err := (*buildingStore).AppendEvent(buildingMioLeaveEvent); err != nil {
		return err
	}

	if o.MioAnimator != nil {
		if err := animator.Draw(o.MioAnimator.GetAnimateSet(), mioLeaveBuildingEvent); err != nil {
			return err
		}
	}
	if o.BuildingAnimator != nil {
		if err := animator.Draw(o.BuildingAnimator.GetAnimateSet(), buildingMioLeaveEvent); err != nil {
			return err
		}
	}
	return nil
}

func (o MioOperator) Act(id uuid.UUID, buildingID uuid.UUID) error {
	mioStore := o.MioAggregator.GetStore()
	buildingStore := o.BuildingAggregator.GetStore()

	mioEvents, err := (*mioStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}
	buildingEvents, err := (*buildingStore).GetEventsByEntityID(id)
	if err != nil {
		return err
	}

	mioActEvent := event.NewMioActEvent(id, buildingID, len(mioEvents)+1)
	buildingMioActEvent := event.NewBuildingEntityActEvent(buildingID, id, len(buildingEvents)+1)

	if err := o.MioAggregator.Aggregate(mioActEvent); err != nil {
		return err
	}
	if err := o.BuildingAggregator.Aggregate(buildingMioActEvent); err != nil {
		return err
	}

	if err := (*mioStore).AppendEvent(mioActEvent); err != nil {
		return err
	}
	if err := (*buildingStore).AppendEvent(buildingMioActEvent); err != nil {
		return err
	}

	if o.MioAnimator != nil {
		if err := animator.Draw(o.MioAnimator.GetAnimateSet(), mioActEvent); err != nil {
			return err
		}
	}
	if o.BuildingAnimator != nil {
		if err := animator.Draw(o.BuildingAnimator.GetAnimateSet(), buildingMioActEvent); err != nil {
			return err
		}
	}
	return nil
}
