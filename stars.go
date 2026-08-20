package telebot

import "time"

type TransactionType = string

const (
	TransactionTypeUser               TransactionType = "user"
	TransactionTypeFragment           TransactionType = "fragment"
	TransactionPartnerTelegramAds     TransactionType = "telegram_ads"
	TransactionPartnerTelegramApi     TransactionType = "telegram_api"
	TransactionTypeAffiliateProgram   TransactionType = "affiliate_program"
	TransactionTypeOther              TransactionType = "other"
	// Bot API 8.3: Transaction partner is a chat
	TransactionTypeChat               TransactionType = "chat"
)

type RevenueState = string

const (
	RevenueStatePending   RevenueState = "pending"
	RevenueStateSucceeded RevenueState = "succeeded"
	RevenueStateFailed    RevenueState = "failed"
)

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
type TransactionPartnerChat struct {
	// Type of the transaction partner, always "chat"
	Type TransactionType `json:"type"`

	// Information about the chat
	Chat *Chat `json:"chat"`
}

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

	// Bot API 8.3: Information about the chat transaction partner
	Chat *TransactionPartnerChat `json:"chat,omitempty"`
}

type RevenueWithdrawal struct {
	// Type of the state
	Type RevenueState `json:"type"`

	// Date the withdrawal was completed in Unix time
	Unixtime int `json:"date,omitempty"`

	// An HTTPS URL that can be used to see transaction details
	URL string `json:"url,omitempty"`
}

// TransactionPartnerAffiliateProgram describes the affiliate program that paid the transaction.
type TransactionPartnerAffiliateProgram struct {
	// Type of the transaction partner, always "affiliate_program"
	Type TransactionType `json:"type"`

	// Sponsor of the affiliate program
	SponsorUser *User `json:"sponsor_user,omitempty"`

	// The number of Telegram Stars received by the bot for each 1000 Telegram Stars received by the affiliate program
	CommissionPerMille int `json:"commission_per_mille"`
}

// VerificationLevel describes the verification level of a user or chat.
type VerificationLevel string

const (
	VerificationLevelNone      VerificationLevel = "none"
	VerificationLevelBasic     VerificationLevel = "basic"
	VerificationLevelDetailed  VerificationLevel = "detailed"
	VerificationLevelFull      VerificationLevel = "full"
)

// VerificationRequirements describes the requirements for verification.
type VerificationRequirements struct {
	// (Optional) Whether the user needs to provide a phone number
	RequirePhoneNumber bool `json:"require_phone_number,omitempty"`

	// (Optional) Whether the user needs to provide their real name
	RequireRealName bool `json:"require_real_name,omitempty"`

	// (Optional) Whether the user needs to provide their date of birth
	RequireBirthDate bool `json:"require_birth_date,omitempty"`

	// (Optional) Whether the user needs to provide their address
	RequireAddress bool `json:"require_address,omitempty"`

	// (Optional) Whether the user needs to provide identification documents
	RequireDocuments bool `json:"require_documents,omitempty"`

	// (Optional) Custom verification requirements
	CustomRequirements []string `json:"custom_requirements,omitempty"`
}

// VerificationInfo contains information about the verification status.
type VerificationInfo struct {
	// Verification level of the user or chat
	Level VerificationLevel `json:"level"`

	// (Optional) Date when verification was granted in Unix time
	VerificationDate int64 `json:"verification_date,omitempty"`

	// (Optional) Expiration date of verification in Unix time
	ExpirationDate int64 `json:"expiration_date,omitempty"`

	// (Optional) Information about the entity that verified this user or chat
	VerifiedBy string `json:"verified_by,omitempty"`

	// (Optional) Additional verification details
	Details string `json:"details,omitempty"`

	// (Optional) Whether the verification is currently active
	IsActive bool `json:"is_active"`
}

// AffiliateInfo contains information about the affiliate program from which the transaction was made.
type AffiliateInfo struct {
	// (Optional) Information about the affiliate program that paid the transaction
	AffiliateProgram *TransactionPartnerAffiliateProgram `json:"affiliate_program,omitempty"`

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

// VerificationTime returns the date when verification was granted.
func (v *VerificationInfo) VerificationTime() time.Time {
	return time.Unix(v.VerificationDate, 0)
}

// ExpirationTime returns the expiration date of the verification.
func (v *VerificationInfo) ExpirationTime() time.Time {
	return time.Unix(v.ExpirationDate, 0)
}

// IsExpired returns whether the verification has expired.
func (v *VerificationInfo) IsExpired() bool {
	if v.ExpirationDate == 0 {
		return false // No expiration date set
	}
	return time.Now().Unix() > v.ExpirationDate
}

// IsValid returns whether the verification is currently valid (active and not expired).
func (v *VerificationInfo) IsValid() bool {
	return v.IsActive && !v.IsExpired()
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
