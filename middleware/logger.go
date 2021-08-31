package middleware

import (
	"encoding/json"
	"log"

	tele "gopkg.in/tucnak/telebot.v3"
)

func Logger(logger *log.Logger) tele.MiddlewareFunc {
	if logger == nil {
		logger = log.Default()
	}
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			data, _ := json.MarshalIndent(c.Update(), "", "  ")
			logger.Println(string(data))
			return next(c)
		}
	}
}
