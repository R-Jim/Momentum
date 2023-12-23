package example

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/R-jim/Momentum/template/animator"
	"github.com/R-jim/Momentum/template/event"
)

func NewSampleAnimator(screen *ebiten.Image, store *event.Store) animator.Animator {
	animateEventSet := map[event.Effect]func(event event.Event) []animator.Frame{
		SampleMoveEffect: func(e event.Event) []animator.Frame {
			return nil
		},
	}

	return animator.NewAnimator(animateEventSet)
}
