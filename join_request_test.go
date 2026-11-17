package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJoinRequestQueries(t *testing.T) {
	b, srv := newFakeBot(t)

	var req *ChatJoinRequest
	b.Handle(OnChatJoinRequest, func(c Context) error {
		req = c.ChatJoinRequest()
		return c.Bot().SendJoinRequestWebApp(req, "https://example.com/check")
	})

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"chat_join_request":{"chat":{"id":-100},
		"from":{"id":7},"user_chat_id":7,"date":1,"query_id":"jq1"}}`), &u))
	b.ProcessUpdate(u)

	require.NotNil(t, req)
	assert.Equal(t, "sendChatJoinRequestWebApp", srv.Last().Method)
	assert.Equal(t, map[string]string{"chat_join_request_query_id": "jq1", "web_app_url": "https://example.com/check"}, srv.Last().Params)

	require.NoError(t, b.AnswerJoinRequest(req, JoinRequestQueue))
	assert.Equal(t, "answerChatJoinRequestQuery", srv.Last().Method)
	assert.Equal(t, map[string]string{"chat_join_request_query_id": "jq1", "result": "queue"}, srv.Last().Params)
}

func TestGuardBotAndBotFlags(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("getChat", json.RawMessage(`{"id":-100,"type":"supergroup","guard_bot":{"id":9,"is_bot":true,"first_name":"G"}}`))

	chat, err := b.ChatByID(-100)
	require.NoError(t, err)
	assert.Equal(t, int64(9), chat.GuardBot.ID)

	var me User
	require.NoError(t, json.Unmarshal([]byte(`{"id":1,"is_bot":true,"first_name":"B",
		"supports_join_request_queries":true,"can_manage_bots":true,"has_topics_enabled":true}`), &me))
	assert.True(t, me.SupportsJoinRequestQueries)
	assert.True(t, me.CanManageBots)
	assert.True(t, me.HasTopicsEnabled)
}
