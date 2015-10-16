package telebot

// User object represents a Telegram user, bot
//
// object represents a group chat if Title is empty.
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`

	LastName string `json:"last_name"`
	Username string `json:"username"`
}

// Chat object represents a Telegram user, bot or group chat.
// Title for channels and group chats
// Type of chat, can be either “private”, or “group”, or “channel”
type Chat struct {
	ID   int    `json:"id"`
	Type string `json:"type"`

	Title     string `json:"title"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// IsGroupChat returns true if chat object represents a group chat.
func (u Chat) IsGroupChat() bool {
	return u.Type == "group"
}

// Update object represents an incoming update.
type Update struct {
	ID      int     `json:"update_id"`
	Payload Message `json:"message"`
}

// Thumbnail object represents a image/sticker of particular size.
type Thumbnail struct {
	File

	Width  int `json:"width"`
	Height int `json:"height"`
}

// Photo object represents a photo with caption.
type Photo struct {
	Thumbnail

	Caption string
}

// Audio object represents an audio file (voice note).
type Audio struct {
	File

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration"`

	// MIME type of the file as defined by sender.
	Mime string `json:"mime_type"`
}

// Document object represents a general file (as opposed to Photo or Audio).
// Telegram users can send files of any type of up to 1.5 GB in size.
type Document struct {
	File

	// Document thumbnail as defined by sender.
	Preview Thumbnail `json:"thumb"`

	// Original filename as defined by sender.
	FileName string `json:"file_name"`

	// MIME type of the file as defined by sender.
	Mime string `json:"mime_type"`
}

// Sticker object represents a WebP image, so-called sticker.
type Sticker struct {
	File

	Width  int `json:"width"`
	Height int `json:"height"`

	// Sticker thumbnail in .webp or .jpg format.
	Preview Thumbnail `json:"thumb"`
}

// Video object represents an MP4-encoded video.
type Video struct {
	Audio

	Width  int `json:"width"`
	Height int `json:"height"`

	// Text description of the video as defined by sender (usually empty).
	Caption string `json:"caption"`

	// Video thumbnail.
	Preview Thumbnail `json:"thumb"`
}

// Contact object represents a contact to Telegram user
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`

	// Contact's username in Telegram (might be empty).
	Username string `json:"user_id"`
}

// Location object represents geographic position.
type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}
