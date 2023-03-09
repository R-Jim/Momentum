package animator

import (
	"github.com/R-jim/Momentum/aggregate/event"
)

type Animator interface {
	GetAnimateSet() map[event.Effect]func(event.Event) error
}

func Draw(animateSet map[event.Effect]func(event.Event) error, event event.Event) error {
	animateFunc, isExist := animateSet[event.Effect]
	if !isExist {
		return ErrEffectNotSupported
	}
	return animateFunc(event)
}
