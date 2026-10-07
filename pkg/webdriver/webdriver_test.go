package webdriver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	assert.NoError(t, webDriver.WaitLoad(page))

	// assert that the page is loaded and idle
	if !page.MustEval(`() => window.performance.timing.loadEventEnd > 0`).Bool() {
		t.Errorf("Page did not finish loading")
	}
	if !page.MustEval(`() => window.requestIdleCallback !== undefined`).Bool() {
		t.Errorf("Page is not idle")
	}
}

// newStallServer serves html at "/" and never answers "/stall" until the test ends.
func newStallServer(t *testing.T, html string) *httptest.Server {
	release := make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, html)
	})
	mux.HandleFunc("/stall", func(_ http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	})
	srv := httptest.NewServer(mux)
	// Cleanups run in reverse order: release the handlers, then close the server.
	t.Cleanup(srv.Close)
	t.Cleanup(func() { close(release) })

	return srv
}

func TestWebDriver_WaitLoad_LoadTimeout(t *testing.T) {
	srv := newStallServer(t, `<html><body><img src="/stall"></body></html>`)

	webDriver := &WebDriver{
		LoadTimeout:  time.Second,
		PageIdleTime: 200 * time.Millisecond,
	}
	webDriver.Connect()
	defer webDriver.Close()

	page := webDriver.UrlToPage(srv.URL)
	defer page.MustClose()

	start := time.Now()
	err := webDriver.WaitLoad(page)

	assert.True(t, errors.Is(err, context.DeadlineExceeded), "expected a deadline error, got %v", err)
	assert.Less(t, time.Since(start), 5*time.Second, "expected WaitLoad to stop at LoadTimeout")
}

func TestWebDriver_WaitLoad_NetworkNeverIdle(t *testing.T) {
	srv := newStallServer(t, `<html><body><script>
		setInterval(() => fetch("/stall").catch(() => {}), 50);
	</script></body></html>`)

	webDriver := &WebDriver{
		LoadTimeout:  time.Second,
		PageIdleTime: 200 * time.Millisecond,
	}
	webDriver.Connect()
	defer webDriver.Close()

	page := webDriver.UrlToPage(srv.URL)
	defer page.MustClose()

	start := time.Now()
	err := webDriver.WaitLoad(page)

	assert.NoError(t, err, "expected the page to be used as it is when the network is not idle")
	assert.Less(t, time.Since(start), 5*time.Second, "expected WaitLoad to stop at LoadTimeout")

	_, err = page.PDF(&proto.PagePrintToPDF{})
	assert.NoError(t, err, "expected the page to render after the timeout")
}
