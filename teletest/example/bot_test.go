package example

import (
	"testing"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/teletest"

	"github.com/stretchr/testify/assert"
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

	expect.Send(t, "Hello!")
	assert.Contains(t, users, int64(1))

	r.Trigger(tele.Update{
		Message: &tele.Message{
			Sender: &tele.User{ID: 1},
			Text:   "echo",
		},
	})

	expect.Reply(t, "echo")
}
