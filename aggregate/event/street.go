package event

import (
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

const (
	StreetInitEffect Effect = "STREET_INIT"

	StreetEntityEnterEffect Effect = "STREET_ENTITY_ENTER"
	StreetEntityLeaveEffect Effect = "STREET_ENTITY_LEAVE"
)

func NewStreetInitEvent(streetID uuid.UUID, headA, headB math.Pos) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: streetID,
		Effect:   StreetInitEffect,
		Data:     []math.Pos{headA, headB},
	}
}

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
