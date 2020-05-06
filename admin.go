package telebot

import (
	"encoding/json"
	"strconv"
	"time"
)

// Rights is a list of privileges available to chat members.
type Rights struct {
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
	CanSendPolls       bool `json:"can_send_polls,omitempty"`
	CanSendOther       bool `json:"can_send_other_messages,omitempty"`
	CanAddPreviews     bool `json:"can_add_web_page_previews,omitempty"`
}

// NoRights is the default Rights{}.
func NoRights() Rights { return Rights{} }

// NoRestrictions should be used when un-restricting or
// un-promoting user.
//
//	   member.Rights = tb.NoRestrictions()
//     bot.Restrict(chat, member)
//
func NoRestrictions() Rights {
	return Rights{
		CanBeEdited:        true,
		CanChangeInfo:      false,
		CanPostMessages:    false,
		CanEditMessages:    false,
		CanDeleteMessages:  false,
		CanInviteUsers:     false,
		CanRestrictMembers: false,
		CanPinMessages:     false,
		CanPromoteMembers:  false,
		CanSendMessages:    true,
		CanSendMedia:       true,
		CanSendPolls:       true,
		CanSendOther:       true,
		CanAddPreviews:     true,
	}
}

// AdminRights could be used to promote user to admin.
func AdminRights() Rights {
	return Rights{
		CanBeEdited:        true,
		CanChangeInfo:      true,
		CanPostMessages:    true,
		CanEditMessages:    true,
		CanDeleteMessages:  true,
		CanInviteUsers:     true,
		CanRestrictMembers: true,
		CanPinMessages:     true,
		CanPromoteMembers:  true,
		CanSendMessages:    true,
		CanSendMedia:       true,
		CanSendPolls:       true,
		CanSendOther:       true,
		CanAddPreviews:     true,
	}
}

// Forever is a Unixtime of "forever" banning.
func Forever() int64 {
	return time.Now().Add(367 * 24 * time.Hour).Unix()
}

// Ban will ban user from chat until `member.RestrictedUntil`.
func (b *Bot) Ban(chat *Chat, member *ChatMember) error {
	params := map[string]string{
		"chat_id":    chat.Recipient(),
		"user_id":    member.User.Recipient(),
		"until_date": strconv.FormatInt(member.RestrictedUntil, 10),
	}

	_, err := b.Raw("kickChatMember", params)
	return err
}

// Unban will unban user from chat, who would have thought eh?
func (b *Bot) Unban(chat *Chat, user *User) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"user_id": user.Recipient(),
	}

	_, err := b.Raw("unbanChatMember", params)
	return err
}

// Restrict lets you restrict a subset of member's rights until
// member.RestrictedUntil, such as:
//
//     * can send messages
//     * can send media
//     * can send other
//     * can add web page previews
//
func (b *Bot) Restrict(chat *Chat, member *ChatMember) error {
	prv, until := member.Rights, member.RestrictedUntil

	params := map[string]string{
		"chat_id":    chat.Recipient(),
		"user_id":    member.User.Recipient(),
		"until_date": strconv.FormatInt(until, 10),
	}
	embedRights(params, prv)

	_, err := b.Raw("restrictChatMember", params)
	return err
}

// Promote lets you update member's admin rights, such as:
//
//     * can change info
//     * can post messages
//     * can edit messages
//     * can delete messages
//     * can invite users
//     * can restrict members
//     * can pin messages
//     * can promote members
//
func (b *Bot) Promote(chat *Chat, member *ChatMember) error {
	prv := member.Rights

	params := map[string]string{
		"chat_id": chat.Recipient(),
		"user_id": member.User.Recipient(),
	}
	embedRights(params, prv)

	_, err := b.Raw("promoteChatMember", params)
	return err
}

// AdminsOf returns a member list of chat admins.
//
// On success, returns an Array of ChatMember objects that
// contains information about all chat administrators except other bots.
// If the chat is a group or a supergroup and
// no administrators were appointed, only the creator will be returned.
func (b *Bot) AdminsOf(chat *Chat) ([]ChatMember, error) {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	data, err := b.Raw("getChatAdministrators", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result []ChatMember
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// Len returns the number of members in a chat.
func (b *Bot) Len(chat *Chat) (int, error) {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	data, err := b.Raw("getChatMembersCount", params)
	if err != nil {
		return 0, err
	}

	var resp struct {
		Result int
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return 0, wrapError(err)
	}
	return resp.Result, nil
}

// SetAdminTitle sets a custom title for an administrator.
// A title should be 0-16 characters length, emoji are not allowed.
func (b *Bot) SetAdminTitle(chat *Chat, user *User, title string) error {
	params := map[string]string{
		"chat_id":      chat.Recipient(),
		"user_id":      user.Recipient(),
		"custom_title": title,
	}

	_, err := b.Raw("setChatAdministratorCustomTitle", params)
	return err
}
