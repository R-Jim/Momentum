package animator

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/R-jim/Momentum/aggregate/spike"
	"github.com/R-jim/Momentum/asset"
	"github.com/hajimehoshi/ebiten/v2"
)

/*
TODO:
enemy has HP
number of enemies to render, radius base on HP
*/

var (
	normalImage *ebiten.Image
	duckImage   *ebiten.Image
)

type SpikeAnimator struct {
	store         spike.Store
	pendingEvents map[string][]spike.Event
}

func initSpikeAnimation() error {
	normalImage = getAssetImage(asset.Spike_png)
	duckImage = getAssetImage(asset.SpikeDuck_png)
	return nil
}

func (sa SpikeAnimator) AppendEvent(event spike.Event) {
	events := sa.pendingEvents[event.ID]
	if events == nil {
		events = []spike.Event{}
	}
	sa.pendingEvents[event.ID] = append(events, event)
}

func (sa SpikeAnimator) resetEventQueue(id string) {
	sa.pendingEvents[id] = []spike.Event{}
}

func (sa SpikeAnimator) animateEvent(screen *ebiten.Image, id string) error {
	if len(sa.pendingEvents) == 0 {
		return ErrNoPendingEvents
	}
	var effectImage *ebiten.Image
	var deathImage *ebiten.Image

	for _, event := range sa.pendingEvents[id] {
		switch event.Effect {
		case spike.MoveEffect:
		case spike.DamageEffect:
			effectImage = HitEffectImage
			currentState, err := spike.GetState(sa.store, id)
			if err != nil {
				return err
			}
			damageValue, _ := event.Data.(int)
			if currentState.Health.Value > 0 && currentState.Health.Value-damageValue <= 0 {
				deathImage = DeathEffectImage
			}

		case spike.StrikeEffect:

		default:
			return fmt.Errorf("[SpikeAnimator][ERROR][%v] err: %v", event.Effect, ErrEffectNotSupported.Error())
		}
		positionState, _ := spike.GetPositionState(sa.store, id)

		if effectImage != nil {
			centerAndRenderImage(screen, effectImage, positionState.X, positionState.Y)
		}
		if deathImage != nil {
			centerAndRenderImage(screen, deathImage, positionState.X, positionState.Y)
		}
	}
	sa.resetEventQueue(id)

	return nil
}

func (sa SpikeAnimator) animateState(screen *ebiten.Image, id string) error {
	// currentState, _ := spike.GetState(ja.store, id)
	var stateImage *ebiten.Image

	// switch currentState.Status {
	// case jet.FlyingStatus:
	// case jet.LandedStatus:
	// 	stateImage = landingImage
	// case jet.EngagingStatus:
	// 	stateImage = engagingImage
	// default:
	// 	return fmt.Errorf("[JetAnimator][ERROR][%v] err: %v", currentState.Status, ErrStateNotSupported.Error())
	// }

	positionState, _ := spike.GetPositionState(sa.store, id)

	frameIndex := rand.Intn(2)
	switch frameIndex {
	case 1:
		stateImage = duckImage
	default:
		stateImage = normalImage
	}

	combatState, _ := spike.GetState(sa.store, id)
	if combatState.Health.Value <= 0 {
		return nil
	}

	if stateImage != nil {
		centerAndRenderImage(screen, stateImage, positionState.X, positionState.Y)
	}
	return nil
}

func (sa SpikeAnimator) drawSpike(screen *ebiten.Image, id string) {
	if err := sa.animateEvent(screen, id); err != nil {
		if !errors.Is(err, ErrNoPendingEvents) {
			fmt.Println(err)
		}
	}
	sa.animateState(screen, id)
}

func (sa SpikeAnimator) Draw(screen *ebiten.Image) {
	if sa.store == nil {
		return
	}

	for id := range sa.store.GetEvents() {
		sa.drawSpike(screen, id)
	}
}
