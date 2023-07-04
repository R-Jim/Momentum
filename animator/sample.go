package animator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type sampleImpl struct {
	animatorImpl *animatorImpl

	screen      *ebiten.Image
	sampleStore *store.Store
}

func NewSampleAnimator(store *store.Store) Animator {
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
