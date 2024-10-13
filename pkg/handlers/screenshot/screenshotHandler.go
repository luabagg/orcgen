// Package screenshot is used to generate screenshots from the rod Page instance.
package screenshot

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/v2/pkg/fileinfo"
	"github.com/luabagg/orcgen/v2/pkg/handlers"
)

// ScreenshotHandler struct.
type ScreenshotHandler struct {
	config   *proto.PageCaptureScreenshot
	fullPage bool
}

// New creates a new ScreenshotHandler instance.
//
// png is the default extension.
func New() handlers.FileHandler[proto.PageCaptureScreenshot] {
	handler := &ScreenshotHandler{
		fullPage: false,
	}
	handler.SetConfig(proto.PageCaptureScreenshot{
		Format:                "png",
		Quality:               new(int),
		FromSurface:           false,
		CaptureBeyondViewport: false,
		OptimizeForSpeed:      false,
	})

	return handler
}

// SetConfig adds the config to the instance.
func (p *ScreenshotHandler) SetConfig(config proto.PageCaptureScreenshot) handlers.FileHandler[proto.PageCaptureScreenshot] {
	p.config = &config

	return p
}

// SetFullPage sets the pages to be converted. If false, only the first page is selected.
//
// Default is false.
func (s *ScreenshotHandler) SetFullPage(fullPage bool) handlers.FileHandler[proto.PageCaptureScreenshot] {
	s.fullPage = fullPage

	return s
}

// GenerateFile converts a rod Page instance to an image file.
func (s *ScreenshotHandler) GenerateFile(page *rod.Page) (*fileinfo.Fileinfo, error) {
	r, err := page.Screenshot(s.fullPage, s.config)
	if err != nil {
		return nil, err
	}

	return &fileinfo.Fileinfo{
		File:     r,
		Filesize: len(r),
	}, nil
}
