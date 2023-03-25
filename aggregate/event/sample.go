package event

import (
	"github.com/google/uuid"
)

const (
	SampleEffect Effect = "SAMPLE_EFFECT"
)

func NewSampleEvent(id uuid.UUID) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: id,
		Version:  0,
		Effect:   SampleEffect,
		Data:     "SAMPLE_DATA",
	}
}
