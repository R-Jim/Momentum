package aggregate

import (
	"fmt"

	"github.com/R-jim/Momentum/template/event"
)

type aggregateImpl struct {
	name         string
	aggregateSet map[event.Effect]map[string]string
}

func NewAggregator(name string, aggregateSet map[event.Effect]map[string]string) Aggregator {
	return aggregateImpl{
		name,
		aggregateSet,
	}
}

type Aggregator interface {
	Aggregate(store *event.Store, inputEvent event.Event) error
}

func (i aggregateImpl) Aggregate(store *event.Store, inputEvent event.Event) error {
	if err := i.aggregate(store, inputEvent); err != nil {
		return fmt.Errorf("[%s][%v][%v][%v] %v", i.name, inputEvent.EntityID, inputEvent.Version, inputEvent.Effect, err)
	}
	fmt.Printf("[%s][%v][%v][%v] aggregated. %v\n", i.name, inputEvent.EntityID, inputEvent.Version, inputEvent.Effect, inputEvent.Data)
	return nil
}

func (i aggregateImpl) aggregate(store *event.Store, inputEvent event.Event) error {
	events, err := (*store).GetEventsByEntityID(inputEvent.EntityID)
	if err != nil {
		return err
	}

	currentState := ""
	for _, event := range events {
		initialStateSet, isExist := i.aggregateSet[event.Effect]
		if !isExist {
			return fmt.Errorf("[COMPOSE_STATE_FAILED] no support for effect [%s], event: %v", event.Effect, event)
		}

		currentState, isExist = initialStateSet[currentState]
		if !isExist {
			return fmt.Errorf("[COMPOSE_STATE_FAILED] no support for state [%s], event: %v", currentState, event)
		}
	}

	initialStateSet, isExist := i.aggregateSet[inputEvent.Effect]
	if !isExist {
		return ErrEffectNotSupported
	}

	_, isExist = initialStateSet[currentState]
	if !isExist {
		return ErrEffectNotSupported
	}

	return nil
}
