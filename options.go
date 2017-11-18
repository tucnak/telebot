package telebot

// Option is a shorcut flag type for certain message features
// (so-called options). It means that instead of passing
// fully-fledged SendOptions* to Send(), you can use these
// flags instead.
//
// Supported options are defined as iota-constants.
type Option int

const (
	// SendOptions.DisableWebPagePreview
	NoPreview Option = iota

	// SendOptions.DisableNotification
	Silent

	// ReplyMarkup.ForceReply
	ForceReply

	// ReplyMarkup.OneTimeKeyboard
	OneTimeKeyboard
)

// SendOptions has most complete control over in what way the message
// must be sent, providing an API-complete set of custom properties
// and options.
//
// Despite its power, SendOptions is rather inconvenient to use all
// the way through bot logic, so you might want to consider storing
// and re-using it somewhere or be using Option flags instead.
type SendOptions struct {
	// If the message is a reply, original message.
	ReplyTo *Message

	// See ReplyMarkup struct definition.
	ReplyMarkup *ReplyMarkup

	// For text messages, disables previews for links in this message.
	DisableWebPagePreview bool

	// Sends the message silently. iOS users will not receive a notification, Android users will receive a notification with no sound.
	DisableNotification bool

	// ParseMode controls how client apps render your message.
	ParseMode ParseMode
}

// ReplyMarkup controls two convenient options for bot-user communications
// such as reply keyboard and inline "keyboard" (a grid of buttons as a part
// of the message).
type ReplyMarkup struct {
	// ForceReply forces Telegram clients to display
	// a reply interface to the user (act as if the user
	// has selected the bot‘s message and tapped "Reply").
	ForceReply bool `json:"force_reply,omitempty"`

	// InlineKeyboard is a grid of InlineButtons displayed in the message.
	//
	// Note: DO NOT confuse with ReplyKeyboard and other keyboard properties!
	InlineKeyboard [][]InlineButton `json:"inline_keyboard,omitempty"`

	// ReplyKeyboard is a grid, consisting of keyboard buttons.
	//
	// Note: you don't need to set HideCustomKeyboard field to show custom keyboard.
	ReplyKeyboard [][]KeyboardButton `json:"keyboard,omitempty"`

	// Requests clients to resize the keyboard vertically for optimal fit
	// (e.g. make the keyboard smaller if there are just two rows of buttons).
	//
	// Defaults to false, in which case the custom keyboard is always of the
	// same height as the app's standard keyboard.
	ResizeReplyKeyboard bool `json:"resize_keyboard,omitempty"`

	// Requests clients to hide the reply keyboard as soon as it's been used.
	//
	// Defaults to false.
	OneTimeKeyboard bool `json:"one_time_keyboard,omitempty"`

	// Use this param if you want to force reply from
	// specific users only.
	//
	// Targets:
	// 1) Users that are @mentioned in the text of the Message object;
	// 2) If the bot's message is a reply (has SendOptions.ReplyTo),
	//       sender of the original message.
	Selective bool `json:"selective,omitempty"`
}

// KeyboardButton represents a button displayed in reply-keyboard.
type KeyboardButton struct {
	Text string `json:"text"`

	Contact  bool `json:"request_contact,omitempty"`
	Location bool `json:"request_location,omitempty"`
}

// InlineKeyboardMarkup represents an inline keyboard that appears
// right next to the message it belongs to.
type InlineKeyboardMarkup struct {
	// Array of button rows, each represented by
	// an Array of KeyboardButton objects.
	InlineKeyboard [][]InlineButton `json:"inline_keyboard,omitempty"`
}

// InlineButton represents a button displayed in the message.
type InlineButton struct {
	Text string `json:"text"`

	URL         string `json:"url,omitempty"`
	Data        string `json:"callback_data,omitempty"`
	InlineQuery string `json:"switch_inline_query,omitempty"`
}
