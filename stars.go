package telebot

import "time"

type TransactionType = string

const (
	TransactionTypeUser           TransactionType = "user"
	TransactionTypeFragment       TransactionType = "fragment"
	TransactionPartnerTelegramAds TransactionType = "telegram_ads"
	TransactionTypeOther          TransactionType = "other"
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

	// Date the transaction was created in Unix time
	Unixtime int64 `json:"date"`

	// (Optional) Source of an incoming transaction (e.g., a user purchasing goods
	// or services, Fragment refunding a failed withdrawal). Only for incoming transactions
	Source TransactionPartner `json:"source"`

	// (Optional) Receiver of an outgoing transaction (e.g., a user for a purchase
	// refund, Fragment for a withdrawal). Only for outgoing transactions
	Receiver TransactionPartner `json:"receiver"`
}

type TransactionPartner struct {
	// Type of the state
	Type    TransactionType `json:"type"`
	User    *User           `json:"user,omitempty"`
	Payload string          `json:"invoice_payload"`

	// (Optional) State of the transaction if the transaction is outgoing$$
	Withdrawal RevenueWithdrawal `json:"withdrawal_state,omitempty"`
}

type RevenueWithdrawal struct {
	// Type of the state
	Type RevenueState `json:"type"`

	// Date the withdrawal was completed in Unix time
	Unixtime int `json:"date,omitempty"`

	// An HTTPS URL that can be used to see transaction details
	URL string `json:"url,omitempty"`
}

// Time returns the date of the transaction.
func (c *StarTransaction) Time() time.Time {
	return time.Unix(c.Unixtime, 0)
}

// Time returns the date of the withdrawal.
func (s *RevenueWithdrawal) Time() time.Time {
	return time.Unix(int64(s.Unixtime), 0)
}
