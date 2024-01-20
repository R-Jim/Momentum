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

func (o Operator) NewRunner(healthBaseValue int, factionValue int, positionValue math.Pos) (Runner, error) {
	runner := Runner{
		id:         uuid.New(),
		faction:    factionValue,
		healthID:   uuid.New(),
		positionID: uuid.New(),
	}

	initRunnerEvent := NewInitEvent(o.RunnerStore, runner.id, runner)
	initHealthEvent := health.NewInitEvent(o.RunnerStore, runner.healthID, healthBaseValue)
	initPositionEvent := position.NewInitEvent(o.PositionStore, runner.positionID, positionValue)

	if err := NewAggregator().Aggregate(o.RunnerStore, initRunnerEvent); err != nil {
		return Runner{}, err
	}
	if err := health.NewAggregator().Aggregate(o.HealthStore, initHealthEvent); err != nil {
		return Runner{}, err
	}
	if err := position.NewAggregator().Aggregate(o.PositionStore, initPositionEvent); err != nil {
		return Runner{}, err
	}

	o.RunnerStore.AppendEvent(initRunnerEvent)
	o.HealthStore.AppendEvent(initHealthEvent)
	o.PositionStore.AppendEvent(initPositionEvent)

	return runner, nil
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
