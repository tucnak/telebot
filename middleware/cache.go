package middleware

import (
	"log"

	tele "gopkg.in/telebot.v3"
)

// CacheContext returns a middleware that store context to cache and retreive data from cache to context
// If no custom cache provided, context store will be rested each iteration.
func CacheContext(cache tele.Cache) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(ctx tele.Context) error {
			for _, key := range cache.Keys() {
				value, err := cache.Get(key)
				if err != nil {
					log.Printf("err: %s was happened, %s -> %s was not got from cache", err, value, key)

					return err
				}

				ctx.Set(key, value)
			}

			defer func() {
				for _, key := range ctx.Keys() {
					value := ctx.Get(key)
					err := cache.Set(key, value)
					if err != nil {
						log.Printf("err: %s was happened, %s -> %s was not stored from context", err, value, key)
					}
				}
			}()

			return next(ctx)
		}
	}
}
