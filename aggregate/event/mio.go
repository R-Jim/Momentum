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

	MioEnterBuildingEffect Effect = "MIO_ENTER_BUILDING"
	MioLeaveBuildingEffect Effect = "MIO_LEAVE_BUILDING"

	MioActEffect Effect = "MIO_ACT"

	MioStreamEffect Effect = "MIO_STREAM_EFFECT"
	MioEatEffect    Effect = "MIO_EAT_EFFECT"
	MioStarveEffect Effect = "MIO_STARVE_EFFECT"
	MioDrinkEffect  Effect = "MIO_DRINK_EFFECT"
	MioSweatEffect  Effect = "MIO_SWEAT_EFFECT"
)

func NewMioInitEvent(entityID uuid.UUID, position math.Pos) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  1,
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

func NewMioEnterBuildingEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   MioEnterBuildingEffect,
		Data:     buildingID,
	}
}

func NewMioLeaveBuildingEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   MioLeaveBuildingEffect,
		Data:     buildingID,
	}
}

func NewMioActEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   MioActEffect,
		Data:     buildingID,
	}
}

func NewMioStreamEvent(entityID uuid.UUID, value, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   MioStreamEffect,
		Data:     value,
	}
}

func NewMioEatEvent(entityID uuid.UUID, value, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   MioEatEffect,
		Data:     value,
	}
}

func NewMioStarveEvent(entityID uuid.UUID, value, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   MioStarveEffect,
		Data:     value,
	}
}

func NewMioDrinkEvent(entityID uuid.UUID, value, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   MioDrinkEffect,
		Data:     value,
	}
}

func NewMioSweatEvent(entityID uuid.UUID, value, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   MioSweatEffect,
		Data:     value,
	}
}
