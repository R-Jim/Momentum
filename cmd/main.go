package main

import (
	"fmt"
	"log"

	"image/color"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/automaton"
	"github.com/R-jim/Momentum/cmd/element"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/operator"
	"github.com/R-jim/Momentum/system"
	"github.com/R-jim/Momentum/ui"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	// entities
	mioID       uuid.UUID
	buildingIDs []uuid.UUID
	workerIDs   []uuid.UUID

	gameMap []math.Path

	// stores
	mioStore      *event.MioStore
	workerStore   *event.WorkerStore
	streetStore   *event.StreetStore
	buildingStore *event.BuildingStore

	// operators
	mioOperator      operator.MioOperator
	workerOperator   operator.WorkerOperator
	streetOperator   operator.StreetOperator
	buildingOperator operator.BuildingOperator

	// animators
	mioAnimator    *animator.Animator
	workerAnimator *animator.Animator

	// automatons
	mioAutomaton     *automaton.MioAutomaton
	workerAutomatons []*automaton.WorkerAutomaton

	automationCounter int

	// ui
	defaultLayer *ebiten.Image
	mioLayer     *ebiten.Image
	workerLayer  *ebiten.Image
	effectLayer  *ebiten.Image

	uiLayer            *ebiten.Image
	workerCommandLayer *ebiten.Image

	clickAbleWorkerCardElements []element.WorkerCardClickableElement

	screenWidth  float64
	screenHeight float64

	workerLayerX      float64
	workerLayerY      float64
	workerLayerWidth  float64
	workerLayerHeight float64

	uiLayerPadding float64

	// game state
	count              int
	selectedBuildingID uuid.UUID
}

func (g *Game) Init() {
	g.initStore()
	g.initOperator()

	g.initEntities()
	g.initAutomaton()

	g.initAnimator()

	g.initUI()

	g.selectedBuildingID = g.buildingIDs[0]
}

func (g *Game) initUI() {
	g.screenWidth = 800
	g.screenHeight = 600

	g.effectLayer = ebiten.NewImage(int(g.screenWidth), int(g.screenHeight))
	g.mioLayer = ebiten.NewImage(int(g.screenWidth), int(g.screenHeight))
	g.workerLayer = ebiten.NewImage(int(g.screenWidth), int(g.screenHeight))
	g.defaultLayer = ebiten.NewImage(int(g.screenWidth), int(g.screenHeight))
	g.uiLayer = ebiten.NewImage(int(g.screenWidth), int(g.screenHeight))

	g.workerLayerWidth = 400
	g.workerLayerHeight = 150

	g.workerCommandLayer = ebiten.NewImage(int(g.workerLayerWidth), int(g.workerLayerHeight))

	g.uiLayerPadding = 20

	g.workerLayerX = (g.screenWidth - g.workerLayerWidth) / 2
	g.workerLayerY = g.screenHeight - g.workerLayerHeight - g.uiLayerPadding
}

