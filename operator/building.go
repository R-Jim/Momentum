package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

type BuildingOperator struct {
	BuildingAggregator aggregator.Aggregator

	BuildingAnimator animator.Animator
}

func (o BuildingOperator) Init(id uuid.UUID, pos math.Pos) error {
	store := o.BuildingAggregator.GetStore()

	event := event.NewBuildingInitEvent(id, pos)

	if err := o.BuildingAggregator.Aggregate(event); err != nil {
		return err
	}

	if err := (*store).AppendEvent(event); err != nil {
		return err
	}

	if o.BuildingAnimator != nil {
		if err := animator.Draw(o.BuildingAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
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

	if o.BuildingAnimator != nil {
		if err := animator.Draw(o.BuildingAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
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

	if o.BuildingAnimator != nil {
		if err := animator.Draw(o.BuildingAnimator.GetAnimateSet(), event); err != nil {
			return err
		}
	}
	return nil
}
