// Package jpeg implements the builder for JPEG files.
package jpeg

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/internal/generator"
)

// JPEGBuilder struct.
type JPEGBuilder struct {
	fullPage bool
}

// GenerateFile converts a rod Page instance to a PNG file.
func (j *JPEGBuilder) GenerateFile(page *rod.Page) ([]byte, error) {
	var quality int = 100
	req := &proto.PageCaptureScreenshot{
		Format:      proto.PageCaptureScreenshotFormatJpeg,
		Quality:     &quality,
		FromSurface: false,
	}

	return page.Screenshot(j.fullPage, req)
}

// SetFullPage sets the pages to be converted. If false, only the first page is selected.
func (j *JPEGBuilder) SetFullPage(fullPage bool) generator.Generator {
	j.fullPage = fullPage

	return j
}
