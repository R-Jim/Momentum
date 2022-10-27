package operator

import (
	"fmt"

	"github.com/R-jim/Momentum/fueltank"
	"github.com/R-jim/Momentum/jet"
)

type jetOperator struct {
	jetAggregator      jet.Aggregator
	fuelTankAggregator fueltank.Aggregator
}

func (j jetOperator) Init(jetID string, fuelTankID string) error {
	jetInitEvent := jet.NewInitEvent(jetID)
	err := j.jetAggregator.Aggregate(jetInitEvent)
	if err != nil {
		return err
	}

	if fuelTankID != "" {
		jetChangeFuelTankEvent := jet.NewChangeFuelTankEvent(jetID, fuelTankID)
		err = j.jetAggregator.Aggregate(jetChangeFuelTankEvent)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j jetOperator) Fly(jetID string, jetFuelTankID string, fuelConsumed int, toPosition jet.PositionState) error {
	jetFuelTankConsumeEvent := fueltank.NewConsumeEvent(jetFuelTankID, fuelConsumed)
	err := j.fuelTankAggregator.Aggregate(jetFuelTankConsumeEvent)
	if err != nil {
		return err
	}

	jetFlyEvent := jet.NewFlyEvent(jetID, toPosition)
	err = j.jetAggregator.Aggregate(jetFlyEvent)
	if err != nil {
		jetFuelTankRefillEvent := fueltank.NewRefillEvent(jetFuelTankID, fuelConsumed)
		newErr := j.fuelTankAggregator.Aggregate(jetFuelTankRefillEvent)
		if newErr != nil {
			fmt.Printf("[%v] rollback failed, err: %v.\n", jetFuelTankRefillEvent.Effect, newErr)
			return err
		}
		return err
	}
	return nil
}

func (j jetOperator) Landing(jetID string) error {
	landingEvent := jet.NewLandingEvent(jetID)
	err := j.jetAggregator.Aggregate(landingEvent)
	if err != nil {
		return err
	}
	return nil
}
