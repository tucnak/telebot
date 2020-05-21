package telebot

import (
	"strconv"
	"time"
)

// Message object represents a message.
type Message struct {
	ID int `json:"message_id"`

	InlineID string `json:"-"`

	// For message sent to channels, Sender will be nil
	Sender *User `json:"from"`

	// Unixtime, use Message.Time() to get time.Time
	Unixtime int64 `json:"date"`

	// Conversation the message belongs to.
	Chat *Chat `json:"chat"`

	// For forwarded messages, sender of the original message.
	OriginalSender *User `json:"forward_from"`

	// For forwarded messages, chat of the original message when
	// forwarded from a channel.
	OriginalChat *Chat `json:"forward_from_chat"`

	// For forwarded messages, unixtime of the original message.
	OriginalUnixtime int `json:"forward_date"`

	// For replies, ReplyTo represents the original message.
	//
	// Note that the Message object in this field will not
	// contain further ReplyTo fields even if it
	// itself is a reply.
	ReplyTo *Message `json:"reply_to_message"`

	// (Optional) Time of last edit in Unix
	LastEdit int64 `json:"edit_date"`

	// AlbumID is the unique identifier of a media message group
	// this message belongs to.
	AlbumID string `json:"media_group_id"`

	// Author signature (in channels).
	Signature string `json:"author_signature"`

	// For a text message, the actual UTF-8 text of the message.
	Text string `json:"text"`

	// For registered commands, will contain the string payload:
	//
	// Ex: `/command <payload>` or `/command@botname <payload>`
	Payload string `json:"-"`

	// For text messages, special entities like usernames, URLs, bot commands,
	// etc. that appear in the text.
	Entities []MessageEntity `json:"entities,omitempty"`

	// Some messages containing media, may as well have a caption.
	Caption string `json:"caption,omitempty"`

	// For messages with a caption, special entities like usernames, URLs,
	// bot commands, etc. that appear in the caption.
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	// For an audio recording, information about it.
	Audio *Audio `json:"audio"`

	// For a gneral file, information about it.
	Document *Document `json:"document"`

	// For a photo, all available sizes (thumbnails).
	Photo *Photo `json:"photo"`

	// For a sticker, information about it.
	Sticker *Sticker `json:"sticker"`

	// For a voice message, information about it.
	Voice *Voice `json:"voice"`

	// For a video note, information about it.
	VideoNote *VideoNote `json:"video_note"`

	// For a video, information about it.
	Video *Video `json:"video"`

	// For a animation, information about it.
	Animation *Animation `json:"animation"`

	// For a contact, contact information itself.
	Contact *Contact `json:"contact"`

	// For a location, its longitude and latitude.
	Location *Location `json:"location"`

	// For a venue, information about it.
	Venue *Venue `json:"venue"`

	// For a poll, information the native poll.
	Poll *Poll `json:"poll"`

	// For a dice, information about it.
	Dice *Dice `json:"dice"`

	// For a service message, represents a user,
	// that just got added to chat, this message came from.
	//
	// Sender leads to User, capable of invite.
	//
	// UserJoined might be the Bot itself.
	UserJoined *User `json:"new_chat_member"`

	// For a service message, represents a user,
	// that just left chat, this message came from.
	//
	// If user was kicked, Sender leads to a User,
	// capable of this kick.
	//
	// UserLeft might be the Bot itself.
	UserLeft *User `json:"left_chat_member"`

	// For a service message, represents a new title
	// for chat this message came from.
	//
	// Sender would lead to a User, capable of change.
	NewGroupTitle string `json:"new_chat_title"`

	// For a service message, represents all available
	// thumbnails of the new chat photo.
	//
	// Sender would lead to a User, capable of change.
	NewGroupPhoto *Photo `json:"new_chat_photo"`

	// For a service message, new members that were added to
	// the group or supergroup and information about them
	// (the bot itself may be one of these members).
	UsersJoined []User `json:"new_chat_members"`

	// For a service message, true if chat photo just
	// got removed.
	//
	// Sender would lead to a User, capable of change.
	GroupPhotoDeleted bool `json:"delete_chat_photo"`

	// For a service message, true if group has been created.
	//
	// You would receive such a message if you are one of
	// initial group chat members.
	//
	// Sender would lead to creator of the chat.
	GroupCreated bool `json:"group_chat_created"`

	// For a service message, true if super group has been created.
	//
	// You would receive such a message if you are one of
	// initial group chat members.
	//
	// Sender would lead to creator of the chat.
	SuperGroupCreated bool `json:"supergroup_chat_created"`

	// For a service message, true if channel has been created.
	//
	// You would receive such a message if you are one of
	// initial channel administrators.
	//
	// Sender would lead to creator of the chat.
	ChannelCreated bool `json:"channel_chat_created"`

	// For a service message, the destination (super group) you
	// migrated to.
	//
	// You would receive such a message when your chat has migrated
	// to a super group.
	//
	// Sender would lead to creator of the migration.
	MigrateTo int64 `json:"migrate_to_chat_id"`

	// For a service message, the Origin (normal group) you migrated
	// from.
	//
	// You would receive such a message when your chat has migrated
	// to a super group.
	//
	// Sender would lead to creator of the migration.
	MigrateFrom int64 `json:"migrate_from_chat_id"`

	// Specified message was pinned. Note that the Message object
	// in this field will not contain further ReplyTo fields even
	// if it is itself a reply.
	PinnedMessage *Message `json:"pinned_message"`

	// Message is an invoice for a payment.
	Invoice *Invoice `json:"invoice"`

	// Message is a service message about a successful payment.
	Payment *Payment `json:"successful_payment"`

	// The domain name of the website on which the user has logged in.
	ConnectedWebsite string `json:"connected_website,omitempty"`

	// Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
}

