package main

import (
	"math"
)

type TexVector struct {
	u, v, w float64
}

func NewTexVector(u, v, w float64) TexVector {
	return TexVector{u, v, w}
}

type Vector4 struct {
	x, y, z, w float64
	originalZ  float64

	texVec TexVector
}

func (v1 Vector4) Add(v2 Vector4) Vector4 {
	return Vector4{
		v1.x + v2.x,
		v1.y + v2.y,
		v1.z + v2.z,
		1,
		v1.originalZ + v2.originalZ,
		NewTexVector(
			v1.texVec.u+v2.texVec.v,
			v1.texVec.v+v2.texVec.v,
			v1.texVec.w+v2.texVec.w,
		),
	}
}

func (v1 Vector4) Sub(v2 Vector4) Vector4 {
	return Vector4{
		v1.x - v2.x,
		v1.y - v2.y,
		v1.z - v2.z,
		1,
		v1.originalZ - v2.originalZ,
		NewTexVector(
			v1.texVec.u-v2.texVec.v,
			v1.texVec.v-v2.texVec.v,
			v1.texVec.w-v2.texVec.w,
		),
	}
}

func (v Vector4) Mul(k float64) Vector4 {
	return Vector4{
		v.x * k,
		v.y * k,
		v.z * k,
		1,
		v.originalZ * k,
		NewTexVector(
			v.texVec.u,
			v.texVec.v,
			v.texVec.w,
		),
	}
}

func (v Vector4) Div(k float64) Vector4 {
	return Vector4{
		v.x / k,
		v.y / k,
		v.z / k,
		1,
		v.originalZ / k,
		NewTexVector(
			v.texVec.u,
			v.texVec.v,
			v.texVec.w,
		),
	}
}

func (v1 Vector4) Dot(v2 Vector4) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func (v Vector4) Len() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector4) Normalise() Vector4 {
	l := v.Len()
	return Vector4{
		v.x / l,
		v.y / l,
		v.z / l,
		1,
		v.originalZ,
		NewTexVector(
			v.texVec.u,
			v.texVec.v,
			v.texVec.w,
		),
	}
}

func (v1 Vector4) CrossProduct(v2 Vector4) Vector4 {
	return Vector4{
		v1.y*v2.z - v1.z*v2.y,
		v1.z*v2.x - v1.x*v2.z,
		v1.x*v2.y - v1.y*v2.x,
		1,
		v1.originalZ,
		NewTexVector(
			v1.texVec.u,
			v1.texVec.v,
			v1.texVec.w,
		),
	}
}

func IntersectPlane(plane_p, plane_n, lineStart, lineEnd Vector4) Vector4 {
	planeN := plane_n.Normalise()
	planeD := -planeN.Dot(plane_p)

	ad := lineStart.Dot(planeN)
	bd := lineEnd.Dot(planeN)
	t := (-planeD - ad) / (bd - ad)

	lineStartToEnd := lineEnd.Sub(lineStart)
	lineToIntersect := lineStartToEnd.Mul(t)

	return lineStart.Add(lineToIntersect)
}

func ClipAgainstPlane(plane_p, plane_n Vector4, in_tri Triangle) []Triangle {
	planeN := plane_n.Normalise()

	dist := func(p Vector4) float64 {
		return planeN.x*p.x + planeN.y*p.y + planeN.z*p.z - planeN.Dot(plane_p)
	}

	var insidePoints, outsidePoints [3]Vector4
	nInsidePointCount, nOutsidePointCount := 0, 0

	for i := 0; i < 3; i++ {
		distance := dist(in_tri.vecs[i])

		if distance >= 0 {
			insidePoints[nInsidePointCount] = in_tri.vecs[i]
			nInsidePointCount++
		} else {
			outsidePoints[nOutsidePointCount] = in_tri.vecs[i]
			nOutsidePointCount++
		}
	}

	// The entire triangle is outside of view. No need to draw it.
	if nInsidePointCount == 0 {
		return make([]Triangle, 0)
	}

	// The entire triangle is completly inside of view. Draw it as is.
	if nInsidePointCount == 3 {
		return []Triangle{in_tri}
	}

	v0 := &in_tri.vecs[0]
	v1 := &in_tri.vecs[1]
	v2 := &in_tri.vecs[2]

	area := EdgeCross(v0, v1, v2)

	updateVector := func(v *Vector4) {
		_, alpha, beta, gamma := BaycentricPointInTriangle(area, v0, v1, v2, v)
		InterpolateVectors(alpha, beta, gamma, v0, v1, v2, v)
	}

	// Two points of the triangle are outside of view. Clip it into a new triangle.
	if nInsidePointCount == 1 && nOutsidePointCount == 2 {
		p1 := IntersectPlane(plane_p, planeN, insidePoints[0], outsidePoints[0])
		p2 := IntersectPlane(plane_p, planeN, insidePoints[0], outsidePoints[1])

		updateVector(&p1)
		updateVector(&p2)

		outTri := Triangle{
			vecs: [3]Vector4{
				insidePoints[0],
				p1,
				p2,
			},
			ilum: in_tri.ilum,
			tex:  in_tri.tex,
		}

		return []Triangle{outTri}
	}

	// One point of the triangle is outside of view. Clip triangle into a quad (two triangles).
	if nInsidePointCount == 2 && nOutsidePointCount == 1 {
		p1 := IntersectPlane(plane_p, planeN, insidePoints[0], outsidePoints[0])
		updateVector(&p1)

		outTri1 := Triangle{
			vecs: [3]Vector4{
				insidePoints[0],
				insidePoints[1],
				p1,
			},
			ilum: in_tri.ilum,
			tex:  in_tri.tex,
		}

		p2 := IntersectPlane(plane_p, planeN, insidePoints[1], outsidePoints[0])
		updateVector(&p2)

		outTri2 := Triangle{
			vecs: [3]Vector4{
				insidePoints[1],
				outTri1.vecs[2],
				p2,
			},
			ilum: in_tri.ilum,
			tex:  in_tri.tex,
		}

		return []Triangle{outTri1, outTri2}
	}

	return make([]Triangle, 0)
}
