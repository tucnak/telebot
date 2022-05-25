package middleware

import tele "github.com/TGeniusFamily/GOFSMtelebot"

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
			if msg := c.Message(); msg != nil && msg.Via != nil {
				return nil
			}
			return next(c)
		}
	}
}
