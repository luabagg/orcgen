package director

import (
	rodutils "github.com/go-rod/rod/lib/utils"
)

// Fileinfo is used for file information control.
type Fileinfo struct {
	File     []byte
	Filesize int
}

// Output saves the file to the informed filepath.
func (f *Fileinfo) Output(filepath string) error {
	return rodutils.OutputFile(filepath, f.File)
}
