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

type StreetStore Store

func NewStreetStore() StreetStore {
	return StreetStore(newStore())
}

func (s StreetStore) NewStreetInitEvent(streetID uuid.UUID, headA, headB math.Pos) Event {
	return Store(s).newEvent(streetID, StreetInitEffect, []math.Pos{headA, headB})
}

func (s StreetStore) NewStreetEntityEnterEvent(streetID uuid.UUID, entityID uuid.UUID) Event {
	return Store(s).newEvent(streetID, StreetEntityEnterEffect, entityID)
}

func (s StreetStore) NewStreetEntityLeaveEvent(streetID uuid.UUID, entityID uuid.UUID) Event {
	return Store(s).newEvent(streetID, StreetEntityLeaveEffect, entityID)
}
