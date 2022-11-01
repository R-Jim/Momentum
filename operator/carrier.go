package operator

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/carrier"
	"github.com/R-jim/Momentum/aggregate/jet"
	"github.com/R-jim/Momentum/animator"
)

type carrierOperator struct {
	jetOperator jetOperator

	carrierAggregator carrier.Aggregator

	animator animator.Animator
}

func (c carrierOperator) Init(carrierID string) error {
	carrierInitEvent := carrier.NewInitEvent(carrierID)
	err := c.carrierAggregator.Aggregate(carrierInitEvent)
	if err != nil {
		return err
	}
	return nil
}

func (c carrierOperator) HouseJet(carrierID string, jetID string) error {
	err := c.jetOperator.Landing(jetID)
	if err != nil {
		return err
	}

	carrierHouseJetEvent := carrier.NewHouseJetEvent(carrierID, jetID)
	err = c.carrierAggregator.Aggregate(carrierHouseJetEvent)
	if err != nil {
		newErr := c.jetOperator.Takeoff(jetID)
		if newErr != nil {
			fmt.Printf("[%v] rollback failed, err: %v.\n", jet.TakeOffEffect, newErr)
		}
		return err
	}
	return nil
}

func (c carrierOperator) LaunchJet(carrierID string, jetID string) error {
	err := c.jetOperator.Takeoff(jetID)
	if err != nil {
		return err
	}

	carrierLaunchJetEvent := carrier.NewLaunchJetEvent(carrierID, jetID)
	err = c.carrierAggregator.Aggregate(carrierLaunchJetEvent)
	if err != nil {
		newErr := c.jetOperator.Landing(jetID)
		if newErr != nil {
			fmt.Printf("[%v] rollback failed, err: %v.\n", jet.LandingEffect, newErr)
		}
		return err
	}
	return nil
}