func (g *Game) initEntities() {
	g.mioID = uuid.New()

	street1ID := uuid.New()
	street2ID := uuid.New()
	street3ID := uuid.New()

	buildingStreetID1 := uuid.New()
	buildingStreetID2 := uuid.New()

	posA := math.NewPos(200, 200)
	posB := math.NewPos(300, 200)
	posC := math.NewPos(500, 200)
	posD := math.NewPos(400, 300)
	buildingPos := math.NewPos(600, 400)

	mapPaths := []math.Path{
		{Start: posA, End: posB, Cost: 10},        // street1
		{Start: posB, End: posC, Cost: 5},         // street2
		{Start: posB, End: posD, Cost: 50},        // street3
		{Start: posC, End: buildingPos, Cost: 15}, // building street1
		{Start: posD, End: buildingPos, Cost: 10}, // building street2
	}

	g.gameMap = mapPaths

	err := g.mioOperator.Init(g.mioID, posA)
	if err != nil {
		log.Fatal(err)
	}

	err = g.streetOperator.Init(street1ID, posA, posB)
	if err != nil {
		log.Fatal(err)
	}
	err = g.streetOperator.Init(street2ID, posB, posC)
	if err != nil {
		log.Fatal(err)
	}
	err = g.streetOperator.Init(street3ID, posB, posD)
	if err != nil {
		log.Fatal(err)
	}
	err = g.streetOperator.Init(buildingStreetID1, posC, buildingPos)
	if err != nil {
		log.Fatal(err)
	}
	err = g.streetOperator.Init(buildingStreetID2, posD, buildingPos)
	if err != nil {
		log.Fatal(err)
	}

	drinkStoreID := uuid.New()
	err = g.buildingOperator.Init(drinkStoreID, event.BuildingTypeDrinkStore, buildingPos)
	if err != nil {
		log.Fatal(err)
	}
	g.buildingIDs = []uuid.UUID{drinkStoreID}

	worker1ID := uuid.New()
	g.workerIDs = []uuid.UUID{worker1ID}
	err = g.workerOperator.Init(worker1ID, posB)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) initAnimator() {
	mioAnimator := animator.NewMioAnimator(g.mioStore)
	workerAnimator := animator.NewWorkerAnimator(g.workerStore)

	g.mioAnimator = &mioAnimator
	g.workerAnimator = &workerAnimator
}

func (g *Game) initAutomaton() {
	g.mioAutomaton = &automaton.MioAutomaton{
		EntityID: g.mioID,
		MapPaths: g.gameMap,

		MioStore:      g.mioStore,
		StreetStore:   g.streetStore,
		BuildingStore: g.buildingStore,

		MioOperator:    g.mioOperator,
		StreetOperator: g.streetOperator,
	}

	for _, workerID := range g.workerIDs {
		g.workerAutomatons = append(g.workerAutomatons, &automaton.WorkerAutomaton{
			EntityID: workerID,
			MapPaths: g.gameMap,

			WorkerStore:   g.workerStore,
			StreetStore:   g.streetStore,
			BuildingStore: g.buildingStore,

			WorkerOperator: g.workerOperator,
			StreetOperator: g.streetOperator,
		})
	}
}

func (g *Game) initOperator() {
	g.streetOperator = operator.NewStreet(g.streetStore)
	g.buildingOperator = operator.NewBuilding(g.buildingStore, nil)
	g.workerOperator = operator.NewWorker(g.workerStore, g.buildingStore)
	g.mioOperator = operator.NewMio(g.mioStore, g.buildingStore)
}

