package main

import (
	"fmt"
	"log"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/operator"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	mioID uuid.UUID

	mioOperator *operator.MioOperator

	mioAnimator *animator.Animator

	p math.Pos

	framesToRender map[int][]animator.Frame
}

func (g *Game) Init() {
	mioID := uuid.New()
	store := store.NewStore()

	mioAnimator := animator.NewMioAnimator(&store)

	mioOperator := operator.MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&store),
		MioAnimator:   mioAnimator,
	}

	err := mioOperator.Init(mioID, math.NewPos(1, 2))
	if err != nil {
		log.Fatal(err)
	}

	g.mioID = mioID
	g.mioAnimator = &mioAnimator
	g.mioOperator = &mioOperator
	g.p = math.NewPos(1, 2)
	g.framesToRender = map[int][]animator.Frame{}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		newPos := math.NewPos(g.p.X+aggregator.MAX_WALK_DISTANT, g.p.Y)
		if err := g.mioOperator.Walk(g.mioID, newPos); err != nil {
			return err
		}
		g.p = newPos
	}

	framesPerEvent := (*g.mioAnimator).Animator().GetFramesPerEvent()
	for _, frames := range framesPerEvent {
		for i, frame := range frames {
			g.framesToRender[i] = append(g.framesToRender[i], frame)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	frames, isExist := g.framesToRender[0]
	if !isExist {
		return
	}

	for _, frame := range frames {
		screen.DrawImage(&frame.Image, &frame.Option)
	}

	updatedFramesToRender := map[int][]animator.Frame{}
	for i := 1; i < len(g.framesToRender); i++ {
		updatedFramesToRender[i-1] = g.framesToRender[i]
	}
	g.framesToRender = updatedFramesToRender
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(400, 200)
	ebiten.SetWindowTitle(fmt.Sprintf("Game test"))

	g := &Game{}
	g.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
