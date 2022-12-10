package artifact

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/common"
)

type aggregateImpl struct {
	store Store
}

type Aggregator interface {
	Aggregate(event Event) error
}

func NewAggregator(store Store) Aggregator {
	return aggregateImpl{
		store: store,
	}
}

/*
*** WHERE THE MAGIC HAPPENS ***
	In case failure at the end of transaction, input reverse event.
*/
func (i aggregateImpl) Aggregate(event Event) error {
	if err := i.aggregate(event); err != nil {
		return fmt.Errorf("[%v] %v", event.Effect, err)
	}
	fmt.Printf("[%v] aggregated.\n", event.Effect)
	return nil
}

func (i aggregateImpl) aggregate(event Event) error {
	if !event.Effect.IsValid() {
		return common.ErrEventNotValid
	}

	events, err := i.store.getEventsByID(event.ID)
	if err != nil {
		return err
	}

	switch event.Effect {
	case InitEffect:
		newState := toState(append(events, event))
		if newState.ID == "" {
			return common.ErrAggregateFail
		}

	case SpawnSpikeEffect:
		currentState := toState(events)
		if currentState.Energy < 10 && currentState.Status != PlantedStatus {
			return common.ErrAggregateFail
		}
	case MoveEffect:
		currentState := toState(events)
		if currentState.Status != PlantedStatus {
			return common.ErrAggregateFail
		}

	case GatherEffect:
		currentState := toState(events)
		if currentState.Status != PlantedStatus {
			return common.ErrAggregateFail
		}
	case DropEFfect:
		currentState := toState(events)
		if currentState.Status != HarvestedStatus {
			return common.ErrAggregateFail
		}

	default:
		return common.ErrEffectNotSupported
	}

	i.store.appendEvent(event)
	return nil
}
