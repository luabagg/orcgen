// package pdf implements the builder for PDF files.
package pdf

import (
	"io"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// PDFBuilder struct.
type PDFBuilder struct {
	fullPage bool
}

// GenerateFile converts a rod Page instance to a PDF file.
func (p *PDFBuilder) GenerateFile(page *rod.Page) ([]byte, error) {
	var pageRanges string
	if !p.fullPage {
		pageRanges = "1"
	}

	r, err := page.PDF(&proto.PagePrintToPDF{
		Landscape:           true,
		DisplayHeaderFooter: true,
		PrintBackground:     true,
		MarginTop:           new(float64),
		MarginBottom:        new(float64),
		MarginLeft:          new(float64),
		MarginRight:         new(float64),
		PageRanges:          pageRanges,
		PreferCSSPageSize:   true,
	})
	if err != nil {
		return nil, err
	}

	return io.ReadAll(r)
}

// SetFullPage sets the pages to be converted. If false, only the first page is selected.
func (p *PDFBuilder) SetFullPage(fullPage bool) {
	p.fullPage = fullPage
}
