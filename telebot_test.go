package telebot

import (
	"testing"
)

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
