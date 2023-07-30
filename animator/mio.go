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

	actSpriteSheet  spriteSheet
	workSpriteSheet spriteSheet
}

func NewMioAnimator(_store *store.Store) Animator {
	mio := mioImpl{
		mioStore: _store,

		mioAsset: mioAsset{
			idleSpriteSheet: newSpriteSheet(asset.Idle_24x24, asset.Idle_24x24_Size),
			walkSpriteSheet: newSpriteSheet(asset.Walk_24x24, asset.Walk_24x24_Size),
			runSpriteSheet:  newSpriteSheet(asset.Walk_24x24, asset.Walk_24x24_Size),

			actSpriteSheet:  newSpriteSheet(asset.Act_24x24, asset.Act_24x24_Size),
			workSpriteSheet: newSpriteSheet(asset.Work_24x24, asset.Work_24x24_Size),
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
		event.MioStreamEffect: func(e event.Event) []frame {
			mioWorkSpriteSheet := mio.mioAsset.workSpriteSheet

			mioWorkSprites := mioWorkSpriteSheet.sprites
			frames := []frame{}
			currentPos := mio.getMioPos(e.EntityID)

			for i := 0; i < len(mioWorkSprites); i++ {
				image := ebiten.NewImageFromImage(mioWorkSprites[i])
				frames = append(frames, frame{
					Image:  image,
					Option: getCenteredDrawImageOptions(mioWorkSprites[i], currentPos),
				})
			}

			return frames
		},
	}

	mio.animatorImpl = newAnimatorImpl()
	mio.animatorImpl.getEventFramesSet = animateEventSet
	mio.animatorImpl.getIdleFramesFunc = mio.getMioIdleFrames
	mio.animatorImpl.subAnimators = []Animator{
		NewMioEffectAnimator(mio.mioStore),
	}
	mio.animatorImpl.defaultRenderLayer = MioRenderLayer
	mio.animatorImpl.store = _store

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
