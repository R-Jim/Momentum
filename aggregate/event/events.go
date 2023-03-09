package event

import "github.com/google/uuid"

type data interface{}

type Effect string

type Event struct {
	ID       uuid.UUID
	EntityID uuid.UUID
	Version  int
	Effect   Effect
	Data     data
}
