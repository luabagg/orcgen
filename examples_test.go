package orcgen_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen"
	"github.com/luabagg/orcgen/pkg/handlers/pdf"
	"github.com/luabagg/orcgen/pkg/handlers/screenshot"
)

// ExampleConvertWebpage gives examples using the ConvertWebpage function.
func ExampleConvertWebpage() {
	filename := "facebook.webp"
	screenshotHandler := screenshot.New(proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatWebp,
	})
	fileinfo, err := orcgen.ConvertWebpage(screenshotHandler, "https://www.facebook.com")
	err = fileinfo.Output(getName(filename))
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	filename = "facebook.pdf"
	pdfHandler := pdf.New(proto.PagePrintToPDF{
		PageRanges: "1",
	})
	fileinfo, err = orcgen.ConvertWebpage(pdfHandler, "https://www.facebook.com")
	err = fileinfo.Output(getName(filename))
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Output:
	// facebook.webp generated successfully
	// facebook.pdf generated successfully
}

// ExampleConvertHTML gives examples using the ConvertHTML function.
func ExampleConvertHTML() {
	filename := "html.webp"
	screenshotHandler := screenshot.New(proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatWebp,
	})
	fileinfo, err := orcgen.ConvertHTML(screenshotHandler, getHTML())
	err = fileinfo.Output(getName(filename))
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	filename = "html.pdf"
	pdfHandler := pdf.New(proto.PagePrintToPDF{
		Landscape:           true,
		DisplayHeaderFooter: true,
		PrintBackground:     true,
		MarginTop:           new(float64),
		MarginBottom:        new(float64),
		MarginLeft:          new(float64),
		MarginRight:         new(float64),
		PreferCSSPageSize:   true,
	})
	fileinfo, err = orcgen.ConvertHTML(pdfHandler, getHTML())
	err = fileinfo.Output(getName(filename))
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Output:
	// html.webp generated successfully
	// html.pdf generated successfully
}

func getHTML() []byte {
	file := filepath.Join(getBasepath(), "testdata/test.html")
	html, _ := os.ReadFile(file)

	return html
}

func getName(name string) string {
	return filepath.Join(getBasepath(), "testdata", name)
}

func getBasepath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}
