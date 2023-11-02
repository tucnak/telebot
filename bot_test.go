package telebot

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// required to test send and edit methods
	token     = os.Getenv("TELEBOT_SECRET")
	chatID, _ = strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	userID, _ = strconv.ParseInt(os.Getenv("USER_ID"), 10, 64)

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

	b, err := NewBot(Settings{Offline: true})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, b.Me)
	assert.Equal(t, DefaultApiURL, b.URL)
	assert.Equal(t, 100, cap(b.Updates))
	assert.NotZero(t, b.client.Timeout)

	pref = defaultSettings()
	client := &http.Client{Timeout: time.Minute}
	pref.URL = "http://api.telegram.org" // not https
	pref.Client = client
	pref.Poller = &LongPoller{Timeout: time.Second}
	pref.Updates = 50
	pref.ParseMode = ModeHTML
	pref.Offline = true

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

	b.Handle("/start", func(c Context) error { return nil })
	assert.Contains(t, b.handlers, "/start")

	reply := ReplyButton{Text: "reply"}
	b.Handle(&reply, func(c Context) error { return nil })

	inline := InlineButton{Unique: "inline"}
	b.Handle(&inline, func(c Context) error { return nil })

	btnReply := (&ReplyMarkup{}).Text("btnReply")
	b.Handle(&btnReply, func(c Context) error { return nil })

	btnInline := (&ReplyMarkup{}).Data("", "btnInline")
	b.Handle(&btnInline, func(c Context) error { return nil })

	assert.Contains(t, b.handlers, btnReply.CallbackUnique())
	assert.Contains(t, b.handlers, btnInline.CallbackUnique())
	assert.Contains(t, b.handlers, reply.CallbackUnique())
	assert.Contains(t, b.handlers, inline.CallbackUnique())
}

func TestBotStart(t *testing.T) {
	if token == "" {
		t.Skip("TELEBOT_SECRET is required")
	}

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
	b.Handle("/start", func(c Context) error {
		assert.Equal(t, c.Text(), "/start")
		tp.done <- struct{}{}
		ok = true
		return nil
	})

	go b.Start()
	<-tp.done
	b.Stop()

	assert.True(t, ok)
}

