package telebot

import "encoding/json"

type photoSize struct {
	File

	Width  int `json:"width"`
	Height int `json:"height"`

	// (Optional)
	Caption string `json:"caption,omitempty"`
}

// Photo object represents a single photo file.
type Photo struct {
	File

	Width  int `json:"width"`
	Height int `json:"height"`

	// (Optional)
	Caption string `json:"caption,omitempty"`
}

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
	Duration int `json:"duration"`

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
	Audio

	Width  int `json:"width"`
	Height int `json:"height"`

	Duration int `json:"duration"`

	// (Optional)
	Thumbnail *Photo `json:"thumb,omitempty"`
	MIME      string `json:"mime_type,omitempty"`
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
	Lat float64 `json:"latitude"`
	// Longitude
	Lng float64 `json:"longitude"`
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
