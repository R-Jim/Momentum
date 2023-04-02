package event

import (
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

const (
	MioInitEffect Effect = "MIO_INIT"

	MioWalkEffect Effect = "MIO_WALK"
	MioRunEffect  Effect = "MIO_RUN"
	MioIdleEffect Effect = "MIO_IDLE"

	MioEnterStreetEffect Effect = "MIO_ENTER_STREET"
)

func NewMioInitEvent(entityID uuid.UUID, position math.Pos) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  0,
		Effect:   MioInitEffect,
		Data:     position,
	}
}

func NewMioWalkEvent(entityID uuid.UUID, version int, newPosition math.Pos) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   MioWalkEffect,
		Data:     newPosition,
	}
}

func NewMioRunEvent(entityID uuid.UUID, version int, newPosition math.Pos) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   MioRunEffect,
		Data:     newPosition,
	}
}

func NewMioIdleEvent(entityID uuid.UUID, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   MioIdleEffect,
	}
}

func NewMioEnterStreetEvent(entityID uuid.UUID, streetID uuid.UUID, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   MioEnterStreetEffect,
		Data:     streetID,
	}
}
