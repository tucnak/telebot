package telebot

import (
	"net/url"
	"os"
)

// File object represents any sort of file.
type File struct {
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size"`

	// file on telegram server https://core.telegram.org/bots/api#file
	FilePath string `json:"file_path"`

	// file on local file system.
	FileLocal string

	// file on the internet
	FileURL *url.URL
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

// FromURL constructs a new file on provided HTTPS URL.
//
// Note, it returns File, not *File for a very good reason:
// in telebot, File is pretty much an embeddable struct,
// so upon uploading media you'll need to set embedded File
// with something. NewFile() returning File makes it a one-liner.
//
//     photo := &tb.Photo{File: tb.FromURL("https://site.com/picture.jpg")}
//
func FromURL(u *url.URL) File {
	return File{FileURL: u}
}

func (f *File) importLocal(g *File) {
	if !g.OnDisk() {
		return
	}

	f.FileLocal = g.FileLocal
}

// InCloud tells whether the file is present on Telegram servers.
func (f *File) InCloud() bool {
	return f.FileID != ""
}

// OnDisk will return true if file is present on disk.
func (f *File) OnDisk() bool {
	if _, err := os.Stat(f.FileLocal); err != nil {
		return false
	}

	return true
}
