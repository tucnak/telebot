package telebot

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	photoID = "AgACAgIAAxkDAAIBV16Ybpg7l2jPgMUiiLJ3WaQOUqTrAAJorjEbh2TBSPSOinaCHfydQO_pki4AAwEAAwIAA3kAA_NQAAIYBA"
)

var (
	// required to test send and edit methods
	token     = os.Getenv("TELEBOT_SECRET")
	chatID, _ = strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	userID, _ = strconv.Atoi(os.Getenv("USER_ID"))

	b, _ = newTestBot()      // cached bot instance to avoid getMe method flooding
	to   = &Chat{ID: chatID} // to chat recipient for send and edit methods
	user = &User{ID: userID} // to user recipient for some special cases
)

func defaultSettings() Settings {
	return Settings{Token: token}
}

func newTestBot() (*Bot, error) {
	return NewBot(defaultSettings())
}

func TestNewBot(t *testing.T) {
	var pref Settings
	_, err := NewBot(pref)
	assert.Error(t, err)

	pref.Token = "BAD TOKEN"
	_, err = NewBot(pref)
	assert.Error(t, err)

	pref.URL = "BAD URL"
	_, err = NewBot(pref)
	assert.Error(t, err)

	b, err := NewBot(Settings{offline: true})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, b.Me)
	assert.Equal(t, DefaultApiURL, b.URL)
	assert.Equal(t, http.DefaultClient, b.client)
	assert.Equal(t, 100, cap(b.Updates))

	pref = defaultSettings()
	client := &http.Client{Timeout: time.Minute}
	pref.URL = "http://api.telegram.org" // not https
	pref.Client = client
	pref.Poller = &LongPoller{Timeout: time.Second}
	pref.Updates = 50
	pref.offline = true

	b, err = NewBot(pref)
	assert.NoError(t, err)
	assert.Equal(t, client, b.client)
	assert.Equal(t, pref.URL, b.URL)
	assert.Equal(t, pref.Poller, b.Poller)
	assert.Equal(t, 50, cap(b.Updates))
}

func TestBotHandle(t *testing.T) {
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}

	b.Handle("/start", func(m *Message) {})
	assert.Contains(t, b.handlers, "/start")

	btn := &InlineButton{Unique: "test"}
	b.Handle(btn, func(c *Callback) {})
	assert.Contains(t, b.handlers, btn.CallbackUnique())

	assert.Panics(t, func() { b.Handle(1, func() {}) })
}

func TestBotStart(t *testing.T) {
	if token == "" {
		t.Skip("TELEBOT_SECRET is required")
	}

	// cached bot has no poller
	assert.Panics(t, func() { b.Start() })

	pref := defaultSettings()
	pref.Poller = &LongPoller{}

	b, err := NewBot(pref)
	if err != nil {
		t.Fatal(err)
	}

	// remove webhook to be sure that bot can poll
	assert.NoError(t, b.RemoveWebhook())

	go b.Start()
	b.Stop()

	tp := newTestPoller()
	go func() {
		tp.updates <- Update{Message: &Message{Text: "/start"}}
	}()

	b, err = NewBot(pref)
	assert.NoError(t, err)
	b.Poller = tp

	var ok bool
	b.Handle("/start", func(m *Message) {
		assert.Equal(t, m.Text, "/start")
		tp.done <- struct{}{}
		ok = true
	})

	go b.Start()
	<-tp.done
	b.Stop()

	assert.True(t, ok)
}

