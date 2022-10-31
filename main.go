package main

import (
	"context"
	"fmt"
	_ "image/png"
	"log"

	"github.com/R-jim/Momentum/aggregate/jet"
	"github.com/R-jim/Momentum/aggregate/storage"
	"github.com/R-jim/Momentum/animator"
	"github.com/R-jim/Momentum/automaton"
	"github.com/R-jim/Momentum/operator"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/sync/errgroup"
)

var (
	opt     operator.Operator
	ani     animator.Animator
	jetAuto automaton.JetAutomaton

	fuelTankID, jetID string
)

// for testing
func initEntities() {
	fuelTankID = "fuel_tank_1"
	jetID = "jet_1"

	err := opt.FuelTank.Init(fuelTankID)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}

	err = opt.Jet.Init(jetID, fuelTankID)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}

	err = opt.FuelTank.Refill(fuelTankID, 150)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}
}

func init() {
	storageStore := storage.NewStore()
	jetStore := jet.NewStore()

	fuelTankAggregator := storage.NewAggregator(storageStore)
	jetAggregator := jet.NewAggregator(jetStore)

	ani = animator.New(jetStore)

	opt = operator.New(
		operator.OperatorAggregator{
			JetAggregator:      jetAggregator,
			FuelTankAggregator: fuelTankAggregator,
		},
		ani,
	)

	jetAuto = automaton.NewJetAutomaton(jetStore, storageStore, opt)
	initEntities()
}

type Game struct {
}

func (g *Game) Update() error {
	err := jetAuto.Auto(jetID)
	if err != nil {
		fmt.Printf("[ERROR]Update: %v\n", err.Error())
	}

	operations := []func() error{}
	operations = append(operations, userInput()...)
	go runConcurrently(operations)
	return nil
}

func userInput() []func() error {
	operations := []func() error{}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		fmt.Println("[Fly]")
		operations = append(operations, func() error {
			return opt.Jet.Fly(jetID, fuelTankID, 5, jet.PositionState{
				X: 1,
				Y: 1,
			})
		})

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		fmt.Println("[Landing]")
		operations = append(operations, func() error {
			return opt.Jet.Landing(jetID)
		})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		fmt.Println("[Landing]")
		operations = append(operations, func() error {
			return opt.Jet.Takeoff(jetID)
		})
	}
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
