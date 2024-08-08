// Package png implements the builder for PNG files.
package png

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/internal/generator"
)

// PNGBuilder struct.
type PNGBuilder struct {
	fullPage bool
	config   generator.Config
}

// GenerateFile converts a rod Page instance to a PNG file.
func (p *PNGBuilder) GenerateFile(page *rod.Page) ([]byte, error) {
	req := &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
	}

	return page.Screenshot(p.fullPage, req)
}

// SetFullPage sets the pages to be converted. If false, only the first page is selected.
func (p *PNGBuilder) SetFullPage(fullPage bool) generator.Generator {
	p.fullPage = fullPage

	return p
}

func (p *PNGBuilder) Configure(c generator.Config) generator.Generator {
	p.config = c

	return p
}
