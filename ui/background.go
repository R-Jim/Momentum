package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func DrawBackground(screen *ebiten.Image) {
	screen.Fill(BACKGROUND)
	// Reset RGB (not Alpha) 0 forcibly
	// op.ColorM.Scale(0, 0, 0, 1)
	// op.ColorM.Translate(float64(0), float64(1), float64(.3), 0)
}
