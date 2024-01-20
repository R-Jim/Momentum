package runner

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/template/event"
)

const (
	InitEffect        event.Effect = "INIT_EFFECT"
	DestroyEffect     event.Effect = "DESTROY_EFFECT"
	UpdateLinksEffect event.Effect = "UPDATE_LINKS_EFFECT"
)

func NewInitEvent(s *event.Store, id uuid.UUID, runner Runner) event.Event {
	return s.NewEvent(id, InitEffect, runner)
}

func NewDestroyEvent(s *event.Store, id uuid.UUID) event.Event {
	return s.NewEvent(id, DestroyEffect, nil)
}

func NewUpdateLinksEvent(s *event.Store, id uuid.UUID, linkIDs []uuid.UUID) event.Event {
	return s.NewEvent(id, UpdateLinksEffect, linkIDs)
}
