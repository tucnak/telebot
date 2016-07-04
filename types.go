package telebot

import "strconv"

// Recipient is basically any possible endpoint you can send
// messages to. It's usually a distinct user or a chat.
type Recipient interface {
	// ID of user or group chat, @Username for channel
	Destination() string
}

// User object represents a Telegram user, bot
//
// object represents a group chat if Title is empty.
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`

	LastName string `json:"last_name"`
	Username string `json:"username"`
}

// Destination is internal user ID.
func (u User) Destination() string {
	return strconv.Itoa(u.ID)
}

// Chat object represents a Telegram user, bot or group chat.
// Title for channels and group chats
// Type of chat, can be either “private”, “group”, "supergroup" or “channel”
type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`

	Title     string `json:"title"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Destination is internal chat ID.
func (c Chat) Destination() string {
	ret := "@" + c.Username
	if c.Type != "channel" {
		ret = strconv.FormatInt(c.ID, 10)
	}
	return ret
}

// IsGroupChat returns true if chat object represents a group chat.
func (c Chat) IsGroupChat() bool {
	return c.Type != "private"
}

// Update object represents an incoming update.
type Update struct {
	ID      int      `json:"update_id"`
	Payload *Message `json:"message"`

	// optional
	Callback *Callback `json:"callback_query"`
	Query    *Query    `json:"inline_query"`
}

// Thumbnail object represents an image/sticker of a particular size.
type Thumbnail struct {
	File

	Width  int `json:"width"`
	Height int `json:"height"`
}

// Photo object represents a photo with caption.
type Photo struct {
	File

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

// KeyboardButton represents a button displayed on in a message.
type KeyboardButton struct {
	Text        string `json:"text"`
	URL         string `json:"url,omitempty"`
	Data        string `json:"callback_data,omitempty"`
	InlineQuery string `json:"switch_inline_query,omitempty"`
}

// InlineKeyboardMarkup represents an inline keyboard that appears right next
// to the message it belongs to.
type InlineKeyboardMarkup struct {
	// Array of button rows, each represented by an Array of KeyboardButton objects.
	InlineKeyboard [][]KeyboardButton `json:"inline_keyboard,omitempty"`
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
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

// Callback object represents a query from a callback button in an
// inline keyboard.
type Callback struct {
	ID string `json:"id"`

	// For message sent to channels, Sender may be empty
	Sender User `json:"from"`

	// Message will be set if the button that originated the query
	// was attached to a message sent by a bot.
	Message Message `json:"message"`

	// MessageID will be set if the button was attached to a message
	// sent via the bot in inline mode.
	MessageID string `json:"inline_message_id"`
	Data      string `json:"data"`
}

// Venue object represents a venue location with name, address and optional foursquare id.
type Venue struct {
	Location Location		`json:"location"`
	Title string			`json:"title"`
	Address string			`json:"address"`
	Foursquare_id string	`json:"foursquare_id",omitempty`
}
