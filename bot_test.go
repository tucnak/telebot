package telebot

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Cached bot instance to avoid getMe method flooding.
var b, _ = newTestBot()

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

func TestHandle(t *testing.T) {
	b.Handle("/start", func(m *Message) {})
	assert.Contains(t, b.handlers, "/start")

	btn := &InlineButton{Unique: "test"}
	b.Handle(btn, func(c *Callback) {})
	assert.Contains(t, b.handlers, btn.CallbackUnique())

	assert.Panics(t, func() { b.Handle(1, func() {}) })
}

func TestStart(t *testing.T) {
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

	pref.Poller = &testPoller{Message: "/start"}
	b, err = NewBot(pref)
	assert.NoError(t, err)

	var ok bool
	b.Handle("/start", func(m *Message) {
		assert.Equal(t, m.Text, "/start")
		ok = true
	})

	time.AfterFunc(50*time.Millisecond, b.Stop)
	b.Start() // stops after some delay
	assert.True(t, ok)
}
