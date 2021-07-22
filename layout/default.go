package layout

import (
	tele "gopkg.in/tucnak/telebot.v3"
)

// DefaultLayout is a simplified layout instance with pre-defined locale by default.
type DefaultLayout struct {
	lt     *Layout
	locale string
}

// Text wraps localized layout function Text using your default locale.
func (dlt *DefaultLayout) Text(k string, args ...interface{}) string {
	return dlt.lt.TextLocale(dlt.locale, k, args...)
}

// Button wraps localized layout function Button using your default locale.
func (dlt *DefaultLayout) Button(k string, args ...interface{}) *tele.Btn {
	return dlt.lt.ButtonLocale(dlt.locale, k, args...)
}

// Markup wraps localized layout function Markup using your default locale.
func (dlt *DefaultLayout) Markup(k string, args ...interface{}) *tele.ReplyMarkup {
	return dlt.lt.MarkupLocale(dlt.locale, k, args...)
}
