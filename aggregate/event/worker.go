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

type WorkerStore Store

func NewWorkerStore() WorkerStore {
	return WorkerStore(newStore())
}

func (s WorkerStore) NewWorkerInitEvent(entityID uuid.UUID, position math.Pos) Event {
	return Store(s).newEvent(entityID, WorkerInitEffect, position)
}

func (s WorkerStore) NewWorkerAssignEvent(entityID uuid.UUID, buildingID uuid.UUID) Event {
	return Store(s).newEvent(entityID, WorkerAssignEffect, buildingID)
}

func (s WorkerStore) NewWorkerUnassignEvent(entityID uuid.UUID, buildingID uuid.UUID) Event {
	return Store(s).newEvent(entityID, WorkerUnassignEffect, buildingID)
}

func (s WorkerStore) NewWorkerActEvent(entityID uuid.UUID, buildingID uuid.UUID) Event {
	return Store(s).newEvent(entityID, WorkerActEffect, buildingID)
}

func (s WorkerStore) NewWorkerMoveEvent(entityID uuid.UUID, position math.Pos) Event {
	return Store(s).newEvent(entityID, WorkerMoveEffect, position)
}

func (s WorkerStore) NewWorkerChangePlannedPoses(entityID uuid.UUID, value []math.Pos) Event {
	return Store(s).newEvent(
		entityID,
		WorkerChangePlannedPoses,
		value,
	)
}
