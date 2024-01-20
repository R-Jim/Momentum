package health

import (
	"github.com/R-jim/Momentum/template/aggregate"
	"github.com/R-jim/Momentum/template/event"
)

func NewAggregator() aggregate.Aggregator {
	return aggregate.NewAggregator(
		"HEALTH",
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
