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
	assert.Equal(t, 10*time.Minute, lt.Duration("dur"))

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
		InlineKeyboard: [][]tele.InlineButton{{{Unique: "stop", Text: "Stop", Data: "1"}}},
	}, lt.MarkupLocale("en", "inline", 1))

	assert.Equal(t, &tele.ArticleResult{
		ResultBase: tele.ResultBase{
			ID:   "1853",
			Type: "article",
		},
		Title:       "Some title",
		Description: "Some description",
		Text:        "The text of the article",
		ThumbURL:    "https://preview.picture",
	}, lt.ResultLocale("en", "article", struct {
		ID          int
		Title       string
		Description string
		Content     string
		PreviewURL  string
	}{
		ID:          1853,
		Title:       "Some title",
		Description: "Some description",
		Content:     "The text of the article",
		PreviewURL:  "https://preview.picture",
	}))
}
