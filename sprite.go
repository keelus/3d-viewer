package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	math "github.com/chewxy/math32"
)

type Texture struct {
	w, h float32
	data [][]color.RGBA
}

func (t Texture) GetColorAt(u, v float32) color.RGBA {
	x := int(math.Round((t.w - 1) * u))
	y := int((t.h - 1) - math.Round((t.h-1)*v))

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x > int(t.w-1) {
		x = int(t.w - 1)
	}
	if y > int(t.h-1) {
		y = int(t.h - 1)
	}
	return t.data[y][x]
}

func LoadTexture(filename string) Texture {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open(filename)

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	defer file.Close()

	texture := Texture{}

	img, _, err := image.Decode(file)

	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()

	w, h := bounds.Max.X, bounds.Max.Y
	texture.w = float32(w)
	texture.h = float32(h)

	for y := 0; y < h; y++ {
		row := []color.RGBA{}
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			r /= 257
			g /= 257
			b /= 257
			a /= 257
			row = append(row, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
		texture.data = append(texture.data, row)
	}

	return texture
}
