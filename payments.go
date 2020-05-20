package telebot

import (
	"encoding/json"
	"math"
)

type Invoice struct {
	// Product name, 1-32 characters.
	Title string `json:"title"`

	// Product description, 1-255 characters.
	Description string `json:"description"`

	// Custom payload, required, 1-128 bytes.
	Payload string `json:"payload"`

	// Unique deep-linking parameter that can be used to
	// generate this invoice when used as a start parameter.
	Start string `json:"start_parameter"`

	// Provider token to use.
	Token string `json:"provider_token"`

	Currency string `json:"currency"`

	Prices []Price `json:"prices"`

	ProviderData string `json:"provider_data"`

	// Processing photo_url, photo_size, photo_width, photo_height fields.
	Photo *Photo

	NeedName            bool `json:"need_name"`
	NeedPhoneNumber     bool `json:"need_phone_number"`
	NeedEmail           bool `json:"need_email"`
	NeedShippingAddress bool `json:"need_shipping_address"`

	SendPhone bool `json:"send_phone_number_to_provider"`
	SendEmail bool `json:"send_email_to_provider"`

	IsFlexible bool `json:"is_flexible"`
}

type Price struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

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

var SupportedCurrencies = map[string]Currency{}

func init() {
	err := json.Unmarshal([]byte(dataSupportedCurrenciesJSON), &SupportedCurrencies)
	if err != nil {
		panic(err)
	}
}
