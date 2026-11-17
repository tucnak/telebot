package telebot

import (
	"encoding/json"
	"strconv"
)

const (
	ReactionTypeEmoji       = "emoji"
	ReactionTypeCustomEmoji = "custom_emoji"
)

// Reaction describes the type of reaction.
// Describes an instance of ReactionTypeCustomEmoji and ReactionTypeEmoji.
//
// Bot API: ReactionType, ReactionTypeEmoji, ReactionTypeCustomEmoji, ReactionTypePaid.
type Reaction struct {
	// Type of the reaction, always “emoji”
	Type string `json:"type"`

	// Reaction emoji.
	Emoji string `json:"emoji,omitempty"`

	// Custom emoji identifier.
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

// ReactionCount represents a reaction added to a message along
// with the number of times it was added.
//
// Bot API: ReactionCount.
type ReactionCount struct {
	// Type of the reaction.
	Type Reaction `json:"type"`

	// Number of times the reaction was added.
	Count int `json:"total_count"`
}

// Reactions represents an object of reaction options.
type Reactions struct {
	// List of reaction types to set on the message.
	Reactions []Reaction `json:"reaction"`

	// Pass True to set the reaction with a big animation.
	Big bool `json:"is_big"`
}

// React changes the chosen reactions on a message. Automatically forwarded messages from a channel to its discussion group have
// the same available reactions as messages in the channel.
// Bot API 8.3+: Reactions are now allowed on most service messages.
//
// Bot API: setMessageReaction.
func (b *Bot) React(to Recipient, msg Editable, r Reactions) error {
	if to == nil {
		return ErrBadRecipient
	}

	msgID, _ := msg.MessageSig()
	params := map[string]string{
		"chat_id":    to.Recipient(),
		"message_id": msgID,
	}

	data, _ := json.Marshal(r.Reactions)
	params["reaction"] = string(data)

	if r.Big {
		params["is_big"] = "true"
	}

	_, err := b.Raw("setMessageReaction", params)
	return err
}

// reactionActor adds the user_id or actor_chat_id of who added reactions.
func reactionActor(params map[string]string, actor Recipient) {
	if chat, ok := actor.(*Chat); ok {
		params["actor_chat_id"] = chat.Recipient()
	} else if actor != nil {
		params["user_id"] = actor.Recipient()
	}
}

// DeleteReaction removes the reaction that actor, a *User or a *Chat, added
// to msg in a group or supergroup. The bot needs the can_delete_messages
// right (Bot API 10.0).
//
// Bot API: deleteMessageReaction.
func (b *Bot) DeleteReaction(msg Editable, actor Recipient) error {
	msgID, chatID := msg.MessageSig()
	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": msgID,
	}
	reactionActor(params, actor)

	_, err := b.Raw("deleteMessageReaction", params)
	return err
}

// DeleteAllReactions removes up to 10000 recent reactions that actor, a
// *User or a *Chat, added in a group or supergroup. The bot needs the
// can_delete_messages right (Bot API 10.0).
//
// Bot API: deleteAllMessageReactions.
func (b *Bot) DeleteAllReactions(chat, actor Recipient) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}
	reactionActor(params, actor)

	_, err := b.Raw("deleteAllMessageReactions", params)
	return err
}
