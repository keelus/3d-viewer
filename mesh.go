package main

type Mesh struct {
	tris                         []Triangle
	triangleAmount, vertexAmount int

	// Lowest and highest vertice values (used to center and offset camera)
	lowestX, highestX float32
	lowestY, highestY float32
	lowestZ, highestZ float32
}
