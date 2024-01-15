package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ContentBlock struct {
	bW, bH   int32
	bX, bY   int32
	bAnchor  Anchor
	bMargin  Margin
	bPadding Padding

	rect            *sdl.Rect
	backgroundColor uint32
}

func NewContentBlock(x, y, w, h int32, margin Margin, padding Padding, anchor Anchor, backgroundColor uint32) ContentBlock {
	return ContentBlock{
		bW:       w,
		bH:       h,
		bX:       x,
		bY:       y,
		bMargin:  margin,
		bPadding: padding,

		bAnchor: anchor,

		rect:            GetFinalRect(x, y, w, h, margin, padding, anchor),
		backgroundColor: backgroundColor,
	}
}

func (cb ContentBlock) Draw(surface *sdl.Surface) {
	surface.FillRect(cb.rect, cb.backgroundColor)
}

func (cb *ContentBlock) UpdateRectToWidth(width int32) {
	cb.rect = GetFinalRect(cb.bX, cb.bY, width, cb.bH, cb.bMargin, cb.bPadding, cb.bAnchor)
}
