package operator

import (
	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/animator"
)

type workerOperator struct {
	workerAggregator aggregator.Aggregator

	workerAnimator animator.Animator
}

func NewWorker(workerAggregator aggregator.Aggregator, workerAnimator animator.Animator) workerOperator {
	return workerOperator{
		workerAggregator: workerAggregator,
		workerAnimator:   workerAnimator,
	}
}
