package telebot

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBusinessAccountMethods(t *testing.T) {
	b, srv := newFakeBot(t)
	call := func(method string, params map[string]string) {
		t.Helper()
		req := srv.Last()
		assert.Equal(t, method, req.Method)
		assert.Equal(t, params, req.Params)
	}

	msg := &Message{ID: 5, Chat: &Chat{ID: 7}}
	require.NoError(t, b.ReadBusinessMessage("bc", msg))
	call("readBusinessMessage", map[string]string{"business_connection_id": "bc", "chat_id": "7", "message_id": "5"})

	require.NoError(t, b.DeleteBusinessMessages("bc", []Editable{msg, &Message{ID: 6, Chat: &Chat{ID: 7}}}))
	call("deleteBusinessMessages", map[string]string{"business_connection_id": "bc", "message_ids": "[5,6]"})

	require.NoError(t, b.SetBusinessAccountName("bc", "Ada", ""))
	call("setBusinessAccountName", map[string]string{"business_connection_id": "bc", "first_name": "Ada"})

	require.NoError(t, b.SetBusinessAccountUsername("bc", ""))
	call("setBusinessAccountUsername", map[string]string{"business_connection_id": "bc", "username": ""})

	require.NoError(t, b.SetBusinessAccountBio("bc", "hi"))
	call("setBusinessAccountBio", map[string]string{"business_connection_id": "bc", "bio": "hi"})

	require.NoError(t, b.SetBusinessAccountProfilePhoto("bc", &InputProfilePhotoStatic{File: FromReader(strings.NewReader("jpg"))}, true))
	req := srv.Last()
	assert.Equal(t, "setBusinessAccountProfilePhoto", req.Method)
	assert.Equal(t, "true", req.Params["is_public"])
	assert.Contains(t, req.Params["photo"], `"type":"static"`)
	assert.Len(t, req.Files, 1)

	require.NoError(t, b.RemoveBusinessAccountProfilePhoto("bc", false))
	call("removeBusinessAccountProfilePhoto", map[string]string{"business_connection_id": "bc"})

	require.NoError(t, b.SetBusinessAccountGiftSettings("bc", true, AcceptedGiftTypes{UniqueGifts: true}))
	req = srv.Last()
	assert.Equal(t, "true", req.Params["show_gift_button"])
	assert.JSONEq(t, `{"unlimited_gifts":false,"limited_gifts":false,"unique_gifts":true,"premium_subscription":false,"gifts_from_channels":false}`,
		req.Params["accepted_gift_types"])

	srv.Reply("getBusinessAccountStarBalance", json.RawMessage(`{"amount":3}`))
	bal, err := b.BusinessAccountStarBalance("bc")
	require.NoError(t, err)
	assert.Equal(t, int64(3), bal.Amount)

	require.NoError(t, b.TransferBusinessAccountStars("bc", 2))
	call("transferBusinessAccountStars", map[string]string{"business_connection_id": "bc", "star_count": "2"})
}

func TestBusinessBotRights(t *testing.T) {
	var conn BusinessConnection
	require.NoError(t, json.Unmarshal([]byte(`{"id":"bc","user":{"id":1},"user_chat_id":1,"date":1,"is_enabled":true,
		"rights":{"can_reply":true,"can_manage_stories":true}}`), &conn))
	assert.True(t, conn.Rights.CanReply)
	assert.True(t, conn.Rights.CanManageStories)
}
