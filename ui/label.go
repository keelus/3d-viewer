package ui

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Label struct {
	bX, bY  int32
	bAnchor Anchor
	bMargin Margin

	rect *sdl.Rect

	textColor sdl.Color
	textValue string
	rendered  *sdl.Surface
	font      *ttf.Font
}

func NewLabel(x, y int32, text string, margin Margin, anchor Anchor, textColor sdl.Color, font *ttf.Font) Label {
	lbl := Label{
		bX:      x,
		bY:      y,
		bAnchor: anchor,
		bMargin: margin,

		textValue: text,
		textColor: textColor,

		font: font,
	}

	lbl.updateRender()

	return lbl
}

func (lbl *Label) updateRender() {
	rendered, err := lbl.font.RenderUTF8Blended(lbl.textValue, lbl.textColor)
	if err != nil {
		log.Fatalf("Error rendering text")
	}

	lbl.rendered = rendered
	lbl.rect = GetFinalRect(lbl.bX, lbl.bY, rendered.W, rendered.H, lbl.bMargin, Padding{0, 0}, lbl.bAnchor)
}

func (lbl *Label) SetText(text string) {
	lbl.textValue = text
	lbl.updateRender()
}

func (lbl Label) Draw(surface *sdl.Surface) {
	if err := lbl.rendered.Blit(nil, surface, lbl.rect); err != nil {
		log.Fatalf("Error rendering a label.")
	}
}

func (lbl Label) GetRectWidth() int32 {
	return lbl.rect.W
}
