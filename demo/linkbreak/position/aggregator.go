package position

import (
	"github.com/R-jim/Momentum/template/aggregate"
	"github.com/R-jim/Momentum/template/event"
)

func NewAggregator() aggregate.Aggregator {
	return aggregate.NewAggregator(
		"POSITION",
		map[event.Effect]map[string]string{
			InitEffect: {
				"": "INIT",
			},
			DestroyEffect: {
				"INIT": "DESTROY",
				"MOVE": "DESTROY",
			},
			MoveEffect: {
				"INIT": "MOVE",
				"MOVE": "MOVE",
			},
		},
	)
}
