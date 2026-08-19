package rod

import (
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
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
func (r *Rod) Connect(opts ...Option) {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.InLambda {
		l := launchInLambda()
		u := l.Preferences(cfg.Preferences).
			WorkingDir(cfg.WorkingDir).
			UserDataDir(cfg.UserDataDir).
			MustLaunch()
		r.Browser = rod.New().ControlURL(u).MustConnect()
		return
	}

	l := launcher.New().
		Bin(cfg.Bin).
		Preferences(cfg.Preferences).
		WorkingDir(cfg.WorkingDir).
		UserDataDir(cfg.UserDataDir)
	u := l.MustLaunch()
	r.Browser = rod.New().ControlURL(u).MustConnect()
}

// Close closes the Browser connection.
func (r *Rod) Close() {
	r.Browser.MustClose()
}

// URLToPage converts the URL to a rod Page instance.
func (r *Rod) URLToPage(url string) *rod.Page {
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

	wait := page.WaitRequestIdle(r.PageIdleTime, nil, nil, nil)
	wait()

	page.CancelTimeout()
}

func launchInLambda() *launcher.Launcher {
	return launcher.New().
		// where lambda runtime stores chromium
		Bin("/opt/chromium").

		// recommended flags to run in serverless environments
		// see https://github.com/alixaxel/chrome-aws-lambda/blob/master/source/index.ts
		Set("allow-running-insecure-content").
		Set("autoplay-policy", "user-gesture-required").
		Set("disable-component-update").
		Set("disable-domain-reliability").
		Set("disable-features", "AudioServiceOutOfProcess", "IsolateOrigins", "site-per-process").
		Set("disable-print-preview").
		Set("disable-setuid-sandbox").
		Set("disable-site-isolation-trials").
		Set("disable-speech-api").
		Set("disable-web-security").
		Set("disk-cache-size", "33554432").
		Set("enable-features", "SharedArrayBuffer").
		Set("hide-scrollbars").
		Set("ignore-gpu-blocklist").
		Set("in-process-gpu").
		Set("mute-audio").
		Set("no-default-browser-check").
		Set("no-pings").
		Set("no-sandbox").
		Set("no-zygote").
		Set("single-process").
		Set("use-gl", "swiftshader").
		Set("window-size", "1920", "1080")
}
