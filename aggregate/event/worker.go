package event

import "github.com/google/uuid"

const (
	WorkerInitEffect Effect = "WORKER_INIT"

	WorkerAssignEffect   Effect = "WORKER_ASSIGN"
	WorkerUnassignEffect Effect = "WORKER_UNASSIGN"

	WorkerActEffect Effect = "WORKER_ACT"
)

func NewWorkerInitEvent(entityID uuid.UUID) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  1,
		Effect:   WorkerInitEffect,
	}
}

func NewWorkerAssignEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   WorkerAssignEffect,
		Data:     buildingID,
	}
}

func NewWorkerUnassignEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   WorkerUnassignEffect,
		Data:     buildingID,
	}
}

func NewWorkerActEvent(entityID uuid.UUID, buildingID uuid.UUID, version int) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   WorkerActEffect,
		Data:     buildingID,
	}
}
