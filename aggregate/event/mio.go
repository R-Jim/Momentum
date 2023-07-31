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

type MioStore Store

func NewMioStore() MioStore {
	return MioStore(newStore())
}

func (s MioStore) NewMioInitEvent(entityID uuid.UUID, position math.Pos) Event {
	return (Store(s)).newEvent(entityID, MioInitEffect, position)
}

func (s MioStore) NewMioWalkEvent(entityID uuid.UUID, newPosition math.Pos) Event {
	return (Store(s)).newEvent(entityID, MioWalkEffect, newPosition)
}

func (s MioStore) NewMioRunEvent(entityID uuid.UUID, newPosition math.Pos) Event {
	return (Store(s)).newEvent(
		entityID,
		MioRunEffect,
		newPosition,
	)
}

func (s MioStore) NewMioIdleEvent(entityID uuid.UUID) Event {
	return (Store(s)).newEvent(
		entityID,
		MioIdleEffect,
		nil,
	)
}

func (s MioStore) NewMioEnterStreetEvent(entityID uuid.UUID, streetID uuid.UUID) Event {
	return (Store(s)).newEvent(
		entityID,
		MioEnterStreetEffect,
		streetID,
	)
}

func (s MioStore) NewMioSelectBuildingEvent(entityID uuid.UUID, buildingID uuid.UUID) Event {
	return (Store(s)).newEvent(
		entityID,
		MioSelectBuildingEffect,
		buildingID,
	)
}

func (s MioStore) NewMioUnselectBuildingEvent(entityID uuid.UUID, buildingID uuid.UUID) Event {
	return (Store(s)).newEvent(
		entityID,
		MioUnselectBuildingEffect,
		buildingID,
	)
}

func (s MioStore) NewMioEnterBuildingEvent(entityID uuid.UUID, buildingID uuid.UUID) Event {
	return (Store(s)).newEvent(
		entityID,
		MioEnterBuildingEffect,
		buildingID,
	)
}

func (s MioStore) NewMioLeaveBuildingEvent(entityID uuid.UUID, buildingID uuid.UUID) Event {
	return (Store(s)).newEvent(
		entityID,
		MioLeaveBuildingEffect,
		buildingID,
	)
}

func (s MioStore) NewMioActEvent(entityID uuid.UUID, buildingID uuid.UUID) Event {
	return (Store(s)).newEvent(
		entityID,
		MioActEffect,
		buildingID,
	)
}

func (s MioStore) NewMioStreamEvent(entityID uuid.UUID, value int) Event {
	return (Store(s)).newEvent(
		entityID,
		MioStreamEffect,
		value,
	)
}

func (s MioStore) NewMioEatEvent(entityID uuid.UUID, value int) Event {
	return (Store(s)).newEvent(
		entityID,
		MioEatEffect,
		value,
	)
}

func (s MioStore) NewMioStarveEvent(entityID uuid.UUID, value int) Event {
	return (Store(s)).newEvent(
		entityID,
		MioStarveEffect,
		value,
	)
}

func (s MioStore) NewMioDrinkEvent(entityID uuid.UUID, value int) Event {
	return (Store(s)).newEvent(
		entityID,
		MioDrinkEffect,
		value,
	)
}

func (s MioStore) NewMioSweatEvent(entityID uuid.UUID, value int) Event {
	return (Store(s)).newEvent(
		entityID,
		MioSweatEffect,
		value,
	)
}

func (s MioStore) NewMioChangePlannedPoses(entityID uuid.UUID, value []math.Pos) Event {
	return (Store(s)).newEvent(
		entityID,
		MioChangePlannedPoses,
		value,
	)
}
