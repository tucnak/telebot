package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneralTopicHiddenMessages(t *testing.T) {
	b, _ := newFakeBot(t)

	var got []string
	b.Handle(OnGeneralTopicHidden, func(c Context) error { got = append(got, "hidden"); return nil })
	b.Handle(OnGeneralTopicUnhidden, func(c Context) error { got = append(got, "unhidden"); return nil })

	for _, field := range []string{"general_forum_topic_hidden", "general_forum_topic_unhidden"} {
		var u Update
		require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":1,"date":1,"chat":{"id":-100},"`+field+`":{}}}`), &u))
		b.ProcessUpdate(u)
	}
	assert.Equal(t, []string{"hidden", "unhidden"}, got)

	var msg Message
	require.NoError(t, json.Unmarshal([]byte(`{"message_id":1,"date":1,"chat":{"id":-100},"forum_topic_closed":{}}`), &msg))
	assert.Equal(t, &ForumTopicClosed{}, msg.TopicClosed)
}
