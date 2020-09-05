package telebot

import "time"

// Layout allows you to control the bot's content,
// that is usually depends on the user context.
type Layout interface {
	With(Context) Layout

	Get(k string) string

	Int(k string) int

	Int64(k string) int64

	Float(k string) float64

	Duration(k string) time.Duration

	Text(k string, args ...interface{}) string

	Markup(k string, args ...interface{}) *ReplyMarkup

	Locale() string
}
