package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBanUsesBanChatMember(t *testing.T) {
	b, srv := newFakeBot(t)

	chat, member := &Chat{ID: -100}, &ChatMember{User: &User{ID: 7}, RestrictedUntil: 99}
	require.NoError(t, b.Ban(chat, member, true))

	req := srv.Last()
	assert.Equal(t, "banChatMember", req.Method)
	assert.Equal(t, map[string]string{
		"chat_id": "-100", "user_id": "7", "until_date": "99", "revoke_messages": "true",
	}, req.Params)
}

func TestLenUsesGetChatMemberCount(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("getChatMemberCount", 12)

	n, err := b.Len(&Chat{ID: -100})
	require.NoError(t, err)
	assert.Equal(t, 12, n)
	assert.Equal(t, "getChatMemberCount", srv.Last().Method)
	assert.Equal(t, "-100", srv.Last().Params["chat_id"])
}

func TestPromoteWelcomeMessages(t *testing.T) {
	b, srv := newFakeBot(t)

	member := &ChatMember{User: &User{ID: 7}, Rights: Rights{CanSendWelcomeMessages: true}}
	require.NoError(t, b.Promote(&Chat{ID: -100}, member))
	assert.Equal(t, "promoteChatMember", srv.Last().Method)
	assert.Equal(t, "true", srv.Last().Params["can_send_welcome_messages"])

	var m ChatMember
	require.NoError(t, json.Unmarshal([]byte(`{"status":"administrator","user":{"id":7},"can_send_welcome_messages":true}`), &m))
	assert.True(t, m.Rights.CanSendWelcomeMessages)
}

func TestChatManagement100(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("getChatAdministrators", json.RawMessage(`[{"status":"administrator","user":{"id":9,"is_bot":true}}]`))

	admins, err := b.AdminsOf(&Chat{ID: -100}, true)
	require.NoError(t, err)
	assert.True(t, admins[0].User.IsBot)
	assert.Equal(t, map[string]string{"chat_id": "-100", "return_bots": "true"}, srv.Last().Params)

	_, err = b.AdminsOf(&Chat{ID: -100})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"chat_id": "-100"}, srv.Last().Params)

	msg := &Message{ID: 5, Chat: &Chat{ID: -100}}
	require.NoError(t, b.DeleteReaction(msg, &User{ID: 7}))
	assert.Equal(t, "deleteMessageReaction", srv.Last().Method)
	assert.Equal(t, map[string]string{"chat_id": "-100", "message_id": "5", "user_id": "7"}, srv.Last().Params)

	require.NoError(t, b.DeleteAllReactions(&Chat{ID: -100}, &Chat{ID: -200}))
	assert.Equal(t, "deleteAllMessageReactions", srv.Last().Method)
	assert.Equal(t, map[string]string{"chat_id": "-100", "actor_chat_id": "-200"}, srv.Last().Params)

	require.NoError(t, b.Restrict(&Chat{ID: -100}, &ChatMember{User: &User{ID: 7}, Rights: Rights{CanReactToMessages: true}}))
	assert.Contains(t, srv.Last().Params["permissions"], `"can_react_to_messages":true`)
}
