package telebot

import "encoding/json"

type Currency struct {
	Code         string `json:"code"`
	Title        string `json:"title"`
	Symbol       string `json:"symbol"`
	Native       string `json:"native"`
	ThousandsSep string `json:"thousands_sep"`
	DecimalSep   string `json:"decimal_sep"`
	SymbolLeft   bool   `json:"symbol_left"`
	SpaceBetween bool   `json:"space_between"`
	Exp          int    `json:"exp"`
	MinAmount    string `json:"min_amount"`
	MaxAmount    string `json:"max_amount"`
}

var SupportedCurrencies = map[string]Currency{}

func init() {
	json.Unmarshal([]byte(dataSupportedCurrenciesJSON), &SupportedCurrencies)
}
