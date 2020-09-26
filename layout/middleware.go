package layout

import (
	tele "gopkg.in/tucnak/telebot.v3"
)

type LocaleFunc func(tele.Recipient) string

func (lt *Layout) Middleware(defaultLocale string, localeFunc ...LocaleFunc) tele.MiddlewareFunc {
	var f LocaleFunc
	if len(localeFunc) > 0 {
		f = localeFunc[0]
	}

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			locale := defaultLocale
			if f != nil {
				if l := f(c.Sender()); l != "" {
					locale = l
				}
			}

			lt.SetLocale(c, locale)
			return next(c)
		}
	}
}
