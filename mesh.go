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
}

func LoadMesh(filename string) Mesh {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		log.Fatalf("Error loading the obj file '%s'", filename)
		return Mesh{}
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
			vertice := Vector4{0, 0, 0, 1}
			for i := 1; i < 4; i++ {
				num, err := strconv.ParseFloat(parts[i], 32)
				if err != nil {
					log.Fatalf("Error parsing the vertice float '%s'", parts[1])
					return mesh
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

			vertices = append(vertices, vertice)
			mesh.vertexAmount++
		} else if parts[0] == "f" {
			triangle := Triangle{}

			for i := 1; i < 4; i++ {
				vIndexString := strings.Split(parts[i], "/")[0]
				vIndex, err := strconv.Atoi(vIndexString)
				if err != nil {
					log.Fatalf("Error parsing the vertice index integer '%s'", parts[i])
					return mesh
				}

				triangle.vecs[i-1] = vertices[vIndex-1]
			}

			mesh.tris = append(mesh.tris, triangle)
			mesh.triangleAmount++
		}
	}

	return mesh
}
