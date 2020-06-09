package telebot

import (
	"io"
	"os"
)

// File object represents any sort of file.
type File struct {
	FileID   string `json:"file_id"`
	UniqueID string `json:"file_unique_id"`
	FileSize int    `json:"file_size"`

	// file on telegram server https://core.telegram.org/bots/api#file
	FilePath string `json:"file_path"`

	// file on local file system.
	FileLocal string `json:"file_local"`

	// file on the internet
	FileURL string `json:"file_url"`

	// file backed with io.Reader
	FileReader io.Reader `json:"-"`
}

// FromDisk constructs a new local (on-disk) file object.
//
// Note, it returns File, not *File for a very good reason:
// in telebot, File is pretty much an embeddable struct,
// so upon uploading media you'll need to set embedded File
// with something. NewFile() returning File makes it a one-liner.
//
//     photo := &tb.Photo{File: tb.FromDisk("chicken.jpg")}
//
func FromDisk(filename string) File {
	return File{FileLocal: filename}
}

// FromURL constructs a new file on provided HTTP URL.
//
// Note, it returns File, not *File for a very good reason:
// in telebot, File is pretty much an embeddable struct,
// so upon uploading media you'll need to set embedded File
// with something. NewFile() returning File makes it a one-liner.
//
//     photo := &tb.Photo{File: tb.FromURL("https://site.com/picture.jpg")}
//
func FromURL(url string) File {
	return File{FileURL: url}
}

// FromReader constructs a new file from io.Reader.
//
// Note, it returns File, not *File for a very good reason:
// in telebot, File is pretty much an embeddable struct,
// so upon uploading media you'll need to set embedded File
// with something. NewFile() returning File makes it a one-liner.
//
//     photo := &tb.Photo{File: tb.FromReader(bytes.NewReader(...))}
//
func FromReader(reader io.Reader) File {
	return File{FileReader: reader}
}

func (f *File) stealRef(g *File) {
	if g.OnDisk() {
		f.FileLocal = g.FileLocal
	}

	if g.FileURL != "" {
		f.FileURL = g.FileURL
	}
}

// InCloud tells whether the file is present on Telegram servers.
func (f *File) InCloud() bool {
	return f.FileID != ""
}

// OnDisk will return true if file is present on disk.
func (f *File) OnDisk() bool {
	_, err := os.Stat(f.FileLocal)
	return err == nil
}
