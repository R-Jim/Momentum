package operator

import (
	"github.com/R-jim/Momentum/aggregate/carrier"
	"github.com/R-jim/Momentum/aggregate/jet"
	"github.com/R-jim/Momentum/aggregate/storage"
	"github.com/R-jim/Momentum/animator"
)

type OperatorAggregator struct {
	JetAggregator      jet.Aggregator
	FuelTankAggregator storage.Aggregator
	CarrierAggregator  carrier.Aggregator
}

type Operator struct {
	Jet      jetOperator
	FuelTank fuelTankOperator
	Carrier  carrierOperator
}

func New(aggregator OperatorAggregator, animator animator.Animator) Operator {
	jetOperator := jetOperator{
		jetAggregator:      aggregator.JetAggregator,
		fuelTankAggregator: aggregator.FuelTankAggregator,

		animator: animator,
	}
	return Operator{
		Jet: jetOperator,
		FuelTank: fuelTankOperator{
			fuelTankAggregator: aggregator.FuelTankAggregator,
		},
		Carrier: carrierOperator{
			jetOperator:       jetOperator,
			carrierAggregator: aggregator.CarrierAggregator,
			animator:          animator,
		},
	}
}
