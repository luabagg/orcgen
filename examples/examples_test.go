package examples

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	orcgen "github.com/luabagg/orcgen/pkg"
)

// import "github.com/luabagg/orcgen"

// Example_New gives examples using the New function from orcgen.
func ExampleNew() {
	// starts the connection.
	gen := orcgen.New(orcgen.PDF)
	defer gen.Close()

	/* using for HTML conversion */

	// this generates a pdf file with the HTML content.
	f, _ := gen.ConvertHTML(getHTML())

	filename := "html.pdf"
	if err := f.Output(filename); err == nil {
		fmt.Printf("%s generated succesfully\n", filename)
	}

	// this generates a png file with the HTML content.
	// notice the SetFullPage use here.
	f, _ = gen.SetExt(orcgen.PNG).
		SetFullPage(true).
		ConvertHTML(getHTML())

	filename = "html.png"
	if err := f.Output(filename); err == nil {
		fmt.Printf("%s generated succesfully\n", filename)
	}

	/* using for URL conversion */

	// this generates a pdf file from www.google.com.
	f, _ = gen.SetExt(orcgen.PDF).
		ConvertWebpage("https://www.google.com")

	filename = "google.pdf"
	if err := f.Output(filename); err == nil {
		fmt.Printf("%s generated succesfully\n", filename)
	}

	// this generates a jpeg file from www.twitter.com.
	// full config example.
	f, _ = gen.SetExt(orcgen.JPEG).
		SetFullPage(true).
		SetLoadTimeout(5 * time.Second).
		SetPageIdleTime(3 * time.Second).
		ConvertWebpage("https://www.twitter.com")

	filename = "twitter.jpeg"
	if err := f.Output(filename); err == nil {
		fmt.Printf("%s generated succesfully\n", filename)
	}

	// Output:
	// html.pdf generated succesfully
	// html.png generated succesfully
	// google.pdf generated succesfully
	// twitter.jpeg generated succesfully
}

// Example_ConvertHTML gives examples using the ConvertHTML function from orcgen.
func ExampleConvertHTML() {
	filename := "html.pdf"
	err := orcgen.ConvertHTML(orcgen.PDF, getHTML(), filename)
	if err == nil {
		fmt.Printf("%s generated succesfully\n", filename)
	}

	// Output:
	// html.pdf generated succesfully
}

// Example_ConvertWebpage gives examples using the ConvertWebpage function from orcgen.
func ExampleConvertWebpage() {
	filename := "github.pdf"
	err := orcgen.ConvertWebpage(orcgen.PDF, "https://www.github.com", filename)
	if err == nil {
		fmt.Printf("%s generated succesfully\n", filename)
	}

	// Output:
	// github.pdf generated succesfully
}

func getHTML() []byte {
	file := filepath.Join(getBasepath(), "test.html")
	html, _ := os.ReadFile(file)

	return html
}

func getBasepath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}
