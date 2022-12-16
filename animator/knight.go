package animator

/*
TODO:
when knight trigger WEAPON_USED event, check for weapon range, is projectile included
*/

import (
	"errors"
	"fmt"

	"github.com/R-jim/Momentum/aggregate/knight"
	"github.com/R-jim/Momentum/asset"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	knightImage *ebiten.Image
)

type KnightAnimator struct {
	store         knight.Store
	pendingEvents map[string][]knight.Event
}

func initKnightAnimation() error {
	knightImage = getAssetImage(asset.MioFa_png)

	return nil
}

func (ja KnightAnimator) AppendEvent(event knight.Event) {
	events := ja.pendingEvents[event.ID]
	if events == nil {
		events = []knight.Event{}
	}
	ja.pendingEvents[event.ID] = append(events, event)
}

func (ja KnightAnimator) resetEventQueue(id string) {
	ja.pendingEvents[id] = []knight.Event{}
}

func (ja KnightAnimator) animateEvent(screen *ebiten.Image, id string) error {
	if len(ja.pendingEvents) == 0 {
		return ErrNoPendingEvents
	}
	var stateImage *ebiten.Image

	for _, event := range ja.pendingEvents[id] {
		switch event.Effect {
		case knight.MoveEffect:
			stateImage = flyingImage
		case knight.ChangeTargetEffect:
		default:
			return fmt.Errorf("[KnightAnimator][ERROR][%v] err: %v", event.Effect, ErrEffectNotSupported.Error())
		}
	}
	ja.resetEventQueue(id)

	positionState, _ := knight.GetPositionState(ja.store, id)

	if stateImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(positionState.X+float64(knightImage.Bounds().Dx()+5), positionState.Y)
		screen.DrawImage(stateImage, op)
	}

	return nil
}

func (ja KnightAnimator) animateState(screen *ebiten.Image, id string) error {
	// currentState, _ := knight.GetState(ja.store, id)
	var stateImage *ebiten.Image
	// switch currentState.Status {
	// // case knight.FlyingStatus:
	// // case knight.LandedStatus:
	// // 	stateImage = landingImage
	// // case knight.EngagingStatus:
	// // 	stateImage = engagingImage
	// default:
	// 	return fmt.Errorf("[KnightAnimator][ERROR][%v] err: %v", currentState.Status, ErrStateNotSupported.Error())
	// }

	positionState, _ := knight.GetPositionState(ja.store, id)
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(positionState.X, positionState.Y)

		// Reset RGB (not Alpha) 0 forcibly
		// op.ColorM.Scale(0, 0, 0, 1)
		// op.ColorM.Translate(float64(0), float64(1), float64(.3), 0)

		ji := ebiten.NewImageFromImage(knightImage)
		screen.DrawImage(ji, op)
	}
	if stateImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(positionState.X+float64(knightImage.Bounds().Dx()+5), positionState.Y)
		screen.DrawImage(stateImage, op)
	}
	return nil
}

func (ja KnightAnimator) drawKnight(screen *ebiten.Image, id string) {
	if err := ja.animateEvent(screen, id); err != nil {
		if !errors.Is(err, ErrNoPendingEvents) {
			fmt.Println(err)
		}
	}
	ja.animateState(screen, id)
}

func (ja KnightAnimator) Draw(screen *ebiten.Image) {
	for id := range ja.store.GetEvents() {
		ja.drawKnight(screen, id)
	}
}
