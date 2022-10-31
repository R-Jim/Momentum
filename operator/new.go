package operator

import (
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/domain/jet"
	"github.com/R-jim/Momentum/domain/storage"
)

type OperatorAggregator struct {
	JetAggregator      jet.Aggregator
	FuelTankAggregator storage.Aggregator
}

type Operator struct {
	Jet      jetOperator
	FuelTank fuelTankOperator
}

func New(aggregator OperatorAggregator, animator animator.Animator) Operator {
	return Operator{
		Jet: jetOperator{
			jetAggregator:      aggregator.JetAggregator,
			fuelTankAggregator: aggregator.FuelTankAggregator,

			animator: animator,
		},
		FuelTank: fuelTankOperator{
			fuelTankAggregator: aggregator.FuelTankAggregator,
		},
	}
}
