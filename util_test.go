package telebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractOk(t *testing.T) {
	data := []byte(`{"ok": true, "result": {}}`)
	require.NoError(t, extractOk(data))

	data = []byte(`{
		"ok": false,
		"error_code": 400,
		"description": "Bad Request: reply message not found"
	}`)
	assert.EqualError(t, extractOk(data), ErrNotFoundToReply.Error())

	data = []byte(`{
		"ok": false,
		"error_code": 429,
		"description": "Too Many Requests: retry after 8",
		"parameters": {"retry_after": 8}
	}`)
	assert.Equal(t, FloodError{
		err:        NewError(429, "Too Many Requests: retry after 8"),
		RetryAfter: 8,
	}, extractOk(data))

	data = []byte(`{
		"ok": false,
		"error_code": 400,
		"description": "Bad Request: group chat was upgraded to a supergroup chat",
		"parameters": {"migrate_to_chat_id": -100123456789}
	}`)
	assert.Equal(t, GroupError{
		err:        ErrGroupMigrated,
		MigratedTo: -100123456789,
	}, extractOk(data))
}

func TestExtractMessage(t *testing.T) {
	data := []byte(`{"ok":true,"result":true}`)
	_, err := extractMessage(data)
	assert.Equal(t, ErrTrueResult, err)

	data = []byte(`{"ok":true,"result":{"foo":"bar"}}`)
	_, err = extractMessage(data)
	require.NoError(t, err)
}

func TestEmbedRights(t *testing.T) {
	rights := NoRestrictions()
	params := map[string]interface{}{
		"chat_id": "1",
		"user_id": "2",
	}
	embedRights(params, rights)

	expected := map[string]interface{}{
		"is_anonymous":              false,
		"chat_id":                   "1",
		"user_id":                   "2",
		"can_be_edited":             true,
		"can_send_messages":         true,
		"can_send_media_messages":   true,
		"can_send_polls":            true,
		"can_send_other_messages":   true,
		"can_add_web_page_previews": true,
		"can_change_info":           false,
		"can_post_messages":         false,
		"can_edit_messages":         false,
		"can_delete_messages":       false,
		"can_invite_users":          false,
		"can_restrict_members":      false,
		"can_pin_messages":          false,
		"can_promote_members":       false,
		"can_manage_video_chats":    false,
		"can_manage_chat":           false,
	}
	assert.Equal(t, expected, params)
}
