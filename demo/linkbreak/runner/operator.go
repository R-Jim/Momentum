package runner

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/demo/linkbreak/health"
	"github.com/R-jim/Momentum/demo/linkbreak/position"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/template/event"
)

type Operator struct {
	RunnerStore   *event.Store
	HealthStore   *event.Store
	PositionStore *event.Store
}

func (o Operator) NewRunner(id uuid.UUID, healthBaseValue int, factionValue int, positionValue math.Pos) error {
	runner := Runner{
		id:         id,
		faction:    factionValue,
		healthID:   uuid.New(),
		positionID: uuid.New(),
	}

	initRunnerEvent := NewInitEvent(o.RunnerStore, runner.id, runner)
	initHealthEvent := health.NewInitEvent(o.RunnerStore, runner.healthID, runner.id, healthBaseValue)
	initPositionEvent := position.NewInitEvent(o.PositionStore, runner.positionID, positionValue)

	if err := NewAggregator().Aggregate(o.RunnerStore, initRunnerEvent); err != nil {
		return err
	}
	if err := health.NewAggregator().Aggregate(o.HealthStore, initHealthEvent); err != nil {
		return err
	}
	if err := position.NewAggregator().Aggregate(o.PositionStore, initPositionEvent); err != nil {
		return err
	}

	o.RunnerStore.AppendEvent(initRunnerEvent)
	o.HealthStore.AppendEvent(initHealthEvent)
	o.PositionStore.AppendEvent(initPositionEvent)

	return nil
}

func (o Operator) MoveRunner(id uuid.UUID, positionValue math.Pos) error {
	var positionID uuid.UUID
	if events, err := o.RunnerStore.GetEventsByEntityID(id); err == nil {
		for _, e := range events {
			switch e.Effect {
			case InitEffect:
				runner, err := event.ParseData[Runner](e)
				if err != nil {
					return err
				}
				positionID = runner.positionID
			}
		}
	} else {
		return err
	}

	if positionID.String() == "" {
		return ErrPositionIDRequired
	}

	movePositionEvent := position.NewMoveEvent(o.PositionStore, positionID, positionValue)
	if err := position.NewAggregator().Aggregate(o.PositionStore, movePositionEvent); err != nil {
		return err
	}
	o.PositionStore.AppendEvent(movePositionEvent)
	return nil
}

func (o Operator) DestroyRunner(id uuid.UUID) error {
	var positionID, healthID uuid.UUID
	if events, err := o.RunnerStore.GetEventsByEntityID(id); err == nil {
		for _, e := range events {
			switch e.Effect {
			case InitEffect:
				runner, err := event.ParseData[Runner](e)
				if err != nil {
					return err
				}
				positionID = runner.positionID
				healthID = runner.healthID
			}
		}
	} else {
		return err
	}

	if positionID.String() == "" {
		return ErrPositionIDRequired
	}
	if healthID.String() == "" {
		return ErrHealthIDRequired
	}

	destroyPositionEvent := position.NewDestroyEvent(o.PositionStore, positionID)
	destroyHealthEvent := health.NewDestroyEvent(o.HealthStore, healthID)
	destroyRunnerEvent := NewDestroyEvent(o.RunnerStore, id)

	if err := position.NewAggregator().Aggregate(o.PositionStore, destroyPositionEvent); err != nil {
		return err
	}
	if err := health.NewAggregator().Aggregate(o.HealthStore, destroyHealthEvent); err != nil {
		return err
	}
	if err := NewAggregator().Aggregate(o.RunnerStore, destroyRunnerEvent); err != nil {
		return err
	}

	o.PositionStore.AppendEvent(destroyPositionEvent)
	o.HealthStore.AppendEvent(destroyHealthEvent)
	o.RunnerStore.AppendEvent(destroyRunnerEvent)
	return nil
}
