package telebot

import (
	"encoding/json"
	"time"
)

type TransactionType = string

const (
	TransactionTypeUser             TransactionType = "user"
	TransactionTypeFragment         TransactionType = "fragment"
	TransactionPartnerTelegramAds   TransactionType = "telegram_ads"
	TransactionPartnerTelegramApi   TransactionType = "telegram_api"
	TransactionTypeAffiliateProgram TransactionType = "affiliate_program"
	TransactionTypeOther            TransactionType = "other"
	// Bot API 8.3: Transaction partner is a chat
	TransactionTypeChat TransactionType = "chat"
)

type RevenueState = string

const (
	RevenueStatePending   RevenueState = "pending"
	RevenueStateSucceeded RevenueState = "succeeded"
	RevenueStateFailed    RevenueState = "failed"
)

// StarTransaction describes a Telegram Star transaction.
//
// Bot API: StarTransaction.
type StarTransaction struct {
	// Unique identifier of the transaction. Coincides with the identifier of the
	// original transaction for refund transactions. Coincides with
	// SuccessfulPayment.telegram_payment_charge_id for successful incoming
	// payments from users.
	ID string `json:"id"`

	// Number of Telegram Stars transferred by the transaction
	Amount int `json:"amount"`

	// (Optional) The number of 0.00000001 TON that were transferred by the transaction
	NanostarAmount int `json:"nanostar_amount,omitempty"`

	// Date the transaction was created in Unix time
	Unixtime int64 `json:"date"`

	// (Optional) Source of an incoming transaction (e.g., a user purchasing goods
	// or services, Fragment refunding a failed withdrawal). Only for incoming transactions
	Source TransactionPartner `json:"source"`

	// (Optional) Receiver of an outgoing transaction (e.g., a user for a purchase
	// refund, Fragment for a withdrawal). Only for outgoing transactions
	Receiver TransactionPartner `json:"receiver"`
}

// TransactionPartnerChat describes a transaction partner that is a chat.
//
// Bot API: TransactionPartnerChat.
type TransactionPartnerChat struct {
	// Type of the transaction partner, always "chat"
	Type TransactionType `json:"type"`

	// Information about the chat
	Chat *Chat `json:"chat"`

	// (Optional) The gift sent to the chat by the bot.
	Gift *Gift `json:"gift,omitempty"`
}

// TransactionPartner describes the source or the recipient of a Telegram
// Star transaction.
//
// Bot API: TransactionPartner, TransactionPartnerUser, TransactionPartnerChat, TransactionPartnerAffiliateProgram, TransactionPartnerFragment, TransactionPartnerTelegramAds, TransactionPartnerTelegramApi, TransactionPartnerOther.
type TransactionPartner struct {
	// Type of the state
	Type    TransactionType `json:"type"`
	User    *User           `json:"user,omitempty"`
	Payload string          `json:"invoice_payload"`

	// (Optional) State of the transaction if the transaction is outgoing$$
	Withdrawal RevenueWithdrawal `json:"withdrawal_state,omitempty"`

	// Bot API 8.0: Star subscriptions
	SubscriptionPeriod int `json:"subscription_period,omitempty"`

	// Bot API 8.0: Gifts
	Gift *Gift `json:"gift,omitempty"`

	// Bot API 8.1: Information about the affiliate program from which the transaction was made
	Affiliate *AffiliateInfo `json:"affiliate,omitempty"`

	// Bot API 8.1: for "affiliate_program" partners, the bot that sponsored
	// the program, and the Stars the bot receives for each 1000 Stars the
	// sponsor receives from referred users.
	SponsorUser        *User `json:"sponsor_user,omitempty"`
	CommissionPerMille int   `json:"commission_per_mille,omitempty"`

	// Bot API 8.3: for "chat" partners, the chat.
	Chat *Chat `json:"chat,omitempty"`

	// Bot API 9.0: for "user" partners, the type of the transaction:
	// "invoice_payment", "paid_media_payment", "gift_purchase",
	// "premium_purchase" or "business_account_transfer".
	TransactionType string `json:"transaction_type,omitempty"`

	// For "paid_media_payment" transactions, the media bought and its
	// payload.
	PaidMedia        []PaidMedia `json:"paid_media,omitempty"`
	PaidMediaPayload string      `json:"paid_media_payload,omitempty"`

	// Bot API 9.0: for "premium_purchase" transactions, the months of
	// Telegram Premium gifted.
	PremiumSubscriptionDuration int `json:"premium_subscription_duration,omitempty"`

	// For "telegram_api" partners, the billed requests over the limits.
	RequestCount int `json:"request_count,omitempty"`
}

