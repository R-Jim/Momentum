package jet

import (
	"fmt"

	"github.com/R-jim/Momentum/common"
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
		currentState := toCombatState(events)
		if currentState.TargetID != "" || currentState.Status != 0 {
			return common.ErrAggregateFail
		}

	case CancelAttackEffect:
		currentState := toCombatState(events)
		if currentState.TargetID != "" {
			return common.ErrAggregateFail
		}

	case EngageEffect:
		currentState := toCombatState(events)
		if currentState.Status != FlyingStatus {
			return common.ErrAggregateFail
		}
	case DisengageEffect:
		currentState := toCombatState(events)
		if currentState.Status != EngagingStatus {
			return common.ErrAggregateFail
		}

	case LandingEffect:
		currentState := toCombatState(events)
		if currentState.Status != FlyingStatus {
			return common.ErrAggregateFail
		}
	case TakeOffEffect:
		currentState := toCombatState(events)
		if currentState.Status != LandedStatus {
			return common.ErrAggregateFail
		}

	case FlyEffect:
		currentState := toCombatState(events)
		if currentState.Status != FlyingStatus {
			return common.ErrAggregateFail
		}

	case FuelTankChangedEffect:
		currentState := toCombatState(events)
		if currentState.Status != LandedStatus {
			return common.ErrAggregateFail
		}
	default:
		return common.ErrEffectNotSupported
	}

	i.store.appendEvent(event)
	return nil
}
