package ui

import "github.com/veandco/go-sdl2/sdl"

type Anchor int

const (
	TOP_LEFT Anchor = iota
	TOP_RIGHT
	TOP_CENTER
	CENTER_CENTER
	BOTTOM_LEFT
	BOTTOM_RIGHT
	BOTTOM_CENTER
)

func GetFinalRect(x, y, w, h int32, margin Margin, padding Padding, anchor Anchor) *sdl.Rect {
	var fx, fy int32
	switch anchor {
	case TOP_LEFT:
		fx, fy = x+margin.x, y+margin.y
	case TOP_RIGHT:
		fx, fy = x-w-padding.x*2-margin.x, y+margin.y
	case TOP_CENTER:
		fx, fy = x-w/2-padding.x-margin.x/2, y+margin.y
	case CENTER_CENTER:
		fx, fy = x-w/2-padding.x-margin.x, y-h/h-padding.y-margin.y
	case BOTTOM_LEFT:
		fx, fy = x+margin.x*2, y-h-padding.y*2-margin.y
	case BOTTOM_RIGHT:
		fx, fy = x-w-padding.x*2-margin.x, y-h-padding.y*2-margin.y
	case BOTTOM_CENTER:
		fx, fy = x-w/2-padding.x-margin.x/2, y-h-padding.y*2-margin.y
	}

	return &sdl.Rect{X: fx, Y: fy, W: w + padding.x*2, H: h + padding.y*2}
}
