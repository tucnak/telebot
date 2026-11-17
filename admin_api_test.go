package telebot

import (
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
