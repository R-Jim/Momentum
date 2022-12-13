package animator

import (
	"errors"
	"fmt"

	"github.com/R-jim/Momentum/aggregate/artifact"
	"github.com/R-jim/Momentum/asset"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	idleImage *ebiten.Image
)

type ArtifactAnimator struct {
	store         artifact.Store
	pendingEvents map[string][]artifact.Event
}

func initArtifactAnimation() error {
	idleImage = getAssetImage(asset.Artifact_png)
	return nil
}

func (aa ArtifactAnimator) AppendEvent(event artifact.Event) {
	events := aa.pendingEvents[event.ID]
	if events == nil {
		events = []artifact.Event{}
	}
	aa.pendingEvents[event.ID] = append(events, event)
}

func (aa ArtifactAnimator) resetEventQueue(id string) {
	aa.pendingEvents[id] = []artifact.Event{}
}

func (aa ArtifactAnimator) animateEvent(screen *ebiten.Image, id string) error {
	// if len(ja.pendingEvents) == 0 {
	// 	return ErrNoPendingEvents
	// }
	// var stateImage *ebiten.Image

	// for _, event := range ja.pendingEvents[id] {
	// 	switch event.Effect {
	// 	case jet.FlyEffect:
	// 		stateImage = flyingImage

	// 	default:
	// 		return fmt.Errorf("[JetAnimator][ERROR][%v] err: %v", event.Effect, ErrEffectNotSupported.Error())
	// 	}
	// }
	// ja.resetEventQueue(id)

	// positionState, _ := jet.GetPositionState(ja.store, id)

	// if stateImage != nil {
	// 	op := &ebiten.DrawImageOptions{}
	// 	op.GeoM.Translate(positionState.X+float64(jetImage.Bounds().Dx()+5), positionState.Y)
	// 	screen.DrawImage(stateImage, op)
	// }

	return nil
}

func (aa ArtifactAnimator) animateState(screen *ebiten.Image, id string) error {
	// currentState, _ := spike.GetState(ja.store, id)
	var stateImage *ebiten.Image
	stateImage = idleImage
	// switch currentState.Status {
	// case jet.FlyingStatus:
	// case jet.LandedStatus:
	// 	stateImage = landingImage
	// case jet.EngagingStatus:
	// 	stateImage = engagingImage
	// default:
	// 	return fmt.Errorf("[JetAnimator][ERROR][%v] err: %v", currentState.Status, ErrStateNotSupported.Error())
	// }

	positionState, _ := artifact.GetPositionState(aa.store, id)

	if stateImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(positionState.X, positionState.Y)
		screen.DrawImage(stateImage, op)
	}
	return nil
}

func (aa ArtifactAnimator) drawArtifact(screen *ebiten.Image, id string) {
	if err := aa.animateEvent(screen, id); err != nil {
		if !errors.Is(err, ErrNoPendingEvents) {
			fmt.Println(err)
		}
	}
	aa.animateState(screen, id)
}

func (aa ArtifactAnimator) Draw(screen *ebiten.Image) {
	if aa.store == nil {
		return
	}

	for id := range aa.store.GetEvents() {
		aa.drawArtifact(screen, id)
	}
}
