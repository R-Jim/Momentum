package health

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/template/event"
)

const (
	InitEffect    event.Effect = "INIT_EFFECT"
	ModifyEffect  event.Effect = "MODIFY_EFFECT"
	DestroyEffect event.Effect = "DESTROY_EFFECT"
)

func NewInitEvent(s *event.Store, id uuid.UUID, baseValue int) event.Event {
	return s.NewEvent(id, InitEffect, baseValue)
}

func NewDestroyEvent(s *event.Store, id uuid.UUID) event.Event {
	return s.NewEvent(id, DestroyEffect, nil)
}

func NewModifyEvent(s *event.Store, id uuid.UUID, value int) event.Event {
	return s.NewEvent(id, ModifyEffect, value)
}
