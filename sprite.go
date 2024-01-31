package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"

	"math"

	"github.com/ncruces/zenity"
)

type Texture struct {
	w, h float64
	data [][]color.RGBA
}

func (t *Texture) GetColorAt(u, v float64) color.RGBA {
	finalU := u - math.Floor(u)
	finalV := v - math.Floor(v)
	x := int(math.Round((t.w - 1) * finalU))
	y := int((t.h - 1) - math.Round((t.h-1)*finalV))

	return t.data[y][x]
}

func LoadTexture(filename string) *Texture {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpg", jpeg.Decode, jpeg.DecodeConfig)

	file, err := os.Open(filename)

	if err != nil {
		zenity.Error(fmt.Sprintf("Error loading the texture '%s'.\n%s", filename, err), zenity.Title("Texture load error"), zenity.ErrorIcon)
		panic(err)
	}

	defer file.Close()

	texture := Texture{}

	img, _, err := image.Decode(file)

	if err != nil {
		zenity.Error(fmt.Sprintf("Error decoding the texture image.\n%s", err), zenity.Title("Texture load error"), zenity.ErrorIcon)
		panic(err)
	}

	bounds := img.Bounds()

	w, h := bounds.Max.X, bounds.Max.Y
	texture.w = float64(w)
	texture.h = float64(h)

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

	return &texture
}
