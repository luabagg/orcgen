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

// Generate generates a file from the given HTML / URL and outputs it to the given path.
//
// There's no checking in the extension type, so make sure to use the correct one.
func Generate[T string | []byte, Config handlers.Config](html T, config Config, output string) error {
	handler := NewHandler(config)

	var fileinfo *fileinfo.Fileinfo
	var err error

	if _, ok := any(html).([]byte); ok {
		fileinfo, err = ConvertHTML(handler, any(html).([]byte))
	} else {
		fileinfo, err = ConvertWebpage(handler, any(html).(string))
	}

	if err != nil {
		return err
	}
	return fileinfo.Output(output)
}

// NewHandler creates a handler from the config.
//
// It checks the config type and instanciates the handler accordingly.
func NewHandler[Config handlers.Config](config Config) handlers.FileHandler[Config] {
	var handler any

	if _, ok := any(config).(proto.PagePrintToPDF); ok {
		handler = pdf.New()
	} else if _, ok := any(config).(proto.PageCaptureScreenshot); ok {
		handler = screenshot.New()
	} else {
		panic("invalid config type provided")
	}

	return any(handler).(handlers.FileHandler[Config]).SetConfig(config)
}

// ConvertHTML converts the bytes using the given handler, and returns a Fileinfo object.
//
// handler is a Handler instance (see pkg/handlers).
// html is the html byte array (if it's a filepath, use os.ReadFile(filepath)).
//
// The connection with the Browser is automatically closed.
func ConvertHTML[Config handlers.Config](handler handlers.FileHandler[Config], html []byte) (*fileinfo.Fileinfo, error) {
	wd := webdriver.FromDefault()
	defer wd.Close()

	page, err := wd.HTMLToPage(html)
	if err != nil {
		return nil, err
	}
	wd.WaitLoad(page)

	return handler.GenerateFile(page)
}

// ConvertWebpage converts the url using the given handler, and returns a Fileinfo object
//
// handler is a Handler instance (see pkg/handlers).
// url will be converted as configured, if you need special treats, check the Webdriver docs.
//
// The connection with the Browser is automatically closed.
func ConvertWebpage[Config handlers.Config](handler handlers.FileHandler[Config], url string) (*fileinfo.Fileinfo, error) {
	wd := webdriver.FromDefault()
	defer wd.Close()

	page := wd.UrlToPage(url)
	wd.WaitLoad(page)

	return handler.GenerateFile(page)
}
