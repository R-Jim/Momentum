package position

import (
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/template/event"
)

func GetPositionProjection(events []event.Event) (math.Pos, error) {
	for i := len(events) - 1; i >= 0; i-- {
		if events[i].Effect == MoveEffect || events[i].Effect == InitEffect {
			lastPosition, err := event.ParseData[math.Pos](events[i])
			if err != nil {
				return math.Pos{}, err
			}
			return lastPosition, nil
		}
	}

	return math.Pos{}, ErrPositionNotFound
}
