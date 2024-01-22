package automaton

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/demo/linkbreak/health"
	"github.com/R-jim/Momentum/demo/linkbreak/link"
	"github.com/R-jim/Momentum/demo/linkbreak/runner"
	"github.com/R-jim/Momentum/template/event"
)

type BreakAutomaton struct {
	linkStore   *event.Store
	runnerStore *event.Store

	healthOperator health.Operator
}

func NewBreakAutomaton(linkStore, runnerStore, healthStore *event.Store) BreakAutomaton {
	return BreakAutomaton{
		linkStore:   linkStore,
		runnerStore: runnerStore,

		healthOperator: health.Operator{
			HealthStore: healthStore,
		},
	}
}

func (b BreakAutomaton) BreakLinkedRunners(requiredLinkStrength int) error {
	breakableLinks := []link.LinkProjection{}

	for _, linkEvents := range b.linkStore.GetEvents() {
		if linkProjection, err := link.GetLinkProjection(linkEvents); err != nil {
			return err
		} else if !linkProjection.IsDestroyed && linkProjection.Strength >= requiredLinkStrength {
			breakableLinks = append(breakableLinks, linkProjection)
		}
	}

	for _, breakableLink := range breakableLinks {
		runnerEvents, err := b.runnerStore.GetEventsByEntityID(breakableLink.TargetID)
		if err != nil {
			return err
		}

		var healthID uuid.UUID
		for _, e := range runnerEvents {
			if e.Effect == runner.InitEffect {
				runner, err := event.ParseData[runner.Runner](e)
				if err != nil {
					return err
				}
				healthID = runner.HealthID()
				break
			}
		}

		if err := b.healthOperator.Modify(healthID, -1); err != nil {
			return err
		}
	}

	return nil
}
