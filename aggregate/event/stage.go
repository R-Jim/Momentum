package event

import "github.com/google/uuid"

const (
	StageInitEffect Effect = "STAGE_INIT"
)

type StageStore Store

func NewStageStore() StageStore {
	return StageStore(newStore())
}

func (s StageStore) NewStageInitEvent(entityID uuid.UUID, goalProductIDs []uuid.UUID) Event {
	return Store(s).newEvent(entityID, StageInitEffect, goalProductIDs)
}
