package event

import (
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

const (
	BuildingInitEffect Effect = "BUILDING_INIT"

	BuildingEntityEnterEffect Effect = "BUILDING_ENTITY_ENTER"
	BuildingEntityLeaveEffect Effect = "BUILDING_ENTITY_LEAVE"

	BuildingEntityActEffect Effect = "BUILDING_ENTITY_ACT"
	BuildingWorkerActEffect Effect = "BUILDING_WORKER_ACT"

	BuildingWorkerAssignEffect   Effect = "BUILDING_WORKER_ASSIGN"
	BuildingWorkerUnassignEffect Effect = "BUILDING_WORKER_UNASSIGN"
)

type BuildingType string

const (
	BuildingTypeMioHouse   BuildingType = "MIO_HOUSE"
	BuildingTypeFoodStore  BuildingType = "FOOD_STORE"
	BuildingTypeDrinkStore BuildingType = "DRINK_STORE"
)

func (t BuildingType) String() string {
	return string(t)
}

type BuildingInitEventData struct {
	Type BuildingType
	Pos  math.Pos
}

type BuildingStore Store

func NewBuildingStore() BuildingStore {
	return BuildingStore(newStore())
}

func (s BuildingStore) NewBuildingInitEvent(buildingID uuid.UUID, buildingType BuildingType, pos math.Pos) Event {
	return (Store(s)).newEvent(
		buildingID,
		BuildingInitEffect,
		BuildingInitEventData{
			Type: buildingType,
			Pos:  pos,
		})
}

func (s BuildingStore) NewBuildingEntityEnterEvent(buildingID uuid.UUID, entityID uuid.UUID) Event {
	return (Store(s)).newEvent(
		buildingID,
		BuildingEntityEnterEffect,
		entityID,
	)
}

func (s BuildingStore) NewBuildingEntityLeaveEvent(buildingID uuid.UUID, entityID uuid.UUID) Event {
	return (Store(s)).newEvent(
		buildingID,
		BuildingEntityLeaveEffect,
		entityID,
	)
}

func (s BuildingStore) NewBuildingEntityActEvent(buildingID uuid.UUID, entityID uuid.UUID) Event {
	return (Store(s)).newEvent(
		buildingID,
		BuildingEntityActEffect,
		entityID,
	)
}

func (s BuildingStore) NewBuildingWorkerActEvent(buildingID uuid.UUID, workerID uuid.UUID) Event {
	return (Store(s)).newEvent(
		buildingID,
		BuildingWorkerActEffect,
		workerID,
	)
}

func (s BuildingStore) NewBuildingWorkerAssignEvent(buildingID uuid.UUID, workerID uuid.UUID) Event {
	return (Store(s)).newEvent(
		buildingID,
		BuildingWorkerAssignEffect,
		workerID,
	)
}

func (s BuildingStore) NewBuildingWorkerUnassignEvent(buildingID uuid.UUID, workerID uuid.UUID) Event {
	return (Store(s)).newEvent(
		buildingID,
		BuildingWorkerUnassignEffect,
		workerID,
	)
}
