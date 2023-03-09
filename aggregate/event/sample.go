package event

import "github.com/google/uuid"

const (
	SampleEffect Effect = "SAMPLE_EFFECT"
)

func NewSampleChangeEvent(id uuid.UUID) Event {
	return Event{
		ID: id,
	}
}
