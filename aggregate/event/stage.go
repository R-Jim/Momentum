package event

import "github.com/google/uuid"

const (
	StageInitEffect Effect = "STAGE_INIT"
)

func NewStageInitEvent(entityID uuid.UUID, goalProductIDs []uuid.UUID) Event {
	return newEvent(entityID, 1, StageInitEffect, goalProductIDs)
}
