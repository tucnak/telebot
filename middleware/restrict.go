package middleware

import tele "gopkg.in/tucnak/telebot.v3"

type RestrictConfig struct {
	Chats   []tele.Recipient
	In, Out tele.HandlerFunc
}

func Restrict(v RestrictConfig) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		if v.In == nil {
			v.In = next
		}
		if v.Out == nil {
			v.Out = next
		}
		return func(c tele.Context) error {
			for _, chat := range v.Chats {
				if chat.Recipient() == c.Sender().Recipient() {
					return v.In(c)
				}
			}
			return v.Out(c)
		}
	}
}

func Whitelist(chats ...tele.Recipient) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return Restrict(RestrictConfig{
			Chats: chats,
			In:    next,
			Out:   func(c tele.Context) error { return tele.ErrSkip },
		})(next)
	}
}

func Blacklist(chats ...tele.Recipient) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return Restrict(RestrictConfig{
			Chats: chats,
			Out:   next,
			In:    func(c tele.Context) error { return tele.ErrSkip },
		})(next)
	}
}
