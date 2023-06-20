package animator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/hajimehoshi/ebiten/v2"
)

type sampleImpl struct {
	animatorImpl *AnimatorImpl

	screen      *ebiten.Image
	sampleStore *store.Store
}

func NewSampleAnimator(store *store.Store) Animator {
	animateEventSet := map[event.Effect]func(event event.Event) []Frame{
		event.SampleEffect: func(e event.Event) []Frame {
			return nil
		},
	}

	return sampleImpl{
		sampleStore: store,

		animatorImpl: &AnimatorImpl{
			Events:         []event.Event{},
			getEventFrames: animateEventSet,
		},
	}
}

func (i sampleImpl) Animator() *AnimatorImpl {
	return i.animatorImpl
}
