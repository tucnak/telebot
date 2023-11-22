package telebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		"can_manage_topics":         false,
		"can_send_audios":           true,
		"can_send_documents":        true,
		"can_send_photos":           true,
		"can_send_videos":           true,
		"can_send_video_notes":      true,
		"can_send_voice_notes":      true,
	}
	assert.Equal(t, expected, params)
}
