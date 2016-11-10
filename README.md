# Telebot
> Telebot is a convenient wrapper to Telegram Bots API, written in Golang.

[![GoDoc](https://godoc.org/github.com/tucnak/telebot?status.svg)](https://godoc.org/github.com/tucnak/telebot) [![Travis](https://travis-ci.org/tucnak/telebot.svg?branch=master)](https://travis-ci.org/tucnak/telebot)

Bots are special Telegram accounts designed to handle messages automatically. Users can interact with bots by sending them command messages in private or group chats. These accounts serve as an interface for code running somewhere on your server.

Telebot offers a convenient wrapper to Bots API, so you shouldn't even
bother about networking at all. You may install it with

	go get github.com/tucnak/telebot

(after setting up your `GOPATH` properly).

We highly recommend you to keep your bot access token outside the code base,
preferably as an environmental variable:

	export BOT_TOKEN=<your token here>

Take a look at the minimal bot setup:
```go
package main

import (
	"log"
	"time"
	"os"
	"github.com/tucnak/telebot"
)

func main() {
	bot, err := telebot.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	messages := make(chan telebot.Message, 100)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {
		if message.Text == "/hi" {
			bot.SendMessage(message.Chat,
				"Hello, "+message.Sender.FirstName+"!", nil)
		}
	}
}
```

Previous example leaves all the logic implementation up to you. Usually you
wouldn't want to do that, so Telebot provides a handy route API. Here is an
example of it:
```go
package main

import (
	"log"
	"time"
	"os"
	"github.com/tucnak/telebot"
)

func main() {
	bot, err := telebot.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	// Routes get compiled into regular expressions.
	bot.Handle("/hi", func (context telebot.Context) {
		bot.SendMessage(context.Message.Chat, "Hi!", nil)
	})

	// Handle passes regex named groups into context variable as Args.
	bot.Handle("/greet (?P<name>[a-z]+)", func(ctx telebot.Context) {
		bot.SendMessage(ctx.Message.Chat, "Hello "+ctx.Args["name"], nil)
	})

    // Poll 100 messages at max every second.
	bot.Serve(100, 1*time.Second)
}
```

## Inline mode
As of January 4, 2016, Telegram added inline mode support for bots. Here's
a nice way to handle both incoming messages and inline queries in the meantime:

```go
package main

import (
	"log"
	"time"
	"os"
	"github.com/tucnak/telebot"
)

func main() {
	bot, err := telebot.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	bot.Messages = make(chan telebot.Message, 100)
	bot.Queries = make(chan telebot.Query, 1000)

	go messages(bot)
	go queries(bot)

	bot.Start(1 * time.Second)
}

func messages(bot *telebot.Bot) {
	for message := range bot.Messages {
		log.Printf("Received a message from %s with the text: %s\n",
			message.Sender.Username, message.Text)
	}
}

func queries(bot *telebot.Bot) {
	for query := range bot.Queries {
		log.Println("--- new query ---")
		log.Println("from:", query.From.Username)
		log.Println("text:", query.Text)

		// Create an article (a link) object to show in results.
		article := &telebot.InlineQueryResultArticle{
			Title: "Telebot",
			URL:   "https://github.com/tucnak/telebot",
			InputMessageContent: &telebot.InputTextMessageContent{
				Text:		   "Telebot is a Telegram bot framework.",
				DisablePreview: false,
			},
		}

		// Build the list of results (make sure to pass pointers!).
		results := []telebot.InlineQueryResult{article}

		// Build a response object to answer the query.
		response := telebot.QueryResponse{
			Results:	results,
			IsPersonal: true,
		}

		// Send it.
		if err := bot.AnswerInlineQuery(&query, &response); err != nil {
			log.Println("Failed to respond to query:", err)
		}
	}
}
```

## Files
Telebot lets you upload files from the file system:

```go
boom, err := telebot.NewFile("boom.ogg")
if err != nil {
	return err
}

audio := telebot.Audio{File: boom}

// Next time you send &audio, telebot won't issue
// an upload, but would re-use existing file.
err = bot.SendAudio(recipient, &audio, nil)
```

## Reply markup
Sometimes you wanna send a little complicated messages with some optional parameters. The third argument of all `Send*` methods accepts `telebot.SendOptions`, capable of defining an advanced reply markup:

```go
// Send a selective force reply message.
bot.SendMessage(user, "pong", &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			ForceReply: true,
			Selective: true,
			CustomKeyboard: [][]string{

				[]string{"1", "2", "3"},
				[]string{"4", "5", "6"},
				[]string{"7", "8", "9"},
				[]string{"*", "0", "#"},
			},
		},
	},
)
```
