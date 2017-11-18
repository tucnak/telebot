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

	_, err := NewBot(token)
	if err != nil {
		t.Fatal("couldn't create bot:", err)
	}
}

func TestRecipient(_ *testing.T) {
	bot := Bot{}
	bot.Send(&User{}, "")
	bot.Send(&Chat{}, "")
}

func TestFile(t *testing.T) {
	file, err := NewFile("telebot.go")
	if err != nil {
		t.Fatal(err)
	}

	if file.Exists() {
		t.Fatal("Newly created file can't exist on Telegram servers!")
	}

	file.FileID = "magic"

	if !file.Exists() {
		t.Fatal("File with defined FileID is supposed to exist, fail.")
	}

	if file.Local() != "telebot.go" {
		t.Fatal("File doesn't preserve its original filename.")
	}
}
