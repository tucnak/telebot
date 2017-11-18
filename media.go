package telebot

// Photo object represents a photo (with or without caption).
type Photo struct {
	File

	Width  int `json:"width"`
	Height int `json:"height"`

	Caption string `json:"caption,omitempty"`
}

// Audio object represents an audio file.
type Audio struct {
	File

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration"`

	// Title (optional) as defined by sender or by audio tags.
	Title string `json:"title"`

	// Performer (optional) is defined by sender or by audio tags.
	Performer string `json:"performer"`

	// MIME type (optional) of the file as defined by sender.
	Mime string `json:"mime_type"`

	Caption string `json:"caption,omitempty"`
}

// Voice object represents a voice note.
type Voice struct {
	File

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration"`

	// MIME type (optional) of the file as defined by sender.
	Mime string `json:"mime_type"`

	Caption string `json:"caption,omitempty"`
}

// Document object represents a general file (as opposed to Photo or Audio).
// Telegram users can send files of any type of up to 1.5 GB in size.
type Document struct {
	File

	// Document thumbnail as defined by sender.
	Preview Photo `json:"thumb"`

	// Original filename as defined by sender.
	FileName string `json:"file_name"`

	// MIME type of the file as defined by sender.
	Mime string `json:"mime_type"`

	Caption string `json:"caption,omitempty"`
}

// Sticker object represents a WebP image, so-called sticker.
type Sticker struct {
	File

	Width  int `json:"width"`
	Height int `json:"height"`

	// Sticker thumbnail in .webp or .jpg format.
	Thumbnail Photo `json:"thumb"`

	// Associated emoji
	Emoji string `json:"emoji"`
}

// Video object represents an MP4-encoded video.
type Video struct {
	Audio

	Width  int `json:"width"`
	Height int `json:"height"`

	// Text description of the video as defined by sender.
	Caption string `json:"caption,omitempty"`

	// Video thumbnail.
	Thumbnail Photo `json:"thumb"`
}

// This object represents a video message (available in Telegram apps
// as of v.4.0).
type VideoNote struct {
	File

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration"`

	// Video note thumbnail.
	Thumbnail Photo `json:"thumb"`
}

// Contact object represents a contact to Telegram user
type Contact struct {
	UserID      int    `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

// Location object represents geographic position.
type Location struct {
	Lat float32 `json:"latitude"`
	Lng float32 `json:"longitude"`
}

// Venue object represents a venue location with name, address and
// optional foursquare ID.
type Venue struct {
	Location     Location `json:"location"`
	Title        string   `json:"title"`
	Address      string   `json:"address"`
	FoursquareID string   `json:"foursquare_id,omitempty"`
}
