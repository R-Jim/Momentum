package animator

import (
	_ "image/png"
	"log"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/asset"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/system"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WALK_STEP = 1
)

type mioImpl struct {
	animatorImpl *animatorImpl

	mioStore *store.Store

	mioAsset mioAsset
}

type mioAsset struct {
	idleSpriteSheet spriteSheet
	walkSpriteSheet spriteSheet
	runSpriteSheet  spriteSheet

	actSpriteSheet spriteSheet
}

func NewMioAnimator(store *store.Store) Animator {
	mio := mioImpl{
		mioStore: store,

		mioAsset: mioAsset{
			idleSpriteSheet: newSpriteSheet(asset.Idle_24x24, asset.Idle_24x24_Size),
			walkSpriteSheet: newSpriteSheet(asset.Walk_24x24, asset.Walk_24x24_Size),
			runSpriteSheet:  newSpriteSheet(asset.Walk_24x24, asset.Walk_24x24_Size),

			actSpriteSheet: newSpriteSheet(asset.Act_24x24, asset.Act_24x24_Size),
		},
	}

	animateEventSet := map[event.Effect]func(event event.Event) []frame{
		event.MioWalkEffect: func(e event.Event) []frame {
			mioWalkSpriteSheet := mio.mioAsset.walkSpriteSheet

			mioWalkSprites := mioWalkSpriteSheet.sprites

			destinationPos, err := event.ParseData[math.Pos](e)
			if err != nil {
				return []frame{}
			}
			{
				spriteMultiplier := int(system.AUTOMATION_TICK_PER_FPS / len(mioWalkSprites))
				fillOutSprites := []*ebiten.Image{}
				for _, sprite := range mioWalkSprites {
					for i := 0; i < spriteMultiplier; i++ {
						fillOutSprites = append(fillOutSprites, sprite)
					}
				}
				mioWalkSprites = fillOutSprites
			}

			frames := []frame{}
			currentPos := mio.getMioPos(e.EntityID)
			desireNumberOfFrames := len(mioWalkSprites)

			for i := 0; i < len(mioWalkSprites); i++ {
				currentPos = math.NewPos(math.GetNextStepXY(currentPos, 0, destinationPos, 0, aggregator.MAX_WALK_DISTANT/float64(desireNumberOfFrames), 180))
				image := ebiten.NewImageFromImage(mioWalkSprites[i])
				frames = append(frames, frame{
					Image:  image,
					Option: getCenteredDrawImageOptions(mioWalkSprites[i], currentPos),
				})
			}

			return frames
		},
		event.MioRunEffect: func(e event.Event) []frame {
			mioRunSpriteSheet := mio.mioAsset.runSpriteSheet

			mioRunSprites := mioRunSpriteSheet.sprites

			destinationPos, err := event.ParseData[math.Pos](e)
			if err != nil {
				return []frame{}
			}
			{
				spriteMultiplier := int(system.AUTOMATION_TICK_PER_FPS / len(mioRunSprites))
				fillOutSprites := []*ebiten.Image{}
				for _, sprite := range mioRunSprites {
					for i := 0; i < spriteMultiplier; i++ {
						fillOutSprites = append(fillOutSprites, sprite)
					}
				}
				mioRunSprites = fillOutSprites
			}

			frames := []frame{}
			currentPos := mio.getMioPos(e.EntityID)
			desireNumberOfFrames := len(mioRunSprites)

			for i := 0; i < len(mioRunSprites); i++ {
				currentPos = math.NewPos(math.GetNextStepXY(currentPos, 0, destinationPos, 0, aggregator.MAX_RUN_DISTANT/float64(desireNumberOfFrames), 180))
				image := ebiten.NewImageFromImage(mioRunSprites[i])
				frames = append(frames, frame{
					Image:  image,
					Option: getCenteredDrawImageOptions(mioRunSprites[i], currentPos),
				})
			}

			return frames
		},
		event.MioActEffect: func(e event.Event) []frame {
			mioActSpriteSheet := mio.mioAsset.actSpriteSheet

			mioActSprites := mioActSpriteSheet.sprites
			frames := []frame{}
			currentPos := mio.getMioPos(e.EntityID)

			for i := 0; i < len(mioActSprites); i++ {
				image := ebiten.NewImageFromImage(mioActSprites[i])
				frames = append(frames, frame{
					Image:  image,
					Option: getCenteredDrawImageOptions(mioActSprites[i], currentPos),
				})
			}

			return frames
		},
	}

	mio.animatorImpl = newAnimatorImpl()
	mio.animatorImpl.getEventFramesSet = animateEventSet
	mio.animatorImpl.getIdleFramesFunc = mio.getMioIdleFrames
	return mio
}

func (i mioImpl) Animator() *animatorImpl {
	return i.animatorImpl
}

func (i mioImpl) getMioIdleFrames(entityIDsWithFrame []uuid.UUID) []frame {
	entityWithoutEvent := map[uuid.UUID][]event.Event{}
	for entityID, events := range (*i.mioStore).GetEvents() {
		isHasFrame := false
		for _, entityIDWithFrame := range entityIDsWithFrame {
			if entityID == entityIDWithFrame {
				isHasFrame = true
				continue
			}
		}

		if !isHasFrame {
			entityWithoutEvent[entityID] = events
		}
	}

	result := []frame{}
	for _, events := range entityWithoutEvent {
		mioState, _ := aggregator.GetMioState(events)
		currentPos := mioState.Position
		image := ebiten.NewImageFromImage(i.mioAsset.idleSpriteSheet.sprites[0])
		result = append(result, frame{
			Image:  image,
			Option: getCenteredDrawImageOptions(image, currentPos),
		})
	}

	return result
}

func (i mioImpl) getMioPos(id uuid.UUID) math.Pos {
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
