package spawner

import (
	"github.com/R-jim/Momentum/template/aggregate"
	"github.com/R-jim/Momentum/template/event"
)

func NewAggregator() aggregate.Aggregator {
	return aggregate.NewAggregator(
		"SPAWNER",
		map[event.Effect]map[string]string{
			InitEffect: {
				"": "INIT",
			},
			CountDownEffect: {
				"INIT":       "COUNT_DOWN",
				"COUNT_DOWN": "COUNT_DOWN",
			},
			SpawnEffect: {
				"INIT":       "SPAWN",
				"COUNT_DOWN": "SPAWN",
			},
		},
	)
}
