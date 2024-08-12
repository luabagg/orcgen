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
type PDFHandler[T proto.PagePrintToPDF] struct {
	config   *proto.PagePrintToPDF
	fullPage bool
}

func New[T proto.PagePrintToPDF](config T) handlers.FileHandler[T] {
	builder := &PDFHandler[T]{
		fullPage: false,
	}
	builder.SetConfig(config)

	return builder
}

func (p *PDFHandler[T]) SetConfig(config T) handlers.FileHandler[T] {
	cfg := proto.PagePrintToPDF(config)
	p.config = &cfg

	return p
}

// SetFullPage sets the pages to be converted. If false, only the first page is selected.
// Default is false.
func (p *PDFHandler[T]) SetFullPage(fullPage bool) handlers.FileHandler[T] {
	p.fullPage = fullPage

	return p
}

// GenerateFile converts a rod Page instance to a PDF file.
func (p *PDFHandler[T]) GenerateFile(page *rod.Page) (*fileinfo.Fileinfo, error) {
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
