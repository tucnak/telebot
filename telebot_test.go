package telebot

import (
	"fmt"
	"os"
	"testing"
)

func TestBot(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	token := os.Getenv("TELEBOT_SECRET")
	if token == "" {
		fmt.Println("ERROR: " +
			"In order to test telebot functionality, you need to set up " +
			"TELEBOT_SECRET environmental variable, which represents an API " +
			"key to a Telegram bot.\n")
		t.Fatal("Could't find TELEBOT_SECRET, aborting.")
	}

	_, err := NewBot(Settings{Token: token})
	if err != nil {
		t.Fatal("couldn't create bot:", err)
	}
}

func TestRecipient(t *testing.T) {
	token := os.Getenv("TELEBOT_SECRET")
	if token == "" {
		fmt.Println("ERROR: " +
			"In order to test telebot functionality, you need to set up " +
			"TELEBOT_SECRET environmental variable, which represents an API " +
			"key to a Telegram bot.\n")
		t.Fatal("Could't find TELEBOT_SECRET, aborting.")
	}

	bot, err := NewBot(Settings{Token: token})
	if err != nil {
		t.Fatal("couldn't create bot:", err)
	}

	bot.Send(&User{}, "")
	bot.Send(&Chat{}, "")
}

func TestFile(t *testing.T) {
	file := FromDisk("telebot.go")

	if file.InCloud() {
		t.Fatal("Newly created file can't exist on Telegram servers!")
	}

	file.FileID = "magic"

	if file.FileLocal != "telebot.go" {
		t.Fatal("File doesn't preserve its original filename.")
	}
}
