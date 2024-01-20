package link

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/template/event"
)

const (
	InitEffect       event.Effect = "INIT_EFFECT"
	DestroyEffect    event.Effect = "DESTROY_EFFECT"
	StrengthenEffect event.Effect = "STRENGTHEN_EFFECT"
)

func NewInitEvent(s *event.Store, id uuid.UUID, link Link) event.Event {
	return s.NewEvent(id, InitEffect, link)
}

func NewDestroyEvent(s *event.Store, id uuid.UUID) event.Event {
	return s.NewEvent(id, DestroyEffect, nil)
}

func NewStrengthenEvent(s *event.Store, id uuid.UUID) event.Event {
	return s.NewEvent(id, StrengthenEffect, nil)
}
