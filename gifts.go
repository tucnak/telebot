package telebot

import (
	"encoding/json"
	"strconv"
)

// Gift represents a gift that can be sent by the bot.
//
// Bot API: Gift.
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

	// (Optional) The gift can only be bought by Telegram Premium subscribers.
	IsPremium bool `json:"is_premium,omitempty"`

	// (Optional) The gift can be used, once upgraded, to customize a user's
	// appearance.
	HasColors bool `json:"has_colors,omitempty"`

	// (Optional) The total and remaining number of gifts of this type the
	// bot can send; for limited gifts only.
	PersonalTotalCount     int `json:"personal_total_count,omitempty"`
	PersonalRemainingCount int `json:"personal_remaining_count,omitempty"`

	// (Optional) Background of the gift.
	Background *GiftBackground `json:"background,omitempty"`

	// (Optional) The number of different unique gifts the gift can be
	// upgraded to.
	UniqueGiftVariantCount int `json:"unique_gift_variant_count,omitempty"`

	// (Optional) The chat that published the gift.
	PublisherChat *Chat `json:"publisher_chat,omitempty"`
}

// GiftBackground is the background of a gift. Colors are RGB.
//
// Bot API: GiftBackground.
type GiftBackground struct {
	CenterColor int `json:"center_color"`
	EdgeColor   int `json:"edge_color"`
	TextColor   int `json:"text_color"`
}

// UniqueGift is a unique gift upgraded from a regular gift.
//
// Bot API: UniqueGift.
type UniqueGift struct {
	// Identifier and name of the regular gift it was upgraded from.
	GiftID   string `json:"gift_id"`
	BaseName string `json:"base_name"`

	// Unique name of the gift, usable in https://t.me/nft/... links.
	Name string `json:"name"`

	// Number of the gift among those upgraded from the same regular gift.
	Number int `json:"number"`

	Model    UniqueGiftModel    `json:"model"`
	Symbol   UniqueGiftSymbol   `json:"symbol"`
	Backdrop UniqueGiftBackdrop `json:"backdrop"`

	// (Optional) The regular gift was for Telegram Premium subscribers only.
	IsPremium bool `json:"is_premium,omitempty"`

	// (Optional) The gift was used to craft another gift.
	IsBurned bool `json:"is_burned,omitempty"`

	// (Optional) The gift comes from the TON blockchain and can't be resold
	// or transferred in Telegram.
	IsFromBlockchain bool `json:"is_from_blockchain,omitempty"`

	// (Optional) Color scheme the owner can use for their name and replies.
	Colors *UniqueGiftColors `json:"colors,omitempty"`

	// (Optional) The chat that published the gift.
	PublisherChat *Chat `json:"publisher_chat,omitempty"`
}

// UniqueGiftModel is the model of a unique gift.
//
// Bot API: UniqueGiftModel.
type UniqueGiftModel struct {
	Name           string  `json:"name"`
	Sticker        Sticker `json:"sticker"`
	RarityPerMille int     `json:"rarity_per_mille"`

	// (Optional) Rarity of a crafted model: "uncommon", "rare", "epic" or
	// "legendary".
	Rarity string `json:"rarity,omitempty"`
}

// UniqueGiftSymbol is the symbol shown on the pattern of a unique gift.
//
// Bot API: UniqueGiftSymbol.
type UniqueGiftSymbol struct {
	Name           string  `json:"name"`
	Sticker        Sticker `json:"sticker"`
	RarityPerMille int     `json:"rarity_per_mille"`
}

// UniqueGiftBackdrop is the backdrop of a unique gift.
//
// Bot API: UniqueGiftBackdrop.
type UniqueGiftBackdrop struct {
	Name           string                   `json:"name"`
	Colors         UniqueGiftBackdropColors `json:"colors"`
	RarityPerMille int                      `json:"rarity_per_mille"`
}

