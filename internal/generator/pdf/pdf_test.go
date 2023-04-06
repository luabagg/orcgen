package pdf

import (
	"testing"

	"github.com/go-rod/rod"
	"github.com/luabagg/orcgen/internal/generator"
	"github.com/stretchr/testify/assert"
)

func TestPDFBuilder_SetFullPage(t *testing.T) {
	tests := []struct {
		name  string
		input bool
	}{
		{
			name:  "simple page",
			input: false,
		},
		{
			name:  "fullpage",
			input: true,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// create a new PDFBuilder instance
			instance := new(PDFBuilder).SetFullPage(tc.input)

			assert.Implements(t, (*generator.Generator)(nil), instance, "expected to be a Generator instance")
		})
	}
}

func TestPDFBuilder_GenerateFile(t *testing.T) {
	// create a new browser instance
	b := rod.New().MustConnect()
	defer b.MustClose()

	// create a new PDFBuilder instance
	pdfBuilder := &PDFBuilder{}

	tests := []struct {
		name     string
		instance generator.Generator
		input    *rod.Page
	}{
		{
			name:     "simple page",
			instance: pdfBuilder.SetFullPage(false),
			input:    b.MustPage("https://www.example.com").MustWaitLoad(),
		},
		{
			name:     "fullpage",
			instance: pdfBuilder.SetFullPage(true),
			input:    b.MustPage("https://www.example.com").MustWaitLoad(),
		},
		{
			instance: pdfBuilder,
			input:    b.MustPage(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// generate the PDF file
			pdfData, err := tc.instance.GenerateFile(tc.input)

			assert.NoError(t, err, "Expected no error")
			assert.NotEmpty(t, pdfData, "Expected bytes")

			tc.input.MustClose()
		})
	}
}
