package event

import (
	"github.com/google/uuid"
)

const (
	SampleEffect Effect = "SAMPLE_EFFECT"
)

func NewSampleEvent(id uuid.UUID) Event {
	return newEvent(id, 0, SampleEffect, "SAMPLE_DATA")
}
