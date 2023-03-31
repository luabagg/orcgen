// package png implements the generator for PNG files.
package png

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orc-generator/internal/utils"
)

// PNGGen struct.
type PNGGen struct {
	FullPage bool
	fileinfo utils.Fileinfo
}

// generatePNG converts a rod Page instance to a PNG file.
func (p *PNGGen) generatePNG(page *rod.Page) error {
	defer page.Close()

	req := &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
	}

	bin, err := page.Screenshot(p.FullPage, req)
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
func (p *PNGGen) ConvertWebpage(url string) error {
	return p.generatePNG(utils.UrlToPage(url))
}

// ConvertWebpage converts from a file.
func (p *PNGGen) ConvertHTML(html []byte) error {
	page, err := utils.ByteToPage(html)
	if err != nil {
		return err
	}

	return p.generatePNG(page)
}

// GetOutput saves the file to the informed filepath.
func (p *PNGGen) GetOutput(filepath string) error {
	return p.fileinfo.Output(filepath)
}

// GetFilesize gets the filesize.
func (p *PNGGen) GetFilesize() int {
	return p.fileinfo.Filesize
}

// Close resets struct and close Browser connection.
func (p *PNGGen) Close() {
	p = new(PNGGen)
	utils.CloseBrowser()
}
