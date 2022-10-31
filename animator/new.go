package animator

import (
	"github.com/R-jim/Momentum/domain/jet"
	"github.com/hajimehoshi/ebiten/v2"
)

type Animator interface {
	Draw(screen *ebiten.Image)
	AppendEvent(event interface{}) error
}

type impl struct {
	JetAnimator JetAnimator
}

func New(jetStore jet.Store) Animator {
	initJetAnimation()
	return impl{
		JetAnimator: JetAnimator{
			store:         jetStore,
			pendingEvents: map[string][]jet.Event{},
		},
	}
}

func (i impl) Draw(screen *ebiten.Image) {
	i.JetAnimator.Draw(screen)
}

func (i impl) AppendEvent(event interface{}) error {
	{
		jetEvent, ok := event.(jet.Event)
		if ok {
			i.JetAnimator.AppendEvent(jetEvent)
			return nil
		}
	}

	return ErrCanNotAppendUnsupportedEvent
}
