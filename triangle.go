package main

import (
	"image/color"
	"log"
	"slices"

	math "github.com/chewxy/math32"
	colorconv "github.com/crazy3lf/colorconv"
)

type Triangle struct {
	vecs [3]Vector4
	ilum float32
}

const (
	TRIANGLE_OUTLINE bool = false
	TRIANGLE_FILL    bool = true
)

var (
	TRIANGLE_OUTLINE_COLOR color.RGBA = color.RGBA{255, 255, 255, 255}
	TRIANGLE_FILL_COLOR    color.RGBA = color.RGBA{255, 255, 255, 255}
)

func getZ(x1, y1, z1, x2, y2, z2, x, y float32) float32 {
	dx := x2 - x1
	dy := y2 - y1
	dz := z2 - z1

	var z float32

	if dx != 0 {
		z = ((x-x1)/dx)*dz + z1
	} else if dy != 0 {
		z = ((y-y1)/dy)*dz + z1
	} else {
		z = 0
	}

	return z
}

func (t Triangle) Draw() {
	var c color.RGBA
	var ilum float32 = 0.3 + (1-0.3)*((t.ilum-0)/(1-0))
	if TRIANGLE_FILL {
		c = color.RGBA{
			uint8(float32(TRIANGLE_FILL_COLOR.R) * ilum),
			uint8(float32(TRIANGLE_FILL_COLOR.G) * ilum),
			uint8(float32(TRIANGLE_FILL_COLOR.B) * ilum),
			TRIANGLE_FILL_COLOR.A,
		}
		vecs := []Vector4{t.vecs[0], t.vecs[1], t.vecs[2]}

		slices.SortFunc(vecs, func(vA, vB Vector4) int {
			if vA.y > vB.y {
				return 1
			} else if vA.y < vB.y {
				return -1
			}

			return 0
		})

		v1 := vecs[0]
		v2 := vecs[1]
		v3 := vecs[2]

		x1, y1, z1 := math.Round(v1.x), math.Round(v1.y), v1.originalZ
		x2, y2, z2 := math.Round(v2.x), math.Round(v2.y), v2.originalZ
		x3, y3, z3 := math.Round(v3.x), math.Round(v3.y), v3.originalZ

		slopeA := (x2 - x1) / (y2 - y1)
		slopeB := (x3 - x1) / (y3 - y1)
		slopeC := (x3 - x2) / (y3 - y2)

		// Draw top triangle:
		for y := y1; y <= y2; y++ {
			xA := x1 + (y-y1)*slopeA
			xB := x1 + (y-y1)*slopeB

			zA := getZ(x1, y1, z1, x2, y2, z2, xA, y)
			zB := getZ(x1, y1, z1, x3, y3, z3, xB, y)

			DrawLine(xA, y, zA, xB, y, zB, c)
		}

		// Draw bottom triangle:
		for y := y2; y <= y3; y++ {
			xB := x3 - (y3-y)*slopeB
			xC := x3 - (y3-y)*slopeC

			zB := getZ(x1, y1, z1, x3, y3, z3, xB, y)
			zC := getZ(x2, y2, z2, x3, y3, z3, xC, y)

			DrawLine(xB, y, zB, xC, y, zC, c)
		}
	}

	c = color.RGBA{255, 255, 255, 255}
	if TRIANGLE_OUTLINE {
		DrawLine(t.vecs[0].x, t.vecs[0].y, t.vecs[0].originalZ, t.vecs[1].x, t.vecs[1].y, t.vecs[1].originalZ, c)
		DrawLine(t.vecs[0].x, t.vecs[0].y, t.vecs[0].originalZ, t.vecs[2].x, t.vecs[2].y, t.vecs[2].originalZ, c)
		DrawLine(t.vecs[1].x, t.vecs[1].y, t.vecs[1].originalZ, t.vecs[2].x, t.vecs[2].y, t.vecs[2].originalZ, c)
	}
}

// Bresenham algorithm
func DrawLine(xf0, yf0, zf0, xf1, yf1, zf1 float32, c color.RGBA) {
	x0, y0, x1, y1 := int(xf0), int(yf0), int(xf1), int(yf1)

	dy := y1 - y0
	dx := x1 - x0

	var stepx, stepy int

	if dy < 0 {
		dy = -dy
		stepy = -1
	} else {
		stepy = 1
	}

	if dx < 0 {
		dx = -dx
		stepx = -1
	} else {
		stepx = 1
	}

	dy <<= 1
	dx <<= 1

	if dx > dy {
		fraction := dy - (dx >> 1)
		for x0 != x1 {
			if fraction >= 0 {
				y0 += stepy
				fraction -= dx
			}
			x0 += stepx
			fraction += dy
			z0 := getZ(xf0, yf0, zf0, xf1, yf1, zf1, float32(x0), float32(y0))
			PutPixel(x0, y0, z0, c)
		}
	} else {
		fraction := dx - (dy >> 1)
		for y0 != y1 {
			if fraction >= 0 {
				x0 += stepx
				fraction -= dy
			}
			y0 += stepy
			fraction += dx
			z0 := getZ(xf0, yf0, zf0, xf1, yf1, zf1, float32(x0), float32(y0))
			PutPixel(x0, y0, z0, c)
		}
	}
}

func PutPixel(x, y int, z float32, c color.RGBA) {
	zIdx := y*int(SCREEN_WIDTH) + x

	if zIdx >= 0 && zIdx < len(depthBuffer) {
		if z < depthBuffer[zIdx] {
			idx := 4 * (y*int(SCREEN_WIDTH) + x)
			//col := colorFromZ(z)
			col := c
			writePixel(idx+0, col.B) // B
			writePixel(idx+1, col.G) // G
			writePixel(idx+2, col.R) // R
			writePixel(idx+3, 255)   // A

			depthBuffer[zIdx] = z
		}
	}
}

func writePixel(idx int, val byte) {
	if idx >= 0 && idx < len(pixelBuffer) {
		pixelBuffer[idx] = val
	}
}

func colorFromZ(z float32) color.RGBA {
	var val float32 = math.Abs(0 + (1-0)*((z-4)/(10-4))*360)

	if val == 360 {
		val = 359
	}

	//log.Print(val)

	r, g, b, err := colorconv.HSLToRGB(float64(val), 1, 0.5)
	if err != nil {
		log.Print(err)
	}

	return color.RGBA{
		r,
		g,
		b,
		255,
	}
}
