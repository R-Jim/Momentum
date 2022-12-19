package animator

import (
	"github.com/R-jim/Momentum/aggregate/artifact"
	"github.com/R-jim/Momentum/aggregate/jet"
	"github.com/R-jim/Momentum/aggregate/knight"
	"github.com/R-jim/Momentum/aggregate/spike"
	"github.com/R-jim/Momentum/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Animator interface {
	Draw(screen *ebiten.Image)
	AppendEvent(event interface{}) error

	DrawStatus(ui ui.UI, screen *ebiten.Image, screenWidth, screenHeight float64)
}

type impl struct {
	JetAnimator      JetAnimator
	SpikeAnimator    SpikeAnimator
	ArtifactAnimator ArtifactAnimator
	KnightAnimator   KnightAnimator
}

type AnimatorStores struct {
	JetStore      jet.Store
	SpikeStore    spike.Store
	ArtifactStore artifact.Store
	KnightStore   knight.Store
}

func New(stores AnimatorStores) Animator {
	initCommonImages()
	initJetAnimation()
	initSpikeAnimation()
	initArtifactAnimation()
	initKnightAnimation()
	return impl{
		JetAnimator: JetAnimator{
			store:         stores.JetStore,
			pendingEvents: map[string][]jet.Event{},
		},
		SpikeAnimator: SpikeAnimator{
			store:         stores.SpikeStore,
			pendingEvents: map[string][]spike.Event{},
		},
		ArtifactAnimator: ArtifactAnimator{
			store:         stores.ArtifactStore,
			pendingEvents: map[string][]artifact.Event{},
		},
		KnightAnimator: KnightAnimator{
			store:         stores.KnightStore,
			pendingEvents: map[string][]knight.Event{},
		},
	}
}

func (i impl) Draw(screen *ebiten.Image) {
	i.JetAnimator.Draw(screen)
	i.SpikeAnimator.Draw(screen)
	i.ArtifactAnimator.Draw(screen)
	i.KnightAnimator.Draw(screen)
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
	{
		knightEvent, ok := event.(knight.Event)
		if ok {
			i.KnightAnimator.AppendEvent(knightEvent)
			return nil
		}
	}

	return ErrCanNotAppendUnsupportedEvent
}

func (i impl) DrawStatus(ui ui.UI, screen *ebiten.Image, screenWidth, screenHeight float64) {
	i.KnightAnimator.DrawStatus(ui, screen, screenWidth, screenHeight)
}
