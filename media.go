package telebot

import (
	"encoding/json"
)

// Album lets you group multiple media (so-called InputMedia)
// into a single messsage.
//
// On older clients albums look like N regular messages.
type Album []InputMedia

// InputMedia is a generic type for all kinds of media you
// can put into an album.
type InputMedia interface {
	// As some files must be uploaded (instead of referencing)
	// outer layers of Telebot require it.
	MediaFile() *File
}

// Photo object represents a single photo file.
type Photo struct {
	File

	Width  int `json:"width"`
	Height int `json:"height"`

	// (Optional)
	Caption string `json:"caption,omitempty"`
}

type photoSize struct {
	File
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Caption string `json:"caption,omitempty"`
}

// MediaFile returns &Photo.File
func (p *Photo) MediaFile() *File {
	return &p.File
}

// UnmarshalJSON is custom unmarshaller required to abstract
// away the hassle of treating different thumbnail sizes.
// Instead, Telebot chooses the hi-res one and just sticks to
// it.
//
// I really do find it a beautiful solution.
func (p *Photo) UnmarshalJSON(jsonStr []byte) error {
	var hq photoSize

	if jsonStr[0] == '{' {
		if err := json.Unmarshal(jsonStr, &hq); err != nil {
			return err
		}
	} else {
		var sizes []photoSize

		if err := json.Unmarshal(jsonStr, &sizes); err != nil {
			return err
		}

		hq = sizes[len(sizes)-1]
	}

	p.File = hq.File
	p.Width = hq.Width
	p.Height = hq.Height

	return nil
}

// Audio object represents an audio file.
type Audio struct {
	File

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration,omitempty"`

	// (Optional)
	Caption   string `json:"caption,omitempty"`
	Title     string `json:"title,omitempty"`
	Performer string `json:"performer,omitempty"`
	MIME      string `json:"mime_type,omitempty"`
}

// Document object represents a general file (as opposed to Photo or Audio).
// Telegram users can send files of any type of up to 1.5 GB in size.
type Document struct {
	File

	// Original filename as defined by sender.
	FileName string `json:"file_name"`

	// (Optional)
	Thumbnail *Photo `json:"thumb,omitempty"`
	Caption   string `json:"caption,omitempty"`
	MIME      string `json:"mime_type"`
}

// Video object represents a video file.
type Video struct {
	File

	Width  int `json:"width"`
	Height int `json:"height"`

	Duration int `json:"duration,omitempty"`

	// (Optional)
	Caption   string `json:"caption,omitempty"`
	Thumbnail *Photo `json:"thumb,omitempty"`
	MIME      string `json:"mime_type,omitempty"`
}

// MediaFile returns &Video.File
func (v *Video) MediaFile() *File {
	return &v.File
}

// Voice object represents a voice note.
type Voice struct {
	File

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration"`

	// (Optional)
	MIME string `json:"mime_type,omitempty"`
}

// VideoNote represents a video message (available in Telegram apps
// as of v.4.0).
type VideoNote struct {
	File

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration"`

	// (Optional)
	Thumbnail *Photo `json:"thumb,omitempty"`
}

// Contact object represents a contact to Telegram user
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`

	// (Optional)
	LastName string `json:"last_name"`
	UserID   int    `json:"user_id,omitempty"`
}

// Location object represents geographic position.
type Location struct {
	// Latitude
	Lat float32 `json:"latitude"`
	// Longitude
	Lng float32 `json:"longitude"`

	// Period in seconds for which the location will be updated
	// (see Live Locations, should be between 60 and 86400.)
	LivePeriod int `json:"live_period,omitempty"`
}

// Venue object represents a venue location with name, address and
// optional foursquare ID.
type Venue struct {
	Location Location `json:"location"`
	Title    string   `json:"title"`
	Address  string   `json:"address"`

	// (Optional)
	FoursquareID string `json:"foursquare_id,omitempty"`
}
