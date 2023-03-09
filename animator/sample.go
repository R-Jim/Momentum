package animator

import (
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/hajimehoshi/ebiten/v2"
)

type sampleImpl struct {
	screen      *ebiten.Image
	sampleStore *store.Store
}

func NewSampleAnimator(store *store.Store) Animator {
	return sampleImpl{
		sampleStore: store,
	}
}

func (i sampleImpl) GetAnimateSet() map[event.Effect]func(event.Event) error {
	return map[event.Effect]func(event.Event) error{
		event.SampleEffect: func(e event.Event) error {
			// op := &ebiten.DrawImageOptions{}
			// image := ebiten.NewImage(w, h)
			// i.screen.DrawImage(image, op)
			return nil
		},
	}
}
