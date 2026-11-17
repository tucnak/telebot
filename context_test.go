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
