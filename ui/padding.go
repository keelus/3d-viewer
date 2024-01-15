package ui

type Padding struct {
	x, y int32
}

func NewPadding(x, y int32) Padding {
	return Padding{x, y}
}
