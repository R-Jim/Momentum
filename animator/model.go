package animator

import (
	"bytes"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type frame struct {
	Image  *ebiten.Image
	Option *ebiten.DrawImageOptions
}

type spriteSheet struct {
	asset      []byte
	spriteSize []int
	sprites    []*ebiten.Image
}

func newSpriteSheet(asset []byte, spriteSize []int) spriteSheet {
	spriteSheetImg, _, err := image.Decode(bytes.NewReader(asset))
	if err != nil {
		log.Fatal(err)
	}

	var spriteSheetXSize, spriteSheetYSize int

	spriteSheetXSize = spriteSheetImg.Bounds().Min.X + spriteSheetImg.Bounds().Max.X
	spriteSheetYSize = spriteSheetImg.Bounds().Min.Y + spriteSheetImg.Bounds().Max.Y
	// log.Printf("sprite sheet size: %vx%v\n", spriteSheetXSize, spriteSheetYSize)

	numberOfRows := (int)(math.Round((float64)(spriteSheetXSize / spriteSize[0])))
	numberOfCols := (int)(math.Round((float64)(spriteSheetYSize / spriteSize[1])))

	sprites := []*ebiten.Image{}

	for col := 0; col < numberOfCols; col++ {
		for row := 0; row < numberOfRows; row++ {
			spriteSheetCopy := ebiten.NewImageFromImage(spriteSheetImg)

			sprites = append(sprites, ebiten.NewImageFromImage(spriteSheetCopy.SubImage(image.Rect(
				row*spriteSize[0], col*spriteSize[1],
				(row+1)*spriteSize[0], (col+1)*spriteSize[1],
			))))
		}
	}

	// log.Printf("sprite sheet size: %vx%v, number of sprites: %v, number of Rows: %v, Cols: %v\n", spriteSheetXSize, spriteSheetYSize, len(sprites), numberOfRows, numberOfCols)

	return spriteSheet{
		asset:      asset,
		spriteSize: spriteSize,
		sprites:    sprites,
	}
}
