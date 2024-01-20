package link

import (
	"github.com/R-jim/Momentum/template/event"
)

func GetAliveLink(events []event.Event) (*Link, error) {
	var l *Link
	for _, e := range events {
		switch e.Effect {
		case InitEffect:
			linkData, err := event.ParseData[Link](e)
			if err != nil {
				return nil, err
			}
			l = &linkData
		case DestroyEffect:
			return nil, nil
		}
	}

	return l, nil
}
