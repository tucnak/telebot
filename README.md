# Telebot
>Telebot is a convenient wrapper to Telegram Bots API, written in Golang.

Bots are special Telegram accounts designed to handle messages automatically. Users can interact with bots by sending them command messages in private or group chats. These accounts serve as an interface for code running somewhere on your server.

Telebot offers a convenient wrapper to Bots API, so you shouldn't even care about networking at all.

```go
import (
    "time"
    "github.com/tucnak/telebot"
)

func main() {
    bot, err := telebot.Create("SECRET TOKEN")
    if err != nil {
        return
    }

    messages := make(chan telebot.Message)
    bot.Listen(messages, 1*time.Second)

    for message := range messages {
        if message.Text == "/hi" {
            bot.SendMessage(message.Chat,
                "Hello, "+message.Sender.FirstName+"!")
        }
    }
}
```

## Features
- Listening to messages from users / group chats
- Replying to users, depending on scope

## TODO
- Message forwarding
- Sophiesticated content types (Document, Contact, Video, etc)
- Reply keyboard styles
