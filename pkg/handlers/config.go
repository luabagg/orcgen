package handlers

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/v2/pkg/fileinfo"
)

type Config interface {
	proto.PageCaptureScreenshot | proto.PagePrintToPDF
}

// BaseHandler is a non-generic interface for file handlers.
// This allows for type-safe handler construction without runtime type assertions.
type BaseHandler interface {
	// GenerateFile converts a rod Page instance to a file.
	GenerateFile(page *rod.Page) (*fileinfo.Fileinfo, error)
}

// FileHandler interface contains the methods used for page conversion.
type FileHandler[T Config] interface {
	BaseHandler
	// SetConfig adds the config to the instance.
	SetConfig(config T) FileHandler[T]
	// SetFullPage sets the pages to be converted. If false, only the first page is selected.
	// Default is false.
	SetFullPage(fullPage bool) FileHandler[T]
}
