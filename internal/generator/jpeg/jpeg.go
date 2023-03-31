// package jpeg implements the generator for JPEG files.
package jpeg

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orc-generator/internal/utils"
)

// JPEGGen struct.
type JPEGGen struct {
	FullPage bool
	fileinfo utils.Fileinfo
}

// generateJPEG converts a rod Page instance to a PNG file.
func (j *JPEGGen) generateJPEG(page *rod.Page) error {
	defer page.Close()

	var quality int = 100
	req := &proto.PageCaptureScreenshot{
		Format:      proto.PageCaptureScreenshotFormatJpeg,
		Quality:     &quality,
		FromSurface: false,
	}

	bin, err := page.Screenshot(j.FullPage, req)
	if err != nil {
		return err
	}

	j.fileinfo = utils.Fileinfo{
		File:     bin,
		Filesize: len(bin),
	}

	return nil
}

// ConvertWebpage converts from a URL.
func (j *JPEGGen) ConvertWebpage(url string) error {
	return j.generateJPEG(utils.UrlToPage(url))
}

// ConvertWebpage converts from a file.
func (j *JPEGGen) ConvertHTML(html []byte) error {
	page, err := utils.ByteToPage(html)
	if err != nil {
		return err
	}

	return j.generateJPEG(page)
}

// GetOutput saves the file to the informed filepath.
func (j *JPEGGen) GetOutput(filepath string) error {
	return j.fileinfo.Output(filepath)
}

// GetFilesize gets the filesize.
func (j *JPEGGen) GetFilesize() int {
	return j.fileinfo.Filesize
}

// Close resets struct and close Browser connection.
func (j *JPEGGen) Close() {
	j = new(JPEGGen)
	utils.CloseBrowser()
}
