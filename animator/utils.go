package animator

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func getAssetImage(asset []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(asset))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func centerAndRenderImage(screen *ebiten.Image, image *ebiten.Image, x, y float64) {
	imageHalfSizeX := float64(image.Bounds().Dx()) / 2
	imageHalfSizeY := float64(image.Bounds().Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x-imageHalfSizeX, y-imageHalfSizeY)
	screen.DrawImage(image, op)
}
