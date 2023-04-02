// orcgen generates files from HTML -
// any static webpage can be informed, or even the HTML itself.
// The file will be generated according the choosen extension.
package orcgen

import (
	"log"

	"github.com/luabagg/orc-generator/internal"
	"github.com/luabagg/orc-generator/internal/generator"
)

const (
	// Valid extension types constants.

	// PDF const.
	PDF = generator.PDF
	// PNG const.
	PNG = generator.PNG
	// JPEG const.
	JPEG = generator.JPEG
)

var (
	// FullPage sets the pages to be converted. If false, only the first page is selected.
	FullPage bool = true
)

// New starts a .
func New(ext generator.Ext) *internal.Director {
	var gen generator.Generator
	if gen = generator.Build(ext, FullPage); gen == nil {
		log.Fatal("Generator not found.")
	}

	return internal.NewDirector(gen).Connect()
}
