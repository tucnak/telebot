package middleware

import tele "gopkg.in/tucnak/telebot.v3"

func AutoRespond() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Callback() != nil {
				defer c.Respond()
			}
			return next(c)
		}
	}
}

func IgnoreVia() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Message().Via != nil {
				return nil
			}
			return next(c)
		}
	}
}
