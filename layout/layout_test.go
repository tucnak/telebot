package layout

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tele "gopkg.in/tucnak/telebot.v3"
)

func TestLayout(t *testing.T) {
	os.Setenv("TOKEN", "TEST")

	lt, err := New("example.yml")
	if err != nil {
		t.Fatal(err)
	}

	pref := lt.Settings()
	assert.Equal(t, "TEST", pref.Token)
	assert.Equal(t, "html", pref.ParseMode)
	assert.Equal(t, &tele.LongPoller{}, pref.Poller)

	assert.Equal(t, "string", lt.String("str"))
	assert.Equal(t, 123, lt.Int("num"))
	assert.Equal(t, int64(123), lt.Int64("num"))
	assert.Equal(t, float64(123), lt.Float("num"))
	assert.Equal(t, 10*time.Minute, lt.Duration("dur"))

	assert.Equal(t, &tele.ReplyMarkup{
		ReplyKeyboard: [][]tele.ReplyButton{
			{{Text: "Help"}},
			{{Text: "Settings"}},
		},
		ResizeKeyboard: true,
	}, lt.Markup(nil, "reply_shortened"))

	assert.Equal(t, &tele.ReplyMarkup{
		ReplyKeyboard:   [][]tele.ReplyButton{{{Text: "Send a contact", Contact: true}}},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}, lt.Markup(nil, "reply_extended"))

	assert.Equal(t, &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{{{Unique: "stop", Text: "Stop", Data: "1"}}},
	}, lt.Markup(nil, "inline", 1))
}
