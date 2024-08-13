package orcgen_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen"
	"github.com/luabagg/orcgen/pkg/handlers/pdf"
	"github.com/luabagg/orcgen/pkg/handlers/screenshot"
	"github.com/luabagg/orcgen/pkg/webdriver"
)

// ExampleCompleteConfiguration is an example of how to use the package structs directly.
func ExampleCompleteConfiguration() {
	screenshotHandler := screenshot.New()
	screenshotHandler.SetFullPage(false)

	wd := webdriver.FromDefault()
	defer wd.Close()

	page := wd.UrlToPage("https://google.com")
	wd.WaitLoad(page)
	page.MustInsertText("github orcgen package golang").Keyboard.Type(input.Enter)
	wd.WaitLoad(page)

	filename := "google.png"
	fileinfo, err := screenshotHandler.GenerateFile(page)
	if err == nil {
		fileinfo.Output(getName(filename))
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Output:
	// google.png generated successfully
}

// ExampleConvertWebpage gives examples using the ConvertWebpage function.
func ExampleGenerate() {
	// Converting the GitHub homepage to a webp file.
	filename := "github.webp"
	err := orcgen.Generate(
		"https://www.github.com",
		proto.PageCaptureScreenshot{
			Format: proto.PageCaptureScreenshotFormatWebp,
		},
		getName(filename),
	)
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Converting the HTML file to a PDF file.
	filename = "html.pdf"
	err = orcgen.Generate(
		getHTML(),
		proto.PagePrintToPDF{
			Landscape:           true,
			DisplayHeaderFooter: true,
			PrintBackground:     true,
			MarginTop:           new(float64),
			MarginBottom:        new(float64),
			MarginLeft:          new(float64),
			MarginRight:         new(float64),
			PreferCSSPageSize:   true,
		},
		getName(filename),
	)
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Output:
	// github.webp generated successfully
	// html.pdf generated successfully
}

// ExampleNewHandler shows how to use ExampleNewHandler function to create a new handler.
func ExampleNewHandler() {
	screenshotHandler := orcgen.NewHandler(
		proto.PageCaptureScreenshot{
			Format: proto.PageCaptureScreenshotFormatWebp,
		},
	)
	screenshotHandler.SetFullPage(true)

	pdfHandler := orcgen.NewHandler(
		proto.PagePrintToPDF{
			PrintBackground: false,
		},
	)
	pdfHandler.SetFullPage(false)
}

// ExampleConvertWebpage gives examples using the ConvertWebpage function.
func ExampleConvertWebpage() {
	// Converting the Faceboox homepage to a PNG file.
	filename := "facebook.png" // png is the default extension for screenshots.
	fileinfo, err := orcgen.ConvertWebpage(
		screenshot.New(), "https://www.facebook.com",
	)

	err = fileinfo.Output(getName(filename))
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Converting the X homepage to a PDF file.
	filename = "x.pdf"
	fileinfo, err = orcgen.ConvertWebpage(
		pdf.New().SetFullPage(true), "https://www.x.com",
	)

	err = fileinfo.Output(getName(filename))
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Output:
	// facebook.png generated successfully
	// x.pdf generated successfully
}

// ExampleConvertHTML gives examples using the ConvertHTML function.
func ExampleConvertHTML() {
	// Converting the HTML file to a JPG file.
	filename := "html.jpg"
	fileinfo, err := orcgen.ConvertHTML(
		screenshot.New().SetConfig(proto.PageCaptureScreenshot{
			Format: proto.PageCaptureScreenshotFormatJpeg,
		}),
		getHTML(),
	)
	err = fileinfo.Output(getName(filename))
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Converting the HTML file to a PDF file.
	filename = "html.pdf"
	fileinfo, err = orcgen.ConvertHTML(pdf.New(), getHTML())
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
