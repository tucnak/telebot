package telebot

import (
	"log"
	"testing"
	"time"
)

const TESTING_TOKEN = "107177593:AAHBJfF3nv3pZXVjXpoowVhv_KSGw56s8zo"

func TestCreate(t *testing.T) {
	_, err := Create(TESTING_TOKEN)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	bot, err := Create(TESTING_TOKEN)
	if err != nil {
		t.Fatal(err)
	}

	bot.AddListener(func(bot *Bot, message Message) {
		if message.Text == "/hi" {
			bot.SendMessage(message.Sender,
				"Hello, "+message.Sender.FirstName+"!")
		}
	})

	log.Println("Listening...")
	bot.Listen(1 * time.Second)
}
