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

	messages := make(chan Message)
	bot.Listen(messages, 1*time.Second)

	log.Println("Listening...")

	for message := range messages {
		if message.Text == "/hi" {
			bot.SendMessage(message.Chat,
				"Hello, "+message.Sender.FirstName+"!")
		}
	}
}
