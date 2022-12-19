package animator

/*
TODO:
when knight trigger WEAPON_USED event, check for weapon range, is projectile included
*/

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/R-jim/Momentum/aggregate/knight"
	"github.com/R-jim/Momentum/asset"
	"github.com/R-jim/Momentum/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	statusPanelWidth  = 100
	statusPanelHeight = 70
	statusPanelMargin = 1
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
	var effectImage *ebiten.Image

	for _, event := range ja.pendingEvents[id] {
		switch event.Effect {
		case knight.MoveEffect:
		case knight.ChangeTargetEffect:
		case knight.DamageEffect:
			effectImage = HitEffectImage
		case knight.StrikeEffect:
		default:
			return fmt.Errorf("[KnightAnimator][ERROR][%v] err: %v", event.Effect, ErrEffectNotSupported.Error())
		}
	}
	ja.resetEventQueue(id)

	positionState, _ := knight.GetPositionState(ja.store, id)

	if effectImage != nil {
		centerAndRenderImage(screen, effectImage, positionState.X, positionState.Y)
	}

	return nil
}

func (ja KnightAnimator) animateState(screen *ebiten.Image, id string) error {
	// currentState, _ := knight.GetState(ja.store, id)
	// var stateImage *ebiten.Image
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
		// Reset RGB (not Alpha) 0 forcibly
		// op.ColorM.Scale(0, 0, 0, 1)
		// op.ColorM.Translate(float64(0), float64(1), float64(.3), 0)
		centerAndRenderImage(screen, knightImage, positionState.X, positionState.Y)
	}
	// if stateImage != nil {
	// 	centerAndRenderImage(screen, stateImage, positionState.X, positionState.Y)
	// }
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

func (ja KnightAnimator) DrawStatus(ui ui.UI, screen *ebiten.Image, screenWidth, screenHeight float64) {
	statusContainerWidth := len(ja.store.GetEntityIDs())*statusPanelWidth + statusPanelMargin*(len(ja.store.GetEntityIDs())+1)
	statusContainerHeight := statusPanelHeight + statusPanelMargin*2
	statusContainer := ebiten.NewImage(statusContainerWidth, statusContainerHeight)
	statusContainer.Fill(color.Black)

	statusPanel := ebiten.NewImage(statusPanelWidth, statusPanelHeight)
	statusPanel.Fill(color.White)

	for i, id := range ja.store.GetEntityIDs() {
		statusPanel := ebiten.NewImageFromImage(statusPanel)

		op := &ebiten.DrawImageOptions{}
		x := i*statusPanelWidth + (i+1)*statusPanelMargin
		y := statusPanelMargin
		op.GeoM.Translate(float64(x), float64(y))

		combatState, _ := knight.GetState(ja.store, id)
		textPadding := float64(5)
		statuses := []string{
			combatState.ID,
			fmt.Sprintf("HP: %v/%v", combatState.Health.Value, combatState.Health.Max),
		}

		for i, status := range statuses {
			textX := (statusPanelWidth - float64(len(status))*ui.SmallFontSize) / 2
			textY := i*int(ui.SmallFontSize) + i*int(textPadding)
			ui.DrawSmallMsg(statusPanel, []string{status}, int(textX), int(textY), false)

		}
		statusContainer.DrawImage(statusPanel, op)
	}
	centerAndRenderImage(screen, statusContainer, screenWidth/2, float64(screenHeight-statusPanelHeight-10))
}
