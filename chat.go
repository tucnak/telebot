package telebot

import "strconv"

// User object represents a Telegram user, bot.
type User struct {
	ID int `json:"id"`

	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	IsBot        bool   `json:"is_bot"`

	// Returns only in getMe
	CanJoinGroups   bool `json:"can_join_groups"`
	CanReadMessages bool `json:"can_read_all_group_messages"`
	SupportsInline  bool `json:"supports_inline_queries"`
}

// Recipient returns user ID (see Recipient interface).
func (u *User) Recipient() string {
	return strconv.Itoa(u.ID)
}

// Chat object represents a Telegram user, bot, group or a channel.
type Chat struct {
	ID int64 `json:"id"`

	// See ChatType and consts.
	Type ChatType `json:"type"`

	// Won't be there for ChatPrivate.
	Title string `json:"title"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`

	// Returns only in getChat
	Photo            *ChatPhoto `json:"photo,omitempty"`
	Description      string     `json:"description,omitempty"`
	InviteLink       string     `json:"invite_link,omitempty"`
	PinnedMessage    *Message   `json:"pinned_message,omitempty"`
	Permissions      *Rights    `json:"permissions,omitempty"`
	SlowMode         int        `json:"slow_mode_delay,omitempty"`
	StickerSet       string     `json:"sticker_set_name,omitempty"`
	CanSetStickerSet bool       `json:"can_set_sticker_set,omitempty"`
}

// ChatPhoto object represents a chat photo.
type ChatPhoto struct {
	// File identifiers of small (160x160) chat photo
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`

	// File identifiers of big (640x640) chat photo
	BigFileID       string `json:"big_file_id"`
	BigFileUniqueID string `json:"big_file_unique_id"`
}

// Recipient returns chat ID (see Recipient interface).
func (c *Chat) Recipient() string {
	return strconv.FormatInt(c.ID, 10)
}

// ChatMember object represents information about a single chat member.
type ChatMember struct {
	Rights

	User  *User        `json:"user"`
	Role  MemberStatus `json:"status"`
	Title string       `json:"custom_title"`

	// Date when restrictions will be lifted for the user, unix time.
	//
	// If user is restricted for more than 366 days or less than
	// 30 seconds from the current time, they are considered to be
	// restricted forever.
	//
	// Use tb.Forever().
	//
	RestrictedUntil int64 `json:"until_date,omitempty"`
}
