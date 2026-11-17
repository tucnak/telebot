package telebot

import "strconv"

// Channel direct messages and suggested posts (Bot API 9.2). A channel's
// direct messages chat is a supergroup with Chat.IsDirectMessages, split
// into topics, one per user. Send to a topic with
// SendOptions.DirectMessagesTopicID, and suggest a post to the channel with
// SendOptions.SuggestedPost.

// DirectMessagesTopic is a topic of a direct messages chat.
//
// Bot API: DirectMessagesTopic.
type DirectMessagesTopic struct {
	ID int64 `json:"topic_id"`

	// (Optional) The user who created the topic; currently always set.
	User *User `json:"user,omitempty"`
}

// Currencies of suggested post prices.
const (
	CurrencyStars = "XTR"
	CurrencyTON   = "TON"
)

// SuggestedPostPrice is the price of a suggested post, in the smallest units
// of its currency: Telegram Stars ("XTR") or nanograms ("TON").
//
// Bot API: SuggestedPostPrice.
type SuggestedPostPrice struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"`
}

// SuggestedPostParameters suggests the sent message as a post to the
// channel; pass it to Send or set it as SendOptions.SuggestedPost.
//
// Bot API: SuggestedPostParameters.
type SuggestedPostParameters struct {
	// (Optional) Proposed price; the post is unpaid if nil.
	Price *SuggestedPostPrice `json:"price,omitempty"`

	// (Optional) Proposed publication time, 5 minutes to 30 days ahead, as a
	// Unix time; zero lets the channel choose.
	SendDate int64 `json:"send_date,omitempty"`
}

// Suggested post states.
const (
	SuggestedPostPending  = "pending"
	SuggestedPostApproved = "approved"
	SuggestedPostDeclined = "declined"
)

// SuggestedPostInfo describes a suggested post.
//
// Bot API: SuggestedPostInfo.
type SuggestedPostInfo struct {
	// SuggestedPostPending, SuggestedPostApproved or SuggestedPostDeclined.
	State string `json:"state"`

	Price    *SuggestedPostPrice `json:"price,omitempty"`
	SendDate int64               `json:"send_date,omitempty"`
}

// SuggestedPostApprovalInfo is a service message about the approval of a
// suggested post.
//
// Bot API: SuggestedPostApproved.
type SuggestedPostApprovalInfo struct {
	SuggestedPostMessage *Message            `json:"suggested_post_message,omitempty"`
	Price                *SuggestedPostPrice `json:"price,omitempty"`
	SendDate             int64               `json:"send_date"`
}

// SuggestedPostApprovalFailed is a service message about a suggested post
// whose approval failed, currently for lack of the user's funds.
//
// Bot API: SuggestedPostApprovalFailed.
type SuggestedPostApprovalFailed struct {
	SuggestedPostMessage *Message           `json:"suggested_post_message,omitempty"`
	Price                SuggestedPostPrice `json:"price"`
}

// SuggestedPostDeclineInfo is a service message about a declined suggested
// post.
//
// Bot API: SuggestedPostDeclined.
type SuggestedPostDeclineInfo struct {
	SuggestedPostMessage *Message `json:"suggested_post_message,omitempty"`
	Comment              string   `json:"comment,omitempty"`
}

// SuggestedPostPaid is a service message about a payment for a suggested
// post.
//
// Bot API: SuggestedPostPaid.
type SuggestedPostPaid struct {
	SuggestedPostMessage *Message `json:"suggested_post_message,omitempty"`

	// CurrencyStars or CurrencyTON.
	Currency string `json:"currency"`

	// (Optional) Nanograms received, for TON payments.
	Amount int64 `json:"amount,omitempty"`

	// (Optional) Stars received, for Star payments.
	StarAmount *StarAmount `json:"star_amount,omitempty"`
}

// SuggestedPostRefunded is a service message about a refunded suggested
// post.
//
// Bot API: SuggestedPostRefunded.
type SuggestedPostRefunded struct {
	SuggestedPostMessage *Message `json:"suggested_post_message,omitempty"`

	// "post_deleted" or "payment_refunded".
	Reason string `json:"reason"`
}

// ApproveSuggestedPost approves the suggested post msg in a direct messages
// chat. sendDate, a Unix time, sets when it is published if the post didn't
// set it; zero leaves it. The bot needs the can_post_messages right in the
// channel.
//
// Bot API: approveSuggestedPost.
func (b *Bot) ApproveSuggestedPost(msg Editable, sendDate int64) error {
	msgID, chatID := msg.MessageSig()
	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": msgID,
	}
	if sendDate != 0 {
		params["send_date"] = strconv.FormatInt(sendDate, 10)
	}

	_, err := b.Raw("approveSuggestedPost", params)
	return err
}

// DeclineSuggestedPost declines the suggested post msg in a direct messages
// chat, with an optional comment of up to 128 characters. The bot needs the
// can_manage_direct_messages right in the channel.
//
// Bot API: declineSuggestedPost.
func (b *Bot) DeclineSuggestedPost(msg Editable, comment string) error {
	msgID, chatID := msg.MessageSig()
	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": msgID,
	}
	if comment != "" {
		params["comment"] = comment
	}

	_, err := b.Raw("declineSuggestedPost", params)
	return err
}
