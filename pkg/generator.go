// orc-generator generates files from HTML -
// any static webpage can be informed, or even the HTML itself.
// The file will be generated according the choosen extension.
package pkg

import (
	"time"

	"github.com/go-rod/rod"
	"github.com/luabagg/orc-generator/internal/generator"
	"github.com/luabagg/orc-generator/internal/generator/jpeg"
	"github.com/luabagg/orc-generator/internal/generator/pdf"
	"github.com/luabagg/orc-generator/internal/generator/png"
	"github.com/luabagg/orc-generator/internal/utils"
)

// Ext enum - valid extension types.
type Ext string

const (
	// PDF enum const.
	PDF Ext = "pdf"

	// PNG enum const.
	PNG Ext = "png"

	// JPEG enum const.
	JPEG Ext = "jpeg"
)

var (
	// LoadTimeout controlls max page load time before context is canceled.
	LoadTimeout time.Duration = 10 * time.Second
	// PageIdleTime sets the wait time after the page stops receiving requests.
	PageIdleTime time.Duration = 200 * time.Millisecond
	// FullPage sets the pages to be converted. If false, only the first page is selected.
	FullPage bool = true
)

// New starts a new generator from the given extension.
func New(ext Ext) generator.Generator {
	utils.Browser = rod.New().MustConnect()
	utils.Timeout = LoadTimeout
	utils.IdleTime = PageIdleTime

	switch ext {
	case PDF:
		return &pdf.PDFGen{
			FullPage: FullPage,
		}
	case PNG:
		return &png.PNGGen{
			FullPage: FullPage,
		}
	case JPEG:
		return &jpeg.JPEGGen{
			FullPage: FullPage,
		}
	}

	return nil
}
