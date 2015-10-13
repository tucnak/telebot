package telebot

import (
	"fmt"
	"os"
)

// File object represents any sort of file.
type File struct {
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size"`

	// Local absolute path to file on local file system.
	filename string
}

// NewFile attempts to create a File object, leading to a real
// file on the file system, that could be uploaded later.
//
// Notice that NewFile doesn't upload file, but only creates
// a descriptor for it.
func NewFile(path string) (File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return File{}, FileError{
			fmt.Sprintf("'%s' does not exist!", path),
		}
	}

	return File{filename: path}, nil
}

// Exists says whether the file presents on Telegram servers or not.
func (f File) Exists() bool {
	return f.FileID != ""
}

// Local returns location of file on local file system, if it's
// actually there, otherwise returns empty string.
func (f File) Local() string {
	return f.filename
}
