package internal

import (
	"bytes"

	"golang.org/x/net/html"
)

// IsValidHTML checks if the provided HTML string is valid.
func IsValidHTML(binHtml []byte) bool {
	reader := bytes.NewReader(binHtml)
	_, err := html.Parse(reader)
	return err == nil
}
