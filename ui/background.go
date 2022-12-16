package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func DrawBackground(screen *ebiten.Image) {
	screenWidth, screenHeight := screen.Size()
	background := ebiten.NewImage(screenWidth, screenHeight)
	background.Fill(BACKGROUND)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)

	// Reset RGB (not Alpha) 0 forcibly
	// op.ColorM.Scale(0, 0, 0, 1)
	// op.ColorM.Translate(float64(0), float64(1), float64(.3), 0)

	screen.DrawImage(background, op)
}
