package telebot

// Callback object represents a query from a callback button in an
// inline keyboard.
type Callback struct {
	ID string `json:"id"`

	// For message sent to channels, Sender may be empty
	Sender *User `json:"from"`

	// Message will be set if the button that originated the query
	// was attached to a message sent by a bot.
	Message *Message `json:"message"`

	// MessageID will be set if the button was attached to a message
	// sent via the bot in inline mode.
	MessageID string `json:"inline_message_id"`

	// Data associated with the callback button. Be aware that
	// a bad client can send arbitrary data in this field.
	Data string `json:"data"`
}

// CallbackResponse builds a response to a Callback query.
//
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
