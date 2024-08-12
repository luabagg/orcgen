package webdriver

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/assert"
)

func TestWebDriver_Connect(t *testing.T) {
	r := &WebDriver{}
	r.Connect()
	defer r.Close()
	page, _ := r.Browser.Page(proto.TargetCreateTarget{})

	assert.NotNil(t, page)
}

func TestWebDriver_Close(t *testing.T) {
	r := &WebDriver{}
	r.Connect()
	r.Close()

	page, _ := r.Browser.Page(proto.TargetCreateTarget{})
	assert.Nil(t, page)
}

func TestWebDriver_UrlToPage(t *testing.T) {
	webDriver := WebDriver{
		LoadTimeout:  10 * time.Second,
		PageIdleTime: 200 * time.Millisecond,
	}
	webDriver.Connect()
	defer webDriver.Browser.Close()

	// Call UrlToPage function to create a page instance
	page := webDriver.UrlToPage("https://www.example.com")
	assert.NotNil(t, page)

	// Check that the page has loaded successfully
	title, _ := page.MustElement("title").Text()
	assert.Equal(t, "Example Domain", title)
}

func GetDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}
	dir := filepath.Dir(filename)
	return dir, nil
}

func TestWebDriver_HTMLToPage(t *testing.T) {
	webDriver := WebDriver{
		LoadTimeout:  10 * time.Second,
		PageIdleTime: 200 * time.Millisecond,
	}
	webDriver.Connect()
	defer webDriver.Close()

	// convert the HTML byte slice to a Rod Page instance
	html := []byte("<html><head><title>ORC gen</title></head><body><h1>Hello, World!</h1></body></html>")
	page, err := webDriver.HTMLToPage(html)
	defer page.MustClose()

	assert.NoError(t, err)

	// get the page content and verifies if it's not empty
	content := page.MustHTML()
	assert.NotEmpty(t, content)

	// Check that the page has loaded successfully
	title, err := page.MustElement("title").Text()

	assert.NoError(t, err)

	assert.Equal(t, "ORC gen", title)
}

func TestWebDriver_WaitLoad(t *testing.T) {
	webDriver := &WebDriver{
		LoadTimeout:  10 * time.Second,
		PageIdleTime: 200 * time.Millisecond,
	}
	webDriver.Connect()
	defer webDriver.Close()

	// get a WebDriver Page instance
	page := webDriver.Browser.MustPage("https://www.example.com")
	defer page.MustClose()

	// call WaitLoad function
	webDriver.WaitLoad(page)

	// assert that the page is loaded and idle
	if !page.MustEval(`() => window.performance.timing.loadEventEnd > 0`).Bool() {
		t.Errorf("Page did not finish loading")
	}
	if !page.MustEval(`() => window.requestIdleCallback !== undefined`).Bool() {
		t.Errorf("Page is not idle")
	}
}
