package pdf

import (
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/luabagg/orcgen/pkg/handlers"
	"github.com/stretchr/testify/assert"
)

func TestPDFHandler_SetConfig(t *testing.T) {
	tests := []struct {
		name  string
		input proto.PagePrintToPDF
	}{
		{
			name:  "valid config",
			input: proto.PagePrintToPDF{},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// create a new PDFHandler instance
			instance := New().SetConfig(tc.input)

			assert.Implements(t, (*handlers.FileHandler[proto.PagePrintToPDF])(nil), instance, "expected to be a Generator instance")
		})
	}
}

func TestPDFHandler_SetFullPage(t *testing.T) {
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
			// create a new PDFHandler instance
			instance := New().SetFullPage(tc.input)

			assert.Implements(t, (*handlers.FileHandler[proto.PagePrintToPDF])(nil), instance, "expected to be a Generator instance")
		})
	}
}

func TestPDFHandler_GenerateFile(t *testing.T) {
	// create a new browser instance
	b := rod.New().MustConnect()
	defer b.MustClose()

	// create a new PDFHandler instance
	pdfHandler := New()

	tests := []struct {
		name     string
		instance handlers.FileHandler[proto.PagePrintToPDF]
		input    *rod.Page
	}{
		{
			name:     "simple page",
			instance: pdfHandler,
			input:    b.MustPage("https://www.example.com").MustWaitLoad(),
		},
		{
			name:     "fullpage",
			instance: pdfHandler.SetFullPage(true),
			input:    b.MustPage("https://www.example.com").MustWaitLoad(),
		},
		{
			instance: pdfHandler,
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
