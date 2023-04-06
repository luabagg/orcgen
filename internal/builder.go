package internal

import (
	"github.com/luabagg/orcgen/internal/generator"
	"github.com/luabagg/orcgen/internal/generator/jpeg"
	"github.com/luabagg/orcgen/internal/generator/pdf"
	"github.com/luabagg/orcgen/internal/generator/png"
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

// builders controls the existing builders.
var builders = map[Ext]generator.Generator{
	PDF:  &pdf.PDFBuilder{},
	PNG:  &png.PNGBuilder{},
	JPEG: &jpeg.JPEGBuilder{},
}

// Build builds a new Generator.
func Build(ext Ext) generator.Generator {
	gen := builders[Ext(ext)]

	return gen
}
