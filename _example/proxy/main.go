package main

import (
	"gopkg.in/telebot.v3"
	"log"
	"time"
)

func main() {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  "token",
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		Proxy:  &telebot.Proxy{Address: "ip:port"},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/hello", func(c telebot.Context) error {
		return c.Send("Hello!")
	})

	bot.Start()
}
