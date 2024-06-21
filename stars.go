package telebot

type RevenueWithdrawalState struct {
	// Type of the state, always “pending”
	Type string `json:"type"`

	// Date the withdrawal was completed in Unix time
	Date int `json:"date,omitempty"`

	// An HTTPS URL that can be used to see transaction details
	URL string `json:"url,omitempty"`
}

type TransactionPartner struct {
	// Type of the state, always “fragment”
	Type string `json:"type"`

	// (Optional) State of the transaction if the transaction is outgoing$$
	WithdrawalState RevenueWithdrawalState `json:"withdrawal_state,omitempty"`

	// Information about the user
	User User `json:"user,omitempty"`
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
	Date int `json:"date"`

	// (Optional) Source of an incoming transaction (e.g., a user purchasing goods
	//or services, Fragment refunding a failed withdrawal). Only for incoming transactions
	Source TransactionPartner `json:"source"`

	// (Optional) Receiver of an outgoing transaction (e.g., a user for a purchase
	//refund, Fragment for a withdrawal). Only for outgoing transactions
	Receiver TransactionPartner `json:"receiver"`
}

type StarTransactions struct {
	// The list of transactions
	Transactions []StarTransaction `json:"transactions"`
}
