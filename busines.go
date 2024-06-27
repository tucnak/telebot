package telebot

import "time"

type BusinessConnection struct {
	// Unique identifier of the business connection
	ID string `json:"id"`

	// Business account user that created the business connection
	User *User `json:"user"`

	// Identifier of a private chat with the user who created the business connection. This
	// number may have more than 32 significant bits and some programming languages may
	// have difficulty/silent defects in interpreting it. But it has at most 52 significant bits,
	// so a 64-bit integer or double-precision float type are safe for storing this identifier.
	UserChatID int64 `json:"user_chat_id"`

	// Unixtime, use BusinessConnection.Time() to get time.Time.
	Unixtime int64 `json:"date"`

	// True, if the bot can act on behalf of the business account in chats that were active in the last 24 hours
	CanReply bool `json:"can_reply"`

	// True, if the connection is active
	IsEnabled bool `json:"is_enabled"`
}

// Time returns the moment of business connection creation in local time.
func (b *BusinessConnection) Time() time.Time {
	return time.Unix(b.Unixtime, 0)
}

type BusinessMessagesDeleted struct {
	// Unique identifier of the business connection
	BusinessConnectionID string `json:"business_connection_id"`

	// Information about a chat in the business account. The bot
	// may not have access to the chat or the corresponding user.
	Chat *Chat `json:"chat"`

	// The list of identifiers of deleted messages in the chat of the business account
	MessageIDs []int `json:"message_ids"`
}
