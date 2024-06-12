package telebot

import (
	"encoding/json"
	"time"
)

type TransactionType = string

const (
	TransactionTypeFragment TransactionType = "fragment"
	TransactionTypeUser     TransactionType = "user"
	TransactionTypeOther    TransactionType = "other"
)

type RevenueState = string

const (
	RevenueStatePending   RevenueState = "pending"
	RevenueStateSucceeded RevenueState = "succeeded"
	RevenueStateFailed    RevenueState = "failed"
)

type TransactionPartner struct {
	// Type of the state
	Type TransactionType `json:"type"`

	// (Optional) State of the transaction if the transaction is outgoing$$
	WithdrawalState RevenueWithdrawalState `json:"withdrawal_state,omitempty"`

	// Information about the user
	Partner *User `json:"user,omitempty"`
}

type RevenueWithdrawalState struct {
	// Type of the state
	Type RevenueState `json:"type"`

	// Date the withdrawal was completed in Unix time
	Date int `json:"date,omitempty"`

	// An HTTPS URL that can be used to see transaction details
	URL string `json:"url,omitempty"`
}

type StarTransaction struct {
	// Unique identifier of the transaction. Coincides with the identifer of the
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

// Date returns the local datetime.
func (c *StarTransaction) Date() time.Time {
	return time.Unix(c.Unixtime, 0)
}

// GetStarTransactions Returns the bot's Telegram Star transactions in chronological order
func (b *Bot) GetStarTransactions(offset, limit int) ([]StarTransaction, error) {
	params := map[string]int{
		"offset": offset,
		"limit":  limit,
	}

	data, err := b.Raw("getStarTransactions", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result struct {
			Transactions []StarTransaction `json:"transactions"`
		}
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result.Transactions, nil
}

// SendStarRefund returns a successful payment in Telegram Stars.
func (b *Bot) SendStarRefund(to Recipient, chargeID string) error {
	params := map[string]string{
		"user_id":                    to.Recipient(),
		"telegram_payment_charge_id": chargeID,
	}

	_, err := b.Raw("refundStarPayment", params)
	return err
}
