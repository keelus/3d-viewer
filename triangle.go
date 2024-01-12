package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Triangle struct {
	vecs [3]Vector4
	ilum float32
}

const (
	TRIANGLE_OUTLINE bool = true
	TRIANGLE_FILL    bool = true
)

var (
	TRIANGLE_OUTLINE_COLOR color.RGBA = color.RGBA{255, 255, 255, 255}
	TRIANGLE_FILL_COLOR    color.RGBA = color.RGBA{255, 255, 255, 255}
)

func (t Triangle) Draw() {
	if TRIANGLE_OUTLINE {
		drawLine := func(v1, v2 Vector4) {
			rl.DrawLineEx(
				rl.NewVector2(v1.x, v1.y),
				rl.NewVector2(v2.x, v2.y),
				2,
				TRIANGLE_OUTLINE_COLOR,
			)
		}

		drawLine(t.vecs[0], t.vecs[1])
		drawLine(t.vecs[0], t.vecs[2])
		drawLine(t.vecs[1], t.vecs[2])
	}

	if TRIANGLE_FILL {
		var ilum float32 = 0.3 + (1-0.3)*((t.ilum-0)/(1-0))

		finalColor := color.RGBA{
			uint8(float32(TRIANGLE_FILL_COLOR.R) * ilum),
			uint8(float32(TRIANGLE_FILL_COLOR.G) * ilum),
			uint8(float32(TRIANGLE_FILL_COLOR.B) * ilum),
			TRIANGLE_FILL_COLOR.A,
		}

		rl.DrawTriangle(
			rl.NewVector2(t.vecs[0].x, t.vecs[0].y),
			rl.NewVector2(t.vecs[1].x, t.vecs[1].y),
			rl.NewVector2(t.vecs[2].x, t.vecs[2].y),
			finalColor,
		)

		// Added an inverted triangle to fix coloring issues on some faces
		// Issue happens on fill, but should print without inversing. Raylib FillTriangle issue?
		rl.DrawTriangle(
			rl.NewVector2(t.vecs[2].x, t.vecs[2].y),
			rl.NewVector2(t.vecs[1].x, t.vecs[1].y),
			rl.NewVector2(t.vecs[0].x, t.vecs[0].y),
			color.RGBA{
				uint8(float32(TRIANGLE_FILL_COLOR.R) * ilum),
				uint8(float32(TRIANGLE_FILL_COLOR.G) * ilum),
				uint8(float32(TRIANGLE_FILL_COLOR.B) * ilum),
				finalColor.B,
			},
		)
	}

}
