# Telebot
>Telebot is a convenient wrapper to Telegram Bots API, written in Golang.

[![GoDoc](https://godoc.org/github.com/tucnak/telebot?status.svg)](https://godoc.org/github.com/tucnak/telebot)
[![Travis](https://travis-ci.org/tucnak/telebot.svg?branch=master)](https://travis-ci.org/tucnak/telebot)

Bots are special Telegram accounts designed to handle messages automatically. Users can interact with bots by sending them command messages in private or group chats. These accounts serve as an interface for code running somewhere on your server.

Telebot offers a convenient wrapper to Bots API, so you shouldn't even care about networking at all. Here is an example "helloworld" bot, written with telebot:
```go
import (
    "time"
    "github.com/tucnak/telebot"
)

func main() {
    bot, err := telebot.NewBot("SECRET TOKEN")
    if err != nil {
        return
    }

    // routes are compiled as regexps
    bot.Handle("/hi", func (context telebot.Context) {
	   bot.SendMessage(context.Message.Chat, "Hi!", nil)
    })

    // named parameters found in routes will get injected in the controller
	bot.Handle("/greet (?P<name>[a-z]+) (?P<last_name>[a-z]+)", func(context telebot.Context) {
	   bot.SendMessage(context.Message.Chat, fmt.Sprintf("Hello %s, %s", context.Args["last_name"], context.Args["name"]), nil)
	})

	bot.Serve()
}
```

## Inline mode
As of January 4, 2016, Telegram added inline mode support for bots. Telebot does support inline mode in a fancy manner. Here's a nice way to handle both incoming messages and inline queries:
```go
import (
	"log"
    "time"

    "github.com/tucnak/telebot"
)

var bot *telebot.Bot

func main() {
    if newBot, err := telebot.NewBot("SECRET TOKEN"); err != nil {
        return
    } else {
		// shadowing, remember?
		bot = newBot
	}

	bot.Messages = make(chan telebot.Message, 1000)
	bot.Queries = make(chan telebot.Query, 1000)

	go messages()
	go queries()

    bot.Start(1 * time.Second)
}

func messages() {
	for message := range bot.Messages {
		// ...
	}
}

func queries() {
	for query := range bot.Queries {
		log.Println("--- new query ---")
		log.Println("from:", query.From)
		log.Println("text:", query.Text)

		// There you build a slice of let's say, article results:
		results := []telebot.Result{...}

		// And finally respond to the query:
		if err := bot.Respond(query, results); err != nil {
			log.Println("ouch:", err)
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
