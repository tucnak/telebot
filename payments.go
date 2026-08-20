package telebot

import (
	"encoding/json"
	"math"
	"strconv"
)

// Stars is a provider token for Telegram Stars.
const Stars = "XTR"

// ShippingQuery contains information about an incoming shipping query.
type ShippingQuery struct {
	Sender  *User           `json:"from"`
	ID      string          `json:"id"`
	Payload string          `json:"invoice_payload"`
	Address ShippingAddress `json:"shipping_address"`
}

// ShippingAddress represents a shipping address.
type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

// ShippingOption represents one shipping option.
type ShippingOption struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Prices []Price `json:"prices"`
}

// Payment contains basic information about a successful payment.
type Payment struct {
	Currency         string `json:"currency"`
	Total            int    `json:"total_amount"`
	Payload          string `json:"invoice_payload"`
	OptionID         string `json:"shipping_option_id"`
	Order            Order  `json:"order_info"`
	TelegramChargeID string `json:"telegram_payment_charge_id"`
	ProviderChargeID string `json:"provider_payment_charge_id"`

	// Bot API 8.0: Star subscriptions
	SubscriptionExpirationDate int64 `json:"subscription_expiration_date,omitempty"`
	IsRecurring                bool  `json:"is_recurring,omitempty"`
	IsFirstRecurring           bool  `json:"is_first_recurring,omitempty"`
}

type RefundedPayment struct {
	Currency         string `json:"currency"`
	Total            int    `json:"total_amount"`
	Payload          string `json:"invoice_payload"`
	TelegramChargeID string `json:"telegram_payment_charge_id"`
	ProviderChargeID string `json:"provider_payment_charge_id"`
}

// PreCheckoutQuery contains information about an incoming pre-checkout query.
type PreCheckoutQuery struct {
	Sender   *User  `json:"from"`
	ID       string `json:"id"`
	Currency string `json:"currency"`
	Payload  string `json:"invoice_payload"`
	Total    int    `json:"total_amount"`
	OptionID string `json:"shipping_option_id"`
	Order    Order  `json:"order_info"`
}

// Order represents information about an order.
type Order struct {
	Name        string          `json:"name"`
	PhoneNumber string          `json:"phone_number"`
	Email       string          `json:"email"`
	Address     ShippingAddress `json:"shipping_address"`
}

// Invoice contains basic information about an invoice.
type Invoice struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Payload     string  `json:"payload"`
	Currency    string  `json:"currency"`
	Prices      []Price `json:"prices"`
	Token       string  `json:"provider_token"`
	Data        string  `json:"provider_data"`

	Photo     *Photo `json:"photo"`
	PhotoSize int    `json:"photo_size"`

	// Unique deep-linking parameter that can be used to
	// generate this invoice when used as a start parameter (0).
	Start string `json:"start_parameter"`

	// Shows the total price in the smallest units of the currency.
	// For example, for a price of US$ 1.45 pass amount = 145.
	Total int `json:"total_amount"`

	MaxTipAmount        int   `json:"max_tip_amount"`
	SuggestedTipAmounts []int `json:"suggested_tip_amounts"`

	NeedName            bool `json:"need_name"`
	NeedPhoneNumber     bool `json:"need_phone_number"`
	NeedEmail           bool `json:"need_email"`
	NeedShippingAddress bool `json:"need_shipping_address"`
	SendPhoneNumber     bool `json:"send_phone_number_to_provider"`
	SendEmail           bool `json:"send_email_to_provider"`
	Flexible            bool `json:"is_flexible"`

	// Bot API 8.0: Star subscriptions
	SubscriptionPeriod int `json:"subscription_period,omitempty"`

	// Bot API 8.0: Business connection support
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
}

