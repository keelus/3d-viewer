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
	tex  *Texture
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
	if TRIANGLE_FILL {
		FillTriangle(&t.vecs[0], &t.vecs[1], &t.vecs[2], t.tex)
		FillTriangle(&t.vecs[2], &t.vecs[1], &t.vecs[0], t.tex)
	}

}

func PutPixel(x, y int, z float32, u, v float32, tex *Texture) {
	zIdx := y*int(SCREEN_WIDTH) + x

	if zIdx >= 0 && zIdx < len(depthBuffer) {
		if z < depthBuffer[zIdx] {
			idx := 4 * (y*int(SCREEN_WIDTH) + x)
			c := color.RGBA{255, 0, 255, 255}
			if tex != nil {
				c = tex.GetColorAt(u, v)
			}
			writePixel(idx+0, c.B) // B
			writePixel(idx+1, c.G) // G
			writePixel(idx+2, c.R) // R
			writePixel(idx+3, c.A) // A

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

func DrawPoint(v *Vector4, tex *Texture) {
	PutPixel(int(v.x), int(v.y), v.originalZ, v.texVec.u, v.texVec.v, tex)
}

func GetSlope(vA, vB Vector4) float32 {
	return (vB.x - vA.x) / (vB.y - vA.y)
}

func EdgeCross(a, b, p *Vector4) float32 {
	return (b.x-a.x)*(p.y-a.y) - (b.y-a.y)*(p.x-a.x)
}

func FillTriangle(v0, v1, v2 *Vector4, tex *Texture) {
	xMin := math.Min(math.Min(v0.x, v1.x), v2.x)
	yMin := math.Min(math.Min(v0.y, v1.y), v2.y)
	xMax := math.Max(math.Max(v0.x, v1.x), v2.x)
	yMax := math.Max(math.Max(v0.y, v1.y), v2.y)

	area := EdgeCross(v0, v1, v2)

	for y := yMin; y <= yMax; y++ {
		for x := xMin; x <= xMax; x++ {
			p := &Vector4{x, y, 0, 0, 0, NewTexVector(0, 0, 0)}

			w0 := EdgeCross(v1, v2, p)
			w1 := EdgeCross(v2, v0, p)
			w2 := EdgeCross(v0, v1, p)

			alpha := w0 / area
			beta := w1 / area
			gamma := w2 / area

			p.z = alpha*v0.z + beta*v1.z + gamma*v2.z
			p.texVec.u = alpha*v0.texVec.u + beta*v1.texVec.u + gamma*v2.texVec.u
			p.texVec.v = alpha*v0.texVec.v + beta*v1.texVec.v + gamma*v2.texVec.v
			p.texVec.w = alpha*v0.texVec.w + beta*v1.texVec.w + gamma*v2.texVec.w
			p.originalZ = alpha*v0.originalZ + beta*v1.originalZ + gamma*v2.originalZ

			p.texVec.u /= p.texVec.w
			p.texVec.v /= p.texVec.w

			isInside := w0 >= 0 && w1 >= 0 && w2 >= 0

			if isInside {
				DrawPoint(p, tex)
			}
		}
	}
}
