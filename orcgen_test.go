package orcgen

import (
	"os"
	"testing"

	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/v2/pkg/fileinfo"
	"github.com/luabagg/orcgen/v2/pkg/handlers"
	"github.com/luabagg/orcgen/v2/pkg/handlers/pdf"
	"github.com/luabagg/orcgen/v2/pkg/handlers/screenshot"
	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	html, _ := os.ReadFile("testdata/test.html")

	tests := []struct {
		name      string
		pdfConfig *PDFConfig
		html      []byte
		testErr   bool
	}{
		{
			"test HTML PDF generation",
			&PDFConfig{},
			html,
			false,
		},
		{
			"test HTML PDF generation",
			nil,
			html,
			true,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.testErr {
				defer func() {
					assert.NotNil(t, recover())
				}()
			}

			err := Generate(tc.html, *tc.pdfConfig, "test.pdf")
			assert.NoError(t, err)
		})
	}
}

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name             string
		pdfConfig        *PDFConfig
		screenshotConfig *ScreenshotConfig
	}{
		{
			"test ScreenshotHandler",
			nil,
			&ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatPng,
			},
		},
		{
			"test PDFHandler",
			&PDFConfig{
				PrintBackground: true,
			},
			nil,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.pdfConfig != nil {
				handler := NewHandler(*tc.pdfConfig)
				assert.IsType(t, &pdf.PDFHandler{}, handler, "expected to be a PDFHandler")
			} else if tc.screenshotConfig != nil {
				handler := NewHandler(*tc.screenshotConfig)
				assert.IsType(t, &screenshot.ScreenshotHandler{}, handler, "expected to be a ScreenshotHandler")
			} else {
				assert.Fail(t, "expected a valid config")
			}
		})
	}
}

func TestConvertHTML(t *testing.T) {
	html, _ := os.ReadFile("testdata/test.html")

	tests := []struct {
		name              string
		pdfHandler        handlers.FileHandler[PDFConfig]
		screenshotHandler handlers.FileHandler[ScreenshotConfig]
		html              []byte
		output            string
	}{
		{
			"test PDF",
			NewHandler(PDFConfig{}),
			nil,
			html,
			"test.pdf",
		},
		{
			"test PNG",
			nil,
			NewHandler(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatPng,
			}),
			html,
			"test.png",
		},
		{
			"test JPEG",
			nil,
			NewHandler(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatJpeg,
			}),
			html,
			"test.jpeg",
		},
		{
			"test WEBP",
			nil,
			NewHandler(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatWebp,
			}),
			html,
			"test.webp",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var err error
			var fi *fileinfo.Fileinfo

			if tc.pdfHandler != nil {
				fi, err = ConvertHTML(tc.pdfHandler, tc.html)
			} else if tc.screenshotHandler != nil {
				fi, err = ConvertHTML(tc.screenshotHandler, tc.html)
			} else {
				assert.Fail(t, "expected a valid handler")
			}

			assert.NoError(t, err)
			assert.IsType(t, &fileinfo.Fileinfo{}, fi, "expected to be a Fileinfo")

			os.Remove(tc.output)
		})
	}
}

func TestConvertWebpage(t *testing.T) {
	url := "https://www.example.com"

	tests := []struct {
		name              string
		pdfHandler        handlers.FileHandler[PDFConfig]
		screenshotHandler handlers.FileHandler[ScreenshotConfig]
		url               string
		output            string
	}{
		{
			"test PDF",
			NewHandler(PDFConfig{}),
			nil,
			url,
			"test.pdf",
		},
		{
			"test PNG",
			nil,
			NewHandler(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatPng,
			}),
			url,
			"test.png",
		},
		{
			"test JPEG",
			nil,
			NewHandler(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatJpeg,
			}),
			url,
			"test.jpeg",
		},
		{
			"test WEBP",
			nil,
			NewHandler(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatWebp,
			}),
			url,
			"test.webp",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var err error
			var fi *fileinfo.Fileinfo

			if tc.pdfHandler != nil {
				fi, err = ConvertWebpage(tc.pdfHandler, tc.url)
			} else if tc.screenshotHandler != nil {
				fi, err = ConvertWebpage(tc.screenshotHandler, tc.url)
			} else {
				assert.Fail(t, "expected a valid handler")
			}

			assert.NoError(t, err)
			assert.IsType(t, &fileinfo.Fileinfo{}, fi, "expected to be a Fileinfo")

			os.Remove(tc.output)
		})
	}
}
