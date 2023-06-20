package animator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type Frame struct {
	Image  ebiten.Image
	Option ebiten.DrawImageOptions
}

type AnimatorImpl struct {
	Events         []event.Event
	getEventFrames map[event.Effect]func(event event.Event) []Frame
}

type Animator interface {
	Animator() *AnimatorImpl
}

func (a *AnimatorImpl) AppendEvent(event event.Event) {
	a.Events = append(a.Events, event)
}

func (a *AnimatorImpl) GetFramesPerEvent() map[uuid.UUID][]Frame {
	framesPerEvent := map[uuid.UUID][]Frame{}

	for _, event := range a.Events {
		getFramesFunc, isExist := a.getEventFrames[event.Effect]
		if !isExist {
			continue
		}
		framesPerEvent[event.ID] = getFramesFunc(event)
	}

	a.Events = []event.Event{}
	return framesPerEvent
}
