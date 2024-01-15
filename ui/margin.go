package ui

type Margin struct {
	x, y int32
}

func NewMargin(x, y int32) Margin {
	return Margin{x, y}
}
