// package utils contains the rod implementation and other useful methods.
package utils

import (
	"os"
	"time"

	"github.com/go-rod/rod"
)

var (
	// Browser is a rod Browser instance.
	Browser *rod.Browser
	// Timeout controlls max page load time before context is canceled.
	Timeout time.Duration
	// IdleTime sets the wait time after the page stops receiving requests.
	IdleTime time.Duration
)

// waitLoad sets a wait time according to the page loading.
func waitLoad(page *rod.Page) {
	page = page.Timeout(Timeout).MustWaitLoad()

	wait := page.WaitRequestIdle(IdleTime, nil, nil)
	wait()

	page.CancelTimeout()
}

// UrlToPage converts the URL to a rod Page instance.
func UrlToPage(url string) *rod.Page {
	page := Browser.MustPage(url)

	waitLoad(page)

	return page
}

// UrlToPage converts the binary to a rod Page instance.
func ByteToPage(bin []byte) (*rod.Page, error) {
	file, err := os.CreateTemp("", "*.html")
	if err != nil {
		return &rod.Page{}, err
	}

	defer os.Remove(file.Name())

	if _, err = file.Write(bin); err != nil {
		return &rod.Page{}, err
	}

	page := Browser.MustPage("file://" + file.Name())

	waitLoad(page)

	return page, nil
}

// CloseBrowser closes the Browser connection.
func CloseBrowser() {
	Browser.MustClose()
}
