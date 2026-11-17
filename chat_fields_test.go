package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatFullInfoFields(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("getChat", json.RawMessage(`{"id":7,"type":"private","active_usernames":["a","b"],
		"join_to_send_messages":true,"join_by_request":true,"paid_message_star_count":10,
		"rating":{"level":3,"rating":120,"current_level_rating":100,"next_level_rating":200},
		"unique_gift_colors":{"model_custom_emoji_id":"m","symbol_custom_emoji_id":"s","light_theme_main_color":1,
			"light_theme_other_colors":[2],"dark_theme_main_color":3,"dark_theme_other_colors":[4]}}`))

	chat, err := b.ChatByID(7)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, chat.ActiveUsernames)
	assert.True(t, chat.JoinToSendMessages)
	assert.True(t, chat.JoinByRequest)
	assert.Equal(t, 10, chat.PaidMessageStarCount)
	assert.Equal(t, &UserRating{Level: 3, Rating: 120, CurrentLevelRating: 100, NextLevelRating: 200}, chat.Rating)
	assert.Equal(t, "m", chat.UniqueGiftColors.ModelCustomEmojiID)
}

func TestPaidMessages(t *testing.T) {
	b, _ := newFakeBot(t)

	var got *Message
	b.Handle(OnPaidMessagePriceChanged, func(c Context) error { got = c.Message(); return nil })

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":1,"date":1,"chat":{"id":-100},
		"paid_star_count":5,"paid_message_price_changed":{"paid_message_star_count":20}}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, got)
	assert.Equal(t, 5, got.PaidStarCount)
	assert.Equal(t, 20, got.PaidMessagePriceChanged.StarCount)

	var topic Topic
	require.NoError(t, json.Unmarshal([]byte(`{"name":"New","icon_color":1,"message_thread_id":3,"is_name_implicit":true}`), &topic))
	assert.True(t, topic.IsNameImplicit)
}
