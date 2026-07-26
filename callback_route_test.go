package telebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Buttons built with a slash-prefixed unique (e.g. menu.Data(text, "/setting/lang", ...))
// must reach the command handler registered under the same path, while plain
// uniques keep going through the "\f"-namespaced lookup.
func TestCallbackRouting(t *testing.T) {
	b, err := NewBot(Settings{Token: "1:offline", Offline: true, Synchronous: true})
	assert.NoError(t, err)

	var (
		args     []string
		sender   int64
		callback bool
		plain    bool
	)

	b.Handle("/setting/lang", func(c Context) error {
		args, sender, callback = c.Args(), c.Sender().ID, c.Callback() != nil
		return nil
	})
	b.Handle("\fplain", func(c Context) error {
		plain = true
		return nil
	})

	b.ProcessUpdate(Update{Callback: &Callback{
		Sender:  &User{ID: 555},
		Message: &Message{Chat: &Chat{ID: 9}, Sender: &User{ID: 42}},
		Data:    "\f/setting/lang|100|zh|save",
	}})

	assert.Equal(t, []string{"100", "zh", "save"}, args)
	assert.Equal(t, int64(555), sender, "sender must be the clicking user, not the bot")
	assert.True(t, callback, "handler must be able to tell it was invoked by a callback")

	b.ProcessUpdate(Update{Callback: &Callback{
		Sender:  &User{ID: 555},
		Message: &Message{Chat: &Chat{ID: 9}},
		Data:    "\fplain|x",
	}})
	assert.True(t, plain, "plain uniques must keep working")
}
