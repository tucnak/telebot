package middleware

import (
	"github.com/sirupsen/logrus"
	tele "gopkg.in/tucnak/telebot.v3"
)

type LoggerFieldsFunc = func(tele.Context) logrus.Fields

func Logger(logger *logrus.Logger, fieldsFunc ...LoggerFieldsFunc) tele.MiddlewareFunc {
	if logger == nil {
		logger = logrus.New()
	}

	var f LoggerFieldsFunc
	if len(fieldsFunc) > 0 && fieldsFunc[0] != nil {
		f = fieldsFunc[0]
	} else {
		f = func(c tele.Context) logrus.Fields { return nil }
	}

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			data := c.Data()
			if clb := c.Callback(); clb != nil {
				data = clb.Unique + "|" + data
			} else if data == "" {
				data = c.Text()
			}

			logger.WithFields(f(c)).Info(data)
			return next(c)
		}
	}
}

func DefaultLogger() tele.MiddlewareFunc {
	return Logger(logrus.New(), func(c tele.Context) logrus.Fields {
		sender := c.Sender()
		if sender == nil {
			return nil
		}

		return logrus.Fields{
			// Default set of fields
			"sender": sender.Recipient(),
		}
	})
}
