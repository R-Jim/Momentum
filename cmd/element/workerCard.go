package element

import (
	"github.com/R-jim/Momentum/operator"
	"github.com/R-jim/Momentum/ui"
	"github.com/google/uuid"
)

type WorkerCardClickableElement struct {
	workerOperator operator.WorkerOperator

	workerID uuid.UUID

	ClickAbleImage ui.ClickAbleImage
}

func NewWorkerCardClickableElement(workerOperator operator.WorkerOperator, workerID uuid.UUID, clickAbleImage ui.ClickAbleImage) WorkerCardClickableElement {
	return WorkerCardClickableElement{
		workerOperator,
		workerID,

		clickAbleImage,
	}
}

func (e WorkerCardClickableElement) OnClick(buildingID uuid.UUID) {
	e.workerOperator.AssignBuilding(e.workerID, buildingID)
}
