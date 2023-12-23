package example

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/template/event"
)

const (
	SampleInitEffect event.Effect = "SAMPLE_EFFECT"
	SampleMoveEffect event.Effect = "MOVE_EFFECT"
)

func NewSampleInitEvent(s event.Store, id uuid.UUID) event.Event {
	return s.NewEvent(id, SampleInitEffect, "SAMPLE_INIT_DATA")
}

func NewSampleMoveEvent(s event.Store, id uuid.UUID) event.Event {
	return s.NewEvent(id, SampleMoveEffect, "SAMPLE_MOVE_DATA")
}
