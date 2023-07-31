package animator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type sampleImpl struct {
	animatorImpl *animatorImpl

	screen      *ebiten.Image
	sampleStore *event.SampleStore
}

func NewSampleAnimator(store *event.SampleStore) Animator {
	animateEventSet := map[event.Effect]func(event event.Event) []frame{
		event.SampleEffect: func(e event.Event) []frame {
			return nil
		},
	}

	return sampleImpl{
		sampleStore: store,

		animatorImpl: &animatorImpl{
			framesToRender:    []map[uuid.UUID][]frame{},
			getEventFramesSet: animateEventSet,
		},
	}
}

func (i sampleImpl) Animator() *animatorImpl {
	return i.animatorImpl
}