// UniqueGiftBackdropColors are the colors of a unique gift's backdrop, in
// RGB.
//
// Bot API: UniqueGiftBackdropColors.
type UniqueGiftBackdropColors struct {
	CenterColor int `json:"center_color"`
	EdgeColor   int `json:"edge_color"`
	SymbolColor int `json:"symbol_color"`
	TextColor   int `json:"text_color"`
}

// UniqueGiftColors is a color scheme based on a unique gift, for a user's
// name, replies and link previews. Colors are RGB.
//
// Bot API: UniqueGiftColors.
type UniqueGiftColors struct {
	ModelCustomEmojiID    string `json:"model_custom_emoji_id"`
	SymbolCustomEmojiID   string `json:"symbol_custom_emoji_id"`
	LightThemeMainColor   int    `json:"light_theme_main_color"`
	LightThemeOtherColors []int  `json:"light_theme_other_colors"`
	DarkThemeMainColor    int    `json:"dark_theme_main_color"`
	DarkThemeOtherColors  []int  `json:"dark_theme_other_colors"`
}

// GiftInfo is a service message about a regular gift sent or received.
//
// Bot API: GiftInfo.
type GiftInfo struct {
	Gift Gift `json:"gift"`

	// (Optional) Identifier of the received gift; for business accounts only.
	OwnedGiftID string `json:"owned_gift_id,omitempty"`

	// (Optional) Stars the receiver can claim by converting the gift.
	ConvertStarCount int `json:"convert_star_count,omitempty"`

	// (Optional) Stars prepaid for upgrading the gift.
	PrepaidUpgradeStarCount int `json:"prepaid_upgrade_star_count,omitempty"`

	// (Optional) The upgrade was bought after the gift was sent.
	IsUpgradeSeparate bool `json:"is_upgrade_separate,omitempty"`

	// (Optional) The gift can be upgraded to a unique gift.
	CanBeUpgraded bool `json:"can_be_upgraded,omitempty"`

	// (Optional) Text added to the gift and its entities.
	Text     string          `json:"text,omitempty"`
	Entities []MessageEntity `json:"entities,omitempty"`

	// (Optional) The sender and text are shown only to the receiver.
	IsPrivate bool `json:"is_private,omitempty"`

	// (Optional) Number reserved for the gift when upgraded.
	UniqueGiftNumber int `json:"unique_gift_number,omitempty"`
}

// UniqueGiftInfo is a service message about a unique gift sent or received.
//
// Bot API: UniqueGiftInfo.
type UniqueGiftInfo struct {
	Gift UniqueGift `json:"gift"`

	// How the gift was obtained: "upgrade", "transfer", "resale", "gifted_upgrade"
	// or "offer".
	Origin string `json:"origin"`

	// (Optional) Text added to the gift and its entities (Bot API 10.3).
	Text     string          `json:"text,omitempty"`
	Entities []MessageEntity `json:"entities,omitempty"`

	// (Optional) The sender and text are shown only to the receiver
	// (Bot API 10.3).
	IsPrivate bool `json:"is_private,omitempty"`

	// (Optional) For resold gifts, the currency ("XTR" or "TON") and amount
	// paid.
	LastResaleCurrency string `json:"last_resale_currency,omitempty"`
	LastResaleAmount   int64  `json:"last_resale_amount,omitempty"`

	// (Optional) Identifier of the received gift; for business accounts only.
	OwnedGiftID string `json:"owned_gift_id,omitempty"`

	// (Optional) Stars needed to transfer the gift, omitted if the bot can't.
	TransferStarCount int `json:"transfer_star_count,omitempty"`

	// (Optional) When the gift can be transferred, as a Unix time.
	NextTransferDate int64 `json:"next_transfer_date,omitempty"`
}

// Gifts represents a list of gifts.
//
// Bot API: Gifts.
type Gifts struct {
	// The list of gifts
	Gifts []Gift `json:"gifts"`
}

// GetAvailableGifts returns the list of gifts that can be sent by the bot to users.
//
// Bot API: getAvailableGifts.
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
//
// Bot API: sendGift.
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
