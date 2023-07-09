package animator

import (
	_ "image/png"
	"log"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/asset"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type mioEffectImpl struct {
	animatorImpl *animatorImpl

	mioStore *store.Store

	mioEffectAsset mioEffectAsset
}

type mioEffectAsset struct {
	effectUpSpriteSheet spriteSheet
}

func NewMioEffectAnimator(store *store.Store) Animator {
	mio := mioEffectImpl{
		mioStore: store,

		mioEffectAsset: mioEffectAsset{
			effectUpSpriteSheet: newSpriteSheet(asset.UpEffect_24x24, asset.UpEffect_24x24_Size),
		},
	}

	animateEventSet := map[event.Effect]func(event event.Event) []frame{
		event.MioEatEffect:   mio.getEffectUpFrames,
		event.MioDrinkEffect: mio.getEffectUpFrames,
	}

	mio.animatorImpl = newAnimatorImpl()
	mio.animatorImpl.getEventFramesSet = animateEventSet
	mio.animatorImpl.defaultRenderLayer = EffectRenderLayer

	return mio
}

func (i mioEffectImpl) getEffectUpFrames(e event.Event) []frame {
	mioEffectUpSpriteSheet := i.mioEffectAsset.effectUpSpriteSheet

	mioEffectUpSprites := mioEffectUpSpriteSheet.sprites
	frames := []frame{}
	currentPos := i.getMioPos(e.EntityID)

	for i := 0; i < len(mioEffectUpSprites); i++ {
		image := ebiten.NewImageFromImage(mioEffectUpSprites[i])
		frames = append(frames, frame{
			Image:  image,
			Option: getCenteredDrawImageOptions(mioEffectUpSprites[i], currentPos),
		})
	}

	return frames
}

func (i mioEffectImpl) Animator() *animatorImpl {
	return i.animatorImpl
}

func (i mioEffectImpl) getMioPos(id uuid.UUID) math.Pos {
	events, err := (*i.mioStore).GetEventsByEntityID(id)
	if err != nil {
		log.Fatalln(err)
	}

	mioState, err := aggregator.GetMioState(events)
	if err != nil {
		log.Fatalln(err)
	}

	return mioState.Position
}
