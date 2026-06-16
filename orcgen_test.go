package orcgen

import (
	"os"
	"testing"

	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/v2/pkg/fileinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPDFBuilder(t *testing.T) {
	config := PDFConfig{
		Landscape:       true,
		PrintBackground: true,
	}

	builder := PDF(config)
	require.NotNil(t, builder)
	require.NotNil(t, builder.handler)
}

func TestScreenshotBuilder(t *testing.T) {
	config := ScreenshotConfig{
		Format: proto.PageCaptureScreenshotFormatPng,
	}

	builder := Screenshot(config)
	require.NotNil(t, builder)
	require.NotNil(t, builder.handler)
}

func TestHandlerBuilderFullPage(t *testing.T) {
	t.Run("PDF builder with FullPage", func(t *testing.T) {
		builder := PDF(PDFConfig{}).FullPage(true)
		require.NotNil(t, builder)
	})

	t.Run("Screenshot builder with FullPage", func(t *testing.T) {
		builder := Screenshot(ScreenshotConfig{}).FullPage(true)
		require.NotNil(t, builder)
	})

	t.Run("Method chaining", func(t *testing.T) {
		builder := PDF(PDFConfig{
			Landscape: true,
		}).FullPage(true).FullPage(false)
		require.NotNil(t, builder)
	})
}

func TestGenerateHTML(t *testing.T) {
	html, err := os.ReadFile("testdata/test.html")
	require.NoError(t, err)

	tests := []struct {
		name    string
		builder *HandlerBuilder
		html    []byte
		output  string
	}{
		{
			"PDF generation from HTML",
			PDF(PDFConfig{
				PrintBackground: true,
			}),
			html,
			"testdata/test_pdf.pdf",
		},
		{
			"PNG screenshot from HTML",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatPng,
			}),
			html,
			"testdata/test_png.png",
		},
		{
			"JPEG screenshot from HTML",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatJpeg,
			}),
			html,
			"testdata/test_jpeg.jpeg",
		},
		{
			"WebP screenshot from HTML",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatWebp,
			}),
			html,
			"testdata/test_webp.webp",
		},
		{
			"PDF with FullPage",
			PDF(PDFConfig{}).FullPage(true),
			html,
			"testdata/test_pdf_fullpage.pdf",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer os.Remove(tc.output)

			err := GenerateHTML(tc.html, tc.builder, tc.output)
			assert.NoError(t, err)

			// Verify file was created
			_, err = os.Stat(tc.output)
			assert.NoError(t, err)
		})
	}
}

func TestGenerateURL(t *testing.T) {
	url := "https://www.example.com"

	tests := []struct {
		name    string
		builder *HandlerBuilder
		url     string
		output  string
	}{
		{
			"PDF from URL",
			PDF(PDFConfig{
				PrintBackground: true,
			}),
			url,
			"testdata/url_test_pdf.pdf",
		},
		{
			"PNG screenshot from URL",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatPng,
			}),
			url,
			"testdata/url_test_png.png",
		},
		{
			"JPEG screenshot from URL",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatJpeg,
			}),
			url,
			"testdata/url_test_jpeg.jpeg",
		},
		{
			"WebP screenshot from URL",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatWebp,
			}),
			url,
			"testdata/url_test_webp.webp",
		},
		{
			"Screenshot with FullPage",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatPng,
			}).FullPage(true),
			url,
			"testdata/url_test_fullpage.png",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer os.Remove(tc.output)

			err := GenerateURL(tc.url, tc.builder, tc.output)
			assert.NoError(t, err)

			// Verify file was created
			_, err = os.Stat(tc.output)
			assert.NoError(t, err)
		})
	}
}

