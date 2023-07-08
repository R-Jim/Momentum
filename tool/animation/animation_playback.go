package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	_ "image/jpeg"
	_ "image/png"
)

func animate(screen *ebiten.Image, sprites []image.Image) {
	for _, sprite := range sprites {
		op := &ebiten.DrawImageOptions{}
		image := ebiten.NewImageFromImage(sprite)
		screen.DrawImage(image, op)
	}
}

type Animator struct {
	sprites []*ebiten.Image
	fps     int
	counter int
}

func (a *Animator) Update() error {
	a.counter++
	return nil
}

func (a *Animator) Draw(screen *ebiten.Image) {
	tick := 60 / a.fps

	index := (a.counter / tick) % len(a.sprites)

	// log.Printf("animate frame, %v\n", index)

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Apply(0, 0)

	image := ebiten.NewImageFromImage(a.sprites[index])
	screen.DrawImage(image, op)
}

func (a *Animator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("missing sprite sheet file name")
	}

	spriteSheetFile := os.Args[1]
	if spriteSheetFile == "" {
		log.Fatalln("missing sprite sheet file name")
	}

	fileNameRegex, err := regexp.Compile(".png|.jpeg")
	if err != nil {
		log.Fatalln(err)
	}
	spriteSheetName := fileNameRegex.Split(spriteSheetFile, 2)[0]

	var sizePerSprite []int // [x,y]
	defaultFPS := 60
	{
		spriteSizeComponentRegex, err := regexp.Compile("-\\d+x\\d+")
		if err != nil {
			log.Fatalln(err)
		}
		sizePerSpriteComponent := spriteSizeComponentRegex.FindString(spriteSheetName)
		spriteSizeRegex, err := regexp.Compile("\\d+")
		if err != nil {
			log.Fatalln(err)
		}
		spriteSizesString := spriteSizeRegex.FindAllString(sizePerSpriteComponent, 2)
		for _, sizeString := range spriteSizesString {
			size, err := strconv.Atoi(sizeString)
			if err != nil {
				log.Fatalln(err)
			}
			sizePerSprite = append(sizePerSprite, int(size))
		}
	}
	{
		fpsComponentRegex, err := regexp.Compile("(?:-\\d+x\\d+-)\\d+")
		if err != nil {
			log.Fatalln(err)
		}
		fpsComponent := fpsComponentRegex.FindString(spriteSheetName)
		fpsRegex, err := regexp.Compile("\\d+")
		if err != nil {
			log.Fatalln(err)
		}
		fpsString := fpsRegex.FindString(fpsComponent)
		if fpsString != "" {
			defaultFPS, err = strconv.Atoi(fpsString)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	log.Printf("animating sprite sheet, %s, sprite size: %vx%v, fps: %v\n", spriteSheetFile, sizePerSprite[0], sizePerSprite[1], defaultFPS)

	sprites, err := getSprites(spriteSheetFile, sizePerSprite[0], sizePerSprite[1])
	if err != nil {
		log.Fatalln(err)
	}

	ebiten.SetWindowSize(sizePerSprite[0], sizePerSprite[1])
	ebiten.SetWindowTitle(fmt.Sprintf("Animation playback (%v fps)", defaultFPS))

	if err := ebiten.RunGame(&Animator{
		sprites: sprites,
		fps:     defaultFPS,
		counter: 0,
	}); err != nil {
		log.Fatal(err)
	}
}

func getSprites(file string, spriteXSize, spriteYSize int) ([]*ebiten.Image, error) {
	var spriteSheetXSize, spriteSheetYSize int
	var spriteSheet ebiten.Image

	if reader, err := os.Open(file); err == nil {
		defer reader.Close()
		spriteSheetImage, _, err := image.Decode(reader)
		if err != nil {
			log.Fatalln(err)
		}
		spriteSheetXSize = spriteSheetImage.Bounds().Min.X + spriteSheetImage.Bounds().Max.X
		spriteSheetYSize = spriteSheetImage.Bounds().Min.Y + spriteSheetImage.Bounds().Max.Y
		log.Printf("sprite sheet size: %vx%v\n", spriteSheetXSize, spriteSheetYSize)

		spriteSheet = *ebiten.NewImageFromImage(spriteSheetImage)
	} else {
		log.Fatalln(err)
	}

	numberOfRows := (int)(math.Round((float64)(spriteSheetXSize / spriteXSize)))
	numberOfCols := (int)(math.Round((float64)(spriteSheetYSize / spriteYSize)))

	sprites := []*ebiten.Image{}

	for col := 0; col < numberOfCols; col++ {
		for row := 0; row < numberOfRows; row++ {
			spriteSheetCopy := ebiten.NewImageFromImage(&spriteSheet)

			sprites = append(sprites, ebiten.NewImageFromImage(spriteSheetCopy.SubImage(image.Rect(
				row*spriteXSize, col*spriteYSize,
				(row+1)*spriteXSize, (col+1)*spriteYSize,
			))))
		}
	}

	log.Printf("sprite sheet size: %vx%v, number of sprites: %v, number of Rows: %v, Cols: %v\n", spriteSheetXSize, spriteSheetYSize, len(sprites), numberOfRows, numberOfCols)

	return sprites, nil
}
