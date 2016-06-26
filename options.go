package telebot

// ParseMode determines the way client applications treat the text of the message
type ParseMode string

// Supported ParseMode
const (
	ModeDefault  ParseMode = ""
	ModeMarkdown ParseMode = "Markdown"
	ModeHTML     ParseMode = "HTML"
)

// SendOptions represents a set of custom options that could
// be appled to messages sent.
type SendOptions struct {
	// If the message is a reply, original message.
	ReplyTo Message

	// See ReplyMarkup struct definition.
	ReplyMarkup ReplyMarkup

	// For text messages, disables previews for links in this message.
	DisableWebPagePreview bool

	// Sends the message silently. iOS users will not receive a notification, Android users will receive a notification with no sound.
	DisableNotification bool

	// ParseMode controls how client apps render your message.
	ParseMode ParseMode
}

// AnswerOptions specifies options about responses to inline results
type AnswerOptions struct {
	// Cache time (default: 300). Set to -1 to set to 0 (no cache)
	CacheTime int

	// Cache personality (cache result for this user or for all)
	// Default: false
	IsPersonal bool

	// Offset to send back if the client wants more results.
	Offset string
}

// ReplyMarkup specifies convenient options for bot-user communications.
type ReplyMarkup struct {
	// ForceReply forces Telegram clients to display
	// a reply interface to the user (act as if the user
	// has selected the bot‘s message and tapped "Reply").
	ForceReply bool `json:"force_reply,omitempty"`

	// CustomKeyboard is Array of button rows, each represented by an Array of Strings.
	//
	// Note: you don't need to set HideCustomKeyboard field to show custom keyboard.
	CustomKeyboard [][]string `json:"keyboard,omitempty"`

	InlineKeyboard [][]KeyboardButton `json:"inline_keyboard,omitempty"`

	// Requests clients to resize the keyboard vertically for optimal fit
	// (e.g., make the keyboard smaller if there are just two rows of buttons).
	// Defaults to false, in which case the custom keyboard is always of the
	// same height as the app's standard keyboard.
	ResizeKeyboard bool `json:"resize_keyboard,omitempty"`
	// Requests clients to hide the keyboard as soon as it's been used. Defaults to false.
	OneTimeKeyboard bool `json:"one_time_keyboard,omitempty"`

	// Requests clients to hide the custom keyboard.
	//
	// Note: You dont need to set CustomKeyboard field to hide custom keyboard.
	HideCustomKeyboard bool `json:"hide_keyboard,omitempty"`

	// Use this param if you want to force reply from
	// specific users only.
	//
	// Targets:
	// 1) Users that are @mentioned in the text of the Message object;
	// 2) If the bot's message is a reply (has SendOptions.ReplyTo),
	//       sender of the original message.
	Selective bool `json:"selective,omitempty"`
}