func TestBotProcessUpdate(t *testing.T) {
	b, err := NewBot(Settings{Synchronous: true, Offline: true})
	if err != nil {
		t.Fatal(err)
	}

	b.Handle(OnMedia, func(c Context) error {
		assert.NotNil(t, c.Message().Photo)
		return nil
	})
	b.ProcessUpdate(Update{Message: &Message{Photo: &Photo{}}})

	b.Handle("/start", func(c Context) error {
		assert.Equal(t, "/start", c.Text())
		return nil
	})
	b.Handle("hello", func(c Context) error {
		assert.Equal(t, "hello", c.Text())
		return nil
	})
	b.Handle(OnText, func(c Context) error {
		assert.Equal(t, "text", c.Text())
		return nil
	})
	b.Handle(OnPinned, func(c Context) error {
		assert.NotNil(t, c.Message())
		return nil
	})
	b.Handle(OnPhoto, func(c Context) error {
		assert.NotNil(t, c.Message().Photo)
		return nil
	})
	b.Handle(OnVoice, func(c Context) error {
		assert.NotNil(t, c.Message().Voice)
		return nil
	})
	b.Handle(OnAudio, func(c Context) error {
		assert.NotNil(t, c.Message().Audio)
		return nil
	})
	b.Handle(OnAnimation, func(c Context) error {
		assert.NotNil(t, c.Message().Animation)
		return nil
	})
	b.Handle(OnDocument, func(c Context) error {
		assert.NotNil(t, c.Message().Document)
		return nil
	})
	b.Handle(OnSticker, func(c Context) error {
		assert.NotNil(t, c.Message().Sticker)
		return nil
	})
	b.Handle(OnVideo, func(c Context) error {
		assert.NotNil(t, c.Message().Video)
		return nil
	})
	b.Handle(OnVideoNote, func(c Context) error {
		assert.NotNil(t, c.Message().VideoNote)
		return nil
	})
	b.Handle(OnContact, func(c Context) error {
		assert.NotNil(t, c.Message().Contact)
		return nil
	})
	b.Handle(OnLocation, func(c Context) error {
		assert.NotNil(t, c.Message().Location)
		return nil
	})
	b.Handle(OnVenue, func(c Context) error {
		assert.NotNil(t, c.Message().Venue)
		return nil
	})
	b.Handle(OnDice, func(c Context) error {
		assert.NotNil(t, c.Message().Dice)
		return nil
	})
	b.Handle(OnInvoice, func(c Context) error {
		assert.NotNil(t, c.Message().Invoice)
		return nil
	})
	b.Handle(OnPayment, func(c Context) error {
		assert.NotNil(t, c.Message().Payment)
		return nil
	})
	b.Handle(OnAddedToGroup, func(c Context) error {
		assert.NotNil(t, c.Message().GroupCreated)
		return nil
	})
	b.Handle(OnUserJoined, func(c Context) error {
		assert.NotNil(t, c.Message().UserJoined)
		return nil
	})
	b.Handle(OnUserLeft, func(c Context) error {
		assert.NotNil(t, c.Message().UserLeft)
		return nil
	})
	b.Handle(OnNewGroupTitle, func(c Context) error {
		assert.Equal(t, "title", c.Message().NewGroupTitle)
		return nil
	})
	b.Handle(OnNewGroupPhoto, func(c Context) error {
		assert.NotNil(t, c.Message().NewGroupPhoto)
		return nil
	})
	b.Handle(OnGroupPhotoDeleted, func(c Context) error {
		assert.True(t, c.Message().GroupPhotoDeleted)
		return nil
	})
	b.Handle(OnMigration, func(c Context) error {
		from, to := c.Migration()
		assert.Equal(t, int64(1), from)
		assert.Equal(t, int64(2), to)
		return nil
	})
	b.Handle(OnEdited, func(c Context) error {
		assert.Equal(t, "edited", c.Message().Text)
		return nil
	})
	b.Handle(OnChannelPost, func(c Context) error {
		assert.Equal(t, "post", c.Message().Text)
		return nil
	})
	b.Handle(OnEditedChannelPost, func(c Context) error {
		assert.Equal(t, "edited post", c.Message().Text)
		return nil
	})
	b.Handle(OnCallback, func(c Context) error {
		if data := c.Callback().Data; data[0] != '\f' {
			assert.Equal(t, "callback", data)
		}
		return nil
	})
	b.Handle("\funique", func(c Context) error {
		assert.Equal(t, "callback", c.Callback().Data)
		return nil
	})
	b.Handle(OnQuery, func(c Context) error {
		assert.Equal(t, "query", c.Data())
		return nil
	})
	b.Handle(OnInlineResult, func(c Context) error {
		assert.Equal(t, "result", c.InlineResult().ResultID)
		return nil
	})
	b.Handle(OnShipping, func(c Context) error {
		assert.Equal(t, "shipping", c.ShippingQuery().ID)
		return nil
	})
	b.Handle(OnCheckout, func(c Context) error {
		assert.Equal(t, "checkout", c.PreCheckoutQuery().ID)
		return nil
	})
	b.Handle(OnPoll, func(c Context) error {
		assert.Equal(t, "poll", c.Poll().ID)
		return nil
	})
	b.Handle(OnPollAnswer, func(c Context) error {
		assert.Equal(t, "poll", c.PollAnswer().PollID)
		return nil
	})

	b.Handle(OnWebApp, func(c Context) error {
		assert.Equal(t, "webapp", c.Message().WebAppData.Data)
		return nil
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
	b.ProcessUpdate(Update{Message: &Message{Invoice: &Invoice{}}})
	b.ProcessUpdate(Update{Message: &Message{Payment: &Payment{}}})
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
	b.ProcessUpdate(Update{InlineResult: &InlineResult{ResultID: "result"}})
	b.ProcessUpdate(Update{ShippingQuery: &ShippingQuery{ID: "shipping"}})
	b.ProcessUpdate(Update{PreCheckoutQuery: &PreCheckoutQuery{ID: "checkout"}})
	b.ProcessUpdate(Update{Poll: &Poll{ID: "poll"}})
	b.ProcessUpdate(Update{PollAnswer: &PollAnswer{PollID: "poll"}})
	b.ProcessUpdate(Update{Message: &Message{WebAppData: &WebAppData{Data: "webapp"}}})
}

func TestBotOnError(t *testing.T) {
	b, err := NewBot(Settings{Synchronous: true, Offline: true})
	if err != nil {
		t.Fatal(err)
	}

	var ok bool
	b.onError = func(err error, c Context) {
		assert.Equal(t, b, c.(*nativeContext).b)
		assert.NotNil(t, err)
		ok = true
	}

	b.runHandler(func(c Context) error {
		return errors.New("not nil")
	}, &nativeContext{b: b})

	assert.True(t, ok)
}

func TestBotMiddleware(t *testing.T) {
	t.Run("calling order", func(t *testing.T) {
		var trace []string

		handler := func(name string) HandlerFunc {
			return func(c Context) error {
				trace = append(trace, name)
				return nil
			}
		}

		middleware := func(name string) MiddlewareFunc {
			return func(next HandlerFunc) HandlerFunc {
				return func(c Context) error {
					trace = append(trace, name+":in")
					err := next(c)
					trace = append(trace, name+":out")
					return err
				}
			}
		}

		b, err := NewBot(Settings{Synchronous: true, Offline: true})
		if err != nil {
			t.Fatal(err)
		}

		b.Use(middleware("global1"), middleware("global2"))
		b.Handle("/a", handler("/a"), middleware("handler1a"), middleware("handler2a"))

		group := b.Group()
		group.Use(middleware("group1"), middleware("group2"))
		group.Handle("/b", handler("/b"), middleware("handler1b"))

		b.ProcessUpdate(Update{
			Message: &Message{Text: "/a"},
		})
		assert.Equal(t, []string{
			"global1:in", "global2:in",
			"handler1a:in", "handler2a:in",
			"/a",
			"handler2a:out", "handler1a:out",
			"global2:out", "global1:out",
		}, trace)

		trace = trace[:0]

		b.ProcessUpdate(Update{
			Message: &Message{Text: "/b"},
		})
		assert.Equal(t, []string{
			"global1:in", "global2:in",
			"group1:in", "group2:in",
			"handler1b:in",
			"/b",
			"handler1b:out",
			"group2:out", "group1:out",
			"global2:out", "global1:out",
		}, trace)
	})

	fatal := func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			t.Fatal("fatal middleware should not be called")
			return nil
		}
	}

	nop := func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			return next(c)
		}
	}

	t.Run("combining with global middleware", func(t *testing.T) {
		b, err := NewBot(Settings{Synchronous: true, Offline: true})
		if err != nil {
			t.Fatal(err)
		}

		// Pre-allocate middleware slice to make sure
		// it has extra capacity after group-level middleware is added.
		b.group.middleware = make([]MiddlewareFunc, 0, 2)
		b.Use(nop)

		b.Handle("/a", func(c Context) error { return nil }, nop)
		b.Handle("/b", func(c Context) error { return nil }, fatal)

		b.ProcessUpdate(Update{Message: &Message{Text: "/a"}})
	})

	t.Run("combining with group middleware", func(t *testing.T) {
		b, err := NewBot(Settings{Synchronous: true, Offline: true})
		if err != nil {
			t.Fatal(err)
		}

		g := b.Group()
		// Pre-allocate middleware slice to make sure
		// it has extra capacity after group-level middleware is added.
		g.middleware = make([]MiddlewareFunc, 0, 2)
		g.Use(nop)

		g.Handle("/a", func(c Context) error { return nil }, nop)
		g.Handle("/b", func(c Context) error { return nil }, fatal)

		b.ProcessUpdate(Update{Message: &Message{Text: "/a"}})
	})
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
		File:    FromURL("https://telegra.ph/file/65c5237b040ebf80ec278.jpg"),
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

		photo2 := *photo
		photo2.Caption = ""

		msgs, err := b.SendAlbum(to, Album{photo, &photo2}, ModeHTML)
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

		sleep()

		edited, err = b.EditCaption(msg, "*new caption with markdown*", ModeMarkdown)
		require.NoError(t, err)
		assert.Equal(t, "new caption with markdown", edited.Caption)
		assert.Equal(t, EntityBold, edited.CaptionEntities[0].Type)

		sleep()

		edited, err = b.EditCaption(msg, "_new caption with markdown \\(V2\\)_", ModeMarkdownV2)
		require.NoError(t, err)
		assert.Equal(t, "new caption with markdown (V2)", edited.Caption)
		assert.Equal(t, EntityItalic, edited.CaptionEntities[0].Type)

		b.parseMode = ModeDefault
	})

	t.Run("Edit(what=Media)", func(t *testing.T) {
		edited, err := b.Edit(msg, photo)
		require.NoError(t, err)
		assert.Equal(t, edited.Photo.UniqueID, photo.UniqueID)

		resp, err := http.Get("https://telegra.ph/file/274e5eb26f348b10bd8ee.mp4")
		require.NoError(t, err)
		defer resp.Body.Close()

		file, err := ioutil.TempFile("", "")
		require.NoError(t, err)

		_, err = io.Copy(file, resp.Body)
		require.NoError(t, err)

		animation := &Animation{
			File:     FromDisk(file.Name()),
			Caption:  t.Name(),
			FileName: "animation.gif",
		}

		msg, err := b.Send(msg.Chat, animation)
		require.NoError(t, err)

		if msg.Animation != nil {
			assert.Equal(t, msg.Animation.FileID, animation.FileID)
		} else {
			assert.Equal(t, msg.Document.FileID, animation.FileID)
		}

		_, err = b.Edit(edited, animation)
		require.NoError(t, err)
	})

	t.Run("Edit(what=Animation)", func(t *testing.T) {})

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
		assert.Equal(t, ErrNotFoundToForward, err)
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
		assert.Nil(t, edited.ReplyMarkup)

		_, err = b.Edit(edited, bad)
		assert.Equal(t, ErrBadButtonData, err)
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

	// Should be after the Edit tests.
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
		var (
			set1 = []Command{{
				Text:        "test1",
				Description: "test command 1",
			}}
			set2 = []Command{{
				Text:        "test2",
				Description: "test command 2",
			}}
			scope = CommandScope{
				Type:   CommandScopeChat,
				ChatID: chatID,
			}
		)

		err := b.SetCommands(set1)
		require.NoError(t, err)

		cmds, err := b.Commands()
		require.NoError(t, err)
		assert.Equal(t, set1, cmds)

		err = b.SetCommands(set2, "en", scope)
		require.NoError(t, err)

		cmds, err = b.Commands()
		require.NoError(t, err)
		assert.Equal(t, set1, cmds)

		cmds, err = b.Commands("en", scope)
		require.NoError(t, err)
		assert.Equal(t, set2, cmds)

		require.NoError(t, b.DeleteCommands("en", scope))
		require.NoError(t, b.DeleteCommands())
	})

	t.Run("InviteLink", func(t *testing.T) {
		inviteLink, err := b.CreateInviteLink(&Chat{ID: chatID}, nil)
		require.NoError(t, err)
		assert.True(t, len(inviteLink.InviteLink) > 0)

		sleep()

		response, err := b.EditInviteLink(&Chat{ID: chatID}, &ChatInviteLink{InviteLink: inviteLink.InviteLink})
		require.NoError(t, err)
		assert.True(t, len(response.InviteLink) > 0)

		sleep()

		response, err = b.RevokeInviteLink(&Chat{ID: chatID}, inviteLink.InviteLink)
		require.Nil(t, err)
		assert.True(t, len(response.InviteLink) > 0)
	})
}

func sleep() {
	time.Sleep(time.Second)
}
