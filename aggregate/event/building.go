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

func NewBuildingInitEvent(buildingID uuid.UUID, buildingType BuildingType, pos math.Pos) Event {
	return newEvent(buildingID,
		1,
		BuildingInitEffect,
		BuildingInitEventData{
			Type: buildingType,
			Pos:  pos,
		})
}

func NewBuildingEntityEnterEvent(buildingID uuid.UUID, version int, entityID uuid.UUID) Event {
	return newEvent(buildingID,
		version,
		BuildingEntityEnterEffect,
		entityID,
	)
}

func NewBuildingEntityLeaveEvent(buildingID uuid.UUID, version int, entityID uuid.UUID) Event {
	return newEvent(buildingID,
		version,
		BuildingEntityLeaveEffect,
		entityID,
	)
}

func NewBuildingEntityActEvent(buildingID uuid.UUID, entityID uuid.UUID, version int) Event {
	return newEvent(buildingID,
		version,
		BuildingEntityActEffect,
		entityID,
	)
}

func NewBuildingWorkerActEvent(buildingID uuid.UUID, workerID uuid.UUID, version int) Event {
	return newEvent(buildingID,
		version,
		BuildingWorkerActEffect,
		workerID,
	)
}

func NewBuildingWorkerAssignEvent(buildingID uuid.UUID, workerID uuid.UUID, version int) Event {
	return newEvent(buildingID,
		version,
		BuildingWorkerAssignEffect,
		workerID,
	)
}

func NewBuildingWorkerUnassignEvent(buildingID uuid.UUID, workerID uuid.UUID, version int) Event {
	return newEvent(buildingID,
		version,
		BuildingWorkerUnassignEffect,
		workerID,
	)
}
