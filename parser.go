package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ncruces/zenity"
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
		zenity.Error(fmt.Sprintf("Error loading the obj file.\n%s", err), zenity.Title("OBJ load error"), zenity.ErrorIcon)
		panic(err)
	}

	// Load .mtl and create a dictionary of: materialName - textureImage
	mtlTex := GetMtlTex(bytes, filename)

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
					zenity.Error(fmt.Sprintf("Error while parsing a texture vertice.\n%s", err), zenity.Title("Texture parsing error"), zenity.ErrorIcon)
					panic(err)
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

func GetMtlTex(bytes []byte, objFilename string) map[string]*Texture {
	basePath := filepath.Dir(objFilename)
	filename := ""
	for _, line := range strings.Split(string(bytes), "\n") {
		cleanLine := strings.TrimSpace(line)
		if cleanLine == "" {
			continue
		}

		parts := strings.Fields(cleanLine)

		if parts[0] == "mtllib" {
			filename = strings.Replace(cleanLine, "mtllib ", "", 1)
			break
		}
	}
	mtlTexDict := make(map[string]*Texture)
	if filename != "" {
		bytes, err := os.ReadFile(filepath.Join(basePath, filename))

		if err != nil {
			zenity.Error(fmt.Sprintf("Error loading the .mtl file '%s'.\n%s", filename, err), zenity.Title("Material load error"), zenity.ErrorIcon)
			panic(err)
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

			if strings.Contains(cleanLine, ".png") ||
				strings.Contains(cleanLine, ".jpg") ||
				strings.Contains(cleanLine, ".jpeg") ||
				strings.Contains(cleanLine, ".jif") {
				prefix := strings.Fields(cleanLine)[0] + " " // Trim the texture prefix (e.g 'map_Ka')
				texFileName := strings.Replace(cleanLine, prefix, "", 1)
				texFilePath := filepath.Join(basePath, texFileName)
				mtlTexDict[mtlKey] = LoadTexture(texFilePath)
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
				isTextured := true

				vParts := strings.Split(parts[i], "/")
				vIndexString := strings.Split(parts[i], "/")[0]

				var vTexIndexString string
				if len(vParts) == 2 {
					vTexIndexString = strings.Split(parts[i], "/")[1]
				} else if len(vParts) == 3 {
					vTexIndexString = strings.Split(parts[i], "/")[1]
				}

				if vTexIndexString == "" {
					isTextured = false
				}

				vIndex, err := strconv.Atoi(vIndexString)
				if err != nil {
					zenity.Error(fmt.Sprintf("Error parsing a vertice.\n%s", err), zenity.Title("Mesh parsing error"), zenity.ErrorIcon)
					panic(err)
				}

				vTexIndex := 0
				if isTextured {
					vTexIndex, err = strconv.Atoi(vTexIndexString)
					if err != nil {
						zenity.Error(fmt.Sprintf("Error parsing a vertice.\n%s", err), zenity.Title("Mesh parsing error"), zenity.ErrorIcon)
						panic(err)
					}
				}

				triangle.tex = lastTexture
				triangle.vecs[i-1] = vertices[vIndex-1]
				if isTextured {
					triangle.vecs[i-1].texVec = texVertices[vTexIndex-1]
				}
			}

			tris = append(tris, triangle)
		}
	}

	return tris
}
