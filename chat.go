package telebot

import "strconv"

// User object represents a Telegram user, bot
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`

	LastName string `json:"last_name"`
	Username string `json:"username"`
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
	User   *User        `json:"user"`
	Status MemberStatus `json:"status"`

	// Due for banned/restricted, Unixtime.
	Until int64 `json:"until_date,omitempty"`

	CanBeEdited        bool `json:"can_be_edited,omitempty"`
	CanChangeInfo      bool `json:"can_change_info,omitempty"`
	CanPostMessages    bool `json:"can_post_messages,omitempty"`
	CanEditMessages    bool `json:"can_edit_messages,omitempty"`
	CanDeleteMessages  bool `json:"can_delete_messages,omitempty"`
	CanInviteUsers     bool `json:"can_invite_users,omitempty"`
	CanRestrictMembers bool `json:"can_restrict_members,omitempty"`
	CanPinMessages     bool `json:"can_pin_messages,omitempty"`
	CanPromoteMembers  bool `json:"can_promote_members,omitempty"`
	CanSendMessages    bool `json:"can_send_messages,omitempty"`
	CanSendMedia       bool `json:"can_send_media_messages,omitempty"`
	CanSendOther       bool `json:"can_send_other_messages,omitempty"`
	CanAddPreviews     bool `json:"can_add_web_page_previews,omitempty"`
}
