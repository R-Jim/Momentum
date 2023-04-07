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

func NewBuildingInitEvent(buildingID uuid.UUID, pos math.Pos) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: buildingID,
		Version:  1,
		Effect:   BuildingInitEffect,
		Data:     pos,
	}
}

func NewBuildingEntityEnterEvent(buildingID uuid.UUID, version int, entityID uuid.UUID) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: buildingID,
		Version:  version,
		Effect:   BuildingEntityEnterEffect,
		Data:     entityID,
	}
}

func NewBuildingEntityLeaveEvent(buildingID uuid.UUID, version int, entityID uuid.UUID) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: buildingID,
		Version:  version,
		Effect:   BuildingEntityLeaveEffect,
		Data:     entityID,
	}
}

func NewBuildingEntityActEvent(buildingID uuid.UUID, entityID uuid.UUID, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: buildingID,
		Version:  version,
		Effect:   BuildingEntityActEffect,
		Data:     entityID,
	}
}

func NewBuildingWorkerActEvent(buildingID uuid.UUID, workerID uuid.UUID, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: buildingID,
		Version:  version,
		Effect:   BuildingWorkerActEffect,
		Data:     workerID,
	}
}

func NewBuildingWorkerAssignEvent(buildingID uuid.UUID, workerID uuid.UUID, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: buildingID,
		Version:  version,
		Effect:   BuildingWorkerAssignEffect,
		Data:     workerID,
	}
}

func NewBuildingWorkerUnassignEvent(buildingID uuid.UUID, workerID uuid.UUID, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: buildingID,
		Version:  version,
		Effect:   BuildingWorkerUnassignEffect,
		Data:     workerID,
	}
}
