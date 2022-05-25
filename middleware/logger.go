package middleware

import (
	"encoding/json"
	"log"

	tele "github.com/TGeniusFamily/GOFSMtelebot"
)

func Logger(logger ...*log.Logger) tele.MiddlewareFunc {
	var l *log.Logger
	if len(logger) > 0 {
		l = logger[0]
	} else {
		l = log.Default()
	}

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			data, _ := json.MarshalIndent(c.Update(), "", "  ")
			l.Println(string(data))
			return next(c)
		}
	}
}
