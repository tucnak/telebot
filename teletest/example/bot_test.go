package example

import (
	"testing"

	"github.com/stretchr/testify/require"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/teletest"
)

var b, _ = NewBot()

func TestBot(t *testing.T) {
	r, expect := teletest.New(b)

	r.Trigger(tele.Update{
		Message: &tele.Message{
			Sender: &tele.User{ID: 1},
			Text:   "/start",
		},
	})

	expect.Reply(t, "Hello!")
	require.Contains(t, users, int64(1))

	r.Trigger(tele.Update{
		Message: &tele.Message{
			Sender: &tele.User{ID: 1},
			Text:   "echo",
		},
	})

	expect.Reply(t, "echo")
}
