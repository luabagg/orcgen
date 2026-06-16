package orcgen_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-rod/rod/lib/input"
	"github.com/luabagg/orcgen/v2"
	"github.com/luabagg/orcgen/v2/pkg/handlers/pdf"
	"github.com/luabagg/orcgen/v2/pkg/handlers/screenshot"
	"github.com/luabagg/orcgen/v2/pkg/webdriver"
)

// ExampleGenerateHTML_new demonstrates the new recommended API using the builder pattern.
// This approach eliminates runtime type assertions and provides compile-time type safety.
func ExampleGenerateHTML_new() {
	// Converting HTML bytes to a PDF file with full page capture
	filename := "html_new.pdf"
	err := orcgen.GenerateHTML(
		getHTML(),
		orcgen.PDF(orcgen.PDFConfig{
			Landscape:           true,
			DisplayHeaderFooter: true,
			PrintBackground:     true,
			PreferCSSPageSize:   true,
		}).FullPage(true),
		getName(filename),
	)
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Converting HTML bytes to a screenshot
	filename = "html_new.png"
	err = orcgen.GenerateHTML(
		getHTML(),
		orcgen.Screenshot(orcgen.ScreenshotConfig{
			Format: "png",
		}),
		getName(filename),
	)
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Output:
	// html_new.pdf generated successfully
	// html_new.png generated successfully
}

// ExampleGenerateURL_new demonstrates the new recommended API for URL conversion.
func ExampleGenerateURL_new() {
	// Converting a URL to a PDF
	filename := "github_new.pdf"
	err := orcgen.GenerateURL(
		"https://www.github.com",
		orcgen.PDF(orcgen.PDFConfig{
			Landscape:       true,
			PrintBackground: true,
		}),
		getName(filename),
	)
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Converting a URL to a WebP screenshot with full page
	filename = "github_new.webp"
	err = orcgen.GenerateURL(
		"https://www.github.com",
		orcgen.Screenshot(orcgen.ScreenshotConfig{
			Format: "webp",
		}).FullPage(true),
		getName(filename),
	)
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Output:
	// github_new.pdf generated successfully
	// github_new.webp generated successfully
}

// Example demonstrates advanced usage with direct webdriver control.
func Example() {
	wd := webdriver.FromDefault()
	defer wd.Close()

	// Using the page directly to search before screenshotting:
	page, err := wd.UrlToPage("https://google.com")
	if err != nil {
		return
	}

	if err := wd.WaitLoad(page); err != nil {
		return
	}

	page.MustInsertText("github orcgen package golang").Keyboard.Type(input.Enter)

	if err := wd.WaitLoad(page); err != nil {
		return
	}

	// Using the handler directly - creates a PNG of the Google search:
	screenshotHandler := screenshot.New().SetConfig(orcgen.ScreenshotConfig{
		Format: "png",
	})

	fileinfo, err := screenshotHandler.GenerateFile(page)
	if err == nil {
		filename := "google.png"
		fileinfo.Output(getName(filename))
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Creates a PDF of the Google search:
	pdfHandler := pdf.New().SetConfig(orcgen.PDFConfig{
		PrintBackground: true,
		PageRanges:      "1,2",
	})

	fileinfo, err = pdfHandler.GenerateFile(page)
	if err == nil {
		filename := "google.pdf"
		fileinfo.Output(getName(filename))
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Output:
	// google.png generated successfully
	// google.pdf generated successfully
}

// ExampleGenerateURL demonstrates URL conversion with the fluent builder API.
func ExampleGenerateURL() {
	// Converting a URL to a WebP screenshot
	filename := "github.webp"
	err := orcgen.GenerateURL(
		"https://www.github.com",
		orcgen.Screenshot(orcgen.ScreenshotConfig{
			Format: "webp",
		}),
		getName(filename),
	)
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Output:
	// github.webp generated successfully
}

// ExampleGenerateHTML demonstrates HTML conversion with the fluent builder API.
func ExampleGenerateHTML() {
	// Converting HTML bytes to a PDF file
	filename := "html.pdf"
	err := orcgen.GenerateHTML(
		getHTML(),
		orcgen.PDF(orcgen.PDFConfig{
			Landscape:           true,
			DisplayHeaderFooter: true,
			PrintBackground:     true,
			MarginTop:           new(float64),
			MarginBottom:        new(float64),
			MarginLeft:          new(float64),
			MarginRight:         new(float64),
			PreferCSSPageSize:   true,
		}),
		getName(filename),
	)
	if err == nil {
		fmt.Printf("%s generated successfully\n", filename)
	}

	// Output:
	// html.pdf generated successfully
}

// ExampleHandlerBuilder demonstrates the builder pattern with method chaining.
func ExampleHandlerBuilder() {
	// Screenshot builder with full page enabled
	screenshotBuilder := orcgen.Screenshot(orcgen.ScreenshotConfig{
		Format: "png",
	}).FullPage(true)

	// PDF builder with full page disabled
	pdfBuilder := orcgen.PDF(orcgen.PDFConfig{
		PrintBackground: true,
	}).FullPage(false)

	// Use the builders
	orcgen.GenerateURL("https://example.com", screenshotBuilder, getName("example.png"))
	orcgen.GenerateHTML(getHTML(), pdfBuilder, getName("example.pdf"))
}

// ExampleConvertURL demonstrates using ConvertURL for more control.
func ExampleConvertURL() {
	// Converting a URL to a PNG file
	filename := "facebook.png"
	fileinfo, err := orcgen.ConvertURL(
		orcgen.Screenshot(orcgen.ScreenshotConfig{
			Format: "png",
		}),
		"https://www.facebook.com",
	)

	if err == nil {
		err = fileinfo.Output(getName(filename))
		if err == nil {
			fmt.Printf("%s generated successfully\n", filename)
		}
	}

	// Converting a URL to a PDF file with full page
	filename = "x.pdf"
	fileinfo, err = orcgen.ConvertURL(
		orcgen.PDF(orcgen.PDFConfig{
			PrintBackground: true,
		}).FullPage(true),
		"https://www.x.com",
	)

	if err == nil {
		err = fileinfo.Output(getName(filename))
		if err == nil {
			fmt.Printf("%s generated successfully\n", filename)
		}
	}

	// Output:
	// facebook.png generated successfully
	// x.pdf generated successfully
}

// ExampleConvertHTML demonstrates using ConvertHTML for more control.
func ExampleConvertHTML() {
	// Converting HTML to a JPG file
	filename := "html.jpg"
	fileinfo, err := orcgen.ConvertHTML(
		orcgen.Screenshot(orcgen.ScreenshotConfig{
			Format: "jpeg",
		}),
		getHTML(),
	)
	if err == nil {
		err = fileinfo.Output(getName(filename))
		if err == nil {
			fmt.Printf("%s generated successfully\n", filename)
		}
	}

	// Converting HTML to a PDF file
	filename = "html.pdf"
	fileinfo, err = orcgen.ConvertHTML(
		orcgen.PDF(orcgen.PDFConfig{
			PrintBackground: true,
		}),
		getHTML(),
	)
	if err == nil {
		err = fileinfo.Output(getName(filename))
		if err == nil {
			fmt.Printf("%s generated successfully\n", filename)
		}
	}

	// Output:
	// html.jpg generated successfully
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
