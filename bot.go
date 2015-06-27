package telebot

import (
	"time"
)

type Bot struct {
	Token string

	// Bot as `User` on API level.
	Identity User
}

func (b *Bot) Listen(subscription chan<- Message, interval time.Duration) {
	updates := make(chan Update)
	pulse := time.NewTicker(interval)
	latest_update := 0

	go func() {
		for range pulse.C {
			go api_getUpdates(b.Token,
				latest_update+1,
				updates)
		}
	}()

	go func() {
		for update := range updates {
			if update.Id > latest_update {
				latest_update = update.Id
			}

			subscription <- update.Payload
		}
	}()
}

func (b *Bot) SendMessage(recipient User, message string) {
	go api_sendMessage(b.Token, recipient, message)
}
