package layout

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
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

	assert.ElementsMatch(t, []tele.Command{{
		Text:        "start",
		Description: "Start the bot",
	}, {
		Text:        "help",
		Description: "How to use the bot",
	}}, lt.Commands())

	assert.Equal(t, "string", lt.String("str"))
	assert.Equal(t, 123, lt.Int("num"))
	assert.Equal(t, int64(123), lt.Int64("num"))
	assert.Equal(t, float64(123), lt.Float("num"))
	assert.Equal(t, tele.ChatID(123), lt.ChatID("num"))

	assert.Equal(t, []string{"abc", "def"}, lt.Strings("strs"))
	assert.Equal(t, []int{123, 456}, lt.Ints("nums"))
	assert.Equal(t, []int64{123, 456}, lt.Int64s("nums"))
	assert.Equal(t, []float64{123, 456}, lt.Floats("nums"))

	obj := lt.Get("obj")
	assert.NotNil(t, obj)

	const dur = 10 * time.Minute
	assert.Equal(t, dur, obj.Duration("dur"))
	assert.True(t, lt.Duration("obj.dur") == obj.Duration("dur"))

	arr := lt.Slice("arr")
	assert.Len(t, arr, 2)

	for _, v := range arr {
		assert.Equal(t, dur, v.Duration("dur"))
	}

	assert.Equal(t, &tele.Btn{
		Unique: "pay",
		Text:   "Pay",
		Data:   "1|100.00|USD",
	}, lt.ButtonLocale("en", "pay", struct {
		UserID   int
		Amount   string
		Currency string
	}{
		UserID:   1,
		Amount:   "100.00",
		Currency: "USD",
	}))

	assert.Equal(t, &tele.ReplyMarkup{
		ReplyKeyboard: [][]tele.ReplyButton{
			{{Text: "Help"}},
			{{Text: "Settings"}},
		},
		ResizeKeyboard: true,
	}, lt.MarkupLocale("en", "reply_shortened"))

	assert.Equal(t, &tele.ReplyMarkup{
		ReplyKeyboard:   [][]tele.ReplyButton{{{Text: "Send a contact", Contact: true}}},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}, lt.MarkupLocale("en", "reply_extended"))

	assert.Equal(t, &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{{
			{
				Unique: "stop",
				Text:   "Stop",
				Data:   "1",
			},
		}},
	}, lt.MarkupLocale("en", "inline", 1))

	assert.Equal(t, &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{{
			{
				Text:   "This is a web app",
				WebApp: &tele.WebApp{URL: "https://google.com"},
			},
		}},
	}, lt.MarkupLocale("en", "web_app"))

	assert.Equal(t, &tele.ArticleResult{
		ResultBase: tele.ResultBase{
			ID:   "1853",
			Type: "article",
		},
		Title:       "Some title",
		Description: "Some description",
		ThumbURL:    "https://preview.picture",
		Text:        "This is an article.",
	}, lt.ResultLocale("en", "article", struct {
		ID          int
		Title       string
		Description string
		PreviewURL  string
	}{
		ID:          1853,
		Title:       "Some title",
		Description: "Some description",
		PreviewURL:  "https://preview.picture",
	}))
}
