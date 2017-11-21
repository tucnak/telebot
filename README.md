# Telebot
>"I never knew creating bots in Telegram was so _easy_!"

[![GoDoc](https://godoc.org/gopkg.in/tucnak/telebot.v2?status.svg)](https://godoc.org/gopkg.in/tucnak/telebot.v2)
[![Travis](https://travis-ci.org/tucnak/telebot.svg?branch=v2)](https://travis-ci.org/tucnak/telebot)

```bash
go get gopkg.in/tucnak/telebot.v2
```

Telebot is a bot framework for Telegram Bots API. This package provides a super convenient API
for command routing, message and inline query requests, as well as callbacks. Actually, I went a
couple steps further and instead of making a 1:1 API wrapper I focused on the beauty of API and
bot performance. All the methods of telebot API are _extremely_ easy to remember and later, get
used to. Telebot is agnostic to the source of updates as long as it implements the Poller interface.
Poller means you can plug your telebot into virtually any bot infrastructure, if you have any. Also,
consider Telebot a highload-ready solution. I'll soon benchmark the most popular actions and if
necessary, optimize against them without sacrificing API quality.

Take a look at the minimal telebot setup:
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

    b.Handle("/hello", func(m *tb.Message) {
        b.Send(m.From, "hello world")
    })

    b.Start()
}
```

Simple, innit? Telebot's routing system takes care of deliviering updates
to their "endpoints", so in order to get handle any meaningful event,
all you have to do is just plug your handler to one of them endpoints
and you're ready to go! You might want to switch-case handle more specific
scenarios later.

```go
b, _ := tb.NewBot(settings)

b.Handle("/help", func (m *Message) {
    // help command handler
})

b.Handle(tb.OnChannelPost, func (m *Message) {
    // channel post messages only
})

b.Handle(tb.Callback, func (c *Callback) {
    // incoming bot callbacks
})
```

Moreover, this API is completely extensible, so new handy endpoints might
appear in the following minor versions of this package.

## Inline mode
Docs TBA.

## Files
Telebot allows to both upload (from disk / by URL) and download (from Telegram)
and files in bot's scope. Telegram allows files up to 20 MB in size.

```go
a := &tb.Audio{File: tb.FromDisk("file.ogg")}

fmt.Println(a.OnDisk()) // true
fmt.Println(a.InCloud()) // false

// Will upload the file from disk and send it to recipient
bot.Send(recipient, a)

// Next time you'll be sending this very *Audio, Telebot won't
// re-upload the same file but rather utilize its Telegram FileID
bot.Send(otherRecipient, a)

fmt.Println(a.OnDisk()) // true
fmt.Println(a.InCloud()) // true
fmt.Println(a.FileID) // <telegram file id: ABC-DEF1234ghIkl-zyx57W2v1u123ew11>
```

You might want to save certain files in order to avoid re-upploading. Feel free
to marshal them into whatever format, `File` only contain public fields, so no
data will ever be lost.

TBA.
