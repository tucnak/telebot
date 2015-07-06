package telebot

// SendOptions represents a set of custom options that could
// be appled to messages sent.
type SendOptions struct {
	// If the message is a reply, original message.
	ReplyTo Message

	// See ForceReply struct definition.
	ForceReply ForceReply

	// For text messages, disables previews for links in this message.
	DisableWebPagePreview bool
}

// ForceReply forces Telegram clients to display
// a reply interface to the user (act as if the user
// has selected the bot‘s message and tapped "Reply").
type ForceReply struct {
	// Enable if intended.
	Require bool `json:"force_reply"`

	// Use this param if you want to force reply from
	// specific users only.
	//
	// Targets:
	// 1) Users that are @mentioned in the text of the Message object;
	// 2) If the bot's message is a reply (has SendOptions.ReplyTo),
	//       sender of the original message.
	Selective bool `json:"selective"`
}
