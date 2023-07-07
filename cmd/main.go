package main

import (
	"fmt"
	"log"

	"image/color"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/automaton"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/operator"
	"github.com/R-jim/Momentum/system"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	mioID      uuid.UUID
	buildingID uuid.UUID

	mioStore *store.Store

	mioOperator *operator.MioOperator

	mioAnimator *animator.Animator

	mioAutomaton *automaton.MioAutomaton

	automationCounter int

	gameMap       []math.Path
	buildingPoses []math.Pos // TODO: should get building pos from building store
}

func (g *Game) Init() {
	mioID := uuid.New()

	mioStore := store.NewStore()
	buildingStore := store.NewStore()

	mioAnimator := animator.NewMioAnimator(&mioStore)

	buildingOperator := operator.BuildingOperator{
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	mioOperator := operator.MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&mioStore),
		MioAnimator:   mioAnimator,

		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	g.mioID = mioID
	g.mioStore = &mioStore
	g.mioAnimator = &mioAnimator
	g.mioOperator = &mioOperator

	posA := math.NewPos(200, 200)
	posB := math.NewPos(300, 200)
	posC := math.NewPos(500, 200)
	posD := math.NewPos(400, 300)
	buildingPos := math.NewPos(600, 400)

	g.buildingPoses = []math.Pos{buildingPos}

	mapPaths := []math.Path{
		{Start: posA, End: posB, Cost: 10},        // street1
		{Start: posB, End: posC, Cost: 5},         // street2
		{Start: posB, End: posD, Cost: 50},        // street3
		{Start: posC, End: buildingPos, Cost: 15}, // building street1
		{Start: posD, End: buildingPos, Cost: 10}, // building street2
	}

	g.gameMap = mapPaths

	mapGraph := math.NewGraph(mapPaths)

	street1ID := uuid.New()
	street2ID := uuid.New()
	street3ID := uuid.New()
	buildingStreetID1 := uuid.New()
	buildingStreetID2 := uuid.New()

	streetStore := store.NewStore()

	streetOperator := operator.NewStreet(aggregator.NewStreetAggregator(&streetStore), nil)

	g.mioAutomaton = &automaton.MioAutomaton{
		EntityID: mioID,
		MapPaths: mapPaths,
		MapGraph: mapGraph,

		MioStore:      &mioStore,
		StreetStore:   &streetStore,
		BuildingStore: &buildingStore,

		MioOperator:    mioOperator,
		StreetOperator: streetOperator,
	}

	err := mioOperator.Init(mioID, buildingPos)
	if err != nil {
		log.Fatal(err)
	}

	err = streetOperator.Init(street1ID, posA, posB)
	if err != nil {
		log.Fatal(err)
	}
	err = streetOperator.Init(street2ID, posB, posC)
	if err != nil {
		log.Fatal(err)
	}
	err = streetOperator.Init(street3ID, posB, posD)
	if err != nil {
		log.Fatal(err)
	}
	err = streetOperator.Init(buildingStreetID1, posC, buildingPos)
	if err != nil {
		log.Fatal(err)
	}
	err = streetOperator.Init(buildingStreetID2, posD, buildingPos)
	if err != nil {
		log.Fatal(err)
	}

	drinkStoreID := uuid.New()
	err = buildingOperator.Init(drinkStoreID, event.BuildingTypeDrinkStore, buildingPos)
	if err != nil {
		log.Fatal(err)
	}
	g.buildingID = drinkStoreID
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		if err := g.mioOperator.SelectBuilding(g.mioID, g.buildingID); err != nil {
			return err
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		if err := g.mioOperator.Act(g.mioID, g.buildingID); err != nil {
			// return err
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		if err := g.mioOperator.EnterBuilding(g.mioID, g.buildingID); err != nil {
			// return err
		}
	}

	g.automationCounter++
	if g.automationCounter >= int(system.AUTOMATION_TICK_PER_FPS) {
		g.mioAutomaton.PathFindingUpdate()
		g.mioAutomaton.Move()
		g.automationCounter = 0
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawMap(screen)
	g.DrawBuilding(screen)

	frames := (*g.mioAnimator).Animator().GetFrames()

	for _, frame := range frames {
		screen.DrawImage(frame.Image, frame.Option)
	}
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
