package telebot

import (
	"encoding/json"
	"fmt"
)

// Option is a shorcut flag type for certain message features
// (so-called options). It means that instead of passing
// fully-fledged SendOptions* to Send(), you can use these
// flags instead.
//
// Supported options are defined as iota-constants.
type Option int

const (
	// NoPreview = SendOptions.DisableWebPagePreview
	NoPreview Option = iota

	// Silent = SendOptions.DisableNotification
	Silent

	// ForceReply = ReplyMarkup.ForceReply
	ForceReply

	// OneTimeKeyboard = ReplyMarkup.OneTimeKeyboard
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

func (og *SendOptions) copy() *SendOptions {
	cp := *og
	if cp.ReplyMarkup != nil {
		cp.ReplyMarkup = cp.ReplyMarkup.copy()
	}

	return &cp
}

// ReplyMarkup controls two convenient options for bot-user communications
// such as reply keyboard and inline "keyboard" (a grid of buttons as a part
// of the message).
type ReplyMarkup struct {
	// InlineKeyboard is a grid of InlineButtons displayed in the message.
	//
	// Note: DO NOT confuse with ReplyKeyboard and other keyboard properties!
	InlineKeyboard [][]InlineButton `json:"inline_keyboard,omitempty"`

	// ReplyKeyboard is a grid, consisting of keyboard buttons.
	//
	// Note: you don't need to set HideCustomKeyboard field to show custom keyboard.
	ReplyKeyboard [][]ReplyButton `json:"keyboard,omitempty"`

	// ForceReply forces Telegram clients to display
	// a reply interface to the user (act as if the user
	// has selected the bot‘s message and tapped "Reply").
	ForceReply bool `json:"force_reply,omitempty"`

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

	// Requests clients to remove the reply keyboard.
	//
	// Dafaults to false.
	ReplyKeyboardRemove bool `json:"remove_keyboard,omitempty"`

	// Use this param if you want to force reply from
	// specific users only.
	//
	// Targets:
	// 1) Users that are @mentioned in the text of the Message object;
	// 2) If the bot's message is a reply (has SendOptions.ReplyTo),
	//       sender of the original message.
	Selective bool `json:"selective,omitempty"`
}

func (og *ReplyMarkup) copy() *ReplyMarkup {
	cp := *og

	cp.ReplyKeyboard = make([][]ReplyButton, len(og.ReplyKeyboard))
	for i, row := range og.ReplyKeyboard {
		cp.ReplyKeyboard[i] = make([]ReplyButton, len(row))
		copy(cp.ReplyKeyboard[i], row)
	}

	cp.InlineKeyboard = make([][]InlineButton, len(og.InlineKeyboard))
	for i, row := range og.InlineKeyboard {
		cp.InlineKeyboard[i] = make([]InlineButton, len(row))
		copy(cp.InlineKeyboard[i], row)
	}

	return &cp
}

// ReplyButton represents a button displayed in reply-keyboard.
//
// Set either Contact or Location to true in order to request
// sensitive info, such as user's phone number or current location.
// (Available in private chats only.)
type ReplyButton struct {
	Text string `json:"text"`

	Contact  bool     `json:"request_contact,omitempty"`
	Location bool     `json:"request_location,omitempty"`
	Poll     PollType `json:"request_poll,omitempty"`

	// Not used anywhere.
	// Will be removed in future releases.
	Action func(*Message) `json:"-"`
}

// InlineKeyboardMarkup represents an inline keyboard that appears
// right next to the message it belongs to.
type InlineKeyboardMarkup struct {
	// Array of button rows, each represented by
	// an Array of KeyboardButton objects.
	InlineKeyboard [][]InlineButton `json:"inline_keyboard,omitempty"`
}

// MarshalJSON implements json.Marshaler. It allows to pass
// PollType as keyboard's poll type instead of KeyboardButtonPollType object.
func (pt PollType) MarshalJSON() ([]byte, error) {
	var aux = struct {
		Type string `json:"type"`
	}{
		Type: string(pt),
	}
	return json.Marshal(&aux)
}

type row []Btn
func (r *ReplyMarkup) Row(many ...Btn) row {
	return many
}

func (r *ReplyMarkup) Inline(rows ...row) {
	inlineKeys := make([][]InlineButton, 0, len(rows))
	for i, row := range rows {
		keys := make([]InlineButton, 0, len(row))
		for j, btn := range row {
			btn := btn.Inline()
			if btn == nil {
				panic(fmt.Sprintf(
					"telebot: button row %d column %d is not an inline button",
					i, j))
			}
			keys = append(keys, *btn)
		}
		inlineKeys = append(inlineKeys, keys)
	}

	r.InlineKeyboard = inlineKeys
}

func (r *ReplyMarkup) Reply(rows ...row) {
	replyKeys := make([][]ReplyButton, 0, len(rows))
	for i, row := range rows {
		keys := make([]ReplyButton, 0, len(row))
		for j, btn := range row {
			btn := btn.Reply()
			if btn == nil {
				panic(fmt.Sprintf(
					"telebot: button row %d column %d is not a reply button",
					i, j))
			}
			keys = append(keys, *btn)
		}
		replyKeys = append(replyKeys, keys)
	}

	r.ReplyKeyboard = replyKeys
}
func(r *ReplyMarkup) Text(unique,text  string) Btn {
	return Btn{Unique: unique, Text: text}
}

func(r *ReplyMarkup) URL(unique,url string) Btn {
	return Btn{Unique: unique, URL: url}
}


func(r *ReplyMarkup) Query(unique string, query string) Btn {
	return Btn{Unique: unique, InlineQuery: query}
}

func(r *ReplyMarkup) QueryChat(unique string, query string) Btn {
	return Btn{Unique: unique, InlineQueryChat: query}
}

func(r *ReplyMarkup) Login(unique,text string,login *Login) Btn {
	return Btn{Unique: unique, Login: login, Text: text}
}

func(r *ReplyMarkup) Contact(text string) Btn {
	return Btn{Contact:true, Text: text}
}

func(r *ReplyMarkup) Location(text string) Btn {
	return Btn{Location:true, Text: text}
}

func(r *ReplyMarkup) Poll(poll PollType) Btn {
	return Btn{Poll: poll}
}

// Btn is a constructor button, which will later become either a reply, or an inline button.
type Btn struct {
	Unique          string
	Text            string
	URL             string
	Data            string
	InlineQuery     string
	InlineQueryChat string
	Contact         bool
	Location        bool
	Poll            PollType
	Login           *Login
}

func (b Btn) Inline() *InlineButton {
	if b.Unique == "" {
		return nil
	}

	return &InlineButton{
		Unique:          b.Unique,
		Text:            b.Text,
		URL:             b.URL,
		Data:            b.Data,
		InlineQuery:     b.InlineQuery,
		InlineQueryChat: b.InlineQueryChat,
		Login:           nil,
	}
}

func (b Btn) Reply() *ReplyButton {
	if b.Unique != "" {
		return nil
	}

	return &ReplyButton{
		Text:     b.Text,
		Contact:  b.Contact,
		Location: b.Location,
		Poll:     b.Poll,
	}
}
