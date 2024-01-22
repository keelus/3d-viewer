package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Mesh struct {
	tris                         []Triangle
	triangleAmount, vertexAmount int

	// Lowest and highest vertice values (used to center and offset camera)
	lowestX, highestX float32
	lowestY, highestY float32
	lowestZ, highestZ float32
}

func LoadMesh(filename string) *Mesh {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		log.Fatalf("Error loading the obj file '%s'", filename)
		return &Mesh{}
	}

	mesh := Mesh{}
	vertices := []Vector4{}
	for _, line := range strings.Split(string(bytes), "\n") {
		cleanLine := strings.TrimSpace(line)
		if cleanLine == "" {
			continue
		}

		parts := strings.Fields(cleanLine)

		if parts[0] == "v" {
			vertice := Vector4{0, 0, 0, 1, -1, 0, 0, 0}
			for i := 1; i < 4; i++ {
				num, err := strconv.ParseFloat(parts[i], 32)
				if err != nil {
					log.Fatalf("Error parsing the vertice float '%s'", parts[1])
					return &mesh
				}

				num32 := float32(num)

				switch i {
				case 1:
					vertice.x = num32
				case 2:
					vertice.y = num32
				case 3:
					vertice.z = num32
				}
			}

			if len(vertices) == 0 {
				mesh.lowestX, mesh.highestX = vertice.x, vertice.x
				mesh.lowestY, mesh.highestY = vertice.y, vertice.y
				mesh.lowestZ, mesh.highestZ = vertice.z, vertice.z
			} else {
				if vertice.x < mesh.lowestX {
					mesh.lowestX = vertice.x
				}
				if vertice.y < mesh.lowestY {
					mesh.lowestY = vertice.y
				}
				if vertice.z < mesh.lowestZ {
					mesh.lowestZ = vertice.z
				}

				if vertice.x > mesh.highestX {
					mesh.highestX = vertice.x
				}
				if vertice.y > mesh.highestY {
					mesh.highestY = vertice.y
				}
				if vertice.z > mesh.highestZ {
					mesh.highestZ = vertice.z
				}
			}

			vertices = append(vertices, vertice)
			mesh.vertexAmount++
		} else if parts[0] == "f" {
			triangle := Triangle{}

			for i := 1; i < 4; i++ {
				vIndexString := strings.Split(parts[i], "/")[0]
				vIndex, err := strconv.Atoi(vIndexString)
				if err != nil {
					log.Fatalf("Error parsing the vertice index integer '%s'", parts[i])
					return &mesh
				}

				triangle.vecs[i-1] = vertices[vIndex-1]
			}

			mesh.tris = append(mesh.tris, triangle)
			mesh.triangleAmount++
		}
	}

	return &mesh
}
