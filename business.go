package telebot

import (
	"encoding/json"
	"strconv"
	"time"
)

// BusinessConnection describes the connection of the bot with a business account.
//
// Bot API: BusinessConnection.
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
	//
	// Deprecated: Telegram no longer sends it; use Rights.CanReply.
	CanReply bool `json:"can_reply"`

	// Bot API 9.0: (Optional) Rights of the bot in the business account.
	Rights *BusinessBotRights `json:"rights,omitempty"`

	// True, if the connection is active
	Enabled bool `json:"is_enabled"`
}

// Time returns the moment of business connection creation in local time.
func (b *BusinessConnection) Time() time.Time {
	return time.Unix(b.Unixtime, 0)
}

// BusinessMessagesDeleted is received when messages are deleted from a connected business account.
//
// Bot API: BusinessMessagesDeleted.
type BusinessMessagesDeleted struct {
	// Unique identifier of the business connection
	BusinessConnectionID string `json:"business_connection_id"`

	// Information about a chat in the business account. The bot
	// may not have access to the chat or the corresponding user.
	Chat *Chat `json:"chat"`

	// The list of identifiers of deleted messages in the chat of the business account
	MessageIDs []int `json:"message_ids"`
}

// BusinessIntro is the intro of a business account.
//
// Bot API: BusinessIntro.
type BusinessIntro struct {
	// (Optional)
	// Title text of the business intro
	Title string `json:"title"`

	// Message text of the business intro
	Message string `json:"message"`

	// Sticker of the business intro
	Sticker *Sticker `json:"sticker"`
}

// BusinessLocation is the location of a business account.
//
// Bot API: BusinessLocation.
type BusinessLocation struct {
	// Address of the business
	Address string `json:"address"`

	// (Optional) Location of the business
	Location *Location `json:"location"`
}

// BusinessOpeningHoursInterval is a time interval when a business is open.
//
// Bot API: BusinessOpeningHoursInterval.
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

// BusinessOpeningHours are the opening hours of a business account.
//
// Bot API: BusinessOpeningHours.
type BusinessOpeningHours struct {
	// Unique name of the time zone for which the opening hours are defined
	Timezone string `json:"time_zone_name"`

	// List of time intervals describing business opening hours
	OpeningHours []BusinessOpeningHoursInterval `json:"opening_hours"`
}

// BusinessConnection returns the information about the connection of the bot with a business account.
//
// Bot API: getBusinessConnection.
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

// DeleteBusinessGift deletes a regular gift owned by a managed
// business account without converting it to Telegram Stars (Bot API 10.4).
// It needs the can_convert_gifts_to_stars business bot right.
//
// Bot API: deleteBusinessAccountGift.
func (b *Bot) DeleteBusinessGift(connectionID, ownedGiftID string) error {
	_, err := b.Raw("deleteBusinessAccountGift", map[string]string{
		"business_connection_id": connectionID,
		"owned_gift_id":          ownedGiftID,
	})
	return err
}

// BusinessBotRights are the rights of a business bot (Bot API 9.0).
//
// Bot API: BusinessBotRights.
type BusinessBotRights struct {
	CanReply                   bool `json:"can_reply,omitempty"`
	CanReadMessages            bool `json:"can_read_messages,omitempty"`
	CanDeleteSentMessages      bool `json:"can_delete_sent_messages,omitempty"`
	CanDeleteAllMessages       bool `json:"can_delete_all_messages,omitempty"`
	CanEditName                bool `json:"can_edit_name,omitempty"`
	CanEditBio                 bool `json:"can_edit_bio,omitempty"`
	CanEditProfilePhoto        bool `json:"can_edit_profile_photo,omitempty"`
	CanEditUsername            bool `json:"can_edit_username,omitempty"`
	CanChangeGiftSettings      bool `json:"can_change_gift_settings,omitempty"`
	CanViewGiftsAndStars       bool `json:"can_view_gifts_and_stars,omitempty"`
	CanConvertGiftsToStars     bool `json:"can_convert_gifts_to_stars,omitempty"`
	CanTransferAndUpgradeGifts bool `json:"can_transfer_and_upgrade_gifts,omitempty"`
	CanTransferStars           bool `json:"can_transfer_stars,omitempty"`
	CanManageStories           bool `json:"can_manage_stories,omitempty"`
}

// AcceptedGiftTypes are the types of gifts a user or chat accepts
// (Bot API 9.0).
//
// Bot API: AcceptedGiftTypes.
type AcceptedGiftTypes struct {
	UnlimitedGifts      bool `json:"unlimited_gifts"`
	LimitedGifts        bool `json:"limited_gifts"`
	UniqueGifts         bool `json:"unique_gifts"`
	PremiumSubscription bool `json:"premium_subscription"`
	GiftsFromChannels   bool `json:"gifts_from_channels"`
}

