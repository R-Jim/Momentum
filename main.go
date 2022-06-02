package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"

	assets "github.com/R-jim/Momentum/asset"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	playerImage *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(assets.Player_png))
	if err != nil {
		log.Fatal(err)
	}
	playerImage = ebiten.NewImageFromImage(img)
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawPlayer(screen)
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

func (g *Game) drawPlayer(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(playerImage, op)
}
