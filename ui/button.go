package ui

import "github.com/veandco/go-sdl2/sdl"

type Button struct {
	bX, bY  int32
	bH, bW  int32
	bAnchor Anchor

	hover bool

	rect *sdl.Rect

	colorIdle, colorHover uint32
}

func NewButton(x, y, w, h int32, margin Margin, anchor Anchor, colorIdle, colorHover uint32) Button {
	btn := Button{
		bX:      x,
		bY:      y,
		bW:      w,
		bH:      h,
		bAnchor: anchor,

		hover:      false,
		colorIdle:  colorIdle,
		colorHover: colorHover,
	}

	btn.rect = GetFinalRect(x, y, w, h, margin, Padding{}, anchor)

	return btn
}

func (b Button) Draw(surface *sdl.Surface) {
	if !b.hover {
		surface.FillRect(b.rect, b.colorIdle)
	} else {
		surface.FillRect(b.rect, b.colorHover)
	}
}

func (b *Button) Hover(x, y int32) bool {
	b.hover = x >= b.rect.X && x <= b.rect.X+b.rect.W && y >= b.rect.Y && y <= b.rect.Y+b.rect.H
	return b.hover
}
