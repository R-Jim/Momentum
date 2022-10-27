package operator

import (
	"github.com/R-jim/Momentum/fueltank"
	"github.com/R-jim/Momentum/jet"
)

type Operator struct {
	Jet      jetOperator
	FuelTank fuelTankOperator
}

func New(jetAggregator jet.Aggregator, fuelTankAggregator fueltank.Aggregator) Operator {
	return Operator{
		Jet: jetOperator{
			jetAggregator:      jetAggregator,
			fuelTankAggregator: fuelTankAggregator,
		},
		FuelTank: fuelTankOperator{
			fuelTankAggregator: fuelTankAggregator,
		},
	}
}
