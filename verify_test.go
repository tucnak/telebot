package telebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerify(t *testing.T) {
	b, srv := newFakeBot(t)

	require.NoError(t, b.Verify(&User{ID: 7}, "<b>Staff</b>", ModeHTML))
	assert.Equal(t, "verifyUser", srv.Last().Method)
	assert.Equal(t, map[string]string{
		"user_id": "7", "custom_description": "<b>Staff</b>", "custom_description_parse_mode": "HTML",
	}, srv.Last().Params)

	require.NoError(t, b.Verify(&Chat{ID: -100}, "Official", Entities{{Type: EntityBold, Length: 8}}))
	assert.Equal(t, "verifyChat", srv.Last().Method)
	assert.Equal(t, "-100", srv.Last().Params["chat_id"])
	assert.JSONEq(t, `[{"type":"bold","offset":0,"length":8}]`, srv.Last().Params["custom_description_entities"])

	require.NoError(t, b.Verify(ChatID(-100), ""))
	assert.Equal(t, map[string]string{"chat_id": "-100"}, srv.Last().Params)

	assert.Equal(t, ErrBadRecipient, b.Verify(nil, ""))
}
