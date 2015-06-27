package telebot

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestTelebot(t *testing.T) {
	token := os.Getenv("TELEBOT_SECRET")
	if token == "" {
		fmt.Println("ERROR: " +
			"In order to test telebot functionality, you need to set up " +
			"TELEBOT_SECRET environmental variable, which represents an API " +
			"key to a Telegram bot.\n")
		t.Fatal("Could't find TELEBOT_SECRET, aborting.")
	}

	bot, err := Create(token)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Uncomment when Telegram fixes behavior for self-messaging

	/*messages := make(chan Message)

	intelligence := "welcome to the jungle"

	bot.SendMessage(bot.Identity, intelligence)
	bot.Listen(messages, 1*time.Second)

	select {
	case message := <-messages:
		{
			if message.Text != intelligence {
				t.Error("Self-handshake failed.")
			}
		}

	case <-time.After(5 * time.Second):
		t.Error("Self-handshake test took too long, aborting.")
	}*/
}
