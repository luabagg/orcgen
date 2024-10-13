package screenshot

import (
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/v2/pkg/handlers"
	"github.com/luabagg/orcgen/v2/pkg/webdriver"
	"github.com/stretchr/testify/assert"
)

func TestScreenshotHandler_SetConfig(t *testing.T) {
	tests := []struct {
		name  string
		input proto.PageCaptureScreenshot
	}{
		{
			name:  "valid config",
			input: proto.PageCaptureScreenshot{},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// create a new ScreenshotHandler instance
			instance := New().SetConfig(tc.input)

			assert.Implements(t, (*handlers.FileHandler[proto.PageCaptureScreenshot])(nil), instance, "expected to be a Generator instance")
		})
	}
}

func TestScreenshotHandler_SetFullPage(t *testing.T) {
	tests := []struct {
		name  string
		input bool
	}{
		{
			name:  "simple page",
			input: false,
		},
		{
			name:  "fullpage",
			input: true,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// create a new ScreenshotHandler instance
			instance := New()

			assert.Implements(t, (*handlers.FileHandler[proto.PageCaptureScreenshot])(nil), instance, "expected to be a Generator instance")
		})
	}
}

func TestScreenshotHandler_GenerateFile(t *testing.T) {
	// create a new browser instance
	wd := webdriver.FromDefault()

	// create a new ScreenshotHandler instance
	screenshotHandler := New()

	tests := []struct {
		name     string
		instance handlers.FileHandler[proto.PageCaptureScreenshot]
		input    *rod.Page
	}{
		{
			name:     "simple page",
			instance: screenshotHandler,
			input:    wd.UrlToPage("https://www.example.com"),
		},
		{
			name:     "fullpage",
			instance: screenshotHandler.SetFullPage(true),
			input:    wd.UrlToPage("https://www.google.com"),
		},
		{
			instance: screenshotHandler,
			input:    wd.Browser.MustPage(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// generate the JPEG file
			wd.WaitLoad(tc.input)
			jpegData, err := tc.instance.GenerateFile(tc.input)

			assert.NoError(t, err, "Expected no error")
			assert.NotEmpty(t, jpegData, "Expected bytes")
		})
	}
}
