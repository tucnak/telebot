package telebot

import (
	"strconv"
	"time"
)

// Rights is a list of privileges available to chat members.
type Rights struct {
	CanBeEdited        bool `json:"can_be_edited,omitempty"`             // 1
	CanChangeInfo      bool `json:"can_change_info,omitempty"`           // 2
	CanPostMessages    bool `json:"can_post_messages,omitempty"`         // 3
	CanEditMessages    bool `json:"can_edit_messages,omitempty"`         // 4
	CanDeleteMessages  bool `json:"can_delete_messages,omitempty"`       // 5
	CanInviteUsers     bool `json:"can_invite_users,omitempty"`          // 6
	CanRestrictMembers bool `json:"can_restrict_members,omitempty"`      // 7
	CanPinMessages     bool `json:"can_pin_messages,omitempty"`          // 8
	CanPromoteMembers  bool `json:"can_promote_members,omitempty"`       // 9
	CanSendMessages    bool `json:"can_send_messages,omitempty"`         // 10
	CanSendMedia       bool `json:"can_send_media_messages,omitempty"`   // 11
	CanSendOther       bool `json:"can_send_other_messages,omitempty"`   // 12
	CanAddPreviews     bool `json:"can_add_web_page_previews,omitempty"` // 13
}

// NoRights is the default Rights{}
func NoRights() Rights { return Rights{} }

// NoRestrictions should be used when un-restricting or
// un-promoting user.
//
//	   member.Rights = NoRestrictions()
//     bot.Restrict(chat, member)
//
func NoRestrictions() Rights {
	return Rights{
		true, false, false, false, false, // 1-5
		false, false, false, false, true, // 6-10
		true, true, true}
}

// AdminRights could be used to promote user to admin.
func AdminRights() Rights {
	return Rights{
		true, true, true, true, true, // 1-5
		true, true, true, true, true, // 6-10
		true, true, true} // 11-13
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

	respJSON, err := b.sendCommand("kickChatMember", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// Unban will unban user from chat, who would have thought eh?
func (b *Bot) Unban(chat *Chat, user *User) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"user_id": user.Recipient(),
	}

	respJSON, err := b.sendCommand("unbanChatMember", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// Restrict let's you restrict a subset of member's rights until
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

	respJSON, err := b.sendCommand("restrictChatMember", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
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

	respJSON, err := b.sendCommand("promoteChatMember", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}
