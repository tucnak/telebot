package layout

import (
	tele "gopkg.in/telebot.v3"
)

// DefaultLayout is a simplified layout instance with pre-defined locale by default.
type DefaultLayout struct {
	locale string
	lt     *Layout

	Config
}

// Settings returns layout settings.
func (dlt *DefaultLayout) Settings() tele.Settings {
	return dlt.lt.Settings()
}

// Text wraps localized layout function Text using your default locale.
func (dlt *DefaultLayout) Text(k string, args ...interface{}) string {
	return dlt.lt.TextLocale(dlt.locale, k, args...)
}

// Callback returns a callback endpoint used to handle buttons.
func (dlt *DefaultLayout) Callback(k string) tele.CallbackEndpoint {
	return dlt.lt.Callback(k)
}

// Button wraps localized layout function Button using your default locale.
func (dlt *DefaultLayout) Button(k string, args ...interface{}) *tele.Btn {
	return dlt.lt.ButtonLocale(dlt.locale, k, args...)
}

// Markup wraps localized layout function Markup using your default locale.
func (dlt *DefaultLayout) Markup(k string, args ...interface{}) *tele.ReplyMarkup {
	return dlt.lt.MarkupLocale(dlt.locale, k, args...)
}

// Result wraps localized layout function Result using your default locale.
func (dlt *DefaultLayout) Result(k string, args ...interface{}) tele.Result {
	return dlt.lt.ResultLocale(dlt.locale, k, args...)
}
