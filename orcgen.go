// Package orcgen generates files from HTML -
// any static webpage can be informed, or even an HTML file.
// The file will be generated according the choosen extension.
package orcgen

import (
	"github.com/luabagg/orcgen/internal"
	"github.com/luabagg/orcgen/pkg/director"
)

// Valid extension types constants.
const (
	// PDF const.
	PDF = internal.PDF
	// PNG const.
	PNG = internal.PNG
	// JPEG const.
	JPEG = internal.JPEG
)

// New starts a new Director - the Director contains the available methods for file conversion.
//
// ext is the extension to be converted to (use the defined constants above).
//
// Connect and Close are used for the Browser connection control.
// ConvertWebpage and ConvertHTML are used for page conversion.
//
// There are a set of setters for specific config.
func New(ext internal.Ext) *director.Director {
	dir := director.NewDirector(ext)
	return dir.Connect()
}

// ConvertWebpage converts the url to the ext format, and saves the file.
//
// ext is the extension to be converted to (use the defined constants above).
// url is the url to convert.
// output is a filepath containing the extension.
//
// The connection is automatically closed.
func ConvertWebpage(ext internal.Ext, url string, output string) error {
	d := New(ext)
	defer d.Close()

	fi, err := d.ConvertWebpage(url)
	if err != nil {
		return err
	}

	return fi.Output(output)
}

// ConvertHTML converts the informed bytes to the ext format, and saves the file.
//
// ext is the extension to be converted to (use the defined constants above).
// html is the html byte array (if it's a filepath, use os.ReadFile(filepath)).
// output is a filepath containing the extension.
//
// The connection is automatically closed.
func ConvertHTML(ext internal.Ext, html []byte, output string) error {
	d := New(ext)
	defer d.Close()

	fi, err := d.ConvertHTML(html)
	if err != nil {
		return err
	}

	return fi.Output(output)
}
