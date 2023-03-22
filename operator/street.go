package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/animator"
)

type streetOperator struct {
	streetAggregator aggregator.Aggregator

	streetAnimator animator.Animator
}

func NewStreet(streetAggregator aggregator.Aggregator, streetAnimator animator.Animator) streetOperator {
	return streetOperator{
		streetAggregator: streetAggregator,
		streetAnimator:   streetAnimator,
	}
}
