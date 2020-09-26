package telebot

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	pref.ParseMode = ModeHTML
	pref.offline = true

	b, err = NewBot(pref)
	require.NoError(t, err)
	assert.Equal(t, client, b.client)
	assert.Equal(t, pref.URL, b.URL)
	assert.Equal(t, pref.Poller, b.Poller)
	assert.Equal(t, 50, cap(b.Updates))
	assert.Equal(t, ModeHTML, b.parseMode)
}

func TestBotHandle(t *testing.T) {
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}

	b.Handle("/start", func(m *Message) {})
	assert.Contains(t, b.handlers, "/start")

	reply := ReplyButton{Text: "reply"}
	b.Handle(&reply, func(m *Message) {})

	inline := InlineButton{Unique: "inline"}
	b.Handle(&inline, func(c *Callback) {})

	btnReply := (&ReplyMarkup{}).Text("btnReply")
	b.Handle(&btnReply, func(m *Message) {})

	btnInline := (&ReplyMarkup{}).Data("", "btnInline")
	b.Handle(&btnInline, func(c *Callback) {})

	assert.Contains(t, b.handlers, btnReply.CallbackUnique())
	assert.Contains(t, b.handlers, btnInline.CallbackUnique())
	assert.Contains(t, b.handlers, reply.CallbackUnique())
	assert.Contains(t, b.handlers, inline.CallbackUnique())

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
	require.NoError(t, b.RemoveWebhook())

	go b.Start()
	b.Stop()

	tp := newTestPoller()
	go func() {
		tp.updates <- Update{Message: &Message{Text: "/start"}}
	}()

	b, err = NewBot(pref)
	require.NoError(t, err)
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
	b.Handle(OnAnimation, func(m *Message) {
		assert.NotNil(t, m.Animation)
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
	b.ProcessUpdate(Update{Message: &Message{Animation: &Animation{}}})
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

	photo := &Photo{
		File:    File{FileID: photoID},
		Caption: t.Name(),
	}
	var msg *Message

	t.Run("Send(what=Sendable)", func(t *testing.T) {
		msg, err = b.Send(to, photo)
		require.NoError(t, err)
		assert.NotNil(t, msg.Photo)
		assert.Equal(t, photo.Caption, msg.Caption)
	})

	t.Run("SendAlbum()", func(t *testing.T) {
		_, err = b.SendAlbum(nil, nil)
		assert.Equal(t, ErrBadRecipient, err)

		_, err = b.SendAlbum(to, nil)
		assert.Error(t, err)

		msgs, err := b.SendAlbum(to, Album{photo, photo})
		require.NoError(t, err)
		assert.Len(t, msgs, 2)
		assert.NotEmpty(t, msgs[0].AlbumID)
	})

	t.Run("EditCaption()+ParseMode", func(t *testing.T) {
		b.parseMode = ModeHTML

		edited, err := b.EditCaption(msg, "<b>new caption with html</b>")
		require.NoError(t, err)
		assert.Equal(t, "new caption with html", edited.Caption)
		assert.Equal(t, EntityBold, edited.CaptionEntities[0].Type)

		edited, err = b.EditCaption(msg, "*new caption with markdown*", ModeMarkdown)
		require.NoError(t, err)
		assert.Equal(t, "new caption with markdown", edited.Caption)
		assert.Equal(t, EntityBold, edited.CaptionEntities[0].Type)

		b.parseMode = ModeDefault
	})

	t.Run("Edit(what=InputMedia)", func(t *testing.T) {
		edited, err := b.Edit(msg, photo)
		require.NoError(t, err)
		assert.Equal(t, edited.Photo.UniqueID, photo.UniqueID)
	})

	t.Run("Send(what=string)", func(t *testing.T) {
		msg, err = b.Send(to, t.Name())
		require.NoError(t, err)
		assert.Equal(t, t.Name(), msg.Text)

		rpl, err := b.Reply(msg, t.Name())
		require.NoError(t, err)
		assert.Equal(t, rpl.Text, msg.Text)
		assert.NotNil(t, rpl.ReplyTo)
		assert.Equal(t, rpl.ReplyTo, msg)
		assert.True(t, rpl.IsReply())

		fwd, err := b.Forward(to, msg)
		require.NoError(t, err)
		assert.NotNil(t, msg, fwd)
		assert.True(t, fwd.IsForwarded())

		fwd.ID += 1 // nonexistent message
		_, err = b.Forward(to, fwd)
		assert.Equal(t, ErrToForwardNotFound, err)
	})

	t.Run("Edit(what=string)", func(t *testing.T) {
		msg, err = b.Edit(msg, t.Name())
		require.NoError(t, err)
		assert.Equal(t, t.Name(), msg.Text)

		_, err = b.Edit(msg, msg.Text)
		assert.Error(t, err) // message is not modified
	})

	t.Run("Edit(what=ReplyMarkup)", func(t *testing.T) {
		good := &ReplyMarkup{
			InlineKeyboard: [][]InlineButton{
				{{
					Data: "btn",
					Text: "Hi Telebot!",
				}},
			},
		}
		bad := &ReplyMarkup{
			InlineKeyboard: [][]InlineButton{
				{{
					Data: strings.Repeat("*", 65),
					Text: "Bad Button",
				}},
			},
		}

		edited, err := b.Edit(msg, good)
		require.NoError(t, err)
		assert.Equal(t, edited.ReplyMarkup.InlineKeyboard, good.InlineKeyboard)

		edited, err = b.EditReplyMarkup(edited, nil)
		require.NoError(t, err)
		assert.Nil(t, edited.ReplyMarkup.InlineKeyboard)

		_, err = b.Edit(edited, bad)
		assert.Equal(t, ErrButtonDataInvalid, err)
	})

	t.Run("Edit(what=Location)", func(t *testing.T) {
		loc := &Location{Lat: 42, Lng: 69, LivePeriod: 60}
		edited, err := b.Send(to, loc)
		require.NoError(t, err)
		assert.NotNil(t, edited.Location)

		loc = &Location{Lat: loc.Lng, Lng: loc.Lat}
		edited, err = b.Edit(edited, *loc)
		require.NoError(t, err)
		assert.NotNil(t, edited.Location)
	})

	// should be the last
	t.Run("Delete()", func(t *testing.T) {
		require.NoError(t, b.Delete(msg))
	})

	t.Run("Notify()", func(t *testing.T) {
		assert.Equal(t, ErrBadRecipient, b.Notify(nil, Typing))
		require.NoError(t, b.Notify(to, Typing))
	})

	t.Run("Answer()", func(t *testing.T) {
		assert.Error(t, b.Answer(&Query{}, &QueryResponse{
			Results: Results{&ArticleResult{}},
		}))
	})

	t.Run("Respond()", func(t *testing.T) {
		assert.Error(t, b.Respond(&Callback{}, &CallbackResponse{}))
	})

	t.Run("Payments", func(t *testing.T) {
		assert.NotPanics(t, func() {
			b.Accept(&PreCheckoutQuery{})
			b.Accept(&PreCheckoutQuery{}, "error")
		})
		assert.NotPanics(t, func() {
			b.Ship(&ShippingQuery{})
			b.Ship(&ShippingQuery{}, "error")
			b.Ship(&ShippingQuery{}, ShippingOption{}, ShippingOption{})
			assert.Equal(t, ErrUnsupportedWhat, b.Ship(&ShippingQuery{}, 0))
		})
	})

	t.Run("Commands", func(t *testing.T) {
		orig := []Command{{
			Text:        "test",
			Description: "test command",
		}}
		require.NoError(t, b.SetCommands(orig))

		cmds, err := b.GetCommands()
		require.NoError(t, err)
		assert.Equal(t, orig, cmds)
	})
}
