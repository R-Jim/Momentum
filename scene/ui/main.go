package main

import (
	"context"
	"fmt"
	_ "image/png"
	"log"

	"github.com/R-jim/Momentum/aggregate/artifact"
	"github.com/R-jim/Momentum/aggregate/carrier"
	"github.com/R-jim/Momentum/aggregate/jet"
	"github.com/R-jim/Momentum/aggregate/knight"
	"github.com/R-jim/Momentum/aggregate/spike"
	"github.com/R-jim/Momentum/aggregate/storage"
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/operator"
	"github.com/R-jim/Momentum/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/sync/errgroup"
)

type modeIndex int

const (
	startModeIndex modeIndex = 1
	playModeIndex  modeIndex = 2
	pauseModeIndex modeIndex = 3
)

const (
	screenWidth  = 640
	screenHeight = 480
	fontSize     = 24
)

var (
	opt      operator.Operator
	ani      animator.Animator
	knightID string
)

// for testing
func initEntities() {
	knightID = "knight_1"

	err := opt.Knight.Init(knightID, knight.Health{Max: 50, Value: 50})
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}

	err = opt.Knight.Move(knightID, knight.PositionState{
		X: 100,
		Y: 100,
	})
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}
}

func (g *Game) init() {
	storageStore := storage.NewStore()
	jetStore := jet.NewStore()
	carrierStore := carrier.NewStore()
	spikeStore := spike.NewStore()
	artifactStore := artifact.NewStore()
	knightStore := knight.NewStore()

	fuelTankAggregator := storage.NewAggregator(storageStore)
	jetAggregator := jet.NewAggregator(jetStore)
	carrierAggregator := carrier.NewAggregator(carrierStore)
	spikeAggregator := spike.NewAggregator(spikeStore)
	artifactAggregator := artifact.NewAggregator(artifactStore)
	knightAggregator := knight.NewAggregator(knightStore)

	ani = animator.New(animator.AnimatorStores{
		JetStore:      jetStore,
		SpikeStore:    spikeStore,
		ArtifactStore: artifactStore,
		KnightStore:   knightStore,
	})

	opt = operator.New(
		operator.OperatorAggregator{
			JetAggregator:      jetAggregator,
			FuelTankAggregator: fuelTankAggregator,
			CarrierAggregator:  carrierAggregator,
			SpikeAggregator:    spikeAggregator,
			ArtifactAggregator: artifactAggregator,
			KnightAggregator:   knightAggregator,
		},
		ani,
	)

	initEntities()
	g.mode = startModeIndex
	g.ui = ui.New(fonts.PressStart2P_ttf, fontSize)
}

type Game struct {
	mode modeIndex
	ui   ui.UI
}

func (g *Game) Update() error {
	// operations := []func() error{}
	// operations = append(operations, userInput()...)
	// go runConcurrently(operations)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch g.mode {
		case startModeIndex:
			g.mode = playModeIndex
		case playModeIndex:
			g.mode = pauseModeIndex
		case pauseModeIndex:
			g.mode = playModeIndex
		}
	}
	return nil
}

func userInput() []func() error {
	operations := []func() error{}
	// if inpututil.IsKeyJustPressed(ebiten.KeyF) {
	// 	fmt.Println("[Fly]")
	// 	operations = append(operations, func() error {
	// 		return opt.Jet.Fly(jet1ID, fuelTank1ID, 5, jet.PositionState{
	// 			X: 1,
	// 			Y: 1,
	// 		})
	// 	})
	// }
	return operations
}

func runConcurrently(operations []func() error) {
	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(1)

	for _, operation := range operations {
		op := operation
		g.Go(func() error {
			return op()
		})
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	ui.DrawBackground(screen)
	var titleTexts string
	var text string
	switch g.mode {
	case startModeIndex:
		titleTexts = "PROJECT NAME HERE"
		text = "PRESS SPACE TO START"
	case playModeIndex, pauseModeIndex:
		if g.mode == pauseModeIndex {
			text = "PAUSED"
		}
		ani.Draw(screen)
	}

	g.ui.DrawTitle(screen, []string{titleTexts}, (screenWidth-len(titleTexts)*int(g.ui.TitleFontSize))/2, 0, false)
	g.ui.DrawMsg(screen, []string{text}, (screenWidth-len(text)*int(g.ui.FontSize))/2, screenHeight, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Momentum")

	game := &Game{}
	game.init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
