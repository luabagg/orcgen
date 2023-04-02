// package internal contains the rod implementation.
package internal

import (
	"time"

	"github.com/luabagg/orcgen/internal/generator"
)

// Director controlls the page conversion methods.
type Director struct {
	generator generator.Generator
	rod       *Rod
}

// NewDirector opens a new Director instance.
func NewDirector(g generator.Generator) *Director {
	return &Director{
		generator: g,
		rod: &Rod{
			LoadTimeout:  200 * time.Millisecond,
			PageIdleTime: 10 * time.Second,
		},
	}
}

// Connect starts the Browser connection.
func (d *Director) Connect() *Director {
	d.rod.Connect()

	return d
}

// SetGenerator sets the specific builder.
func (d *Director) SetGenerator(g generator.Generator) *Director {
	d.generator = g

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

// ConvertWebpage converts from a URL.
func (d *Director) ConvertWebpage(url string) (*Fileinfo, error) {
	page := d.rod.UrlToPage(url)

	defer page.Close()

	b, err := d.generator.GenerateFile(page)
	if err != nil {
		return nil, err
	}

	return &Fileinfo{
		File:     b,
		Filesize: len(b),
	}, nil
}

// ConvertHTML converts from a file.
func (d *Director) ConvertHTML(html []byte) (*Fileinfo, error) {
	page, err := d.rod.ByteToPage(html)
	if err != nil {
		return nil, err
	}

	defer page.Close()

	b, err := d.generator.GenerateFile(page)
	if err != nil {
		return nil, err
	}

	return &Fileinfo{
		File:     b,
		Filesize: len(b),
	}, nil
}

// Close resets struct and close Browser connection.
func (d *Director) Close() {
	d.generator = nil
	d.rod.Close()
}
