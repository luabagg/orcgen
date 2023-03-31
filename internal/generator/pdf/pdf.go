// package pdf implements the generator for PDF files.
package pdf

import (
	"io"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orc-generator/internal/utils"
)

// PDFGen struct.
type PDFGen struct {
	FullPage bool
	fileinfo utils.Fileinfo
}

// generatePDF converts a rod Page instance to a PDF file.
func (p *PDFGen) generatePDF(page *rod.Page) error {
	defer page.Close()

	var pageRanges string
	if !p.FullPage {
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
		return err
	}

	bin, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	p.fileinfo = utils.Fileinfo{
		File:     bin,
		Filesize: len(bin),
	}

	return nil
}

// ConvertWebpage converts from a URL.
func (p *PDFGen) ConvertWebpage(url string) error {
	page := utils.UrlToPage(url)

	return p.generatePDF(page)
}

// ConvertWebpage converts from a file.
func (p *PDFGen) ConvertHTML(html []byte) error {
	page, err := utils.ByteToPage(html)
	if err != nil {
		return err
	}

	return p.generatePDF(page)
}

// GetOutput saves the file to the informed filepath.
func (p *PDFGen) GetOutput(filepath string) error {
	return p.fileinfo.Output(filepath)
}

// GetFilesize gets the filesize.
func (p *PDFGen) GetFilesize() int {
	return p.fileinfo.Filesize
}

// Close resets struct and close Browser connection.
func (p *PDFGen) Close() {
	p = new(PDFGen)
	utils.CloseBrowser()
}
