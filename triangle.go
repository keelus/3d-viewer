package main

import (
	"image/color"
	"log"

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

		FillTriangle(t.vecs[2], t.vecs[1], t.vecs[0], c)

	}

}

func PutPixel(x, y int, z float32, u, v float32) {
	zIdx := y*int(SCREEN_WIDTH) + x

	if zIdx >= 0 && zIdx < len(depthBuffer) {
		if z < depthBuffer[zIdx] {
			idx := 4 * (y*int(SCREEN_WIDTH) + x)
			//col := colorFromZ(z)
			c := bricks.GetColorAt(u, v)
			writePixel(idx+0, c.B) // B
			writePixel(idx+1, c.G) // G
			writePixel(idx+2, c.R) // R
			writePixel(idx+3, 255) // A

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

func DrawPoint(v Vector4, c color.RGBA) {
	PutPixel(int(v.x), int(v.y), v.originalZ, v.u, v.v)
}

func GetSlope(vA, vB Vector4) float32 {
	return (vB.x - vA.x) / (vB.y - vA.y)
}

func edge_cross(a, b, p Vector4) float32 {
	ab := Vector4{b.x - a.x, b.y - a.y, 0, 0, 0, 0, 0}
	ap := Vector4{p.x - a.x, p.y - a.y, 0, 0, 0, 0, 0}
	return ab.x*ap.y - ab.y*ap.x
}

func FillTriangle(v0, v1, v2 Vector4, c color.RGBA) {
	xMin := math.Min(math.Min(v0.x, v1.x), v2.x)
	yMin := math.Min(math.Min(v0.y, v1.y), v2.y)
	xMax := math.Max(math.Max(v0.x, v1.x), v2.x)
	yMax := math.Max(math.Max(v0.y, v1.y), v2.y)

	area := edge_cross(v0, v1, v2)

	for y := yMin; y <= yMax; y++ {
		for x := xMin; x <= xMax; x++ {
			p := Vector4{x, y, 0, 0, 0, 0, 0}

			w0 := edge_cross(v1, v2, p)
			w1 := edge_cross(v2, v0, p)
			w2 := edge_cross(v0, v1, p)

			alpha := w0 / area
			beta := w1 / area
			gamma := w2 / area

			p.z = alpha*v0.z + beta*v1.z + gamma*v2.z
			p.u = math.Abs(alpha*v0.u + beta*v1.u + gamma*v2.u)
			p.v = math.Abs(alpha*v0.v + beta*v1.v + gamma*v2.v)

			isInside := w0 >= 0 && w1 >= 0 && w2 >= 0

			if isInside {
				DrawPoint(p, c)
			}
		}
	}
}
