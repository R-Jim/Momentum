package event

import (
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

const (
	WorkerInitEffect Effect = "WORKER_INIT"

	WorkerAssignEffect   Effect = "WORKER_ASSIGN"
	WorkerUnassignEffect Effect = "WORKER_UNASSIGN"

	WorkerActEffect Effect = "WORKER_ACT"

	WorkerMoveEffect Effect = "WORKER_MOVE"

	WorkerChangePlannedPoses Effect = "WORKER_CHANGE_PLANNED_POSITIONS"
)

func NewWorkerInitEvent(entityID uuid.UUID, position math.Pos) Event {
	return newEvent(entityID, 1, WorkerInitEffect, position)
}

func NewWorkerAssignEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return newEvent(entityID, version, WorkerAssignEffect, buildingID)
}

func NewWorkerUnassignEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return newEvent(entityID, version, WorkerUnassignEffect, buildingID)
}

func NewWorkerActEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return newEvent(entityID, version, WorkerActEffect, buildingID)
}

func NewWorkerMoveEvent(entityID uuid.UUID, version int, position math.Pos) Event {
	return newEvent(entityID, version, WorkerMoveEffect, position)
}

func NewWorkerChangePlannedPoses(entityID uuid.UUID, version int, value []math.Pos) Event {
	return newEvent(
		entityID,
		version,
		WorkerChangePlannedPoses,
		value,
	)
}
