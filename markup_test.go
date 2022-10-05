package telebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBtn(t *testing.T) {
	r := &ReplyMarkup{}

	assert.Equal(t, &ReplyButton{Text: "T"}, r.Text("T").Reply())
	assert.Equal(t, &ReplyButton{Text: "T", Contact: true}, r.Contact("T").Reply())
	assert.Equal(t, &ReplyButton{Text: "T", Location: true}, r.Location("T").Reply())
	assert.Equal(t, &ReplyButton{Text: "T", Poll: PollAny}, r.Poll("T", PollAny).Reply())

	assert.Nil(t, r.Data("T", "u").Reply())
	assert.Equal(t, &InlineButton{Unique: "u", Text: "T"}, r.Data("T", "u").Inline())
	assert.Equal(t, &InlineButton{Unique: "u", Text: "T", Data: "1|2"}, r.Data("T", "u", "1", "2").Inline())
	assert.Equal(t, &InlineButton{Text: "T", URL: "url"}, r.URL("T", "url").Inline())
	assert.Equal(t, &InlineButton{Text: "T", InlineQuery: "q"}, r.Query("T", "q").Inline())
	assert.Equal(t, &InlineButton{Text: "T", InlineQueryChat: "q"}, r.QueryChat("T", "q").Inline())
	assert.Equal(t, &InlineButton{Text: "T", Login: &Login{Text: "T"}}, r.Login("T", &Login{Text: "T"}).Inline())
	assert.Equal(t, &InlineButton{Text: "T", WebApp: &WebApp{URL: "url"}}, r.WebApp("T", &WebApp{URL: "url"}).Inline())
}

func TestOptions(t *testing.T) {
	r := &ReplyMarkup{}
	r.Reply(
		r.Row(r.Text("Menu")),
		r.Row(r.Text("Settings")),
	)

	assert.Equal(t, [][]ReplyButton{
		{{Text: "Menu"}},
		{{Text: "Settings"}},
	}, r.ReplyKeyboard)

	i := &ReplyMarkup{}
	i.Inline(i.Row(
		i.Data("Previous", "prev"),
		i.Data("Next", "next"),
	))

	assert.Equal(t, [][]InlineButton{{
		{Unique: "prev", Text: "Previous"},
		{Unique: "next", Text: "Next"},
	}}, i.InlineKeyboard)

	assert.Panics(t, func() {
		r.Reply(r.Row(r.Data("T", "u")))
		i.Inline(i.Row(i.Text("T")))
	})

	assert.Equal(t, r.copy(), r)
	assert.Equal(t, i.copy(), i)

	o := &SendOptions{ReplyMarkup: r}
	assert.Equal(t, o.copy(), o)

	data, err := PollQuiz.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte(`{"type":"quiz"}`), data)
}