func TestBotProcessUpdate(t *testing.T) {
	b, err := NewBot(Settings{Synchronous: true, offline: true})
	if err != nil {
		t.Fatal(err)
	}

	b.Handle("/start", func(m *Message) {
		assert.Equal(t, "/start", m.Text)
	})
	b.Handle("hello", func(m *Message) {
		assert.Equal(t, "hello", m.Text)
	})
	b.Handle(OnText, func(m *Message) {
		assert.Equal(t, "text", m.Text)
	})
	b.Handle(OnPinned, func(m *Message) {
		assert.NotNil(t, m.PinnedMessage)
	})
	b.Handle(OnPhoto, func(m *Message) {
		assert.NotNil(t, m.Photo)
	})
	b.Handle(OnVoice, func(m *Message) {
		assert.NotNil(t, m.Voice)
	})
	b.Handle(OnAudio, func(m *Message) {
		assert.NotNil(t, m.Audio)
	})
	b.Handle(OnDocument, func(m *Message) {
		assert.NotNil(t, m.Document)
	})
	b.Handle(OnSticker, func(m *Message) {
		assert.NotNil(t, m.Sticker)
	})
	b.Handle(OnVideo, func(m *Message) {
		assert.NotNil(t, m.Video)
	})
	b.Handle(OnVideoNote, func(m *Message) {
		assert.NotNil(t, m.VideoNote)
	})
	b.Handle(OnContact, func(m *Message) {
		assert.NotNil(t, m.Contact)
	})
	b.Handle(OnLocation, func(m *Message) {
		assert.NotNil(t, m.Location)
	})
	b.Handle(OnVenue, func(m *Message) {
		assert.NotNil(t, m.Venue)
	})
	b.Handle(OnAddedToGroup, func(m *Message) {
		assert.NotNil(t, m.GroupCreated)
	})
	b.Handle(OnUserJoined, func(m *Message) {
		assert.NotNil(t, m.UserJoined)
	})
	b.Handle(OnUserLeft, func(m *Message) {
		assert.NotNil(t, m.UserLeft)
	})
	b.Handle(OnNewGroupTitle, func(m *Message) {
		assert.Equal(t, "title", m.NewGroupTitle)
	})
	b.Handle(OnNewGroupPhoto, func(m *Message) {
		assert.NotNil(t, m.NewGroupPhoto)
	})
	b.Handle(OnGroupPhotoDeleted, func(m *Message) {
		assert.True(t, m.GroupPhotoDeleted)
	})
	b.Handle(OnMigration, func(from, to int64) {
		assert.Equal(t, int64(1), from)
		assert.Equal(t, int64(2), to)
	})
	b.Handle(OnEdited, func(m *Message) {
		assert.Equal(t, "edited", m.Text)
	})
	b.Handle(OnChannelPost, func(m *Message) {
		assert.Equal(t, "post", m.Text)
	})
	b.Handle(OnEditedChannelPost, func(m *Message) {
		assert.Equal(t, "edited post", m.Text)
	})
	b.Handle(OnCallback, func(c *Callback) {
		if c.Data[0] != '\f' {
			assert.Equal(t, "callback", c.Data)
		}
	})
	b.Handle("\funique", func(c *Callback) {
		assert.Equal(t, "callback", c.Data)
	})
	b.Handle(OnQuery, func(q *Query) {
		assert.Equal(t, "query", q.Text)
	})
	b.Handle(OnChosenInlineResult, func(r *ChosenInlineResult) {
		assert.Equal(t, "result", r.ResultID)
	})
	b.Handle(OnCheckout, func(pre *PreCheckoutQuery) {
		assert.Equal(t, "checkout", pre.ID)
	})
	b.Handle(OnPoll, func(p *Poll) {
		assert.Equal(t, "poll", p.ID)
	})
	b.Handle(OnPollAnswer, func(pa *PollAnswer) {
		assert.Equal(t, "poll", pa.PollID)
	})

	b.ProcessUpdate(Update{Message: &Message{Text: "/start"}})
	b.ProcessUpdate(Update{Message: &Message{Text: "/start@other_bot"}})
	b.ProcessUpdate(Update{Message: &Message{Text: "hello"}})
	b.ProcessUpdate(Update{Message: &Message{Text: "text"}})
	b.ProcessUpdate(Update{Message: &Message{PinnedMessage: &Message{}}})
	b.ProcessUpdate(Update{Message: &Message{Photo: &Photo{}}})
	b.ProcessUpdate(Update{Message: &Message{Voice: &Voice{}}})
	b.ProcessUpdate(Update{Message: &Message{Audio: &Audio{}}})
	b.ProcessUpdate(Update{Message: &Message{Document: &Document{}}})
	b.ProcessUpdate(Update{Message: &Message{Sticker: &Sticker{}}})
	b.ProcessUpdate(Update{Message: &Message{Video: &Video{}}})
	b.ProcessUpdate(Update{Message: &Message{VideoNote: &VideoNote{}}})
	b.ProcessUpdate(Update{Message: &Message{Contact: &Contact{}}})
	b.ProcessUpdate(Update{Message: &Message{Location: &Location{}}})
	b.ProcessUpdate(Update{Message: &Message{Venue: &Venue{}}})
	b.ProcessUpdate(Update{Message: &Message{Dice: &Dice{}}})
	b.ProcessUpdate(Update{Message: &Message{GroupCreated: true}})
	b.ProcessUpdate(Update{Message: &Message{UserJoined: &User{ID: 1}}})
	b.ProcessUpdate(Update{Message: &Message{UsersJoined: []User{{ID: 1}}}})
	b.ProcessUpdate(Update{Message: &Message{UserLeft: &User{}}})
	b.ProcessUpdate(Update{Message: &Message{NewGroupTitle: "title"}})
	b.ProcessUpdate(Update{Message: &Message{NewGroupPhoto: &Photo{}}})
	b.ProcessUpdate(Update{Message: &Message{GroupPhotoDeleted: true}})
	b.ProcessUpdate(Update{Message: &Message{Chat: &Chat{ID: 1}, MigrateTo: 2}})
	b.ProcessUpdate(Update{EditedMessage: &Message{Text: "edited"}})
	b.ProcessUpdate(Update{ChannelPost: &Message{Text: "post"}})
	b.ProcessUpdate(Update{ChannelPost: &Message{PinnedMessage: &Message{}}})
	b.ProcessUpdate(Update{EditedChannelPost: &Message{Text: "edited post"}})
	b.ProcessUpdate(Update{Callback: &Callback{MessageID: "inline", Data: "callback"}})
	b.ProcessUpdate(Update{Callback: &Callback{Data: "callback"}})
	b.ProcessUpdate(Update{Callback: &Callback{Data: "\funique|callback"}})
	b.ProcessUpdate(Update{Query: &Query{Text: "query"}})
	b.ProcessUpdate(Update{ChosenInlineResult: &ChosenInlineResult{ResultID: "result"}})
	b.ProcessUpdate(Update{PreCheckoutQuery: &PreCheckoutQuery{ID: "checkout"}})
	b.ProcessUpdate(Update{Poll: &Poll{ID: "poll"}})
	b.ProcessUpdate(Update{PollAnswer: &PollAnswer{PollID: "poll"}})
}

