package ui

import (
	"log"

	"github.com/veandco/go-sdl2/ttf"
)

func LoadFont(fpath string, size int) *ttf.Font {
	if font, err := ttf.OpenFont(fpath, size); err != nil {
		log.Fatalf("Error loading the font '%s'", fpath)
	} else {
		return font
	}

	return nil
}
