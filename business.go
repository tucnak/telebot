package telebot

import (
	"encoding/json"
	"time"
)

type BusinessConnection struct {
	// Unique identifier of the business connection
	ID string `json:"id"`

	// Business account user that created the business connection
	Sender *User `json:"user"`

	// Identifier of a private chat with the user who created the business connection. This
	// number may have more than 32 significant bits and some programming languages may
	// have difficulty/silent defects in interpreting it. But it has at most 52 significant bits,
	// so a 64-bit integer or double-precision float type are safe for storing this identifier.
	UserChatID int64 `json:"user_chat_id"`

	// Unixtime, use BusinessConnection.Time() to get time.Time.
	Unixtime int64 `json:"date"`

	// True, if the bot can act on behalf of the business account in chats that were active in the last 24 hours
	CanReply bool `json:"can_reply"`

	// True, if the connection is active
	Enabled bool `json:"is_enabled"`
}

// Time returns the moment of business connection creation in local time.
func (b *BusinessConnection) Time() time.Time {
	return time.Unix(b.Unixtime, 0)
}

type BusinessMessagesDeleted struct {
	// Unique identifier of the business connection
	BusinessConnectionID string `json:"business_connection_id"`

	// Information about a chat in the business account. The bot
	// may not have access to the chat or the corresponding user.
	Chat *Chat `json:"chat"`

	// The list of identifiers of deleted messages in the chat of the business account
	MessageIDs []int `json:"message_ids"`
}

type BusinessIntro struct {
	// (Optional)
	// Title text of the business intro
	Title string `json:"title"`

	// Message text of the business intro
	Message string `json:"message"`

	// Sticker of the business intro
	Sticker *Sticker `json:"sticker"`
}

type BusinessLocation struct {
	// Address of the business
	Address string `json:"address"`

	// (Optional) Location of the business
	Location *Location `json:"location"`
}

type BusinessOpeningHoursInterval struct {
	// The minute's sequence number in a week, starting on Monday,
	// marking the start of the time interval during which the business
	// is open; 0 - 7 * 24 * 60
	OpeningMinute int `json:"opening_minute"`

	// The minute's sequence number in a week, starting on Monday,
	// marking the start of the time interval during which the business
	// is open; 0 - 7 * 24 * 60
	ClosingMinute int `json:"closing_minute"`
}

type BusinessOpeningHours struct {
	// Unique name of the time zone for which the opening hours are defined
	Timezone string `json:"time_zone_name"`

	// List of time intervals describing business opening hours
	OpeningHours []BusinessOpeningHoursInterval `json:"opening_hours"`
}

// BusinessConnection returns the information about the connection of the bot with a business account.
func (b *Bot) BusinessConnection(id string) (*BusinessConnection, error) {
	params := map[string]string{
		"business_connection_id": id,
	}

	data, err := b.Raw("getBusinessConnection", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *BusinessConnection
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}