func TestBot(t *testing.T) {
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}
	if chatID == 0 {
		t.Skip("CHAT_ID is required for Bot methods test")
	}

	_, err := b.Send(to, nil)
	assert.Equal(t, ErrUnsupportedWhat, err)
	_, err = b.Edit(&Message{Chat: &Chat{}}, nil)
	assert.Equal(t, ErrUnsupportedWhat, err)

	_, err = b.Send(nil, "")
	assert.Equal(t, ErrBadRecipient, err)
	_, err = b.Forward(nil, nil)
	assert.Equal(t, ErrBadRecipient, err)

	t.Run("Send(what=Sendable)", func(t *testing.T) {
		photo := &Photo{
			File:    File{FileID: photoID},
			Caption: t.Name(),
		}

		msg, err := b.Send(to, photo)
		assert.NoError(t, err)
		assert.NotNil(t, msg.Photo)
		assert.Equal(t, photo.Caption, msg.Caption)
	})

	var msg *Message

	t.Run("Send(what=string)", func(t *testing.T) {
		msg, err = b.Send(to, t.Name())
		assert.NoError(t, err)
		assert.Equal(t, t.Name(), msg.Text)

		rpl, err := b.Reply(msg, t.Name())
		assert.NoError(t, err)
		assert.Equal(t, rpl.Text, msg.Text)
		assert.NotNil(t, rpl.ReplyTo)
		assert.Equal(t, rpl.ReplyTo, msg)
		assert.True(t, rpl.IsReply())

		fwd, err := b.Forward(to, msg)
		assert.NoError(t, err)
		assert.NotNil(t, msg, fwd)
		assert.True(t, fwd.IsForwarded())

		fwd.ID += 1 // nonexistent message
		fwd, err = b.Forward(to, fwd)
		assert.Equal(t, ErrToForwardNotFound, err)
	})

	t.Run("Edit(what=string)", func(t *testing.T) {
		msg, err = b.Edit(msg, t.Name())
		assert.NoError(t, err)
		assert.Equal(t, t.Name(), msg.Text)

		_, err = b.Edit(msg, msg.Text)
		assert.Error(t, err) // message is not modified
	})

	t.Run("Edit(what=Location)", func(t *testing.T) {
		loc := &Location{Lat: 42, Lng: 69, LivePeriod: 60}
		msg, err := b.Send(to, loc)
		assert.NoError(t, err)
		assert.NotNil(t, msg.Location)

		loc = &Location{Lat: loc.Lng, Lng: loc.Lat}
		msg, err = b.Edit(msg, *loc)
		assert.NoError(t, err)
		assert.NotNil(t, msg.Location)
	})

	t.Run("EditReplyMarkup()", func(t *testing.T) {
		markup := &ReplyMarkup{
			InlineKeyboard: [][]InlineButton{{{
				Data: "btn",
				Text: "Hi Telebot!",
			}}},
		}
		badMarkup := &ReplyMarkup{
			InlineKeyboard: [][]InlineButton{{{
				Data: strings.Repeat("*", 65),
				Text: "Bad Button",
			}}},
		}

		msg, err := b.EditReplyMarkup(msg, markup)
		assert.NoError(t, err)
		assert.Equal(t, msg.ReplyMarkup.InlineKeyboard, markup.InlineKeyboard)

		msg, err = b.EditReplyMarkup(msg, nil)
		assert.NoError(t, err)
		assert.Nil(t, msg.ReplyMarkup.InlineKeyboard)

		_, err = b.EditReplyMarkup(msg, badMarkup)
		assert.Equal(t, ErrButtonDataInvalid, err)
	})

	t.Run("Commands", func(t *testing.T) {
		orig := []Command{{
			Text:        "test",
			Description: "test command",
		}}
		assert.NoError(t, b.SetCommands(orig))

		cmds, err := b.GetCommands()
		assert.NoError(t, err)
		assert.Equal(t, orig, cmds)
	})
}
