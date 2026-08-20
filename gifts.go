package telebot

import (
	"encoding/json"
	"strconv"
)

// Gift represents a gift that can be sent by the bot.
type Gift struct {
	// Unique identifier of the gift
	ID string `json:"id"`

	// The sticker that represents the gift
	Sticker *Sticker `json:"sticker"`

	// The number of Telegram Stars that must be paid to send the sticker
	StarCount int `json:"star_count"`

	// (Optional) The number of Telegram Stars that must be paid to upgrade the gift to a unique one
	UpgradeStarCount int `json:"upgrade_star_count,omitempty"`

	// (Optional) The total number of the gifts of this type that can be sent; for limited gifts only
	TotalCount int `json:"total_count,omitempty"`

	// (Optional) The number of remaining gifts of this type that can be sent; for limited gifts only
	RemainingCount int `json:"remaining_count,omitempty"`
}

// Gifts represents a list of gifts.
type Gifts struct {
	// The list of gifts
	Gifts []Gift `json:"gifts"`
}

// GetAvailableGifts returns the list of gifts that can be sent by the bot to users.
func (b *Bot) GetAvailableGifts() ([]Gift, error) {
	data, err := b.Raw("getAvailableGifts", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result Gifts
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result.Gifts, nil
}

// SendGift sends a gift to the given user or chat member. The gift can't be converted to Telegram Stars by the user.
// Additional text can be passed as a string option. PayForUpgrade can be passed as a bool to pay for upgrading the gift.
// For Bot API 8.3+, you can pass a chat as the first parameter to send gifts to specific chat members.
func (b *Bot) SendGift(to Recipient, giftID string, opts ...interface{}) error {
	params := map[string]string{
		"gift_id": giftID,
	}

	// Check if recipient is a user or chat
	switch to.(type) {
	case *User:
		// Send to user directly
		params["user_id"] = to.Recipient()
	case *Chat:
		// Bot API 8.3: Send to chat member (requires additional chat_id parameter)
		params["chat_id"] = to.Recipient()
	default:
		// Default to user_id for backward compatibility
		params["user_id"] = to.Recipient()
	}

	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			// Text for the gift
			params["text"] = v
		case *SendOptions:
			if v.ParseMode != ModeDefault {
				params["text_parse_mode"] = v.ParseMode
			}
			if v.Entities != nil {
				if data, err := json.Marshal(v.Entities); err == nil {
					params["text_entities"] = string(data)
				}
			}
		case bool:
			// Pay for upgrade option
			params["pay_for_upgrade"] = strconv.FormatBool(v)
		default:
			// Handle other variadic options if needed
		}
	}

	_, err := b.Raw("sendGift", params)
	return err
}
