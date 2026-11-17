package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDraftStop(t *testing.T) {
	b, srv := newFakeBot(t)
	stop := &SendOptions{CanStop: true, KeepOnStop: true, ThreadID: 3}

	require.NoError(t, b.SendDraft(ChatID(2), 9, "partial", stop))
	assert.Equal(t, map[string]string{
		"chat_id": "2", "draft_id": "9", "text": "partial", "message_thread_id": "3",
		"can_stop": "true", "keep_on_stop": "true",
	}, srv.Last().Params)

	require.NoError(t, b.SendRichDraft(ChatID(2), 9, &InputRichMessage{Markdown: "x"}, stop))
	req := srv.Last()
	assert.Equal(t, "sendRichMessageDraft", req.Method)
	assert.Equal(t, "true", req.Params["can_stop"])
	assert.Equal(t, "true", req.Params["keep_on_stop"])

	require.NoError(t, b.SendDraft(ChatID(2), 9, "partial"))
	assert.NotContains(t, srv.Last().Params, "can_stop")
}

func TestGenerationStopped(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendMessage", json.RawMessage(`{"message_id":1,"chat":{"id":2}}`))

	var got *MessageGenerationStopped
	b.Handle(OnGenerationStopped, func(c Context) error {
		got = c.GenerationStopped()
		return c.Send("stopped")
	})

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"stopped_message_generation":{
		"chat":{"id":2,"type":"private"},"message_thread_id":3,"draft_id":9}}`), &u))
	b.ProcessUpdate(u)

	require.NotNil(t, got)
	assert.Equal(t, 9, got.DraftID)
	assert.Equal(t, 3, got.ThreadID)
	assert.Equal(t, "2", srv.Last().Params["chat_id"])
	assert.Contains(t, AllowedUpdates, "stopped_message_generation")
}
