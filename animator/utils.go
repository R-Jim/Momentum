package animator

import (
	"bytes"
	"image"
	"log"

	"github.com/R-jim/Momentum/math"
	"github.com/hajimehoshi/ebiten/v2"
)

func getAssetImage(asset []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(asset))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func getCenteredDrawImageOptions(image *ebiten.Image, pos math.Pos) *ebiten.DrawImageOptions {
	imageHalfSizeX := float64(image.Bounds().Dx()) / 2
	imageHalfSizeY := float64(image.Bounds().Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X-imageHalfSizeX, pos.Y-imageHalfSizeY)
	return op
}
