package internal

import (
	"os"
	"time"

	"github.com/go-rod/rod"
)

// Rod is responsible for browsser operations.
type Rod struct {
	// Browser is a rod Browser instance.
	Browser *rod.Browser
	// LoadTimeout controlls max page load time before context is canceled.
	LoadTimeout time.Duration
	// PageIdleTime sets the wait time after the page stops receiving requests.
	PageIdleTime time.Duration
}

// Connect starts the Browser connection.
func (r *Rod) Connect() {
	r.Browser = rod.New().MustConnect()
}

// Close closes the Browser connection.
func (r *Rod) Close() {
	r.Browser.MustClose()
}

// UrlToPage converts the URL to a rod Page instance.
func (r *Rod) UrlToPage(url string) *rod.Page {
	return r.Browser.MustPage(url)
}

// ByteToPage converts the binary to a rod Page instance.
func (r *Rod) ByteToPage(bin []byte) (*rod.Page, error) {
	file, err := os.CreateTemp("", "*.html")
	if err != nil {
		return &rod.Page{}, err
	}

	defer os.Remove(file.Name())

	if _, err = file.Write(bin); err != nil {
		return &rod.Page{}, err
	}

	page := r.Browser.MustPage("file://" + file.Name())

	return page, nil
}

// WaitLoad sets a wait time according to the page loading.
func (r *Rod) WaitLoad(page *rod.Page) {
	page = page.Timeout(r.LoadTimeout).MustWaitLoad()

	wait := page.WaitRequestIdle(r.PageIdleTime, nil, nil)
	wait()

	page.CancelTimeout()
}
