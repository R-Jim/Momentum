package animator

import (
	_ "image/png"
	"log"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/asset"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WORKER_STEP = 2
)

type workerImpl struct {
	animatorImpl *animatorImpl

	workerStore *event.WorkerStore

	workerAsset workerAsset
}

type workerAsset struct {
	moveSpriteSheet spriteSheet
}

func NewWorkerAnimator(_store *event.WorkerStore) Animator {
	worker := workerImpl{
		workerStore: _store,

		workerAsset: workerAsset{
			moveSpriteSheet: newSpriteSheet(asset.Walk_24x24, asset.Walk_24x24_Size),
		},
	}

	animateEventSet := map[event.Effect]func(event event.Event) []frame{
		event.WorkerMoveEffect: func(e event.Event) []frame {
			workerWalkSpriteSheet := worker.workerAsset.moveSpriteSheet

			workerWalkSprites := workerWalkSpriteSheet.sprites

			destinationPos, err := event.ParseData[math.Pos](e)
			if err != nil {
				return []frame{}
			}

			frames := []frame{}
			currentPos := worker.getWorkerPos(e.EntityID)
			desireNumberOfFrames := len(workerWalkSprites)

			for i := 0; i < len(workerWalkSprites); i++ {
				currentPos = math.NewPos(math.GetNextStepXY(currentPos, 0, destinationPos, 0, aggregator.MAX_WORKER_MOVE_DISTANT/float64(desireNumberOfFrames), 180))
				image := ebiten.NewImageFromImage(workerWalkSprites[i])
				frames = append(frames, frame{
					Image:  image,
					Option: getCenteredDrawImageOptions(workerWalkSprites[i], currentPos),
				})
			}

			return frames
		},
	}

	worker.animatorImpl = newAnimatorImpl()
	worker.animatorImpl.getEventFramesSet = animateEventSet
	worker.animatorImpl.defaultRenderLayer = WorkerRenderLayer
	s := event.Store(*_store)
	worker.animatorImpl.store = &s

	return worker
}

func (i workerImpl) Animator() *animatorImpl {
	return i.animatorImpl
}

func (i workerImpl) getWorkerPos(id uuid.UUID) math.Pos {
	events, err := event.Store(*i.workerStore).GetEventsByEntityID(id)
	if err != nil {
		log.Fatalln(err)
	}

	workerState, err := aggregator.GetWorkerState(events)
	if err != nil {
		log.Fatalln(err)
	}

	return workerState.Position
}
