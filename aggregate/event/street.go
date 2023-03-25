package event

import (
	"github.com/google/uuid"
)

const (
	StreetEntityEnterEffect Effect = "STREET_ENTITY_ENTER"
)

func NewStreetEntityEnterEvent(streetID uuid.UUID, version int, entityID uuid.UUID) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: streetID,
		Version:  version,
		Effect:   StreetEntityEnterEffect,
		Data:     entityID,
	}
}
