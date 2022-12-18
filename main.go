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
	"github.com/R-jim/Momentum/automaton"
	"github.com/R-jim/Momentum/operator"
	"github.com/R-jim/Momentum/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/sync/errgroup"
)

var (
	opt        operator.Operator
	ani        animator.Animator
	knightAuto automaton.KnightAutomaton
	spikeAuto  automaton.SpikeAutomaton

	knightID string
	spikeID  string

	isPaused  bool
	turnCount int
)

// for testing
func initEntities() {
	knightID = "knight_1"

	err := opt.Knight.Init(knightID, knight.Health{Max: 50, Value: 50})
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}

	err = opt.Knight.Move(knightID, knight.PositionState{
		X: 150,
		Y: 100,
	})
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}

	spikeID = "spike_1"
	err = opt.Spike.Init(spikeID, "")
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}

	err = opt.Spike.Move(spikeID, spike.PositionState{
		X: 100,
		Y: 100,
	})
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}

	err = opt.Knight.ChangeTarget(knightID, knight.Target{
		ID: spikeID,
		Position: knight.PositionState{
			X: 100,
			Y: 100,
		},
	})
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}
}

func init() {
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

	knightAuto = automaton.NewKnightAutomaton(knightStore, spikeStore, opt)
	spikeAuto = automaton.NewSpikeAutomaton(carrierStore, knightStore, spikeStore, opt)
	initEntities()

	isPaused = true
}

type Game struct {
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		isPaused = false
	}
	if isPaused {
		return nil
	}
	{
		fmt.Printf("\n--------- Turn %v ----------\n", turnCount)
		err := knightAuto.Auto(knightID)
		if err != nil {
			return err
		}
		err = spikeAuto.Auto(spikeID)
		if err != nil {
			return err
		}
		turnCount++
	}

	isPaused = true
	// operations := []func() error{}
	// operations = append(operations, userInput()...)
	// go runConcurrently(operations)
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
	ani.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 400, 400
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Momentum")

	game := &Game{}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
