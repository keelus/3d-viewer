package main

import (
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
)

func parseVector(parts []string) Vector4 {
	vertice := Vector4{0, 0, 0, 1, -1, NewTexVector(0, 0, 0)}

	for i := 1; i < 4; i++ {
		num, err := strconv.ParseFloat(parts[i], 32)
		if err != nil {
			log.Fatalf("Error parsing the vertice float '%s'", parts[1])
		}

		num32 := float64(num)

		switch i {
		case 1:
			vertice.x = num32
		case 2:
			vertice.y = num32
		case 3:
			vertice.z = num32
		}
	}

	return vertice
}

func ParseObj(filename string) *Mesh {
	mesh := Mesh{}

	bytes, err := os.ReadFile(filename)

	if err != nil {
		log.Fatalf("Error loading the obj file '%s'", filename)
		return &Mesh{}
	}

	// Load .mtl and create a dictionary of: materialName - textureImage
	mtlTex := GetMtlTex(filename)

	// Get the texture vertices
	texVertices := GetTexVerts(bytes)

	// Get the model vertices
	vertices, lowests, highests := GetVerts(bytes)

	mesh.lowestX = lowests[0]
	mesh.lowestY = lowests[1]
	mesh.lowestZ = lowests[2]
	mesh.highestX = highests[0]
	mesh.highestY = highests[1]
	mesh.highestZ = highests[2]

	// Get the triangles, using previous values
	triangles := GetTriangles(bytes, mtlTex, vertices, texVertices)

	mesh.tris = triangles

	mesh.vertexAmount = len(vertices)
	mesh.triangleAmount = len(triangles)

	return &mesh
}

func GetTexVerts(bytes []byte) []TexVector {
	texVerts := []TexVector{}

	for _, line := range strings.Split(string(bytes), "\n") {
		cleanLine := strings.TrimSpace(line)
		if cleanLine == "" {
			continue
		}

		parts := strings.Fields(cleanLine)

		if parts[0] == "vt" {
			texVertice := NewTexVector(0, 0, 0)
			for i := 1; i < 3; i++ {
				num, err := strconv.ParseFloat(parts[i], 32)
				if err != nil {
					log.Fatalf("Error parsing the texture vertice float '%s'", parts[1])
					return texVerts
				}

				num32 := float64(num)

				switch i {
				case 1:
					texVertice.u = num32
				case 2:
					texVertice.v = num32
				}
			}

			texVerts = append(texVerts, texVertice)
		}
	}

	return texVerts
}

func GetVerts(bytes []byte) ([]Vector4, [3]float64, [3]float64) {
	verts := []Vector4{}

	lowests := [3]float64{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	highests := [3]float64{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}

	for _, line := range strings.Split(string(bytes), "\n") {
		cleanLine := strings.TrimSpace(line)
		if cleanLine == "" {
			continue
		}

		parts := strings.Fields(cleanLine)

		if parts[0] == "v" {
			vertice := parseVector(parts)

			if vertice.x < lowests[0] {
				lowests[0] = vertice.x
			}
			if vertice.y < lowests[1] {
				lowests[1] = vertice.y
			}
			if vertice.z < lowests[2] {
				lowests[2] = vertice.z
			}

			if vertice.x > highests[0] {
				highests[0] = vertice.x
			}
			if vertice.y > highests[1] {
				highests[1] = vertice.y
			}
			if vertice.z > highests[2] {
				highests[2] = vertice.z
			}

			verts = append(verts, vertice)
		}
	}

	return verts, lowests, highests
}

var lastTexture *Texture

func GetMtlTex(filename string) map[string]*Texture {
	mtlTexDict := make(map[string]*Texture)
	if strings.HasSuffix(filename, ".obj") {
		basePath := path.Dir(filename)
		newFilename := strings.Replace(filename, ".obj", ".mtl", 1)
		bytes, err := os.ReadFile(newFilename)

		if err != nil {
			log.Fatalf("Error loading the obj file '%s'", filename)
			return mtlTexDict
		}

		var mtlKey string

		for _, line := range strings.Split(string(bytes), "\n") {
			cleanLine := strings.TrimSpace(line)
			if cleanLine == "" {
				continue
			}

			if strings.Contains(cleanLine, "newmtl") {
				mtlKey = strings.Fields(cleanLine)[1]
			}

			if strings.Contains(cleanLine, ".png") {
				texFilename := path.Join(basePath, strings.Fields(cleanLine)[1])
				mtlTexDict[mtlKey] = LoadTexture(texFilename)
			}

		}
	}

	return mtlTexDict
}

func GetTriangles(bytes []byte, mtlTex map[string]*Texture, vertices []Vector4, texVertices []TexVector) []Triangle {
	tris := []Triangle{}

	for _, line := range strings.Split(string(bytes), "\n") {
		cleanLine := strings.TrimSpace(line)
		if cleanLine == "" {
			continue
		}

		parts := strings.Fields(cleanLine)

		if parts[0] == "usemtl" {
			lastTexture = mtlTex[parts[1]]
		} else if parts[0] == "f" {
			triangle := Triangle{}

			for i := 1; i < 4; i++ {
				vParts := strings.Split(parts[i], "/")
				vIndexString := strings.Split(parts[i], "/")[0]
				var vTexIndexString string // TODO: Handle non textured
				if len(vParts) == 2 {
					vTexIndexString = strings.Split(parts[i], "/")[1]
				} else if len(vParts) == 3 {
					vTexIndexString = strings.Split(parts[i], "/")[1]
				}

				vIndex, err := strconv.Atoi(vIndexString)
				if err != nil {
					log.Fatalf("Error parsing the vertice index integer '%s'", parts[i])
				}
				vTexIndex, err := strconv.Atoi(vTexIndexString)
				if err != nil {
					log.Fatalf("Error parsing the vertice index integer '%s'", parts[i])
				}

				triangle.tex = lastTexture
				triangle.vecs[i-1] = vertices[vIndex-1]
				triangle.vecs[i-1].texVec = texVertices[vTexIndex-1]
			}

			tris = append(tris, triangle)
		}
	}

	return tris
}
