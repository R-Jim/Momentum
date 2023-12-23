package example

import (
	"github.com/R-jim/Momentum/template/aggregate"
	"github.com/R-jim/Momentum/template/event"
)

func NewAggregator() aggregate.Aggregator {
	return aggregate.NewAggregator(
		"SAMPLE",
		map[event.Effect]map[string]string{
			SampleInitEffect: {
				"": "INIT",
			},
			SampleMoveEffect: {
				"INIT": "MOVE",
			},
		},
	)
}
