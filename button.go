package main

import "github.com/veandco/go-sdl2/sdl"

type Button struct {
	x, y, w, h int32

	hover bool

	rect *sdl.Rect

	colorIdle, colorHover uint32
}

func NewButton(x, y, w, h int32, colorIdle, colorHover uint32) Button {
	return Button{
		x, y, w, h,
		false,
		&sdl.Rect{x, y, w, h},
		colorIdle, colorHover,
	}
}

func (b Button) Draw() {
	if !b.hover {
		sur.FillRect(b.rect, b.colorIdle)
	} else {
		sur.FillRect(b.rect, b.colorHover)
	}
}

func (b *Button) Hover(x, y int32) bool {
	b.hover = x >= b.x && x <= b.x+b.w && y >= b.y && y <= b.y+b.h
	return b.hover
}
