package automaton

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/demo/linkbreak/health"
	"github.com/R-jim/Momentum/demo/linkbreak/runner"
	"github.com/R-jim/Momentum/template/event"
)

type DestroyRunnerAutomaton struct {
	healthStore *event.Store

	runnerOperator runner.Operator
}

func NewDestroyRunnerAutomaton(runnerStore, positionStore, healthStore *event.Store) DestroyRunnerAutomaton {
	return DestroyRunnerAutomaton{
		healthStore: healthStore,

		runnerOperator: runner.Operator{
			RunnerStore:   runnerStore,
			PositionStore: positionStore,
			HealthStore:   healthStore,
		},
	}
}

func (d DestroyRunnerAutomaton) DestroyEmptyHealthRunner() error {
	destroyableRunnerIDs := []uuid.UUID{}

	for _, healthEvents := range d.healthStore.GetEvents() {
		runnerID, isEmptyHealth, err := isEmptyHealth(healthEvents)
		if err != nil {
			return err
		}
		if isEmptyHealth {
			destroyableRunnerIDs = append(destroyableRunnerIDs, runnerID)
		}
	}

	for _, runnerID := range destroyableRunnerIDs {
		if err := d.runnerOperator.DestroyRunner(runnerID); err != nil {
			return err
		}
	}

	return nil
}

func isEmptyHealth(events []event.Event) (uuid.UUID, bool, error) {
	var runnerID uuid.UUID
	var healthValue int

	for _, e := range events {
		switch e.Effect {
		case health.InitEffect:
			if e.Effect == health.InitEffect {
				health, err := event.ParseData[health.Health](e)
				if err != nil {
					return uuid.UUID{}, false, err
				}
				runnerID = health.OwnerID()
				healthValue = health.BaseValue()
			}
		case health.DestroyEffect:
			return uuid.UUID{}, false, nil
		case health.ModifyEffect:
			value, err := event.ParseData[int](e)
			if err != nil {
				return uuid.UUID{}, false, err
			}
			healthValue += value
		}
	}

	return runnerID, healthValue <= 0, nil
}
