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
	ID      int64    `json:"update_id"`
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

	// FileSize (optional) of the audio file.
	FileSize int `json:"file_size"`

	// Title (optional) as defined by sender or by audio tags.
	Title string `json:"title"`

	// Performer (optional) is defined by sender or by audio tags.
	Performer string `json:"performer"`

	// MIME type (optional) of the file as defined by sender.
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
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
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

// CallbackResponse builds a response to an Callback query.
// See also: https://core.telegram.org/bots/api#answerCallbackQuery
type CallbackResponse struct {
	// The ID of the callback to which this is a response.
	// It is not necessary to specify this field manually.
	CallbackID string `json:"callback_query_id"`

	// Text of the notification. If not specified, nothing will be shown to the user.
	Text string `json:"text,omitempty"`

	// (Optional) If true, an alert will be shown by the client instead
	// of a notification at the top of the chat screen. Defaults to false.
	ShowAlert bool `json:"show_alert,omitempty"`

	// (Optional) URL that will be opened by the user's client.
	// If you have created a Game and accepted the conditions via @Botfather
	// specify the URL that opens your game
	// note that this will only work if the query comes from a callback_game button.
	// Otherwise, you may use links like telegram.me/your_bot?start=XXXX that open your bot with a parameter.
	URL string `json:"url,omitempty"`
}

// Venue object represents a venue location with name, address and optional foursquare id.
type Venue struct {
	Location      Location `json:"location"`
	Title         string   `json:"title"`
	Address       string   `json:"address"`
	Foursquare_id string   `json:"foursquare_id",omitempty`
}

// MessageEntity
// This object represents one special entity in a text message.
// For example, hashtags, usernames, URLs, etc
type MessageEntity struct {

	// type Type of the entity. Can be mention (@username), hashtag,
	// bot_command, url, email, bold (bold text), italic (italic text),
	// code (monowidth string), pre (monowidth block), text_link (for clickable text URLs),
	// text_mention (for users without usernames)
	Type string `json:"type"`

	// offset Offset in UTF-16 code units to the start of the entity
	Offset int `json:"offset"`

	//length Length of the entity in UTF-16 code units
	Length int `json:"length"`

	//url	Optional. For “text_link” only, url that will be opened after user taps on the text
	Url string `json:"url",omitempty`

	//user	Optional. For “text_mention” only, the mentioned user
	User User `json:"user",omitempty`
}

// ChatMember ,
// This struct contains information about one member of the chat.
type ChatMember struct {
	User   User   `json:"user"`
	Status string `json:"status"`
}

// UserProfilePhotos ,
// This struct represent a user's profile pictures.
//
// Count : Total number of profile pictures the target user has
//
// Photos : Array of Array of PhotoSize	, Requested profile pictures (in up to 4 sizes each)
type UserProfilePhotos struct {
	Count  int       `json:"total_count"`
	Photos [][]Photo `json:"photos"`
}
