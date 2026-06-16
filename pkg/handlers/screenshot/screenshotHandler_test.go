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
	defer wd.Close()

	// create a new ScreenshotHandler instance
	screenshotHandler := New()

	// Create test pages
	page1, err := wd.UrlToPage("https://www.example.com")
	if err != nil {
		t.Skip("Skipping test: browser not available")
	}

	page2, err := wd.UrlToPage("https://www.google.com")
	if err != nil {
		t.Skip("Skipping test: browser not available")
	}

	page3, err := wd.Browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		t.Skip("Skipping test: browser not available")
	}

	tests := []struct {
		name     string
		instance handlers.FileHandler[proto.PageCaptureScreenshot]
		input    *rod.Page
	}{
		{
			name:     "simple page",
			instance: screenshotHandler,
			input:    page1,
		},
		{
			name:     "fullpage",
			instance: screenshotHandler.SetFullPage(true),
			input:    page2,
		},
		{
			instance: screenshotHandler,
			input:    page3,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// generate the screenshot file
			if err := wd.WaitLoad(tc.input); err != nil {
				t.Skip("Skipping test: wait load failed")
			}

			jpegData, err := tc.instance.GenerateFile(tc.input)

			assert.NoError(t, err, "Expected no error")
			assert.NotEmpty(t, jpegData, "Expected bytes")
		})
	}
}
