package png

import (
	"testing"

	"github.com/go-rod/rod"
	"github.com/luabagg/orcgen/internal/generator"
	"github.com/stretchr/testify/assert"
)

func TestPNGBuilder_SetFullPage(t *testing.T) {
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
			// create a new PNGBuilder instance
			instance := new(PNGBuilder).SetFullPage(tc.input)

			assert.Implements(t, (*generator.Generator)(nil), instance, "expected to be a Generator instance")
		})
	}
}

func TestPNGBuilder_GenerateFile(t *testing.T) {
	// create a new browser instance
	b := rod.New().MustConnect()
	defer b.MustClose()

	// create a new PNGBuilder instance
	pngBuilder := &PNGBuilder{}

	tests := []struct {
		name     string
		instance generator.Generator
		input    *rod.Page
	}{
		{
			name:     "simple page",
			instance: pngBuilder.SetFullPage(false),
			input:    b.MustPage("https://www.example.com").MustWaitLoad(),
		},
		{
			name:     "fullpage",
			instance: pngBuilder.SetFullPage(true),
			input:    b.MustPage("https://www.example.com").MustWaitLoad(),
		},
		{
			instance: pngBuilder,
			input:    b.MustPage(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// generate the PNG file
			pngData, err := tc.instance.GenerateFile(tc.input)

			assert.NoError(t, err, "Expected no error")
			assert.NotEmpty(t, pngData, "Expected bytes")

			tc.input.MustClose()
		})
	}
}
