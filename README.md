# Telebot
>Telebot is a Telegram bot framework in Go.

[![GoDoc](https://godoc.org/github.com/tucnak/telebot?status.svg)](https://godoc.org/github.com/tucnak/telebot)
[![Travis](https://travis-ci.org/tucnak/telebot.svg?branch=master)](https://travis-ci.org/tucnak/telebot)

Bots are special Telegram accounts designed to handle messages automatically.
Users can interact with bots by sending them command messages in private or
via group chats / channels. These accounts serve as an interface to your code.

Telebot offers a pretty convenient interface to Bots API and uses default HTTP
client. Ideally, you wouldn't need to worry about actual networking at all.

```bash
go get gopkg.in/tucnak/telebot.v2
```

(after setting up your `GOPATH` properly).

We highly recommend you to keep your bot access token outside the code base,
preferably in an environmental variable:

```bash
export BOT_TOKEN=<your token here>
```

Take a look at a minimal functional bot setup:
```go
import (
    "time"
    tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
    b, err := tb.NewBot(tb.Settings{
        Token: "TOKEN_HERE",
        Poller: &tb.LongPoller{10 * time.Second},
    })

    if err != nil {
        return
    }

    b.Handle(tb.OnMessage, func(m *tb.Message) {
        b.Send(m.From, "hello world")
    }

    b.Start()
}
```

## Inline mode
As of January 4, 2016, Telegram added inline mode support for bots. Here's
a nice way to handle both incoming messages and inline queries:

```go
import (
    "time"
    tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
    b, err := tb.NewBot(tb.Settings{
        Token: "TOKEN_HERE",
        Poller: &tb.LongPoller{10 * time.Second},
    })

    if err != nil {
        return
    }

    b.Handle(tb.OnMessage, func(m *tb.Message) {
        b.Send(m.From, "hello world")
    }

    b.Handle(tb.OnQuery, func(q *tb.Query) {
        b.Answer(q, ...)
    }

    b.Start()
}
```

## Files
Telebot allows to both upload and download certain files.

```go
a := &tb.Audio{File: tb.FromDisk("file.ogg")}

fmt.Println(a.OnDisk()) // true
fmt.Println(a.InCloud()) // false

// Next time you'll be sending this very *Audio, Telebot won't
// re-upload the same file but rather use the copy from the
// server.
bot.Send(recipient, a)

fmt.Println(a.OnDisk()) // true
fmt.Println(a.InCloud()) // true
fmt.Println(a.FileID) // <telegram file id: ABC-DEF1234ghIkl-zyx57W2v1u123ew11>
```

You might want to save certain files in order to avoid re-upploading. Feel free
to marshal them into whatever format, Files only contain public fields, so no
data will be lost.

## Reply markup
```go
// Send a selective force reply message.
bot.Send(user, "pong", &tb.ReplyMarkup{
    ForceReply: true,
    Selective: true,

    ReplyKeyboard: keys,
})
```
