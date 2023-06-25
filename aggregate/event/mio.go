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

	MioSelectBuildingEffect   Effect = "MIO_SELECT_EFFECT"
	MioUnselectBuildingEffect Effect = "MIO_UNSELECT_EFFECT"

	MioEnterBuildingEffect Effect = "MIO_ENTER_BUILDING"
	MioLeaveBuildingEffect Effect = "MIO_LEAVE_BUILDING"

	MioActEffect Effect = "MIO_ACT"

	MioStreamEffect Effect = "MIO_STREAM_EFFECT"
	MioEatEffect    Effect = "MIO_EAT_EFFECT"
	MioStarveEffect Effect = "MIO_STARVE_EFFECT"
	MioDrinkEffect  Effect = "MIO_DRINK_EFFECT"
	MioSweatEffect  Effect = "MIO_SWEAT_EFFECT"

	MioChangePlannedPoses Effect = "MIO_CHANGE_PLANNED_POSITIONS"
)

func NewMioInitEvent(entityID uuid.UUID, position math.Pos) Event {
	return newEvent(entityID, 1, MioInitEffect, position)
}

func NewMioWalkEvent(entityID uuid.UUID, version int, newPosition math.Pos) Event {
	return newEvent(entityID, version, MioWalkEffect, newPosition)
}

func NewMioRunEvent(entityID uuid.UUID, version int, newPosition math.Pos) Event {
	return newEvent(
		entityID,
		version,
		MioRunEffect,
		newPosition,
	)
}

func NewMioIdleEvent(entityID uuid.UUID, version int) Event {
	return newEvent(
		entityID,
		version,
		MioIdleEffect,
		nil,
	)
}

func NewMioEnterStreetEvent(entityID uuid.UUID, streetID uuid.UUID, version int) Event {
	return newEvent(
		entityID,
		version,
		MioEnterStreetEffect,
		streetID,
	)
}

func NewMioSelectBuildingEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return newEvent(
		entityID,
		version,
		MioSelectBuildingEffect,
		buildingID,
	)
}

func NewMioUnselectBuildingEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return newEvent(
		entityID,
		version,
		MioUnselectBuildingEffect,
		buildingID,
	)
}

func NewMioEnterBuildingEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return newEvent(
		entityID,
		version,
		MioEnterBuildingEffect,
		buildingID,
	)
}

func NewMioLeaveBuildingEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return newEvent(
		entityID,
		version,
		MioLeaveBuildingEffect,
		buildingID,
	)
}

func NewMioActEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return newEvent(
		entityID,
		version,
		MioActEffect,
		buildingID,
	)
}

func NewMioStreamEvent(entityID uuid.UUID, value, version int) Event {
	return newEvent(
		entityID,
		version,
		MioStreamEffect,
		value,
	)
}

func NewMioEatEvent(entityID uuid.UUID, value, version int) Event {
	return newEvent(
		entityID,
		version,
		MioEatEffect,
		value,
	)
}

func NewMioStarveEvent(entityID uuid.UUID, value, version int) Event {
	return newEvent(
		entityID,
		version,
		MioStarveEffect,
		value,
	)
}

func NewMioDrinkEvent(entityID uuid.UUID, value, version int) Event {
	return newEvent(
		entityID,
		version,
		MioDrinkEffect,
		value,
	)
}

func NewMioSweatEvent(entityID uuid.UUID, value, version int) Event {
	return newEvent(
		entityID,
		version,
		MioSweatEffect,
		value,
	)
}

func NewMioChangePlannedPoses(entityID uuid.UUID, value []math.Pos, version int) Event {
	return newEvent(
		entityID,
		version,
		MioChangePlannedPoses,
		value,
	)
}
