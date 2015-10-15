package telebot

import (
	"time"
)

// Message object represents a message.
type Message struct {
	ID int `json:"message_id"`
	// For message sent to channels, Sender maybe empty
	Sender   User `json:"from"`
	Unixtime int  `json:"date"`

	// For forwarded messages, sender of the original message.
	OriginalSender User `json:"forward_from"`

	// For forwarded messages, unixtime of the original message.
	OriginalUnixtime int `json:"forward_date"`

	// For replies, ReplyTo represents the original message.
	// Note that the Message object in this field will not
	// contain further ReplyTo fields even if it
	// itself is a reply.
	ReplyTo *Message `json:"reply_to_message"`

	// For a text message, the actual UTF-8 text of the message
	Text string `json:"text"`

	// For an audio recording, information about it.
	Audio Audio `json:"audio"`

	// For a general file, information about it.
	Document Document `json:"document"`

	// For a photo, available thumbnails.
	Photo []Thumbnail `json:"photo"`

	// For a sticker, information about it.
	Sticker Sticker `json:"sticker"`

	// For a video, information about it.
	Video Video `json:"video"`

	// For a contact, contact information itself.
	Contact Contact `json:"contact"`

	// For a location, its longitude and latitude.
	Location Location `json:"location"`

	// A group chat message belongs to, empty if personal.
	Chat Chat `json:"chat"`

	// For a service message, represents a user,
	// that just got added to chat, this message came from.
	//
	// Sender leads to User, capable of invite.
	//
	// UserJoined might be the Bot itself.
	UserJoined User `json:"new_chat_participant"`

	// For a service message, represents a user,
	// that just left chat, this message came from.
	//
	// If user was kicked, Sender leads to a User,
	// capable of this kick.
	//
	// UserLeft might be the Bot itself.
	UserLeft User `json:"left_chat_participant"`

	// For a service message, represents a new title
	// for chat this message came from.
	//
	// Sender would lead to a User, capable of change.
	NewChatTitle string `json:"new_chat_title"`

	// For a service message, represents all available
	// thumbnails of new chat photo.
	//
	// Sender would lead to a User, capable of change.
	NewChatPhoto []Thumbnail `json:"new_chat_photo"`

	// For a service message, true if chat photo just
	// got removed.
	//
	// Sender would lead to a User, capable of change.
	ChatPhotoDeleted bool `json:"delete_chat_photo"`

	// For a service message, true if group has been created.
	//
	// You would recieve such a message if you are one of
	// initial group chat members.
	//
	// Sender would lead to creator of the chat.
	ChatCreated bool `json:"group_chat_created"`
}

// Origin returns an origin of message: group chat / personal.
func (m Message) Origin() User {
	// if m.IsPersonal() {
	// 	return m.Chat
	// }

	return m.Sender
}

// Time returns the moment of message creation in local time.
func (m Message) Time() time.Time {
	return time.Unix(int64(m.Unixtime), 0)
}

// IsForwarded says whether message is forwarded copy of another
// message or not.
func (m Message) IsForwarded() bool {
	return m.OriginalSender != User{}
}

// IsReply says whether message is reply to another message or not.
func (m Message) IsReply() bool {
	return m.ReplyTo != nil
}

// IsPersonal returns true, if message is a personal message,
// returns false if sent to group chat.
func (m Message) IsPersonal() bool {
	return !m.Chat.IsGroupChat()
}

// IsService returns true, if message is a service message,
// returns false otherwise.
//
// Service messages are automatically sent messages, which
// typically occur on some global action. For instance, when
// anyone leaves the chat or chat title changes.
func (m Message) IsService() bool {
	service := false

	if (m.UserJoined != User{}) {
		service = true
	}

	if (m.UserLeft != User{}) {
		service = true
	}

	if m.NewChatTitle != "" {
		service = true
	}

	if len(m.NewChatPhoto) > 0 {
		service = true
	}

	if m.ChatPhotoDeleted {
		service = true
	}

	if m.ChatCreated {
		service = true
	}

	return service
}
