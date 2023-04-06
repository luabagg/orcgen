package director

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileinfo_Output(t *testing.T) {
	file := "test.txt"

	tests := []struct {
		name     string
		fileData []byte
		filepath string
	}{
		{
			name:     "save file successfully",
			fileData: []byte("this is a test file"),
			filepath: file,
		},
		{
			name:     "overwrite existing file",
			fileData: []byte("new file data"),
			filepath: file,
		},
		{
			name:     "invalid filepath",
			fileData: []byte("invalid filepath test"),
			filepath: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new fileinfo instance
			fi := &Fileinfo{
				File:     tc.fileData,
				Filesize: len(tc.fileData),
			}

			// Output the file
			err := fi.Output(tc.filepath)

			if tc.filepath != "" {
				// Assert that the file exists and has the expected data
				assert.FileExists(t, tc.filepath)
				data, err := ioutil.ReadFile(tc.filepath)
				assert.NoError(t, err)
				assert.Equal(t, tc.fileData, data)
				assert.NoError(t, err)

				os.Remove(tc.filepath)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
