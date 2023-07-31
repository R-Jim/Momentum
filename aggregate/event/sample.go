package event

import (
	"github.com/google/uuid"
)

const (
	SampleEffect Effect = "SAMPLE_EFFECT"
)

type SampleStore Store

func NewSampleStore() SampleStore {
	return SampleStore(newStore())
}

func (s SampleStore) NewSampleEvent(id uuid.UUID) Event {
	return Store(s).newEvent(id, SampleEffect, "SAMPLE_DATA")
}
