package animator

import (
	"github.com/R-jim/Momentum/aggregate/jet"
	"github.com/R-jim/Momentum/aggregate/spike"
	"github.com/hajimehoshi/ebiten/v2"
)

type Animator interface {
	Draw(screen *ebiten.Image)
	AppendEvent(event interface{}) error
}

type impl struct {
	JetAnimator   JetAnimator
	SpikeAnimator SpikeAnimator
}

type AnimatorStores struct {
	JetStore   jet.Store
	SpikeStore spike.Store
}

func New(stores AnimatorStores) Animator {
	initJetAnimation()
	initSpikeAnimation()
	return impl{
		JetAnimator: JetAnimator{
			store:         stores.JetStore,
			pendingEvents: map[string][]jet.Event{},
		},
		SpikeAnimator: SpikeAnimator{
			store:         stores.SpikeStore,
			pendingEvents: map[string][]spike.Event{},
		},
	}
}

func (i impl) Draw(screen *ebiten.Image) {
	i.JetAnimator.Draw(screen)
	i.SpikeAnimator.Draw(screen)
}

func (i impl) AppendEvent(event interface{}) error {
	{
		jetEvent, ok := event.(jet.Event)
		if ok {
			i.JetAnimator.AppendEvent(jetEvent)
			return nil
		}
	}
	{
		spikeEvent, ok := event.(spike.Event)
		if ok {
			i.SpikeAnimator.AppendEvent(spikeEvent)
			return nil
		}
	}

	return ErrCanNotAppendUnsupportedEvent
}
