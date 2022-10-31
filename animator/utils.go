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
