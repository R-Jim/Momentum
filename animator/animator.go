package animator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/system"
	"github.com/google/uuid"
)

type animatorImpl struct {
	framesToRender []map[uuid.UUID][]frame

	getEventFramesSet map[event.Effect]func(event event.Event) []frame
	getIdleFramesFunc func(entityIDsWithEvent []uuid.UUID) []frame

	subAnimators []Animator
}

func newAnimatorImpl() *animatorImpl {
	return &animatorImpl{
		framesToRender: make([]map[uuid.UUID][]frame, system.DEFAULT_FPS),
	}
}

type Animator interface {
	Animator() *animatorImpl
}

func (a *animatorImpl) ProcessEvent(event event.Event) {
	getFramesFunc := a.getEventFramesSet[event.Effect]
	if getFramesFunc != nil {
		for index, _frame := range getFramesFunc(event) {
			if len(a.framesToRender) == index {
				a.framesToRender = append(
					a.framesToRender,
					map[uuid.UUID][]frame{event.EntityID: {_frame}})
			} else {
				a.framesToRender[index][event.EntityID] = append(a.framesToRender[index][event.EntityID], _frame)
			}
		}
	}

	for _, subAnimator := range a.subAnimators {
		subAnimator.Animator().ProcessEvent(event)
	}
}

func (a *animatorImpl) GetFrames() []frame {
	frames := []frame{}
	entityIDsWithEvent := []uuid.UUID{}

	if len(a.framesToRender) > 0 {
		for entityID, framesByEntityID := range a.framesToRender[0] {
			frames = append(frames, framesByEntityID...)

			entityIDsWithEvent = append(entityIDsWithEvent, entityID)
		}
		a.framesToRender = a.framesToRender[1:]
	}

	if a.getIdleFramesFunc != nil {
		frames = append(frames, a.getIdleFramesFunc(entityIDsWithEvent)...)
	}

	for _, subAnimator := range a.subAnimators {
		frames = append(frames, subAnimator.Animator().GetFrames()...)
	}

	return frames
}
