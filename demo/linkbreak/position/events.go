package position

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/template/event"
)

const (
	InitEffect    event.Effect = "INIT_EFFECT"
	DestroyEffect event.Effect = "DESTROY_EFFECT"
	MoveEffect    event.Effect = "MOVE_EFFECT"
)

func NewInitEvent(s *event.Store, id uuid.UUID, destination math.Pos) event.Event {
	return s.NewEvent(id, InitEffect, destination)
}

func NewDestroyEvent(s *event.Store, id uuid.UUID) event.Event {
	return s.NewEvent(id, DestroyEffect, nil)
}

func NewMoveEvent(s *event.Store, id uuid.UUID, destination math.Pos) event.Event {
	return s.NewEvent(id, MoveEffect, destination)
}
