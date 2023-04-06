package director

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/internal"
	"github.com/luabagg/orcgen/internal/generator"
	"github.com/luabagg/orcgen/internal/generator/jpeg"
	"github.com/luabagg/orcgen/internal/generator/pdf"
	"github.com/luabagg/orcgen/internal/generator/png"
	"github.com/stretchr/testify/assert"
)

func Example() {
	// starts a new connection and converts the byte array to PNG
	fi, _ := NewDirector(internal.PNG).Connect().ConvertHTML([]byte("Example"))
	fmt.Printf("File size: %d bytes\n", fi.Filesize)

	// saves the file output to "test.png"
	filepath := "test.png"
	fi.Output(filepath)

	// gets file info
	info, _ := os.Stat(filepath)
	fmt.Printf("New file saved with %d bytes", info.Size())

	// Output:
	// File size: 7395 bytes
	// New file saved with 7395 bytes

	os.Remove(filepath)
}

func TestNewDirector(t *testing.T) {
	tests := []struct {
		name string
		ext  internal.Ext
	}{
		{
			"test PDF",
			internal.PDF,
		},
		{
			"test PNG",
			internal.PNG,
		},
		{
			"test JPEG",
			internal.JPEG,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d := NewDirector(tc.ext)

			assert.NotNil(t, d.rod)
			assert.Implements(t, (*generator.Generator)(nil), d.generator, "expected to be a Generator instance")
		})
	}
}

func TestDirector_Connect(t *testing.T) {
	d := NewDirector(internal.Ext(0))
	d.Connect()
	defer d.rod.Browser.Close()

	assert.NotNil(t, d.rod.Browser)
}

func TestDirector_SetExt(t *testing.T) {
	testCases := []struct {
		name    string
		ext     internal.Ext
		builder generator.Generator
	}{
		{
			name:    "test PDF ext",
			ext:     internal.PDF,
			builder: &pdf.PDFBuilder{},
		},
		{
			name:    "test PNG ext",
			ext:     internal.PNG,
			builder: &png.PNGBuilder{},
		},
		{
			name:    "test JPEG ext",
			ext:     internal.JPEG,
			builder: &jpeg.JPEGBuilder{},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d := NewDirector(internal.Ext(0))
			dcopy := *d.SetExt(tc.ext)

			assert.IsType(t, Director{}, dcopy)
			assert.IsType(t, tc.builder, d.generator)
		})
	}
}

func TestDirector_SetFullPage(t *testing.T) {
	d := NewDirector(internal.Ext(0))
	d = d.SetFullPage(true)

	assert.IsType(t, new(Director), d)
}

func TestDirector_SetLoadTimeout(t *testing.T) {
	time := 1 * time.Second
	d := NewDirector(internal.Ext(0))
	d = d.SetLoadTimeout(time)

	assert.IsType(t, new(Director), d)
	assert.Equal(t, time, d.rod.LoadTimeout)
}

func TestDirector_SetPageIdleTime(t *testing.T) {
	time := 1 * time.Second
	d := NewDirector(internal.Ext(0))
	d = d.SetPageIdleTime(time)

	assert.IsType(t, new(Director), d)
	assert.Equal(t, time, d.rod.PageIdleTime)
}

func TestDirector_convert(t *testing.T) {
	d := NewDirector(internal.Ext(0)).Connect()
	d.generator = &MockGenerator{}

	pageMock := d.rod.Browser.MustPage()
	defer pageMock.Close()

	mockFileInfo := &Fileinfo{File: []byte("mock file"), Filesize: 9}

	fileInfo, err := d.convert(pageMock)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	assert.Equal(t, mockFileInfo.File, fileInfo.File)
	assert.Equal(t, mockFileInfo.Filesize, fileInfo.Filesize)
}

func TestDirector_ConvertWebpage(t *testing.T) {
	d := NewDirector(internal.Ext(0)).Connect()

	fileInfo, err := d.ConvertWebpage("https://google.com")

	assert.NoError(t, err)
	assert.NotNil(t, fileInfo.File)
	assert.True(t, fileInfo.Filesize > 0)
}

func TestDirector_ConvertHTML(t *testing.T) {
	d := NewDirector(internal.Ext(0)).Connect()

	html := []byte("<html><head><title>ORC gen</title></head><body><h1>Hello, World!</h1></body></html>")
	fileInfo, err := d.ConvertHTML(html)

	assert.NoError(t, err)
	assert.NotNil(t, fileInfo.File)
	assert.True(t, fileInfo.Filesize > 0)
}

func TestDirector_Close(t *testing.T) {
	d := NewDirector(internal.Ext(0)).Connect()
	d.Close()

	page, _ := d.rod.Browser.Page(proto.TargetCreateTarget{})
	assert.Nil(t, page)
}

type MockGenerator struct{}

func (m *MockGenerator) GenerateFile(page *rod.Page) ([]byte, error) {
	return []byte("mock file"), nil
}

func (m *MockGenerator) SetFullPage(fullPage bool) generator.Generator {
	return m
}
