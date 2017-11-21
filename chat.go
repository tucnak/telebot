package telebot

import "strconv"

// User object represents a Telegram user, bot
type User struct {
	ID int `json:"id"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Recipient returns user ID (see Recipient interface).
func (u *User) Recipient() string {
	return strconv.Itoa(u.ID)
}

// Chat object represents a Telegram user, bot, group or a channel.
type Chat struct {
	ID int64 `json:"id"`

	// See telebot.ChatType and consts.
	Type ChatType `json:"type"`

	// Won't be there for ChatPrivate.
	Title string `json:"title"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Recipient returns chat ID (see Recipient interface).
func (c *Chat) Recipient() string {
	if c.Type == ChatChannel {
		return "@" + c.Username
	}
	return strconv.FormatInt(c.ID, 10)
}

// ChatMember object represents information about a single chat member.
type ChatMember struct {
	Rights

	User *User        `json:"user"`
	Role MemberStatus `json:"status"`

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
