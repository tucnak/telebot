package telebot

import (
	"time"
)

type Bot struct {
	Token string

	// Bot as `User` on API level.
	Identity User

	listeners []Listener
}

type Listener func(*Bot, Message)

func (b *Bot) Listen(interval time.Duration) {
	updates := make(chan Update, 1000)
	var latest_update int

	pulse := time.NewTicker(interval)
	go func() {
		for range pulse.C {
			go api_getUpdates(b.Token,
				latest_update+1,
				updates)
		}
	}()

	for update := range updates {
		if update.Id > latest_update {
			latest_update = update.Id
		}

		for _, ear := range b.listeners {
			go ear(b, update.Payload)
		}
	}
}

func (b *Bot) AddListener(ear Listener) {
	b.listeners = append(b.listeners, ear)
}
