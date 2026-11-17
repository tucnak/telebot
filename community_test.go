package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommunityServiceMessages(t *testing.T) {
	b, _ := newFakeBot(t)

	var added, removed *Message
	b.Handle(OnCommunityChatAdded, func(c Context) error { added = c.Message(); return nil })
	b.Handle(OnCommunityChatRemoved, func(c Context) error { removed = c.Message(); return nil })

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":1,"date":1,
		"chat":{"id":-100,"type":"supergroup"},
		"community_chat_added":{"community":{"id":4503599627370495,"name":"Makers"}}}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, added)
	assert.Equal(t, Community{ID: 4503599627370495, Name: "Makers"}, added.CommunityChatAdded.Community)

	u = Update{}
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":2,"message":{"message_id":2,"date":1,
		"chat":{"id":-100,"type":"supergroup"},"community_chat_removed":{}}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, removed)
	assert.NotNil(t, removed.CommunityChatRemoved)
}

func TestChatCommunity(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("getChat", json.RawMessage(`{"id":-100,"type":"supergroup","title":"g","community":{"id":9,"name":"Makers"}}`))

	chat, err := b.ChatByID(-100)
	require.NoError(t, err)
	assert.Equal(t, &Community{ID: 9, Name: "Makers"}, chat.Community)

	info, err := b.ChatFullInfo(ChatID(-100))
	require.NoError(t, err)
	assert.Equal(t, &Community{ID: 9, Name: "Makers"}, info.Community)
}

func TestSubscriptionUpdate(t *testing.T) {
	b, _ := newFakeBot(t)

	var got *BotSubscriptionUpdated
	var sender *User
	b.Handle(OnSubscription, func(c Context) error {
		got, sender = c.Subscription(), c.Sender()
		return nil
	})

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"subscription":{
		"user":{"id":7,"first_name":"U"},"invoice_payload":"pro","state":"canceled"}}`), &u))
	b.ProcessUpdate(u)

	require.NotNil(t, got)
	assert.Equal(t, "pro", got.Payload)
	assert.Equal(t, SubscriptionCanceled, got.State)
	assert.Equal(t, int64(7), sender.ID)
	assert.Contains(t, AllowedUpdates, "subscription")
}
