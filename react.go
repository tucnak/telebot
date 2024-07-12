package telebot

import (
	"encoding/json"
)

const (
	ReactionTypeEmoji       = "emoji"
	ReactionTypeCustomEmoji = "custom_emoji"
)

// Reaction describes the type of reaction.
// Describes an instance of ReactionTypeCustomEmoji and ReactionTypeEmoji.
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
type ReactionCount struct {
	// Type of the reaction.
	Type Reaction `json:"type"`

	// Number of times the reaction was added.
	Count int `json:"total_count"`
}

// ReactionOptions represents an object of reaction options.
type ReactionOptions struct {
	// List of reaction types to set on the message.
	Reactions []Reaction `json:"reaction"`

	// Pass True to set the reaction with a big animation.
	Big bool `json:"is_big"`
}

// React changes the chosen reactions on a message. Service messages can't be
// reacted to. Automatically forwarded messages from a channel to its discussion group have
// the same available reactions as messages in the channel.
func (b *Bot) React(to Recipient, msg Editable, opts ReactionOptions) error {
	if to == nil {
		return ErrBadRecipient
	}
	msgID, _ := msg.MessageSig()

	params := map[string]string{
		"chat_id":    to.Recipient(),
		"message_id": msgID,
	}

	reactions := make([]Reaction, 0, len(opts.Reactions))

	for _, reaction := range opts.Reactions {
		// Type is required, use ReactionTypeEmoji as default
		if reaction.Type == "" {
			reaction.Type = ReactionTypeEmoji
		}

		reactions = append(reactions, reaction)
	}

	data, _ := json.Marshal(reactions)
	params["reaction"] = string(data)

	if opts.Big {
		params["is_big"] = "true"
	}

	_, err := b.Raw("setMessageReaction", params)

	return err
}
