package main

import (
	"context"
	"fmt"
	_ "image/png"
	"log"

	"github.com/R-jim/Momentum/fueltank"
	"github.com/R-jim/Momentum/info"
	"github.com/R-jim/Momentum/jet"
	"github.com/R-jim/Momentum/operator"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/sync/errgroup"
)

var (
	// playerImage *ebiten.Image
	// player      *entity.Core

	fuelTankStore      fueltank.Store
	fuelTankAggregator fueltank.Aggregator

	opt operator.Operator

	fuelTankID, jetID string
)

func initOperator() {
	fuelTankStore = fueltank.NewStore()
	fuelTankAggregator = fueltank.NewAggregator(fuelTankStore)

	jetStore := jet.NewStore()
	jetAggregator := jet.NewAggregator(jetStore)

	opt = operator.New(jetAggregator, fuelTankAggregator)
}

// for testing
func initEntities() {
	fuelTankID = "fuel_tank_1"
	jetID = "jet_1"

	err := opt.FuelTank.Init(fuelTankID)
	if err != nil {
		// TODO: log error
	}

	err = opt.Jet.Init(jetID, fuelTankID)
	if err != nil {
		// TODO: log error
	}

	err = opt.FuelTank.Refill(fuelTankID, 15)
	if err != nil {
		// TODO: log error
	}

	// TEST REROLL BACK
	// jetOperator.Landing(jetID)
}

func init() {
	initOperator()
	initEntities()
	// img, _, err := image.Decode(bytes.NewReader(assets.Player_png))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// playerImage = ebiten.NewImageFromImage(img)

	// player = &entity.Core{
	// 	Position: valueobject.Position{
	// 		X: 10,
	// 		Y: 10,
	// 	},
	// 	Momentum: valueobject.Momentum{
	// 		X: 100,
	// 		Y: 200,
	// 	},
	// }
}

type Game struct{}

func (g *Game) Update() error {
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
	// g.drawPlayer(screen)
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

// func (g *Game) drawPlayer(screen *ebiten.Image) {
// 	op := &ebiten.DrawImageOptions{}
// 	op.GeoM.Translate(float64(player.Position.X), float64(player.Position.Y))
// 	screen.DrawImage(playerImage, op)
// }

//GameData
type GameData struct {
}

//Entity
type Entity struct {
	Info     info.Info
	GameData GameData
}