// MessageEntity object represents "special" parts of text messages,
// including hashtags, usernames, URLs, etc.
type MessageEntity struct {
	// Specifies entity type.
	Type EntityType `json:"type"`

	// Offset in UTF-16 code units to the start of the entity.
	Offset int `json:"offset"`

	// Length of the entity in UTF-16 code units.
	Length int `json:"length"`

	// (Optional) For EntityTextLink entity type only.
	//
	// URL will be opened after user taps on the text.
	URL string `json:"url,omitempty"`

	// (Optional) For EntityTMention entity type only.
	User *User `json:"user,omitempty"`

	// (Optional) For EntityCodeBlock entity type only.
	Language string `json:"language,omitempty"`
}

// MessageSig satisfies Editable interface (see Editable.)
func (m *Message) MessageSig() (string, int64) {
	if m.InlineID != "" {
		return m.InlineID, 0
	}
	return strconv.Itoa(m.ID), m.Chat.ID
}

// Time returns the moment of message creation in local time.
func (m *Message) Time() time.Time {
	return time.Unix(m.Unixtime, 0)
}

// LastEdited returns time.Time of last edit.
func (m *Message) LastEdited() time.Time {
	return time.Unix(m.LastEdit, 0)
}

// IsForwarded says whether message is forwarded copy of another
// message or not.
func (m *Message) IsForwarded() bool {
	return m.OriginalSender != nil || m.OriginalChat != nil
}

// IsReply says whether message is a reply to another message.
func (m *Message) IsReply() bool {
	return m.ReplyTo != nil
}

// Private returns true, if it's a personal message.
func (m *Message) Private() bool {
	return m.Chat.Type == ChatPrivate
}

// FromGroup returns true, if message came from a group OR
// a super group.
func (m *Message) FromGroup() bool {
	return m.Chat.Type == ChatGroup || m.Chat.Type == ChatSuperGroup
}

// FromChannel returns true, if message came from a channel.
func (m *Message) FromChannel() bool {
	return m.Chat.Type == ChatChannel
}

// IsService returns true, if message is a service message,
// returns false otherwise.
//
// Service messages are automatically sent messages, which
// typically occur on some global action. For instance, when
// anyone leaves the chat or chat title changes.
func (m *Message) IsService() bool {
	fact := false

	fact = fact || m.UserJoined != nil
	fact = fact || len(m.UsersJoined) > 0
	fact = fact || m.UserLeft != nil
	fact = fact || m.NewGroupTitle != ""
	fact = fact || m.NewGroupPhoto != nil
	fact = fact || m.GroupPhotoDeleted
	fact = fact || m.GroupCreated || m.SuperGroupCreated
	fact = fact || (m.MigrateTo != m.MigrateFrom)

	return fact
}
