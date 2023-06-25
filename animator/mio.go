package animator

import (
	"bytes"
	"image"
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
	animatorImpl *AnimatorImpl

	mioStore *store.Store

	mioIdle *ebiten.Image
	mioWalk *ebiten.Image

	entityLastPosition *map[uuid.UUID]math.Pos
}

func NewMioAnimator(store *store.Store) Animator {
	mioWalkImg, _, err := image.Decode(bytes.NewReader(asset.MioWalk_png))
	if err != nil {
		log.Fatal(err)
	}

	mio := mioImpl{
		mioStore: store,

		mioWalk: ebiten.NewImageFromImage(mioWalkImg),

		entityLastPosition: &map[uuid.UUID]math.Pos{},
	}

	animateEventSet := map[event.Effect]func(event event.Event) []Frame{
		event.MioWalkEffect: func(e event.Event) []Frame {
			destinationPos, err := event.ParseData[math.Pos](e)
			if err != nil {
				return []Frame{}
			}

			frames := []Frame{}
			currentPos := (*mio.entityLastPosition)[e.EntityID] // TODO: should add a func to compose all of Mio's position effected effects
			desireNumberOfFrames := 24.0

			for i := 0; i < int(desireNumberOfFrames); i++ {
				currentPos = math.NewPos(math.GetNextStepXY(currentPos, 0, destinationPos, 0, aggregator.MAX_WALK_DISTANT/desireNumberOfFrames, 180))
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(currentPos.X, currentPos.Y)
				image := ebiten.NewImageFromImage(mio.mioWalk)
				frames = append(frames, Frame{
					Image:  *image,
					Option: *op,
				})
			}

			(*mio.entityLastPosition)[e.EntityID] = destinationPos
			return frames
		},
	}

	mio.animatorImpl = &AnimatorImpl{
		Events:         []event.Event{},
		getEventFrames: animateEventSet,
	}

	return mio
}

func (i mioImpl) Animator() *AnimatorImpl {
	return i.animatorImpl
}

func (i mioImpl) getDrawEventSet() map[event.Effect]func(event.Event) error {
	return map[event.Effect]func(event.Event) error{}
}
