// Package webdriver provides a wrapper for the rod library to perform browser operations.
package webdriver

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/luabagg/orcgen/v2/internal"
)

type WebDriverConfig struct {
	LoadTimeout  time.Duration
	PageIdleTime time.Duration
}

// WebDriver is a wrapper of the rod library.
type WebDriver struct {
	// Browser is a rod Browser instance.
	Browser *rod.Browser
	// LoadTimeout controlls max page load time before context is canceled.
	LoadTimeout time.Duration
	// PageIdleTime sets the wait time after the page stops receiving requests.
	PageIdleTime time.Duration
}

func FromDefault() *WebDriver {
	wd := &WebDriver{
		LoadTimeout:  30 * time.Second,
		PageIdleTime: 2 * time.Second,
	}
	wd.Connect()

	return wd
}

func FromConfig(config *WebDriverConfig) *WebDriver {
	wd := &WebDriver{
		LoadTimeout:  config.LoadTimeout,
		PageIdleTime: config.PageIdleTime,
	}
	wd.Connect()

	return wd
}

// Connect starts the Browser connection.
func (r *WebDriver) Connect() {
	r.Browser = rod.New().MustConnect()
}

// Close closes the Browser connection.
func (r *WebDriver) Close() {
	r.Browser.MustClose()
}

// UrlToPage converts the URL to a rod Page instance.
func (r *WebDriver) UrlToPage(url string) *rod.Page {
	return r.Browser.MustPage(url)
}

// HTMLToPage converts the binary html to a rod Page instance.
func (r *WebDriver) HTMLToPage(html []byte) (*rod.Page, error) {
	if !internal.IsValidHTML(html) {
		return &rod.Page{}, errors.New("the provided html must be valid")
	}

	file, err := os.CreateTemp("", "*.html")
	if err != nil {
		return &rod.Page{}, err
	}
	defer os.Remove(file.Name())

	if _, err = file.Write(html); err != nil {
		return &rod.Page{}, err
	}

	page := r.Browser.MustPage("file://" + file.Name())

	return page, nil
}

// WaitLoad waits for the page load event, then for the network to be idle.
//
// LoadTimeout limits the two waits together.
// It returns an error if the page does not load before the timeout.
// If the network is not idle before the timeout, it returns nil and the page is used as it is.
func (r *WebDriver) WaitLoad(page *rod.Page) error {
	page = page.Timeout(r.LoadTimeout)
	defer page.CancelTimeout()

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("wait for page load: %w", err)
	}

	wait := page.WaitRequestIdle(r.PageIdleTime, nil, nil, nil)
	wait()

	return nil
}
