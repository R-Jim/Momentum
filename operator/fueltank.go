package operator

import (
	"github.com/R-jim/Momentum/aggregate/storage"
)

type fuelTankOperator struct {
	fuelTankAggregator storage.Aggregator
}

func (ft fuelTankOperator) Init(fuelTankID string) error {
	fuelTankInitEvent := storage.NewInitEvent(fuelTankID)
	err := ft.fuelTankAggregator.Aggregate(fuelTankInitEvent)
	if err != nil {
		return err
	}
	return nil
}

func (ft fuelTankOperator) Refill(fuelTankID string, quantity int) error {
	fuelTankRefillEvent := storage.NewRefillEvent(fuelTankID, quantity)
	err := ft.fuelTankAggregator.Aggregate(fuelTankRefillEvent)
	if err != nil {
		return err
	}
	return nil
}