func TestConvertHTML(t *testing.T) {
	html, err := os.ReadFile("testdata/test.html")
	require.NoError(t, err)

	tests := []struct {
		name    string
		builder *HandlerBuilder
		html    []byte
	}{
		{
			"PDF conversion",
			PDF(PDFConfig{}),
			html,
		},
		{
			"PNG screenshot",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatPng,
			}),
			html,
		},
		{
			"JPEG screenshot",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatJpeg,
			}),
			html,
		},
		{
			"WebP screenshot",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatWebp,
			}),
			html,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fi, err := ConvertHTML(tc.builder, tc.html)
			assert.NoError(t, err)
			assert.IsType(t, &fileinfo.Fileinfo{}, fi, "expected to be a Fileinfo")
			assert.NotEmpty(t, fi.File, "expected file content")
			assert.Greater(t, fi.Filesize, 0, "expected positive filesize")
		})
	}
}

func TestConvertURL(t *testing.T) {
	url := "https://www.example.com"

	tests := []struct {
		name    string
		builder *HandlerBuilder
		url     string
	}{
		{
			"PDF from URL",
			PDF(PDFConfig{}),
			url,
		},
		{
			"PNG screenshot from URL",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatPng,
			}),
			url,
		},
		{
			"JPEG screenshot from URL",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatJpeg,
			}),
			url,
		},
		{
			"WebP screenshot from URL",
			Screenshot(ScreenshotConfig{
				Format: proto.PageCaptureScreenshotFormatWebp,
			}),
			url,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fi, err := ConvertURL(tc.builder, tc.url)
			assert.NoError(t, err)
			assert.IsType(t, &fileinfo.Fileinfo{}, fi, "expected to be a Fileinfo")
			assert.NotEmpty(t, fi.File, "expected file content")
			assert.Greater(t, fi.Filesize, 0, "expected positive filesize")
		})
	}
}

func TestBuilderMethodChaining(t *testing.T) {
	t.Run("PDF builder chaining", func(t *testing.T) {
		builder := PDF(PDFConfig{
			Landscape:       true,
			PrintBackground: true,
		}).FullPage(true)

		require.NotNil(t, builder)
		require.NotNil(t, builder.handler)
	})

	t.Run("Screenshot builder chaining", func(t *testing.T) {
		builder := Screenshot(ScreenshotConfig{
			Format: proto.PageCaptureScreenshotFormatPng,
		}).FullPage(false)

		require.NotNil(t, builder)
		require.NotNil(t, builder.handler)
	})
}

func TestBuilderWithEmptyConfig(t *testing.T) {
	t.Run("PDF with empty config", func(t *testing.T) {
		builder := PDF(PDFConfig{})
		require.NotNil(t, builder)
		require.NotNil(t, builder.handler)
	})

	t.Run("Screenshot with empty config", func(t *testing.T) {
		builder := Screenshot(ScreenshotConfig{})
		require.NotNil(t, builder)
		require.NotNil(t, builder.handler)
	})
}

func TestBuilderReuseability(t *testing.T) {
	html, err := os.ReadFile("testdata/test.html")
	require.NoError(t, err)

	// Create a builder and use it multiple times
	builder := PDF(PDFConfig{
		PrintBackground: true,
	})

	// First use
	fi1, err := ConvertHTML(builder, html)
	assert.NoError(t, err)
	assert.NotNil(t, fi1)

	// Second use - builder should still work
	fi2, err := ConvertHTML(builder, html)
	assert.NoError(t, err)
	assert.NotNil(t, fi2)
}

func TestDifferentConfigurations(t *testing.T) {
	html, err := os.ReadFile("testdata/test.html")
	require.NoError(t, err)

	tests := []struct {
		name   string
		config PDFConfig
	}{
		{
			"Landscape mode",
			PDFConfig{Landscape: true},
		},
		{
			"Portrait mode",
			PDFConfig{Landscape: false},
		},
		{
			"With background",
			PDFConfig{PrintBackground: true},
		},
		{
			"Without background",
			PDFConfig{PrintBackground: false},
		},
		{
			"Custom page ranges",
			PDFConfig{PageRanges: "1"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			builder := PDF(tc.config)
			fi, err := ConvertHTML(builder, html)
			assert.NoError(t, err)
			assert.NotNil(t, fi)
			assert.Greater(t, fi.Filesize, 0)
		})
	}
}
