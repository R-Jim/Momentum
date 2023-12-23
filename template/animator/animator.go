package animator

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/system"
	"github.com/R-jim/Momentum/template/event"
)

/*
An animator produces animation Frame based on the newest event after checking cache counter from store.
If no new Frames for the entity, animator can choose to produce Idle Frames for that entity.
*/
type animatorImpl struct {
	store *event.Store

	defaultRenderLayer RenderLayer

	counterSet map[uuid.UUID]int

	framesToRender []map[uuid.UUID][]Frame

	eventFramesSet    map[event.Effect]func(event event.Event) []Frame
	getIdleFramesFunc func(entityIDsWithEvent []uuid.UUID) []Frame
}

func NewAnimator(eventFramesSet map[event.Effect]func(event event.Event) []Frame) Animator {
	return &animatorImpl{
		counterSet: map[uuid.UUID]int{},

		eventFramesSet: eventFramesSet,
		framesToRender: make([]map[uuid.UUID][]Frame, system.DEFAULT_FPS),
	}

}

type Animator interface {
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
	getFramesFunc := a.eventFramesSet[event.Effect]
	if getFramesFunc != nil {
		for index, _frame := range getFramesFunc(event) {
			if len(a.framesToRender) == index {
				a.framesToRender = append(
					a.framesToRender,
					map[uuid.UUID][]Frame{event.EntityID: {_frame}})
			} else {
				a.framesToRender[index][event.EntityID] = append(a.framesToRender[index][event.EntityID], _frame)
			}
		}
	}
}

func (a *animatorImpl) GetFrames() []Frame {
	a.pullNewEventFromStore()

	frames := []Frame{}
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

	return frames
}
