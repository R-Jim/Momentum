package main

import (
	"context"
	"fmt"
	_ "image/png"
	"log"

	"github.com/R-jim/Momentum/aggregate/carrier"
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

	carrierID, fuelTank1ID, fuelTank2ID, jet1ID, jet2ID string
)

// for testing
func initEntities() {
	carrierID = "carrier_1"

	fuelTank1ID = "fuel_tank_1"
	fuelTank2ID = "fuel_tank_2"

	jet1ID = "jet_1"
	jet2ID = "jet_2"

	err := opt.FuelTank.Init(fuelTank1ID)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}
	err = opt.FuelTank.Init(fuelTank2ID)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}

	err = opt.Jet.Init(jet1ID, fuelTank1ID)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}
	err = opt.Jet.Init(jet2ID, fuelTank2ID)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}

	err = opt.FuelTank.Refill(fuelTank1ID, 1050)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}
	err = opt.FuelTank.Refill(fuelTank2ID, 150)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}

	err = opt.Jet.Takeoff(jet1ID)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}
	err = opt.Jet.Fly(jet1ID, fuelTank1ID, 0, jet.PositionState{
		X: float64(100),
		Y: float64(100),
	})
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}

	err = opt.Carrier.Init(carrierID)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}

	err = opt.Carrier.HouseJet(carrierID, jet1ID)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}
	err = opt.Carrier.HouseJet(carrierID, jet2ID)
	if err != nil {
		fmt.Printf("[ERROR]initEntities: %v\n", err.Error())
	}
}

func init() {
	storageStore := storage.NewStore()
	jetStore := jet.NewStore()
	carrierStore := carrier.NewStore()

	fuelTankAggregator := storage.NewAggregator(storageStore)
	jetAggregator := jet.NewAggregator(jetStore)
	carrierAggregator := carrier.NewAggregator(carrierStore)

	ani = animator.New(jetStore)

	opt = operator.New(
		operator.OperatorAggregator{
			JetAggregator:      jetAggregator,
			FuelTankAggregator: fuelTankAggregator,
			CarrierAggregator:  carrierAggregator,
		},
		ani,
	)

	jetAuto = automaton.NewJetAutomaton(jetStore, storageStore, opt)
	initEntities()
}

type Game struct {
}

func (g *Game) Update() error {
	err := jetAuto.Auto(jet1ID)
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
			return opt.Jet.Fly(jet1ID, fuelTank1ID, 5, jet.PositionState{
				X: 1,
				Y: 1,
			})
		})

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		fmt.Println("[Landing]")
		operations = append(operations, func() error {
			return opt.Jet.Landing(jet1ID)
		})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		fmt.Println("[Launch Jet]")
		operations = append(operations, func() error {
			return opt.Carrier.LaunchJet(carrierID, jet1ID)
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
