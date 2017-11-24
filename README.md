# Telebot
>"I never knew creating Telegram bots could so _sexy_!"

[![GoDoc](https://godoc.org/gopkg.in/tucnak/telebot.v2?status.svg)](https://godoc.org/gopkg.in/tucnak/telebot.v2)
[![Travis](https://travis-ci.org/tucnak/telebot.svg?branch=v2)](https://travis-ci.org/tucnak/telebot)

```bash
go get gopkg.in/tucnak/telebot.v2
```

Telebot is a bot framework for [Telegram](https://telegram.org) [Bot API](https://core.telegram.org/bots/api).
This package provides the best of its kind API for command routing, inline query requests and keyboards, as well
as callbacks. Actually, I went a couple steps further, so instead of making a 1:1 API wrapper I chose to focus on
the beauty of API and performance. All the methods of telebot API are _extremely_ easy to memorize and get
used to. Telebot is agnostic to the source of updates as long as the source implements the `Poller` interface.
`Poller` means you can plug your telebot into virtually any existing bot infrastructure, if you have any. Also,
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
        b.Send(m.Sender, "hello world")
    })

    b.Start()
}
```

Simple, innit? Telebot's routing system takes care of deliviering updates
to their endpoints, so in order to get to handle any meaningful event,
all you have to do is just plug your function to one of them endpoints
and you're ready to go!

```go
b, _ := tb.NewBot(settings)

b.Handle(tb.OnText, func(m *Message) {
    // all text messages that weren't captured
    // by existing handlers
}

b.Handle(tb.OnPhoto, func(m *Message) {
    // photos only
}

b.Handle("/help", func (m *Message) {
    // help command handler
})

b.Handle(tb.OnChannelPost, func (m *Message) {
    // channel post messages only
})

b.Handle(tb.Callback, func (c *Callback) {
    // incoming bot callbacks that weren't
    // captured by specific callback handlers.
})
```

Now there's a dozen of supported endpoints (see package consts). Let me know
if you'd like to see some endpoint or endpoint idea implemented. This system
is completely extensible, so I can introduce them without braking
backwards-compatibity.

## Message CRUD: `Send()`, `Edit()`, `Delete()`
These are the three most important functions for manipulating Telebot messages.
`Send()` takes a Recipient (could be user, chat, channel) and a Sendable. All
telebot-provided media types (Photo, Audio, Video, etc.) are Sendable.

```go
// Sendable is any object that can send itself.
//
// This is pretty cool, since it lets bots implement
// custom Sendables for complex kind of media or
// chat objects spanning across multiple messages.
type Sendable interface {
    Send(*Bot, Recipient, *SendOptions) (*Message, error)
}
```

If you want to edit some existing message, you don't really need to store the
original `*Message` object. In fact, upon edit, Telegram only requires two IDs:
ChatID and MessageID. And it doesn't really require the whole Message. Also you
might want to store references to certain messages in the database, so for me it
made sense for *any* Go struct to be editable as Telegram message, to implement
Editable interface:
```go
// Editable is an interface for all objects that
// provide "message signature", a pair of 32-bit
// message ID and 64-bit chat ID, both required
// for edit operations.
//
// Use case: DB model struct for messages to-be
// edited with, say two collums: msg_id,chat_id
// could easily implement MessageSig() making
// instances of stored messages editable.
type Editable interface {
	// MessageSig is a "message signature".
	//
	// For inline messages, return chatID = 0.
	MessageSig() (messageID int, chatID int64)
}
```

For example, `Message` type is Editable. Here is an implementation of `StoredMessage`
type, provided by telebot:
```go
// StoredMessage is an example struct suitable for being
// stored in the database as-is or being embedded into
// a larger struct, which is often the case (you might
// want to store some metadata alongside, or might not.)
type StoredMessage struct {
	MessageID int   `sql:"message_id" json:"message_id"`
	ChatID    int64 `sql:"chat_id" json:"chat_id"`
}

func (x StoredMessage) MessageSig() (int, int64) {
	return x.MessageID, x.ChatID
}
```

Why bother at all? Well, it allows you to do things like this:
```go
// just two integer columns in the database
var msgs []StoredMessage
db.Find(&msgs) // gorm syntax

for _, msg := range msgs {
    bot.Edit(&msg, "Updated text.")
    // or
    bot.Delete(&msg)
}
```

I find it incredibly neat. Worth noting, at this point of time there exists
another method in the Edit family, `EditCaption()` which is of a pretty
rare use, so I didn't bother including it to `Edit()`, which would inevitably
lead to unnecessary complications.
```go
var m *Message

// change caption of a photo, audio, etc.
bot.EditCaption(m, "new caption")
```

## Inline mode
Docs TBA.

## Files
>Telegram allows files up to 20 MB in size.

Telebot allows to both upload (from disk / by URL) and download (from Telegram)
and files in bot's scope. Also, sending any kind of media with a File created
from disk will upload the file to Telegram automatically:
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

You might want to save certain `File`s in order to avoid re-uploading. Feel free
to marshal them into whatever format, `File` only contain public fields, so no
data will ever be lost.

TBA.
