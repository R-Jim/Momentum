package link

import (
	"github.com/R-jim/Momentum/template/event"
	"github.com/google/uuid"
)

type LinkProjection struct {
	ID          uuid.UUID
	SourceID    uuid.UUID
	TargetID    uuid.UUID
	Strength    int
	IsDestroyed bool
}

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

func GetLinkProjection(events []event.Event) (LinkProjection, error) {
	var l LinkProjection
	for _, e := range events {
		switch e.Effect {
		case InitEffect:
			linkData, err := event.ParseData[Link](e)
			if err != nil {
				return l, err
			}
			l = LinkProjection{
				ID:       e.EntityID,
				SourceID: linkData.source,
				TargetID: linkData.target,
				Strength: 1,
			}
		case StrengthenEffect:
			l.Strength++
		case DestroyEffect:
			l.IsDestroyed = true
			return l, nil
		}
	}

	return l, nil
}
