package telebot

import "encoding/json"

// Chat share theme types.
const (
	ShareThemeEmoji      = "emoji"
	ShareThemeUniqueGift = "unique_gift"
)

// ChatShareTheme is the theme of the screen for sharing a link to a user or
// a chat (Bot API 10.4): emoji on a background ("emoji"), or a unique gift
// ("unique_gift").
//
// Bot API: ChatShareTheme, ChatShareThemeEmoji, ChatShareThemeUniqueGift.
type ChatShareTheme struct {
	// ShareThemeEmoji or ShareThemeUniqueGift.
	Type string `json:"type"`

	// For "emoji", 0-6 regular or custom emoji.
	Emojis []EmojiType `json:"emojis,omitempty"`

	// For "emoji", the backdrop used as the background; zero for the
	// default background. See Bot.ShareThemeBackdrops.
	BackgroundBackdropID int `json:"background_backdrop_id,omitempty"`

	// For "unique_gift", the gift.
	Gift *UniqueGift `json:"gift,omitempty"`
}

// EmojiType is a regular emoji ("emoji") or a custom emoji ("custom_emoji")
// (Bot API 10.4).
//
// Bot API: EmojiType, EmojiTypeEmoji, EmojiTypeCustomEmoji.
type EmojiType struct {
	// "emoji" or "custom_emoji".
	Type string `json:"type"`

	// The emoji, for "emoji".
	Emoji string `json:"emoji,omitempty"`

	// The custom emoji, for "custom_emoji". See Bot.CustomEmojiStickers.
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

// ChatShareThemeBackdrop is a backdrop usable as the background of a chat
// share theme.
//
// Bot API: ChatShareThemeBackdrop.
type ChatShareThemeBackdrop struct {
	ID     int                      `json:"id"`
	Name   string                   `json:"name"`
	Colors UniqueGiftBackdropColors `json:"colors"`
}

// ChatShareThemeBackdrops is a list of chat share theme backdrops.
//
// Bot API: ChatShareThemeBackdrops.
type ChatShareThemeBackdrops struct {
	Backdrops []ChatShareThemeBackdrop `json:"backdrops"`
}

// ShareThemeBackdrops returns the backdrops usable as the background of a
// chat share theme (Bot API 10.4).
//
// Bot API: getChatShareThemeBackdrops.
func (b *Bot) ShareThemeBackdrops() ([]ChatShareThemeBackdrop, error) {
	data, err := b.Raw("getChatShareThemeBackdrops", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result ChatShareThemeBackdrops
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result.Backdrops, nil
}
