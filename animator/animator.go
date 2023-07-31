package animator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/system"
	"github.com/google/uuid"
)

type animatorImpl struct {
	store *event.Store

	defaultRenderLayer RenderLayer

	counterSet map[uuid.UUID]int

	framesToRender []map[uuid.UUID][]frame

	getEventFramesSet map[event.Effect]func(event event.Event) []frame
	getIdleFramesFunc func(entityIDsWithEvent []uuid.UUID) []frame

	subAnimators []Animator
}

func newAnimatorImpl() *animatorImpl {
	return &animatorImpl{
		counterSet: map[uuid.UUID]int{},

		framesToRender: make([]map[uuid.UUID][]frame, system.DEFAULT_FPS),
	}
}

type Animator interface {
	Animator() *animatorImpl
}

func (a *animatorImpl) pullNewEventFromStore() {
	totalCounter := 0
	for _, counter := range a.counterSet {
		totalCounter += counter
	}

	if totalCounter != (*a.store).GetCounter() {
		for id, events := range (*a.store).GetEvents() {
			index, isExist := a.counterSet[id]
			if isExist {
				events = events[index:]
			}

			for _, event := range events {
				a.processEvent(event)
			}
			a.counterSet[id] = index + len(events)
		}
	}

}

func (a *animatorImpl) processEvent(event event.Event) {
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
		subAnimator.Animator().processEvent(event)
	}
}

func (a *animatorImpl) GetFrames() []frame {
	a.pullNewEventFromStore()

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

	for i, _frame := range frames {
		if _frame.RenderLayer.String() == "" {
			_frame.RenderLayer = a.defaultRenderLayer
			frames[i] = _frame
		}
	}

	for _, subAnimator := range a.subAnimators {
		frames = append(frames, subAnimator.Animator().GetFrames()...)
	}

	return frames
}
