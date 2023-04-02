package internal

import (
	"os"
	"time"

	"github.com/go-rod/rod"
)

type Rod struct {
	// Browser is a rod Browser instance.
	browser *rod.Browser
	// LoadTimeout controlls max page load time before context is canceled.
	LoadTimeout time.Duration
	// PageIdleTime sets the wait time after the page stops receiving requests.
	PageIdleTime time.Duration
}

// Connect starts the Browser connection.
func (r *Rod) Connect() {
	r.browser = rod.New().MustConnect()
}

// Close closes the Browser connection.
func (r *Rod) Close() {
	r.browser.MustClose()
}

// waitLoad sets a wait time according to the page loading.
func (r *Rod) waitLoad(page *rod.Page) {
	page = page.Timeout(r.LoadTimeout).MustWaitLoad()

	wait := page.WaitRequestIdle(r.PageIdleTime, nil, nil)
	wait()

	page.CancelTimeout()
}

// UrlToPage converts the URL to a rod Page instance.
func (r *Rod) UrlToPage(url string) *rod.Page {
	page := r.browser.MustPage(url)

	r.waitLoad(page)

	return page
}

// UrlToPage converts the binary to a rod Page instance.
func (r *Rod) ByteToPage(bin []byte) (*rod.Page, error) {
	file, err := os.CreateTemp("", "*.html")
	if err != nil {
		return &rod.Page{}, err
	}

	defer os.Remove(file.Name())

	if _, err = file.Write(bin); err != nil {
		return &rod.Page{}, err
	}

	page := r.browser.MustPage("file://" + file.Name())

	r.waitLoad(page)

	return page, nil
}