func (g *Game) initStore() {
	mioStore := event.NewMioStore()
	streetStore := event.NewStreetStore()
	buildingStore := event.NewBuildingStore()
	workerStore := event.NewWorkerStore()

	g.mioStore = &mioStore
	g.workerStore = &workerStore
	g.streetStore = &streetStore
	g.buildingStore = &buildingStore
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		if err := g.mioOperator.SelectBuilding(g.mioID, g.selectedBuildingID); err != nil {
			return err
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		if err := g.mioOperator.Act(g.mioID, g.selectedBuildingID); err != nil {
			// return err
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		if err := g.mioOperator.EnterBuilding(g.mioID, g.selectedBuildingID); err != nil {
			// return err
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		if err := g.mioOperator.Eat(g.mioID, 20); err != nil {
			// return err
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if err := g.mioOperator.Stream(g.mioID, 20); err != nil {
			// return err
		}
	}

	isCursorPointer := false
	for _, element := range g.clickAbleWorkerCardElements {
		if element.ClickAbleImage.In(ebiten.CursorPosition()) {
			isCursorPointer = true
			break
		}
	}

	if isCursorPointer {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, element := range g.clickAbleWorkerCardElements {
			if element.ClickAbleImage.In(ebiten.CursorPosition()) {
				element.OnClick(g.selectedBuildingID)
				break
			}
		}
	}

	g.automationCounter++
	if g.automationCounter >= int(system.AUTOMATION_TICK_PER_FPS) {
		g.mioAutomaton.PathFindingUpdate()
		g.mioAutomaton.Move()

		for _, workerAutomaton := range g.workerAutomatons {
			workerAutomaton.PathFindingUpdate()
			workerAutomaton.Move()
		}

		g.automationCounter = 0
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawMap(screen)
	g.DrawBuilding(screen)

	if g.count == int(60/system.DEFAULT_FPS) {
		g.effectLayer.Clear()
		g.mioLayer.Clear()
		g.workerLayer.Clear()
		g.defaultLayer.Clear()

		frames := (*g.mioAnimator).Animator().GetFrames()
		frames = append(frames, (*g.workerAnimator).Animator().GetFrames()...)

		for _, frame := range frames {
			switch frame.RenderLayer {
			case animator.EffectRenderLayer:
				g.effectLayer.DrawImage(frame.Image, frame.Option)
			case animator.MioRenderLayer:
				g.mioLayer.DrawImage(frame.Image, frame.Option)
			case animator.WorkerRenderLayer:
				g.workerLayer.DrawImage(frame.Image, frame.Option)
			default:
				g.defaultLayer.DrawImage(frame.Image, frame.Option)
			}
		}

		g.count = 0
	}

	g.DrawWorkerCommandPanel()

	screen.DrawImage(g.effectLayer, &ebiten.DrawImageOptions{})
	screen.DrawImage(g.mioLayer, &ebiten.DrawImageOptions{})
	screen.DrawImage(g.workerLayer, &ebiten.DrawImageOptions{})
	screen.DrawImage(g.defaultLayer, &ebiten.DrawImageOptions{})
	screen.DrawImage(g.uiLayer, &ebiten.DrawImageOptions{})

	g.count++
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

	for _, buildingID := range g.buildingIDs {
		events, err := event.Store(*g.buildingStore).GetEventsByEntityID(buildingID)
		if err != nil {
			log.Fatalln(err)
		}

		building, err := aggregator.GetBuildingState(events)
		if err != nil {
			log.Fatalln(err)
		}

		drawBuilding(building.Pos)
	}
}

func (g *Game) DrawWorkerCommandPanel() {
	numberOfWorkers := len(g.workerIDs)

	clickAbleWorkerCardElements := []element.WorkerCardClickableElement{}

	ebitenutil.DrawRect(g.workerCommandLayer, 0, 0, g.workerLayerWidth, g.workerLayerHeight, color.RGBA{0xff, 0xff, 0x0, 0xff})
	padding := float64(10)
	cardWidth := float64(100)
	cardHeight := g.workerLayerHeight - padding*2

	holderWidth := cardWidth*float64(numberOfWorkers) + padding*float64(numberOfWorkers-1)
	holderX := (g.workerLayerWidth - holderWidth) / 2
	workerCardHolder := ebiten.NewImage(int(holderWidth), int(g.workerLayerHeight-padding))

	for i := 0; i < numberOfWorkers; i++ {
		x := float64(i) * cardWidth
		if i != 0 {
			x += padding
		}

		workerCard := ebiten.NewImage(int(cardWidth), int(cardHeight))
		ebitenutil.DrawRect(workerCard, 0, 0, cardWidth, cardHeight, color.RGBA{0x0, 0xff, 0xff, 0xff})

		clickAbleImage := ui.NewClickAbleImage(workerCard, int(g.workerLayerX+holderX+x), int(g.workerLayerY))
		clickAbleWorkerCardElements = append(clickAbleWorkerCardElements, element.NewWorkerCardClickableElement(
			g.workerOperator, g.workerIDs[i], clickAbleImage,
		))

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, 0)

		workerCardHolder.DrawImage(workerCard, op)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(holderX, padding)

	g.workerCommandLayer.DrawImage(workerCardHolder, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.workerLayerX, g.workerLayerY)
	g.uiLayer.DrawImage(g.workerCommandLayer, op)

	g.clickAbleWorkerCardElements = clickAbleWorkerCardElements
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
