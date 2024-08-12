// Package screenshot is used to generate screenshots from the rod Page instance.
package screenshot

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/pkg/fileinfo"
	"github.com/luabagg/orcgen/pkg/handlers"
)

// ScreenshotHandler struct.
type ScreenshotHandler[T proto.PageCaptureScreenshot] struct {
	config   *proto.PageCaptureScreenshot
	fullPage bool
}

func New[T proto.PageCaptureScreenshot](config T) handlers.FileHandler[T] {
	builder := &ScreenshotHandler[T]{
		fullPage: false,
	}
	builder.SetConfig(config)

	return builder
}

func (s *ScreenshotHandler[T]) SetConfig(config T) handlers.FileHandler[T] {
	if cfg, ok := any(config).(proto.PageCaptureScreenshot); ok {
		s.config = &cfg
	} else {
		panic("invalid config type. Expected proto.PageCaptureScreenshot")
	}

	return s
}

// SetFullPage sets the pages to be converted. If false, only the first page is selected.
// Default is false.
func (s *ScreenshotHandler[T]) SetFullPage(fullPage bool) handlers.FileHandler[T] {
	s.fullPage = fullPage

	return s
}

// GenerateFile converts a rod Page instance to a JPEG file.
func (s *ScreenshotHandler[T]) GenerateFile(page *rod.Page) (*fileinfo.Fileinfo, error) {
	r, err := page.Screenshot(s.fullPage, s.config)
	if err != nil {
		return nil, err
	}

	return &fileinfo.Fileinfo{
		File:     r,
		Filesize: len(r),
	}, nil
}
