// Package orcgen generates files from HTML / URLs -
// any webpage can be informed, or even an HTML file.
//
// The file will be generated accordingly the configured handler.
// You can also configure the webdriver to control the page before saving the file.
package orcgen

import (
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/v2/pkg/fileinfo"
	"github.com/luabagg/orcgen/v2/pkg/handlers"
	"github.com/luabagg/orcgen/v2/pkg/handlers/pdf"
	"github.com/luabagg/orcgen/v2/pkg/handlers/screenshot"
	"github.com/luabagg/orcgen/v2/pkg/webdriver"
)

// Aliases:
type ScreenshotConfig = proto.PageCaptureScreenshot
type PDFConfig = proto.PagePrintToPDF

// HandlerBuilder provides a fluent API for building file handlers.
// This eliminates the need for runtime type assertions.
type HandlerBuilder struct {
	handler handlers.BaseHandler
}

// PDF creates a handler builder for PDF generation.
// Example:
//
//	err := orcgen.GenerateURL("https://example.com",
//	    orcgen.PDF(orcgen.PDFConfig{Landscape: true}),
//	    "output.pdf")
func PDF(config PDFConfig) *HandlerBuilder {
	handler := pdf.New().SetConfig(config)
	return &HandlerBuilder{handler: handler}
}

// Screenshot creates a handler builder for screenshot generation.
// Example:
//
//	err := orcgen.GenerateHTML(htmlBytes,
//	    orcgen.Screenshot(orcgen.ScreenshotConfig{Format: "png"}),
//	    "output.png")
func Screenshot(config ScreenshotConfig) *HandlerBuilder {
	handler := screenshot.New().SetConfig(config)
	return &HandlerBuilder{handler: handler}
}

// FullPage configures whether to capture the full page or just the first page.
// Returns the builder for method chaining.
func (hb *HandlerBuilder) FullPage(enabled bool) *HandlerBuilder {
	// Type assert to set full page on the underlying handler
	switch h := hb.handler.(type) {
	case handlers.FileHandler[proto.PagePrintToPDF]:
		hb.handler = h.SetFullPage(enabled)
	case handlers.FileHandler[proto.PageCaptureScreenshot]:
		hb.handler = h.SetFullPage(enabled)
	}
	return hb
}

// GenerateHTML generates a file from HTML bytes and outputs it to the given path.
//
// This is the recommended way to generate files from HTML.
// Example:
//
//	err := orcgen.GenerateHTML(htmlBytes,
//	    orcgen.PDF(orcgen.PDFConfig{Landscape: true}).FullPage(true),
//	    "output.pdf")
func GenerateHTML(html []byte, builder *HandlerBuilder, output string) error {
	wd := webdriver.FromDefault()
	defer wd.Close()

	page, err := wd.HTMLToPage(html)
	if err != nil {
		return err
	}
	wd.WaitLoad(page)

	fileinfo, err := builder.handler.GenerateFile(page)
	if err != nil {
		return err
	}

	return fileinfo.Output(output)
}

// GenerateURL generates a file from a URL and outputs it to the given path.
//
// Example:
//
//	err := orcgen.GenerateURL("https://example.com",
//	    orcgen.Screenshot(orcgen.ScreenshotConfig{Format: "png"}),
//	    "output.png")
func GenerateURL(url string, builder *HandlerBuilder, output string) error {
	wd := webdriver.FromDefault()
	defer wd.Close()

	page := wd.UrlToPage(url)
	wd.WaitLoad(page)

	fileinfo, err := builder.handler.GenerateFile(page)
	if err != nil {
		return err
	}

	return fileinfo.Output(output)
}

// ConvertHTML converts HTML bytes using the given handler builder, and returns a Fileinfo object.
//
// The connection with the Browser is automatically closed.
//
// Example:
//
//	fi, err := orcgen.ConvertHTML(orcgen.PDF(config), htmlBytes)
func ConvertHTML(builder *HandlerBuilder, html []byte) (*fileinfo.Fileinfo, error) {
	wd := webdriver.FromDefault()
	defer wd.Close()

	page, err := wd.HTMLToPage(html)
	if err != nil {
		return nil, err
	}
	wd.WaitLoad(page)

	return builder.handler.GenerateFile(page)
}

// ConvertURL converts a URL using the given handler builder, and returns a Fileinfo object.
//
// The connection with the Browser is automatically closed.
//
// Example:
//
//	fi, err := orcgen.ConvertURL(orcgen.Screenshot(config), "https://example.com")
func ConvertURL(builder *HandlerBuilder, url string) (*fileinfo.Fileinfo, error) {
	wd := webdriver.FromDefault()
	defer wd.Close()

	page := wd.UrlToPage(url)
	wd.WaitLoad(page)

	return builder.handler.GenerateFile(page)
}
