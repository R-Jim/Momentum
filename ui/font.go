package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

func drawMsg(screen *ebiten.Image, msg []string, font font.Face, fontSize, x, y int) {
	for i, l := range msg {
		text.Draw(screen, l, font, x, y+(i+4)*fontSize, color.NRGBA{
			R: 0x00 + MAIN_TEXT.R,
			G: 0x00 + MAIN_TEXT.G,
			B: 0x00 + MAIN_TEXT.B,
			A: 0x00 + MAIN_TEXT.A,
		})
	}
}

func drawMsgFromBottom(screen *ebiten.Image, msg []string, font font.Face, fontSize, x, y int) {
	for i, l := range msg {
		text.Draw(screen, l, font, x, y-(i+4)*fontSize, color.NRGBA{
			R: 0x00 + MAIN_TEXT.R,
			G: 0x00 + MAIN_TEXT.G,
			B: 0x00 + MAIN_TEXT.B,
			A: 0x00 + MAIN_TEXT.A,
		})
	}
}

func (ui UI) DrawMsg(screen *ebiten.Image, msg []string, x, y int, fromBottom bool) {
	if fromBottom {
		drawMsgFromBottom(screen, msg, ui.font, int(ui.FontSize), x, y)
	} else {
		drawMsg(screen, msg, ui.font, int(ui.FontSize), x, y)
	}
}

func (ui UI) DrawTitle(screen *ebiten.Image, msg []string, x, y int, fromBottom bool) {
	if fromBottom {
		drawMsgFromBottom(screen, msg, ui.titleFont, int(ui.TitleFontSize), x, y)
	} else {
		drawMsg(screen, msg, ui.titleFont, int(ui.TitleFontSize), x, y)
	}
}

func (ui UI) DrawSmallMsg(screen *ebiten.Image, msg []string, x, y int, fromBottom bool) {
	if fromBottom {
		drawMsgFromBottom(screen, msg, ui.smallFont, int(ui.SmallFontSize), x, y)
	} else {
		drawMsg(screen, msg, ui.smallFont, int(ui.SmallFontSize), x, y)
	}
}
