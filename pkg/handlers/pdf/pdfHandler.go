// Package pdf is used to generate PDFs from the rod Page instance.
package pdf

import (
	"io"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/pkg/fileinfo"
	"github.com/luabagg/orcgen/pkg/handlers"
)

// PDFHandler struct.
type PDFHandler struct {
	config   *proto.PagePrintToPDF
	fullPage bool
}

// New creates a new PDFHandler instance.
func New() handlers.FileHandler[proto.PagePrintToPDF] {
	handler := &PDFHandler{
		fullPage: false,
	}
	handler.SetConfig(proto.PagePrintToPDF{
		Landscape:           true,
		DisplayHeaderFooter: true,
		PrintBackground:     true,
		MarginTop:           new(float64),
		MarginBottom:        new(float64),
		MarginLeft:          new(float64),
		MarginRight:         new(float64),
		PreferCSSPageSize:   true,
	})

	return handler
}

// SetConfig adds the config to the instance.
func (p *PDFHandler) SetConfig(config proto.PagePrintToPDF) handlers.FileHandler[proto.PagePrintToPDF] {
	p.config = &config

	return p
}

// SetFullPage sets the pages to be converted. If false, only the first page is selected.
// Default is false.
func (p *PDFHandler) SetFullPage(fullPage bool) handlers.FileHandler[proto.PagePrintToPDF] {
	p.fullPage = fullPage

	return p
}

// GenerateFile converts a rod Page instance to a PDF file.
func (p *PDFHandler) GenerateFile(page *rod.Page) (*fileinfo.Fileinfo, error) {
	if !p.fullPage && p.config.PageRanges == "" {
		p.config.PageRanges = "1"
	}

	r, err := page.PDF(
		p.config,
	)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &fileinfo.Fileinfo{
		File:     bytes,
		Filesize: len(bytes),
	}, nil
}