// ReadBusinessMessage marks msg, received in a chat active in the last 24
// hours, as read on behalf of the business account (Bot API 9.0).
//
// Bot API: readBusinessMessage.
func (b *Bot) ReadBusinessMessage(connectionID string, msg Editable) error {
	msgID, chatID := msg.MessageSig()
	_, err := b.Raw("readBusinessMessage", map[string]string{
		"business_connection_id": connectionID,
		"chat_id":                strconv.FormatInt(chatID, 10),
		"message_id":             msgID,
	})
	return err
}

// DeleteBusinessMessages deletes 1-100 messages of one chat on behalf of the
// business account (Bot API 9.0).
//
// Bot API: deleteBusinessMessages.
func (b *Bot) DeleteBusinessMessages(connectionID string, msgs []Editable) error {
	ids := make([]json.Number, len(msgs))
	for i, m := range msgs {
		id, _ := m.MessageSig()
		ids[i] = json.Number(id)
	}
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	_, err = b.Raw("deleteBusinessMessages", map[string]string{
		"business_connection_id": connectionID,
		"message_ids":            string(data),
	})
	return err
}

// SetBusinessName changes the first and last name of the business
// account (Bot API 9.0).
//
// Bot API: setBusinessAccountName.
func (b *Bot) SetBusinessName(connectionID, firstName, lastName string) error {
	params := map[string]string{
		"business_connection_id": connectionID,
		"first_name":             firstName,
	}
	if lastName != "" {
		params["last_name"] = lastName
	}
	_, err := b.Raw("setBusinessAccountName", params)
	return err
}

// SetBusinessUsername changes the username of the business account;
// an empty username removes it (Bot API 9.0).
//
// Bot API: setBusinessAccountUsername.
func (b *Bot) SetBusinessUsername(connectionID, username string) error {
	_, err := b.Raw("setBusinessAccountUsername", map[string]string{
		"business_connection_id": connectionID,
		"username":               username,
	})
	return err
}

// SetBusinessBio changes the bio of the business account; an empty
// bio removes it (Bot API 9.0).
//
// Bot API: setBusinessAccountBio.
func (b *Bot) SetBusinessBio(connectionID, bio string) error {
	_, err := b.Raw("setBusinessAccountBio", map[string]string{
		"business_connection_id": connectionID,
		"bio":                    bio,
	})
	return err
}

// SetBusinessProfilePhoto changes the profile photo of the business
// account; with public, the public photo shown when the main one is hidden
// (Bot API 9.0).
//
// Bot API: setBusinessAccountProfilePhoto.
func (b *Bot) SetBusinessProfilePhoto(connectionID string, photo InputProfilePhoto, public bool) error {
	files := make(map[string]File)
	data, err := photo.inputProfilePhoto(files)
	if err != nil {
		return err
	}
	params := map[string]string{
		"business_connection_id": connectionID,
		"photo":                  string(data),
	}
	if public {
		params["is_public"] = "true"
	}
	_, err = b.sendFiles("setBusinessAccountProfilePhoto", files, params)
	return err
}

// RemoveBusinessProfilePhoto removes the profile photo of the
// business account, or with public its public photo (Bot API 9.0).
//
// Bot API: removeBusinessAccountProfilePhoto.
func (b *Bot) RemoveBusinessProfilePhoto(connectionID string, public bool) error {
	params := map[string]string{"business_connection_id": connectionID}
	if public {
		params["is_public"] = "true"
	}
	_, err := b.Raw("removeBusinessAccountProfilePhoto", params)
	return err
}

// SetBusinessGiftSettings changes which gifts the business account
// accepts and whether the gift button is always shown (Bot API 9.0).
//
// Bot API: setBusinessAccountGiftSettings.
func (b *Bot) SetBusinessGiftSettings(connectionID string, showGiftButton bool, types AcceptedGiftTypes) error {
	data, err := json.Marshal(types)
	if err != nil {
		return err
	}
	_, err = b.Raw("setBusinessAccountGiftSettings", map[string]string{
		"business_connection_id": connectionID,
		"show_gift_button":       strconv.FormatBool(showGiftButton),
		"accepted_gift_types":    string(data),
	})
	return err
}

// BusinessStarBalance returns the Telegram Stars of the business
// account (Bot API 9.0).
//
// Bot API: getBusinessAccountStarBalance.
func (b *Bot) BusinessStarBalance(connectionID string) (*StarAmount, error) {
	data, err := b.Raw("getBusinessAccountStarBalance", map[string]string{
		"business_connection_id": connectionID,
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *StarAmount
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// TransferBusinessStars moves 1-10000 Telegram Stars from the business
// account to the bot (Bot API 9.0).
//
// Bot API: transferBusinessAccountStars.
func (b *Bot) TransferBusinessStars(connectionID string, stars int) error {
	_, err := b.Raw("transferBusinessAccountStars", map[string]string{
		"business_connection_id": connectionID,
		"star_count":             strconv.Itoa(stars),
	})
	return err
}
