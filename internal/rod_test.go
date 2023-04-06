package internal

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/assert"
)

func TestRod_Connect(t *testing.T) {
	r := &Rod{}
	r.Connect()
	defer r.Close()
	page, _ := r.Browser.Page(proto.TargetCreateTarget{})

	assert.NotNil(t, page)
}

func TestRod_Close(t *testing.T) {
	r := &Rod{}
	r.Connect()
	r.Close()

	page, _ := r.Browser.Page(proto.TargetCreateTarget{})
	assert.Nil(t, page)
}

func TestRod_UrlToPage(t *testing.T) {
	rod := Rod{
		LoadTimeout:  10 * time.Second,
		PageIdleTime: 200 * time.Millisecond,
	}
	rod.Connect()
	defer rod.Browser.Close()

	// Call UrlToPage function to create a page instance
	page := rod.UrlToPage("https://www.example.com")
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

func TestRod_ByteToPage(t *testing.T) {
	rod := Rod{
		LoadTimeout:  10 * time.Second,
		PageIdleTime: 200 * time.Millisecond,
	}
	rod.Connect()
	defer rod.Close()

	// convert the HTML byte slice to a rod Page instance
	html := []byte("<html><head><title>ORC gen</title></head><body><h1>Hello, World!</h1></body></html>")
	page, err := rod.ByteToPage(html)
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

func TestRod_WaitLoad(t *testing.T) {
	rod := &Rod{
		LoadTimeout:  10 * time.Second,
		PageIdleTime: 200 * time.Millisecond,
	}
	rod.Connect()
	defer rod.Close()

	// get a rod Page instance
	page := rod.Browser.MustPage("https://www.example.com")
	defer page.MustClose()

	// call WaitLoad function
	rod.WaitLoad(page)

	// assert that the page is loaded and idle
	if !page.MustEval(`() => window.performance.timing.loadEventEnd > 0`).Bool() {
		t.Errorf("Page did not finish loading")
	}
	if !page.MustEval(`() => window.requestIdleCallback !== undefined`).Bool() {
		t.Errorf("Page is not idle")
	}
}
