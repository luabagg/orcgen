// orcgen generates files from HTML -
// any static webpage can be informed, or even the HTML itself.
// The file will be generated according the choosen extension.
package pkg

import (
	"github.com/luabagg/orcgen/internal"
)

const (
	// Valid extension types constants.

	// PDF const.
	PDF = internal.PDF
	// PNG const.
	PNG = internal.PNG
	// JPEG const.
	JPEG = internal.JPEG
)

// New starts a new orcgen - the Director contains the available methods.
//
// ext is the extension to be converted to (use the defined constants above).
//
// Connect and Close are used for the Browser connection controll.
// ConvertWebpage and ConvertHTML are used for page conversion.
//
// There are a set of setters for specific config.
func New(ext internal.Ext) *internal.Director {
	return internal.NewDirector(ext).Connect()
}
