package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ Context = (*nativeContext)(nil)

func TestContext(t *testing.T) {
	t.Run("Get,Set", func(t *testing.T) {
		var c Context
		c = new(nativeContext)
		c.Set("name", "Jon Snow")
		assert.Equal(t, "Jon Snow", c.Get("name"))
	})
}

func TestContextBusinessMessage(t *testing.T) {
	b, _ := newFakeBot(t)

	for _, kind := range []string{"business_message", "edited_business_message"} {
		var u Update
		require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"`+kind+`":{"message_id":3,"date":1,
			"business_connection_id":"bc","chat":{"id":7,"type":"private"},"from":{"id":5},"text":"hi"}}`), &u))

		c := b.NewContext(u)
		require.NotNil(t, c.Message(), kind)
		assert.Equal(t, int64(7), c.Chat().ID, kind)
		assert.Equal(t, int64(5), c.Sender().ID, kind)
	}
}

func TestContextSendIgnoreThread(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendMessage", json.RawMessage(`{"message_id":1,"chat":{"id":7}}`))

	c := b.NewContext(Update{Message: &Message{ID: 1, ThreadID: 9, TopicMessage: true, Chat: &Chat{ID: 7}}})
	require.NoError(t, c.Send("in the thread"))
	assert.Equal(t, "9", srv.Last().Params["message_thread_id"])

	require.NoError(t, c.Send("outside it", IgnoreThread))
	assert.NotContains(t, srv.Last().Params, "message_thread_id")
}
