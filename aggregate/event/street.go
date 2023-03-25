package event

import (
	"github.com/google/uuid"
)

const (
	StreetEntityEnterEffect Effect = "STREET_ENTITY_ENTER"
	StreetEntityLeaveEffect Effect = "STREET_ENTITY_LEAVE"
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

func NewStreetEntityLeaveEvent(streetID uuid.UUID, version int, entityID uuid.UUID) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: streetID,
		Version:  version,
		Effect:   StreetEntityLeaveEffect,
		Data:     entityID,
	}
}
