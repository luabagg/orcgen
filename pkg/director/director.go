package director

import (
	"time"

	gorod "github.com/go-rod/rod"
	"github.com/luabagg/orcgen/internal"
	"github.com/luabagg/orcgen/internal/configurator"
	"github.com/luabagg/orcgen/internal/generator"
	"github.com/luabagg/orcgen/internal/rod"
	"github.com/luabagg/orcgen/pkg/fileinfo"
)

// Director controls the page conversion methods.
type Director struct {
	generator    generator.Generator
	configurator configurator.Configurator
	rod          *rod.Rod
}

// NewDirector opens a new Director instance.
func NewDirector(ext internal.Ext) *Director {
	return &Director{
		generator: internal.Build(ext),
		rod: &rod.Rod{
			LoadTimeout:  10 * time.Second,
			PageIdleTime: 200 * time.Millisecond,
		},
	}
}

// Connect starts the Browser connection.
func (d *Director) Connect() *Director {
	d.rod.Connect()

	return d
}

// SetExt sets the extension to be converted to.
func (d *Director) SetExt(ext internal.Ext) *Director {
	d.generator = internal.Build(ext)

	return d
}

// SetFullPage sets the pages to be converted. If false, only the first page is selected.
func (d *Director) SetFullPage(fullPage bool) *Director {
	d.generator.SetFullPage(fullPage)

	return d
}

func (d *Director) Configure(c generator.Config) *Director {
	d.configurator.Configure(c)

	return d
}

// SetLoadTimeout controlls max page load time before context is canceled.
func (d *Director) SetLoadTimeout(t time.Duration) *Director {
	d.rod.LoadTimeout = t

	return d
}

// SetPageIdleTime sets the wait time after the page stops receiving requests.
func (d *Director) SetPageIdleTime(t time.Duration) *Director {
	d.rod.PageIdleTime = t

	return d
}

// convert converts a rod Page to a Fileinfo instance.
func (d *Director) convert(page *gorod.Page) (*fileinfo.Fileinfo, error) {
	defer page.Close()

	d.rod.WaitLoad(page)

	b, err := d.generator.GenerateFile(page)
	if err != nil {
		return nil, err
	}

	return &fileinfo.Fileinfo{
		File:     b,
		Filesize: len(b),
	}, nil
}

// ConvertWebpage converts from an URL.
func (d *Director) ConvertWebpage(url string) (*fileinfo.Fileinfo, error) {
	page := d.rod.UrlToPage(url)

	return d.convert(page)
}

// ConvertHTML converts from a file.
func (d *Director) ConvertHTML(html []byte) (*fileinfo.Fileinfo, error) {
	page, err := d.rod.ByteToPage(html)
	if err != nil {
		return nil, err
	}

	return d.convert(page)
}

// Close closes Browser connection.
func (d *Director) Close() {
	d.rod.Close()
}
