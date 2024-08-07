package example

import (
	tele "gopkg.in/telebot.v3"
)

var users = make(map[int64]bool)

func NewBot() (*tele.Bot, error) {
	pref := tele.Settings{
		Synchronous: true,
		Offline:     true,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	b.Handle("/start", func(c tele.Context) error {
		id := c.Sender().ID
		if !users[id] {
			users[id] = true
		}
		return c.Reply("Hello!")
	})
	b.Handle(tele.OnText, func(c tele.Context) error {
		return c.Reply(c.Text())
	})

	return b, nil
}
