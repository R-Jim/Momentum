package knight

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
	case MoveEffect:
		currentState := toState(events)
		if currentState.Health.Value <= 0 {
			return common.ErrAggregateFail
		}

	case StrikeEffect:
		currentState := toState(events)
		if currentState.Health.Value <= 0 {
			return common.ErrAggregateFail
		}

	case GatherArtifactEffect:
		currentState := toState(events)
		if currentState.HarvestedArtifactID != "" {
			return common.ErrAggregateFail
		}
		if currentState.Health.Value <= 0 {
			return common.ErrAggregateFail
		}
	case DropArtifactEFfect:
		currentState := toState(events)
		if currentState.HarvestedArtifactID == "" {
			return common.ErrAggregateFail
		}
	case ChangeTargetEffect:
	default:
		return common.ErrEffectNotSupported
	}

	i.store.appendEvent(event)
	return nil
}
