package main

import (
	math "github.com/chewxy/math32"
)

type Vector4 struct {
	x, y, z, w float32
}

func (v1 Vector4) Add(v2 Vector4) Vector4 {
	return Vector4{
		v1.x + v2.x,
		v1.y + v2.y,
		v1.z + v2.z,
		1,
	}
}

func (v1 Vector4) Sub(v2 Vector4) Vector4 {
	return Vector4{
		v1.x - v2.x,
		v1.y - v2.y,
		v1.z - v2.z,
		1,
	}
}

func (v Vector4) Mul(k float32) Vector4 {
	return Vector4{
		v.x * k,
		v.y * k,
		v.z * k,
		1,
	}
}

func (v Vector4) Div(k float32) Vector4 {
	return Vector4{
		v.x / k,
		v.y / k,
		v.z / k,
		1,
	}
}

func (v1 Vector4) Dot(v2 Vector4) float32 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func (v Vector4) Len() float32 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector4) Normalise() Vector4 {
	l := v.Len()
	return Vector4{
		v.x / l,
		v.y / l,
		v.z / l,
		1,
	}
}

func (v1 Vector4) CrossProduct(v2 Vector4) Vector4 {
	return Vector4{
		v1.y*v2.z - v1.z*v2.y,
		v1.z*v2.x - v1.x*v2.z,
		v1.x*v2.y - v1.y*v2.x,
		1,
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

	dist := func(p Vector4) float32 {
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

	// Two points of the triangle are outside of view. Clip it into a new triangle.
	if nInsidePointCount == 1 && nOutsidePointCount == 2 {
		outTri := Triangle{
			vecs: [3]Vector4{
				insidePoints[0],
				IntersectPlane(plane_p, planeN, insidePoints[0], outsidePoints[0]),
				IntersectPlane(plane_p, planeN, insidePoints[0], outsidePoints[1]),
			},
			ilum: in_tri.ilum,
		}

		return []Triangle{outTri}
	}

	// One point of the triangle is outside of view. Clip triangle into a quad (two triangles).
	if nInsidePointCount == 2 && nOutsidePointCount == 1 {
		outTri1 := Triangle{
			vecs: [3]Vector4{
				insidePoints[0],
				insidePoints[1],
				IntersectPlane(plane_p, planeN, insidePoints[0], outsidePoints[0]),
			},
			ilum: in_tri.ilum,
		}
		outTri2 := Triangle{
			vecs: [3]Vector4{
				insidePoints[1],
				outTri1.vecs[2],
				IntersectPlane(plane_p, planeN, insidePoints[1], outsidePoints[0]),
			},
			ilum: in_tri.ilum,
		}

		return []Triangle{outTri1, outTri2}
	}

	return make([]Triangle, 0)
}
