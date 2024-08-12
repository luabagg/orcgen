// Package orcgen generates files from HTML -
// any static webpage can be informed, or even an HTML file.
// The file will be generated according the configured handler.
package orcgen

import (
	"github.com/luabagg/orcgen/pkg/fileinfo"
	"github.com/luabagg/orcgen/pkg/handlers"
	"github.com/luabagg/orcgen/pkg/webdriver"
)

// ConvertHTML converts the bytes using the given handler, and returns a Fileinfo object.
//
// handler is a Handler instance (see pkg/handlers).
// html is the html byte array (if it's a filepath, use os.ReadFile(filepath)).
//
// The connection with the Browser is automatically closed.
func ConvertHTML[T handlers.Config](handler handlers.FileHandler[T], html []byte) (*fileinfo.Fileinfo, error) {
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
func ConvertWebpage[T handlers.Config](handler handlers.FileHandler[T], url string) (*fileinfo.Fileinfo, error) {
	wd := webdriver.FromDefault()
	defer wd.Close()

	page := wd.UrlToPage(url)
	wd.WaitLoad(page)

	return handler.GenerateFile(page)
}
