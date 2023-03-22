package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/animator"
)

type buildingOperator struct {
	buildingAggregator aggregator.Aggregator

	buildingAnimator animator.Animator
}

func NewBuilding(buildingAggregator aggregator.Aggregator, buildingAnimator animator.Animator) buildingOperator {
	return buildingOperator{
		buildingAggregator: buildingAggregator,
		buildingAnimator:   buildingAnimator,
	}
}
