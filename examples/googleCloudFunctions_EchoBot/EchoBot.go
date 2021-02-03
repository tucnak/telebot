package EchoTelegramBot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
)

// EchoBot receive news updates and replies to.
func EchoBot(w http.ResponseWriter, r *http.Request) {

	// First setup the bot
	bot, err := initializeBot()
	if err != nil {
		log.Fatal(err)
	}

	// Create a handle to respond any text message
	bot.Handle(tb.OnText, func(msg *tb.Message) {
		sendSticker(msg, bot) // Optional
		bot.Send(msg.Chat, msg.Text)
	})

	// Parse and process every new update received
	update := parseUpdate(w, r)
	bot.ProcessUpdate(update)

	fmt.Fprint(w, "\"status\": 200")
}


func initializeBot() (*tb.Bot, error) {

	// Remember to put any sensible data in env.yaml file
	bot, err := tb.NewBot(tb.Settings{
		Token:       os.Getenv("TELEBOT_ECHO_BOT_TOKEN"),
		Synchronous: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	return bot, err
}

func parseUpdate(w http.ResponseWriter, r *http.Request) tb.Update {

	var update tb.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("json.NewDecoder: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	return update
}

// This is totally optional and is just for
// demonstration. Use if you would to.
func sendSticker(msg *tb.Message, bot *tb.Bot) (*tb.Message, error) {
	// To better understand what is going on check the following documentation
	// https://core.telegram.org/bots/api#stickers
	// https://core.telegram.org/bots/api#file

	// You can get the FileID and UniqueID of a sticker by send one to https://t.me/JsonDumpBot
	sticker := &tb.Sticker{
		File: tb.File{
			FileID:   "CAACAgIAAxkBAAEEruRgF0Do01rTyFhM_nTG3GlnBZECzwACIgADTlzSKWF0vv5zFvwUHgQ",
			UniqueID: "AgADIgADTlzSKQ",
		},
		Width:    512,
		Height:   512,
		Animated: true,
	}
	return bot.Send(msg.Chat, sticker)
}