func (i Invoice) params() map[string]string {
	params := map[string]string{
		"title":                         i.Title,
		"description":                   i.Description,
		"start_parameter":               i.Start,
		"payload":                       i.Payload,
		"provider_token":                i.Token,
		"provider_data":                 i.Data,
		"currency":                      i.Currency,
		"max_tip_amount":                strconv.Itoa(i.MaxTipAmount),
		"need_name":                     strconv.FormatBool(i.NeedName),
		"need_phone_number":             strconv.FormatBool(i.NeedPhoneNumber),
		"need_email":                    strconv.FormatBool(i.NeedEmail),
		"need_shipping_address":         strconv.FormatBool(i.NeedShippingAddress),
		"send_phone_number_to_provider": strconv.FormatBool(i.SendPhoneNumber),
		"send_email_to_provider":        strconv.FormatBool(i.SendEmail),
		"is_flexible":                   strconv.FormatBool(i.Flexible),
	}
	if i.Photo != nil {
		if i.Photo.FileURL != "" {
			params["photo_url"] = i.Photo.FileURL
		}
		if i.PhotoSize > 0 {
			params["photo_size"] = strconv.Itoa(i.PhotoSize)
		}
		if i.Photo.Width > 0 {
			params["photo_width"] = strconv.Itoa(i.Photo.Width)
		}
		if i.Photo.Height > 0 {
			params["photo_height"] = strconv.Itoa(i.Photo.Height)
		}
	}
	if len(i.Prices) > 0 {
		data, _ := json.Marshal(i.Prices)
		params["prices"] = string(data)
	}
	if len(i.SuggestedTipAmounts) > 0 {
		var amounts []string
		for _, n := range i.SuggestedTipAmounts {
			amounts = append(amounts, strconv.Itoa(n))
		}

		data, _ := json.Marshal(amounts)
		params["suggested_tip_amounts"] = string(data)
	}
	// Bot API 8.0: Star subscriptions
	if i.SubscriptionPeriod > 0 {
		params["subscription_period"] = strconv.Itoa(i.SubscriptionPeriod)
	}
	// Bot API 8.0: Business connection support
	if i.BusinessConnectionID != "" {
		params["business_connection_id"] = i.BusinessConnectionID
	}
	return params
}

// Price represents a portion of the price for goods or services.
type Price struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

// Currency contains information about supported currency for payments.
type Currency struct {
	Code         string      `json:"code"`
	Title        string      `json:"title"`
	Symbol       string      `json:"symbol"`
	Native       string      `json:"native"`
	ThousandsSep string      `json:"thousands_sep"`
	DecimalSep   string      `json:"decimal_sep"`
	SymbolLeft   bool        `json:"symbol_left"`
	SpaceBetween bool        `json:"space_between"`
	Exp          int         `json:"exp"`
	MinAmount    interface{} `json:"min_amount"`
	MaxAmount    interface{} `json:"max_amount"`
}

func (c Currency) FromTotal(total int) float64 {
	return float64(total) / math.Pow(10, float64(c.Exp))
}

func (c Currency) ToTotal(total float64) int {
	return int(total) * int(math.Pow(10, float64(c.Exp)))
}

// CreateInvoiceLink creates a link for a payment invoice.
func (b *Bot) CreateInvoiceLink(i Invoice) (string, error) {
	data, err := b.Raw("createInvoiceLink", i.params())
	if err != nil {
		return "", err
	}

	var resp struct {
		Result string
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", wrapError(err)
	}
	return resp.Result, nil
}

// RefundStars returns a successful payment in Telegram Stars.
func (b *Bot) RefundStars(to Recipient, chargeID string) error {
	params := map[string]string{
		"user_id":                    to.Recipient(),
		"telegram_payment_charge_id": chargeID,
	}

	_, err := b.Raw("refundStarPayment", params)
	if err != nil {
		return err
	}

	return nil
}

// EditUserStarSubscription allows the bot to cancel or re-enable extension of a subscription paid in Telegram Stars.
func (b *Bot) EditUserStarSubscription(user Recipient, chargeID string, isCanceled bool) error {
	params := map[string]string{
		"user_id":                    user.Recipient(),
		"telegram_payment_charge_id": chargeID,
		"is_canceled":                strconv.FormatBool(isCanceled),
	}

	_, err := b.Raw("editUserStarSubscription", params)
	return err
}

// VerifyUser verifies a user with the specified level and requirements.
func (b *Bot) VerifyUser(user Recipient, level VerificationLevel, requirements *VerificationRequirements) error {
	params := map[string]string{
		"user_id": user.Recipient(),
		"level":   string(level),
	}

	if requirements != nil {
		if data, err := json.Marshal(requirements); err == nil {
			params["requirements"] = string(data)
		}
	}

	_, err := b.Raw("verifyUser", params)
	return err
}

// VerifyChat verifies a chat with the specified level and requirements.
func (b *Bot) VerifyChat(chat *Chat, level VerificationLevel, requirements *VerificationRequirements) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"level":   string(level),
	}

	if requirements != nil {
		if data, err := json.Marshal(requirements); err == nil {
			params["requirements"] = string(data)
		}
	}

	_, err := b.Raw("verifyChat", params)
	return err
}

// RemoveUserVerification removes verification from a user.
func (b *Bot) RemoveUserVerification(user Recipient) error {
	params := map[string]string{
		"user_id": user.Recipient(),
	}

	_, err := b.Raw("removeUserVerification", params)
	return err
}

// RemoveChatVerification removes verification from a chat.
func (b *Bot) RemoveChatVerification(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("removeChatVerification", params)
	return err
}
