package main

import (
	"context"
	"fmt"
	_ "image/png"
	"log"

	"github.com/R-jim/Momentum/aggregate/carrier"
	"github.com/R-jim/Momentum/aggregate/jet"
	"github.com/R-jim/Momentum/aggregate/spike"
	"github.com/R-jim/Momentum/aggregate/storage"
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/operator"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/sync/errgroup"
)

var (
	opt operator.Operator
	ani animator.Animator

	spikeID string
)

// for testing
func initEntities() {
	spikeID = "spike_1"

	err := opt.Spike.Init(spikeID, "artifact")
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}
}

func init() {
	storageStore := storage.NewStore()
	jetStore := jet.NewStore()
	carrierStore := carrier.NewStore()
	spikeStore := spike.NewStore()

	fuelTankAggregator := storage.NewAggregator(storageStore)
	jetAggregator := jet.NewAggregator(jetStore)
	carrierAggregator := carrier.NewAggregator(carrierStore)
	spikeAggregator := spike.NewAggregator(spikeStore)

	ani = animator.New(animator.AnimatorStores{
		JetStore:   jetStore,
		SpikeStore: spikeStore,
	})

	opt = operator.New(
		operator.OperatorAggregator{
			JetAggregator:      jetAggregator,
			FuelTankAggregator: fuelTankAggregator,
			CarrierAggregator:  carrierAggregator,
			SpikeAggregator:    spikeAggregator,
		},
		ani,
	)

	initEntities()
}

type Game struct {
}

func (g *Game) Update() error {
	operations := []func() error{}
	operations = append(operations, userInput()...)
	go runConcurrently(operations)
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
	ani.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 400, 400
}

func main() {
	ebiten.SetWindowSize(1000, 800)
	ebiten.SetWindowTitle("Momentum")

	game := &Game{}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
