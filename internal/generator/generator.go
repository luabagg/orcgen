// package generator contains the implementation of the generators used for page conversion (example HTML to PDF).
package generator

// Generator interface contains the methods used for the page conversion.
type Generator interface {
	ConvertWebpage(url string) error
	ConvertHTML(html []byte) error
	GetOutput(filepath string) error
	GetFilesize() int
	Close()
}
