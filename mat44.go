package main

import (
	"math"
)

type mat44 struct {
	m [4][4]float64
}

func (m mat44) multiplyVector(i Vector4) Vector4 {
	return Vector4{
		i.x*m.m[0][0] + i.y*m.m[1][0] + i.z*m.m[2][0] + i.w*m.m[3][0],
		i.x*m.m[0][1] + i.y*m.m[1][1] + i.z*m.m[2][1] + i.w*m.m[3][1],
		i.x*m.m[0][2] + i.y*m.m[1][2] + i.z*m.m[2][2] + i.w*m.m[3][2],
		i.x*m.m[0][3] + i.y*m.m[1][3] + i.z*m.m[2][3] + i.w*m.m[3][3], -1,
		NewTexVector(
			i.texVec.u, i.texVec.v,
			i.texVec.w,
		),
	}
}

func (m mat44) multiplyTriangle(t Triangle) Triangle {
	return Triangle{
		vecs: [3]Vector4{
			m.multiplyVector(t.vecs[0]),
			m.multiplyVector(t.vecs[1]),
			m.multiplyVector(t.vecs[2]),
		},
	}
}

func identityMatrix() mat44 {
	return mat44{
		m: [4][4]float64{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		},
	}
}

func rotationMatrix(rotationRads Vector4) mat44 {
	gamma := rotationRads.x
	beta := rotationRads.y
	alpha := rotationRads.z

	mat := mat44{
		m: [4][4]float64{
			{math.Cos(alpha) * math.Cos(beta), math.Cos(alpha)*math.Sin(beta)*math.Sin(gamma) - math.Sin(alpha)*math.Cos(gamma), math.Cos(alpha)*math.Sin(beta)*math.Cos(gamma) + math.Sin(alpha)*math.Sin(gamma), 0},
			{math.Sin(alpha) * math.Cos(beta), math.Sin(alpha)*math.Sin(beta)*math.Sin(gamma) + math.Cos(alpha)*math.Cos(gamma), math.Sin(alpha)*math.Sin(beta)*math.Cos(gamma) - math.Cos(alpha)*math.Sin(gamma), 0},
			{-math.Sin(beta), math.Cos(beta) * math.Sin(gamma), math.Cos(beta) * math.Cos(gamma), 0},
			{0, 0, 0, 1},
		},
	}

	return mat
}

func MakeTranslation(x, y, z float64) mat44 {
	return mat44{
		m: [4][4]float64{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 1, 0},
			{x, y, z, 1},
		},
	}
}

func projectionMatrix(aspectRatio, fovDeg, nearDist, farDist float64) mat44 {
	fovRad := 1 / math.Tan(degToRad(fovDeg/2))
	return mat44{
		m: [4][4]float64{
			{aspectRatio * fovRad, 0, 0, 0},
			{0, fovRad, 0, 0},
			{0, 0, farDist / (farDist - nearDist), -1},
			{0, 0, (-farDist * nearDist) / (farDist - nearDist), 0},
		},
	}
}

func (m1 mat44) multiplyMatrix(m2 mat44) mat44 {
	mat := mat44{}

	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			mat.m[r][c] = m1.m[r][0]*m2.m[0][c] + m1.m[r][1]*m2.m[1][c] + m1.m[r][2]*m2.m[2][c] + m1.m[r][3]*m2.m[3][c]
		}
	}

	return mat
}
