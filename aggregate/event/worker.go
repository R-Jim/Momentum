package event

import "github.com/google/uuid"

const (
	WorkerInitEffect Effect = "WORKER_INIT"

	WorkerAssignEffect   Effect = "WORKER_ASSIGN"
	WorkerUnassignEffect Effect = "WORKER_UNASSIGN"

	WorkerActEffect Effect = "WORKER_ACT"
)

func NewWorkerInitEvent(entityID uuid.UUID) Event {
	return newEvent(entityID, 1, WorkerInitEffect, nil)
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
