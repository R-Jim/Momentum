package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	player *ebiten.Image
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(player, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 400, 400
}

func main() {
	ebiten.SetWindowSize(1000, 800)
	ebiten.SetWindowTitle("Momentum")

	player = ebiten.NewImage(100, 100)
	player.Fill(color.White)

	game := &Game{}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
