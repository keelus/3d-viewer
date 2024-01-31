package main

import (
	"image/color"

	"math"
)

type Triangle struct {
	vecs [3]Vector4
	ilum float64
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

func getZ(x1, y1, z1, x2, y2, z2, x, y float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	dz := z2 - z1

	var z float64

	if dx != 0 {
		z = ((x-x1)/dx)*dz + z1
	} else if dy != 0 {
		z = ((y-y1)/dy)*dz + z1
	} else {
		z = 0
	}

	return z
}

func IsClockWise(a, b, c *Vector4) bool {
	return EdgeCross(a, b, c) < 0
}

func (t Triangle) Draw() {
	if TRIANGLE_FILL {
		// Fill triangle expects the order to be counter-clockwise (because of EdgeCross order)
		if IsClockWise(&t.vecs[0], &t.vecs[1], &t.vecs[2]) {
			FillTriangle(&t.vecs[2], &t.vecs[1], &t.vecs[0], t.tex)
		} else {
			FillTriangle(&t.vecs[0], &t.vecs[1], &t.vecs[2], t.tex)
		}
	}

	if TRIANGLE_OUTLINE {
	}
}

func PutPixel(p *Vector4, tex *Texture) {
	fx, fy := int((p.x)), int((p.y))

	zIdx := fy*RENDER_WIDTH + fx

	if zIdx >= 0 && zIdx < depthBufferLength {
		if p.originalZ < depthBuffer[zIdx] {
			idx := 4 * (fy*RENDER_WIDTH + fx)
			c := color.RGBA{255, 0, 255, 255}
			if tex != nil {
				c = tex.GetColorAt(p.texVec.u, p.texVec.v)
				if c.A == 0 {
					return
				}
			}

			renderBuffer[idx+0] = c.B
			renderBuffer[idx+1] = c.G
			renderBuffer[idx+2] = c.R
			renderBuffer[idx+3] = c.A

			depthBuffer[zIdx] = p.originalZ
		}
	}
}

func DrawPoint(v *Vector4, tex *Texture) {
	PutPixel(v, tex)
}

func GetSlope(vA, vB Vector4) float64 {
	return (vB.x - vA.x) / (vB.y - vA.y)
}

func EdgeCross(a, b, p *Vector4) float64 {
	return (b.x-a.x)*(p.y-a.y) - (b.y-a.y)*(p.x-a.x)
}

func BaycentricPointInTriangle(area float64, v0, v1, v2, p *Vector4) (bool, float64, float64, float64) {
	w0 := EdgeCross(v1, v2, p)
	w1 := EdgeCross(v2, v0, p)
	w2 := EdgeCross(v0, v1, p)

	alpha := w0 / area
	beta := w1 / area
	gamma := w2 / area

	isInside := w0 >= 0 && w1 >= 0 && w2 >= 0

	return isInside, alpha, beta, gamma
}

func isTopLeft(start, end *Vector4) bool {
	edge := Vector4{x: end.x - start.x, y: end.y - start.y}

	isTopEdge := edge.y == 0 && edge.x > 0
	isLeftEdge := edge.y < 0

	return isTopEdge || isLeftEdge
}

func FillTriangle(v0, v1, v2 *Vector4, tex *Texture) {
	xMin := math.Floor(math.Min(math.Min(v0.x, v1.x), v2.x))
	yMin := math.Floor(math.Min(math.Min(v0.y, v1.y), v2.y))
	xMax := math.Ceil(math.Max(math.Max(v0.x, v1.x), v2.x))
	yMax := math.Ceil(math.Max(math.Max(v0.y, v1.y), v2.y))

	deltaW0Col := v1.y - v2.y
	deltaW1Col := v2.y - v0.y
	deltaW2Col := v0.y - v1.y

	deltaW0Row := v2.x - v1.x
	deltaW1Row := v0.x - v2.x
	deltaW2Row := v1.x - v0.x

	var bias0, bias1, bias2 float64
	if isTopLeft(v1, v2) {
		bias0 = 0
	} else {
		bias0 = -0.00001
	}
	if isTopLeft(v2, v0) {
		bias1 = 0
	} else {
		bias1 = -0.00001
	}
	if isTopLeft(v0, v1) {
		bias2 = 0
	} else {
		bias2 = -0.00001
	}

	area := EdgeCross(v0, v1, v2)

	p := &Vector4{xMin + 0.5, yMin + 0.5, 0, 0, 0, NewTexVector(0, 0, 0)}

	w0Row := EdgeCross(v1, v2, p) + bias0
	w1Row := EdgeCross(v2, v0, p) + bias1
	w2Row := EdgeCross(v0, v1, p) + bias2

	for y := yMin; y <= yMax; y++ {
		w0 := w0Row
		w1 := w1Row
		w2 := w2Row
		for x := xMin; x <= xMax; x++ {
			isInside := w0 >= 0 && w1 >= 0 && w2 >= 0

			if isInside {
				p.x = x
				p.y = y

				alpha := w0 / area
				beta := w1 / area
				gamma := w2 / area

				p.z = alpha*v0.z + beta*v1.z + gamma*v2.z
				p.texVec.u = alpha*v0.texVec.u + beta*v1.texVec.u + gamma*v2.texVec.u
				p.texVec.v = alpha*v0.texVec.v + beta*v1.texVec.v + gamma*v2.texVec.v
				p.texVec.w = alpha*v0.texVec.w + beta*v1.texVec.w + gamma*v2.texVec.w
				p.originalZ = alpha*v0.originalZ + beta*v1.originalZ + gamma*v2.originalZ

				if p.texVec.w != 0 {
					p.texVec.u /= p.texVec.w
					p.texVec.v /= p.texVec.w
				}

				PutPixel(p, tex)
			}
			w0 += deltaW0Col
			w1 += deltaW1Col
			w2 += deltaW2Col
		}
		w0Row += deltaW0Row
		w1Row += deltaW1Row
		w2Row += deltaW2Row
	}
}

func InterpolateVectors(alpha, beta, gamma float64, v0, v1, v2, target *Vector4) {
	target.z = alpha*v0.z + beta*v1.z + gamma*v2.z
	target.texVec.u = alpha*v0.texVec.u + beta*v1.texVec.u + gamma*v2.texVec.u
	target.texVec.v = alpha*v0.texVec.v + beta*v1.texVec.v + gamma*v2.texVec.v
	target.texVec.w = alpha*v0.texVec.w + beta*v1.texVec.w + gamma*v2.texVec.w
	target.originalZ = alpha*v0.originalZ + beta*v1.originalZ + gamma*v2.originalZ
}
