// package generator contains builders used for page conversion (example HTML to PDF).
package generator

import (
	"github.com/go-rod/rod"
	"github.com/luabagg/orc-generator/internal/generator/jpeg"
	"github.com/luabagg/orc-generator/internal/generator/pdf"
	"github.com/luabagg/orc-generator/internal/generator/png"
)

// Generator interface contains the methods used for page conversion.
type Generator interface {
	// GenerateFile converts a rod Page instance to a file.
	GenerateFile(page *rod.Page) ([]byte, error)
	// SetFullPage sets the pages to be converted. If false, only the first page is selected.
	SetFullPage(fullPage bool)
}

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
var builders = map[Ext]Generator{
	PDF:  &pdf.PDFBuilder{},
	PNG:  &png.PNGBuilder{},
	JPEG: &jpeg.JPEGBuilder{},
}

// Build builds a new Generator.
func Build(ext Ext, fullPage bool) Generator {
	gen := builders[Ext(ext)]
	gen.SetFullPage(fullPage)

	return gen
}
