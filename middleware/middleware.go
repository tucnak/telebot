package middleware

import (
	"errors"

	tele "gopkg.in/telebot.v3"
)

// AutoRespond returns a middleware that automatically responds
// to every callback.
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

// IgnoreVia returns a middleware that ignores all the
// "sent via" messages.
func IgnoreVia() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if msg := c.Message(); msg != nil && msg.Via != nil {
				return nil
			}
			return next(c)
		}
	}
}

type RecoverFunc = func(error, tele.Context)

// Recover returns a middleware that recovers a panic happened in
// the handler.
func Recover(onError ...RecoverFunc) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			var f RecoverFunc
			if len(onError) > 0 {
				f = onError[0]
			} else {
				f = func(err error, c tele.Context) {
					c.Bot().OnError(err, c)
				}
			}

			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						f(err, c)
					} else if s, ok := r.(string); ok {
						f(errors.New(s), c)
					}
				}
			}()

			return next(c)
		}
	}
}
