package layout

import (
	tele "gopkg.in/telebot.v3"
)

// LocaleFunc is the function used to fetch the locale of the recipient.
// Returned locale will be remembered and linked to the corresponding context.
type LocaleFunc func(tele.Recipient) string

// Middleware builds a telebot middleware to make localization work.
//
// Usage:
//		b.Use(lt.Middleware("en", func(r tele.Recipient) string {
//			loc, _ := db.UserLocale(r.Recipient())
//			return loc
//		}))
//
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

			defer func() {
				lt.mu.Lock()
				delete(lt.ctxs, c)
				lt.mu.Unlock()
			}()

			return next(c)
		}
	}
}

// Middleware wraps ordinary layout middleware with your default locale.
func (dlt *DefaultLayout) Middleware() tele.MiddlewareFunc {
	return dlt.lt.Middleware(dlt.locale)
}
