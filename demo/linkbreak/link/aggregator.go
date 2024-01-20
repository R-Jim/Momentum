package link

import (
	"github.com/R-jim/Momentum/template/aggregate"
	"github.com/R-jim/Momentum/template/event"
)

func NewAggregator() aggregate.Aggregator {
	return aggregate.NewAggregator(
		"LINK",
		map[event.Effect]map[string]string{
			InitEffect: {
				"": "INIT",
			},
			DestroyEffect: {
				"INIT":       "DESTROY",
				"STRENGTHEN": "DESTROY",
			},
			StrengthenEffect: {
				"INIT":       "STRENGTHEN",
				"STRENGTHEN": "STRENGTHEN",
			},
		},
	)
}
