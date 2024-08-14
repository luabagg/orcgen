package handlers

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/pkg/fileinfo"
)

type Config interface {
	proto.PageCaptureScreenshot | proto.PagePrintToPDF
}

// FileHandler interface contains the methods used for page conversion.
type FileHandler[T Config] interface {
	// SetConfig adds the config to the instance.
	SetConfig(config T) FileHandler[T]
	// SetFullPage sets the pages to be converted. If false, only the first page is selected.
	// Default is false.
	SetFullPage(fullPage bool) FileHandler[T]
	// GenerateFile converts a rod Page instance to a file.
	GenerateFile(page *rod.Page) (*fileinfo.Fileinfo, error)
}
