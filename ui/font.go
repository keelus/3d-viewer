package ui

import (
	"embed"
	"fmt"

	"github.com/ncruces/zenity"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//go:embed font.ttf
var f embed.FS

func LoadFont(fpath string, size int) *ttf.Font {
	data, _ := f.ReadFile("font.ttf")
	loadedFont, _ := sdl.RWFromMem(data)
	if font, err := ttf.OpenFontRW(loadedFont, 1, size); err != nil {
		zenity.Error(fmt.Sprintf("Error loading the font '%s'.\n%s", fpath, err), zenity.Title("Font load error"), zenity.ErrorIcon)
		panic(err)
	} else {
		return font
	}
}
