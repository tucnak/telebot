package telebot

// EmojiType defines emoji types.
type EmojiType = string

// Reaction describes the type of reaction.
// Describes an instance of ReactionTypeCustomEmoji and ReactionTypeEmoji.
type Reaction struct {
	// Type of the reaction, always “emoji”
	Type string `json:"type"`

	// Reaction emoji.
	Emoji EmojiType `json:"emoji,omitempty"`

	// Custom emoji identifier.
	CustomEmoji string `json:"custom_emoji_id,omitempty"`
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
