package telebot

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
