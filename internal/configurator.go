package internal

import (
	"github.com/luabagg/orcgen/internal/configurator"
	"github.com/luabagg/orcgen/internal/generator/jpeg"
	"github.com/luabagg/orcgen/internal/generator/pdf"
	"github.com/luabagg/orcgen/internal/generator/png"
)

// configurators configures different types of generators
var configurators = map[Ext]configurator.Configurator{
	PDF:  &pdf.PDFBuilder{},
	PNG:  &png.PNGBuilder{},
	JPEG: &jpeg.JPEGBuilder{},
}

// Set sets up a new Generator.
func Set(ext Ext) configurator.Configurator {
	gen := configurators[Ext(ext)]

	return gen
}
