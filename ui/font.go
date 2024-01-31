package ui

import (
	"embed"
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//go:embed font.ttf
var f embed.FS

func LoadFont(fpath string, size int) *ttf.Font {
	data, _ := f.ReadFile("font.ttf")
	loadedFont, _ := sdl.RWFromMem(data)
	if font, err := ttf.OpenFontRW(loadedFont, 1, size); err != nil {
		log.Fatalf("Error loading the font '%s'", fpath)
	} else {
		return font
	}

	return nil
}
