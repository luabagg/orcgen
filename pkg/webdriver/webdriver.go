// Package webdriver provides a wrapper for the rod library to perform browser operations.
package webdriver

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/v2/internal"
)

type WebDriverConfig struct {
	LoadTimeout  time.Duration
	PageIdleTime time.Duration
	Headless     bool
	NoSandbox    bool // Useful for Docker/CI environments
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

	if err := wd.Connect(); err != nil {
		// For backward compatibility, panic if connection fails
		// Users can use FromConfig with error handling if needed
		panic(err)
	}

	return wd
}

func FromConfig(config *WebDriverConfig) (*WebDriver, error) {
	wd := &WebDriver{
		LoadTimeout:  config.LoadTimeout,
		PageIdleTime: config.PageIdleTime,
	}

	err := wd.ConnectWithConfig(config)
	if err != nil {
		return nil, err
	}

	return wd, nil
}

// Connect starts the Browser connection with default settings.
func (r *WebDriver) Connect() error {
	browser := rod.New()
	if err := browser.Connect(); err != nil {
		return err
	}
	r.Browser = browser
	return nil
}

// ConnectWithConfig starts the Browser connection with custom configuration.
func (r *WebDriver) ConnectWithConfig(config *WebDriverConfig) error {
	l := launcher.New()

	if config.Headless {
		l = l.Headless(true)
	}

	if config.NoSandbox {
		l = l.NoSandbox(true)
	}

	url, err := l.Launch()
	if err != nil {
		return err
	}

	browser := rod.New().ControlURL(url)
	if err := browser.Connect(); err != nil {
		return err
	}

	r.Browser = browser
	return nil
}

// Close closes the Browser connection safely.
func (r *WebDriver) Close() error {
	if r.Browser == nil {
		return nil
	}
	return r.Browser.Close()
}

// UrlToPage converts the URL to a rod Page instance.
func (r *WebDriver) UrlToPage(url string) (*rod.Page, error) {
	return r.Browser.Page(proto.TargetCreateTarget{URL: url})
}

// HTMLToPage converts the binary html to a rod Page instance.
func (r *WebDriver) HTMLToPage(html []byte) (*rod.Page, error) {
	if !internal.IsValidHTML(html) {
		return nil, errors.New("the provided html must be valid")
	}

	file, err := os.CreateTemp("", "*.html")
	if err != nil {
		return nil, err
	}

	// Ensure file is closed properly
	fileName := file.Name()
	defer os.Remove(fileName)

	if _, err = file.Write(html); err != nil {
		file.Close()
		return nil, err
	}

	// Close file before browser tries to read it
	if err = file.Close(); err != nil {
		return nil, err
	}

	page, err := r.Browser.Page(proto.TargetCreateTarget{URL: "file://" + fileName})
	if err != nil {
		return nil, err
	}

	return page, nil
}

// WaitLoad sets a wait time according to the page loading.
func (r *WebDriver) WaitLoad(page *rod.Page) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.LoadTimeout)
	defer cancel()

	if err := page.Context(ctx).WaitLoad(); err != nil {
		return err
	}

	wait := page.WaitRequestIdle(r.PageIdleTime, nil, nil, nil)
	wait()

	return nil
}
