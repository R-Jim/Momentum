package ui

import (
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type UI struct {
	titleFont font.Face
	font      font.Face
	smallFont font.Face

	TitleFontSize float64
	FontSize      float64
	SmallFontSize float64
}

func New(fontStyle []byte, fontSize float64) UI {
	tt, err := opentype.Parse(fontStyle)
	if err != nil {
		log.Fatal(err)
	}

	titleFontSize := fontSize * 1.5
	smallFontSize := fontSize / 3

	const dpi = 72
	titleFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    titleFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	smallFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    smallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	font, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return UI{
		titleFont: titleFont,
		font:      font,
		smallFont: smallFont,

		TitleFontSize: titleFontSize,
		FontSize:      fontSize,
		SmallFontSize: smallFontSize,
	}
}
