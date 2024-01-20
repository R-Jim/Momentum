package runner

import (
	"github.com/R-jim/Momentum/template/aggregate"
	"github.com/R-jim/Momentum/template/event"
)

func NewAggregator() aggregate.Aggregator {
	return aggregate.NewAggregator(
		"RUNNER",
		map[event.Effect]map[string]string{
			InitEffect: {
				"": "INIT",
			},
			DestroyEffect: {
				"INIT": "DESTROY",
			},
		},
	)
}
