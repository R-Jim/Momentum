package operator

import (
	"github.com/R-jim/Momentum/fueltank"
)

type fuelTankOperator struct {
	fuelTankAggregator fueltank.Aggregator
}

func (ft fuelTankOperator) Init(fuelTankID string) error {
	fuelTankInitEvent := fueltank.NewInitEvent(fuelTankID)
	err := ft.fuelTankAggregator.Aggregate(fuelTankInitEvent)
	if err != nil {
		return err
	}
	return nil
}

func (ft fuelTankOperator) Refill(fuelTankID string, quantity int) error {
	fuelTankRefillEvent := fueltank.NewRefillEvent(fuelTankID, quantity)
	err := ft.fuelTankAggregator.Aggregate(fuelTankRefillEvent)
	if err != nil {
		return err
	}
	return nil
}
