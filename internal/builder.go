package internal

import (
	"github.com/luabagg/orc-generator/internal/generator"
	"github.com/luabagg/orc-generator/internal/generator/jpeg"
	"github.com/luabagg/orc-generator/internal/generator/pdf"
	"github.com/luabagg/orc-generator/internal/generator/png"
)

// Ext enum - valid extension types.
type Ext int

const (
	// PDF enum const.
	PDF Ext = iota
	// PNG enum const.
	PNG
	// JPEG enum const.
	JPEG
)

// builders controlls the existing builders.
var builders = map[Ext]generator.Generator{
	PDF:  &pdf.PDFBuilder{},
	PNG:  &png.PNGBuilder{},
	JPEG: &jpeg.JPEGBuilder{},
}

// Build builds a new Generator.
func Build(ext Ext, fullPage bool) generator.Generator {
	gen := builders[Ext(ext)]
	gen.SetFullPage(fullPage)

	return gen
}
