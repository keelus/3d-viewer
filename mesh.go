package main

type Mesh struct {
	tris                         []Triangle
	triangleAmount, vertexAmount int

	// Lowest and highest vertice values (used to center and offset camera)
	lowestX, highestX float64
	lowestY, highestY float64
	lowestZ, highestZ float64
}
