package orcgen

// import (
// 	"os"
// 	"testing"

// 	"github.com/luabagg/orcgen/internal"
// 	"github.com/luabagg/orcgen/pkg/director"
// 	"github.com/stretchr/testify/assert"
// )

// func TestNew(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		ext  internal.Ext
// 	}{
// 		{
// 			"test PDF",
// 			PDF,
// 		},
// 		{
// 			"test PNG",
// 			PNG,
// 		},
// 		{
// 			"test JPEG",
// 			JPEG,
// 		},
// 	}
// 	for _, tc := range tests {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()

// 			d := New(tc.ext)
// 			defer d.Close()
// 			assert.IsType(t, new(director.Director), d, "expected to be a Director")
// 		})
// 	}
// }

// func TestConvertWebpage(t *testing.T) {
// 	url := "https://www.example.com"

// 	tests := []struct {
// 		name   string
// 		ext    internal.Ext
// 		url    string
// 		output string
// 	}{
// 		{
// 			"test PDF",
// 			PDF,
// 			url,
// 			"test.pdf",
// 		},
// 		{
// 			"test PNG",
// 			PNG,
// 			url,
// 			"test.png",
// 		},
// 		{
// 			"test JPEG",
// 			JPEG,
// 			url,
// 			"test.jpeg",
// 		},
// 	}
// 	for _, tc := range tests {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()

// 			err := ConvertWebpage(tc.ext, tc.url, tc.output)
// 			assert.NoError(t, err)

// 			os.Remove(tc.output)
// 		})
// 	}
// }

// func TestConvertHTML(t *testing.T) {
// 	html, _ := os.ReadFile("testdata/test.html")

// 	tests := []struct {
// 		name   string
// 		ext    internal.Ext
// 		html   []byte
// 		output string
// 	}{
// 		{
// 			"test PDF",
// 			PDF,
// 			html,
// 			"test.pdf",
// 		},
// 		{
// 			"test PNG",
// 			PNG,
// 			html,
// 			"test.png",
// 		},
// 		{
// 			"test JPEG",
// 			JPEG,
// 			html,
// 			"test.jpeg",
// 		},
// 	}
// 	for _, tc := range tests {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()

// 			err := ConvertHTML(tc.ext, tc.html, tc.output)
// 			assert.NoError(t, err)

// 			os.Remove(tc.output)
// 		})
// 	}
// }
