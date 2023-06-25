package main

import (
	"fmt"
	"log"

	"image/color"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/operator"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	mioID uuid.UUID

	mioOperator *operator.MioOperator

	mioAnimator *animator.Animator

	p math.Pos

	framesToRender map[int][]animator.Frame
	gameMap       []math.Path
	buildingPoses []math.Pos // TODO: should get building pos from building store
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

	posA := math.NewPos(200, 200)
	posB := math.NewPos(450, 200)
	posC := math.NewPos(450, 100)
	posD := math.NewPos(200, 100)
	posE := math.NewPos(300, 300)
	building1Pos := math.NewPos(400, 400)
	building2Pos := math.NewPos(300, 150)

	g.gameMap = []math.Path{
		{Start: posA, End: posB, Cost: 2},
		{Start: posA, End: posD, Cost: 1},
		{Start: posA, End: posE, Cost: 3},

		{Start: posB, End: posC, Cost: 2},
		{Start: posB, End: building1Pos, Cost: 8},
		{Start: posB, End: posE, Cost: 3},

		{Start: posC, End: posD, Cost: 4},
		{Start: posC, End: building2Pos, Cost: 2},

		{Start: posD, End: building2Pos, Cost: 1},

		{Start: posE, End: building1Pos, Cost: 5},
	}

	g.buildingPoses = []math.Pos{building1Pos, building2Pos}
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
	g.DrawMap(screen)
	g.DrawBuilding(screen)

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

func (g *Game) DrawMap(screen *ebiten.Image) {
	// TODO: separate to path draw layer
	drawedPosSet := map[math.Pos]bool{}
	pointRadius := 20.0

	drawPoint := func(pos math.Pos) {
		ebitenutil.DrawRect(screen, pos.X-pointRadius/2, pos.Y-pointRadius/2, pointRadius, pointRadius, color.White)
	}

	for _, path := range g.gameMap {
		if !drawedPosSet[path.Start] {
			drawPoint(path.Start)
			drawedPosSet[path.Start] = true
		}
		if !drawedPosSet[path.End] {
			drawPoint(path.End)
			drawedPosSet[path.End] = true
		}
		ebitenutil.DrawLine(
			screen,
			path.Start.X, path.Start.Y,
			path.End.X, path.End.Y,
			color.White)
	}
}

func (g *Game) DrawBuilding(screen *ebiten.Image) {
	// TODO: separate to building draw layer
	pointRadius := 20.0

	drawBuilding := func(pos math.Pos) {
		ebitenutil.DrawRect(screen, pos.X-pointRadius/2, pos.Y-pointRadius/2, pointRadius, pointRadius, color.RGBA{0x0, 0xff, 0xff, 0xff})
	}

	for _, pos := range g.buildingPoses {
		drawBuilding(pos)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle(fmt.Sprintf("Game test"))

	g := &Game{}
	g.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
