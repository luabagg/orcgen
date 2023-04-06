package generator

import (
	"github.com/go-rod/rod"
)

// Generator interface contains the methods used for page conversion.
type Generator interface {
	// GenerateFile converts a rod Page instance to a file.
	GenerateFile(page *rod.Page) ([]byte, error)
	// SetFullPage sets the pages to be converted. If false, only the first page is selected.
	SetFullPage(fullPage bool) Generator
}
