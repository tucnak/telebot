package telebot

import (
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	photoID = "AgACAgIAAxkDAAIBV16Ybpg7l2jPgMUiiLJ3WaQOUqTrAAJorjEbh2TBSPSOinaCHfydQO_pki4AAwEAAwIAA3kAA_NQAAIYBA"
)

var (
	// required to test send and edit methods
	chatID, _ = strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)

	b, _ = newTestBot()      // cached bot instance to avoid getMe method flooding
	to   = &Chat{ID: chatID} // to recipient for send and edit methods
)

func defaultSettings() Settings {
	return Settings{Token: os.Getenv("TELEBOT_SECRET")}
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

	b, err := newTestBot()
	assert.NoError(t, err)
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

	b, err = NewBot(pref)
	assert.NoError(t, err)
	assert.Equal(t, client, b.client)
	assert.Equal(t, pref.URL, b.URL)
	assert.Equal(t, pref.Poller, b.Poller)
	assert.Equal(t, 50, cap(b.Updates))
}

func TestBotHandle(t *testing.T) {
	b.Handle("/start", func(m *Message) {})
	assert.Contains(t, b.handlers, "/start")

	btn := &InlineButton{Unique: "test"}
	b.Handle(btn, func(c *Callback) {})
	assert.Contains(t, b.handlers, btn.CallbackUnique())

	assert.Panics(t, func() { b.Handle(1, func() {}) })
}

func TestBotStart(t *testing.T) {
	// cached bot has no poller
	assert.Panics(t, func() { b.Start() })

	pref := defaultSettings()
	pref.Poller = &LongPoller{}

	b, err := NewBot(pref)
	assert.NoError(t, err)

	// remove webhook to be sure that bot can poll
	assert.NoError(t, b.RemoveWebhook())

	time.AfterFunc(50*time.Millisecond, b.Stop)
	b.Start() // stops after some delay
	assert.Empty(t, b.stop)

	tp := &testPoller{updates: make(chan Update, 1)}
	go func() {
		tp.updates <- Update{Message: &Message{Text: "/start"}}
	}()

	b, err = NewBot(pref)
	assert.NoError(t, err)
	b.Poller = tp

	var ok bool
	b.Handle("/start", func(m *Message) {
		assert.Equal(t, m.Text, "/start")
		ok = true
	})

	time.AfterFunc(100*time.Millisecond, b.Stop)
	b.Start() // stops after some delay
	assert.True(t, ok)
}

func TestBotIncomingUpdate(t *testing.T) {
	b, err := newTestBot()
	if err != nil {
		t.Fatal(err)
	}

	tp := &testPoller{updates: make(chan Update, 1)}
	b.Poller = tp

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

	go func() {
		tp.updates <- Update{Message: &Message{Text: "/start"}}
		tp.updates <- Update{Message: &Message{Text: "/start@other_bot"}}
		tp.updates <- Update{Message: &Message{Text: "hello"}}
		tp.updates <- Update{Message: &Message{Text: "text"}}
		tp.updates <- Update{Message: &Message{PinnedMessage: &Message{}}}
		tp.updates <- Update{Message: &Message{Photo: &Photo{}}}
		tp.updates <- Update{Message: &Message{Voice: &Voice{}}}
		tp.updates <- Update{Message: &Message{Audio: &Audio{}}}
		tp.updates <- Update{Message: &Message{Document: &Document{}}}
		tp.updates <- Update{Message: &Message{Sticker: &Sticker{}}}
		tp.updates <- Update{Message: &Message{Video: &Video{}}}
		tp.updates <- Update{Message: &Message{VideoNote: &VideoNote{}}}
		tp.updates <- Update{Message: &Message{Contact: &Contact{}}}
		tp.updates <- Update{Message: &Message{Location: &Location{}}}
		tp.updates <- Update{Message: &Message{Venue: &Venue{}}}
		tp.updates <- Update{Message: &Message{GroupCreated: true}}
		tp.updates <- Update{Message: &Message{UserJoined: &User{}}}
		tp.updates <- Update{Message: &Message{UsersJoined: []User{{}}}}
		tp.updates <- Update{Message: &Message{UserLeft: &User{}}}
		tp.updates <- Update{Message: &Message{NewGroupTitle: "title"}}
		tp.updates <- Update{Message: &Message{NewGroupPhoto: &Photo{}}}
		tp.updates <- Update{Message: &Message{GroupPhotoDeleted: true}}
		tp.updates <- Update{Message: &Message{Chat: &Chat{ID: 1}, MigrateTo: 2}}
		tp.updates <- Update{EditedMessage: &Message{Text: "edited"}}
		tp.updates <- Update{ChannelPost: &Message{Text: "post"}}
		tp.updates <- Update{ChannelPost: &Message{PinnedMessage: &Message{}}}
		tp.updates <- Update{EditedChannelPost: &Message{Text: "edited post"}}
		tp.updates <- Update{Callback: &Callback{MessageID: "inline", Data: "callback"}}
		tp.updates <- Update{Callback: &Callback{Data: "callback"}}
		tp.updates <- Update{Callback: &Callback{Data: "\funique|callback"}}
		tp.updates <- Update{Query: &Query{Text: "query"}}
		tp.updates <- Update{ChosenInlineResult: &ChosenInlineResult{ResultID: "result"}}
		tp.updates <- Update{PreCheckoutQuery: &PreCheckoutQuery{ID: "checkout"}}
		tp.updates <- Update{Poll: &Poll{ID: "poll"}}
		tp.updates <- Update{PollAnswer: &PollAnswer{PollID: "poll"}}
	}()

	time.AfterFunc(100*time.Millisecond, b.Stop)
	b.Start() // stops after some delay
}

func TestBotSend(t *testing.T) {
	if chatID == 0 {
		t.Skip("CHAT_ID is required for Send method test")
	}

	_, err := b.Send(to, nil)
	assert.Equal(t, ErrUnsupportedSendable, err)

	_, err = b.Send(nil, t.Name())
	assert.Error(t, err)

	t.Run("what=string", func(t *testing.T) {
		msg, err := b.Send(to, t.Name())
		assert.NoError(t, err)
		assert.Equal(t, t.Name(), msg.Text)
	})

	t.Run("what=Sendable", func(t *testing.T) {
		photo := &Photo{
			File:    File{FileID: photoID},
			Caption: t.Name(),
		}

		msg, err := b.Send(to, photo)
		assert.NoError(t, err)
		assert.NotNil(t, msg.Photo)
		assert.Equal(t, photo.Caption, msg.Caption)
	})
}
