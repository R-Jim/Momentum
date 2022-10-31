package animator

import (
	"errors"
	"fmt"

	"github.com/R-jim/Momentum/asset"
	"github.com/R-jim/Momentum/domain/jet"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	flyingImage   *ebiten.Image
	landingImage  *ebiten.Image
	engagingImage *ebiten.Image

	jetImage *ebiten.Image
)

type JetAnimator struct {
	store         jet.Store
	pendingEvents map[string][]jet.Event
}

func initJetAnimation() error {
	jetImage = getAssetImage(asset.Jet_png)

	flyingImage = getAssetImage(asset.Flying_png)
	landingImage = getAssetImage(asset.Landing_png)
	engagingImage = getAssetImage(asset.Engaging_png)

	return nil
}

func (ja JetAnimator) AppendEvent(event jet.Event) {
	events := ja.pendingEvents[event.ID]
	if events == nil {
		events = []jet.Event{}
	}
	ja.pendingEvents[event.ID] = append(events, event)
}

func (ja JetAnimator) resetEventQueue(id string) {
	ja.pendingEvents[id] = []jet.Event{}
}

func (ja JetAnimator) animateEvent(screen *ebiten.Image, id string) error {
	if len(ja.pendingEvents) == 0 {
		return ErrNoPendingEvents
	}
	var stateImage *ebiten.Image

	for _, event := range ja.pendingEvents[id] {
		switch event.Effect {
		case jet.FlyEffect:
			stateImage = flyingImage

		default:
			return fmt.Errorf("[JetAnimator][ERROR][%v] err: %v", event.Effect, ErrEffectNotSupported.Error())
		}
	}
	ja.resetEventQueue(id)

	positionState, _ := jet.GetPositionState(ja.store, id)

	if stateImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(positionState.X+jetImage.Bounds().Dx()+5), float64(positionState.Y))
		screen.DrawImage(stateImage, op)
	}

	return nil
}

func (ja JetAnimator) animateState(screen *ebiten.Image, id string) error {
	currentState, _ := jet.GetCombatState(ja.store, id)
	var stateImage *ebiten.Image

	switch currentState.Status {
	case jet.FlyingStatus:
	case jet.LandedStatus:
		stateImage = landingImage
	case jet.EngagingStatus:
		stateImage = engagingImage
	default:
		return fmt.Errorf("[JetAnimator][ERROR][%v] err: %v", currentState.Status, ErrStateNotSupported.Error())
	}

	positionState, _ := jet.GetPositionState(ja.store, id)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(positionState.X), float64(positionState.Y))
	screen.DrawImage(jetImage, op)

	if stateImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(positionState.X+jetImage.Bounds().Dx()+5), float64(positionState.Y))
		screen.DrawImage(stateImage, op)
	}
	return nil
}

func (ja JetAnimator) drawJet(screen *ebiten.Image, id string) {
	if err := ja.animateEvent(screen, id); err != nil {
		if !errors.Is(err, ErrNoPendingEvents) {
			fmt.Println(err)
		}
	}
	ja.animateState(screen, id)
}

func (ja JetAnimator) Draw(screen *ebiten.Image) {
	for id := range ja.store.GetEvents() {
		ja.drawJet(screen, id)
	}
}
