package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type ClickAbleImage interface {
	In(x, y int) bool
}

type clickAbleImage struct {
	image *ebiten.Image
	x     int
	y     int
}

func NewClickAbleImage(image *ebiten.Image, x, y int) ClickAbleImage {
	return clickAbleImage{
		image,
		x,
		y,
	}
}

func (s clickAbleImage) In(x, y int) bool {
	// example from https://ebitengine.org/en/examples/drag.html#Code
	return s.image.At(x-s.x, y-s.y).(color.RGBA).A > 0
}
