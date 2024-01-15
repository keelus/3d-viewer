package ui

import "github.com/veandco/go-sdl2/sdl"

type Button struct {
	bX, bY  int32
	bH, bW  int32
	bAnchor Anchor

	pressed bool
	hovered bool

	rect *sdl.Rect

	colorIdle, colorHover, colorPress uint32
}

func NewButton(x, y, w, h int32, margin Margin, anchor Anchor, colorIdle, colorHover, colorPress uint32) Button {
	btn := Button{
		bX:      x,
		bY:      y,
		bW:      w,
		bH:      h,
		bAnchor: anchor,

		pressed:    false,
		hovered:    false,
		colorIdle:  colorIdle,
		colorHover: colorHover,
		colorPress: colorPress,
	}

	btn.rect = GetFinalRect(x, y, w, h, margin, Padding{}, anchor)

	return btn
}

func (b Button) Draw(surface *sdl.Surface) {
	if b.pressed {
		surface.FillRect(b.rect, b.colorPress)
	} else if b.hovered {
		surface.FillRect(b.rect, b.colorHover)
	} else {
		surface.FillRect(b.rect, b.colorIdle)
	}
}

func (b *Button) UpdateAndGetStatus(x, y int32, pressing bool) bool {
	b.hovered = x >= b.rect.X && x <= b.rect.X+b.rect.W && y >= b.rect.Y && y <= b.rect.Y+b.rect.H
	pressed := pressing && b.hovered

	if b.pressed && pressed { // Still clicking, wait until release
		pressed = false
	} else {
		b.pressed = pressed
	}

	return pressed
}