// RevenueWithdrawal describes the state of a revenue withdrawal operation.
//
// Bot API: RevenueWithdrawalState, RevenueWithdrawalStatePending, RevenueWithdrawalStateSucceeded, RevenueWithdrawalStateFailed.
type RevenueWithdrawal struct {
	// Type of the state
	Type RevenueState `json:"type"`

	// Date the withdrawal was completed in Unix time
	Unixtime int `json:"date,omitempty"`

	// An HTTPS URL that can be used to see transaction details
	URL string `json:"url,omitempty"`
}

// TransactionPartnerAffiliateProgram describes the affiliate program that paid the transaction.
//
// Bot API: TransactionPartnerAffiliateProgram.
type TransactionPartnerAffiliateProgram struct {
	// Type of the transaction partner, always "affiliate_program"
	Type TransactionType `json:"type"`

	// Sponsor of the affiliate program
	SponsorUser *User `json:"sponsor_user,omitempty"`

	// The number of Telegram Stars received by the bot for each 1000 Telegram Stars received by the affiliate program
	CommissionPerMille int `json:"commission_per_mille"`
}

// AffiliateInfo contains information about the affiliate program from which the transaction was made.
//
// Bot API: AffiliateInfo.
type AffiliateInfo struct {
	// (Optional) Information about the affiliate program that paid the transaction
	AffiliateProgram *TransactionPartnerAffiliateProgram `json:"affiliate_program,omitempty"`

	// (Optional) The bot or user, or the chat, that received the commission.
	AffiliateUser *User `json:"affiliate_user,omitempty"`
	AffiliateChat *Chat `json:"affiliate_chat,omitempty"`

	// (Optional) The number of Telegram Stars received by the affiliate program for each 1000 Telegram Stars received by the bot
	CommissionPerMille int `json:"commission_per_mille"`

	// (Optional) Monetary amount of the commission
	Amount int `json:"amount,omitempty"`

	// (Optional) Nanostar amount of the commission
	NanostarAmount int `json:"nanostar_amount,omitempty"`
}

// IsAffiliateProgram returns whether this is an affiliate program transaction.
func (a *AffiliateInfo) IsAffiliateProgram() bool {
	return a.AffiliateProgram != nil
}

// CommissionPercentage returns the commission rate as a percentage.
func (a *AffiliateInfo) CommissionPercentage() float64 {
	return float64(a.CommissionPerMille) / 10.0
}

// NanostarAmountAsTON returns the nanostar commission amount as TON value.
func (a *AffiliateInfo) NanostarAmountAsTON() float64 {
	return float64(a.NanostarAmount) / 100000000.0
}

// HasNanostarAmount returns whether the affiliate info has a nanostar amount.
func (a *AffiliateInfo) HasNanostarAmount() bool {
	return a.NanostarAmount > 0
}

// Time returns the date of the transaction.
func (c *StarTransaction) Time() time.Time {
	return time.Unix(c.Unixtime, 0)
}

// Time returns the date of the withdrawal.
func (s *RevenueWithdrawal) Time() time.Time {
	return time.Unix(int64(s.Unixtime), 0)
}

// NanostarAmountAsTON returns the nanostar amount as TON value.
func (c *StarTransaction) NanostarAmountAsTON() float64 {
	return float64(c.NanostarAmount) / 100000000.0
}

// HasNanostarAmount returns whether the transaction has a nanostar amount.
func (c *StarTransaction) HasNanostarAmount() bool {
	return c.NanostarAmount > 0
}

// StarAmount is an amount of Telegram Stars (Bot API 9.1).
//
// Bot API: StarAmount.
type StarAmount struct {
	// Whole Stars, rounded towards zero; can be negative.
	Amount int64 `json:"amount"`

	// (Optional) Billionths of a Star, from -999999999 to 999999999.
	NanostarAmount int64 `json:"nanostar_amount,omitempty"`
}

// StarBalance returns the bot's current balance of Telegram Stars
// (Bot API 9.1).
//
// Bot API: getMyStarBalance.
func (b *Bot) StarBalance() (*StarAmount, error) {
	data, err := b.Raw("getMyStarBalance", nil)
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
