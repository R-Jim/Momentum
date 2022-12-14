package operator

import (
	"github.com/R-jim/Momentum/aggregate/artifact"
	"github.com/R-jim/Momentum/aggregate/carrier"
	"github.com/R-jim/Momentum/aggregate/jet"
	"github.com/R-jim/Momentum/aggregate/knight"
	"github.com/R-jim/Momentum/aggregate/spike"
	"github.com/R-jim/Momentum/aggregate/storage"
	"github.com/R-jim/Momentum/animator"
)

type OperatorAggregator struct {
	JetAggregator      jet.Aggregator
	FuelTankAggregator storage.Aggregator
	CarrierAggregator  carrier.Aggregator
	SpikeAggregator    spike.Aggregator
	ArtifactAggregator artifact.Aggregator
	KnightAggregator   knight.Aggregator
}

type Operator struct {
	Jet      jetOperator
	FuelTank fuelTankOperator
	Carrier  carrierOperator
	Spike    spikeOperator
	Artifact artifactOperator
	Knight   knightOperator
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
		Spike: spikeOperator{
			artifactAggregator: aggregator.ArtifactAggregator,
			spikeAggregator:    aggregator.SpikeAggregator,
			animator:           animator,
		},
		Artifact: artifactOperator{
			artifactAggregator: aggregator.ArtifactAggregator,
			spikeAggregator:    aggregator.SpikeAggregator,
			animator:           animator,
		},
		Knight: knightOperator{
			knightAggregator: aggregator.KnightAggregator,
			animator:         animator,
		},
	}
}
