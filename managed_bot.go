package telebot

import (
	"encoding/json"
)

// Managed bots (Bot API 9.6) are bots that users create to be controlled by
// a bot with User.CanManageBots. Ask a user to create one with a ReplyButton
// whose ManagedBot is set.

// ManagedBotCreated is a service message about a user creating a bot to be
// managed by the current bot.
//
// Bot API: ManagedBotCreated.
type ManagedBotCreated struct {
	// The new bot. Get its token with ManagedBotToken.
	Bot *User `json:"bot"`
}

// ManagedBotUpdated is received, on OnManagedBot, when a managed bot is
// created or its token or owner changes.
//
// Bot API: ManagedBotUpdated.
type ManagedBotUpdated struct {
	// The user that created the bot.
	User *User `json:"user"`

	// The bot. Get its token with ManagedBotToken.
	Bot *User `json:"bot"`
}

// ReplyManagedBot asks the user to create a bot managed by the current bot;
// set it as ReplyButton.ManagedBot.
//
// Bot API: KeyboardButtonRequestManagedBot.
type ReplyManagedBot struct {
	// Signed 32-bit identifier of the request, unique within the message.
	ID int32 `json:"request_id"`

	// (Optional) Suggested name and username for the bot. The username must
	// end with "bot", or for Telegram Premium users also with "agent", "ai",
	// "app", "game" or "gay" (Bot API 10.4).
	SuggestedName     string `json:"suggested_name,omitempty"`
	SuggestedUsername string `json:"suggested_username,omitempty"`
}

// BotAccessSettings are the access settings of a managed bot (Bot API
// 10.0).
//
// Bot API: BotAccessSettings.
type BotAccessSettings struct {
	// Only the owner and AddedUsers can access the bot.
	IsAccessRestricted bool `json:"is_access_restricted"`

	// (Optional) Users other than the owner who can access a restricted bot.
	// When setting, only their IDs are used, up to 10.
	AddedUsers []User `json:"added_users,omitempty"`
}

// ManagedBotToken returns the token of a managed bot.
//
// Bot API: getManagedBotToken.
func (b *Bot) ManagedBotToken(bot *User) (string, error) {
	return b.managedBotToken("getManagedBotToken", bot)
}

// ReplaceManagedBotToken revokes the token of a managed bot and returns a
// new one.
//
// Bot API: replaceManagedBotToken.
func (b *Bot) ReplaceManagedBotToken(bot *User) (string, error) {
	return b.managedBotToken("replaceManagedBotToken", bot)
}

func (b *Bot) managedBotToken(method string, bot *User) (string, error) {
	data, err := b.Raw(method, map[string]string{"user_id": bot.Recipient()})
	if err != nil {
		return "", err
	}

	var resp struct {
		Result string
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", wrapError(err)
	}
	return resp.Result, nil
}

// ManagedBotAccessSettings returns the access settings of a managed bot
// (Bot API 10.0).
//
// Bot API: getManagedBotAccessSettings.
func (b *Bot) ManagedBotAccessSettings(bot *User) (*BotAccessSettings, error) {
	data, err := b.Raw("getManagedBotAccessSettings", map[string]string{"user_id": bot.Recipient()})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *BotAccessSettings
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// SetManagedBotAccessSettings changes the access settings of a managed bot
// (Bot API 10.0).
//
// Bot API: setManagedBotAccessSettings.
func (b *Bot) SetManagedBotAccessSettings(bot *User, settings BotAccessSettings) error {
	params := map[string]interface{}{
		"user_id":              bot.ID,
		"is_access_restricted": settings.IsAccessRestricted,
	}
	if len(settings.AddedUsers) > 0 {
		ids := make([]int64, len(settings.AddedUsers))
		for i, u := range settings.AddedUsers {
			ids[i] = u.ID
		}
		params["added_user_ids"] = ids
	}

	_, err := b.Raw("setManagedBotAccessSettings", params)
	return err
}

// PreparedKeyboardButton is a keyboard button saved for a user of a Mini App.
//
// Bot API: PreparedKeyboardButton.
type PreparedKeyboardButton struct {
	ID string `json:"id"`
}

// PrepareButton stores a keyboard button that user can use in a
// Mini App. The button must request users, a chat or a managed bot
// (Bot API 9.6).
//
// Bot API: savePreparedKeyboardButton.
func (b *Bot) PrepareButton(user Recipient, button ReplyButton) (*PreparedKeyboardButton, error) {
	data, err := json.Marshal(button)
	if err != nil {
		return nil, err
	}

	data, err = b.Raw("savePreparedKeyboardButton", map[string]string{
		"user_id": user.Recipient(),
		"button":  string(data),
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *PreparedKeyboardButton
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}
