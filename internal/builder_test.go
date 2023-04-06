package internal

import (
	"testing"

	"github.com/luabagg/orcgen/internal/generator"
	"github.com/luabagg/orcgen/internal/generator/jpeg"
	"github.com/luabagg/orcgen/internal/generator/pdf"
	"github.com/luabagg/orcgen/internal/generator/png"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	testCases := []struct {
		name    string
		ext     Ext
		builder generator.Generator
	}{
		{
			name:    "test PDF build",
			ext:     PDF,
			builder: &pdf.PDFBuilder{},
		},
		{
			name:    "test PNG build",
			ext:     PNG,
			builder: &png.PNGBuilder{},
		},
		{
			name:    "test JPEG build",
			ext:     JPEG,
			builder: &jpeg.JPEGBuilder{},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// builds from ext
			builder := Build(tc.ext)

			assert.Implements(t, (*generator.Generator)(nil), builder)
			assert.IsType(t, tc.builder, builder)
		})
	}
}
