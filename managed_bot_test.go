package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagedBotUpdates(t *testing.T) {
	b, _ := newFakeBot(t)

	var upd *ManagedBotUpdated
	var created *ManagedBotCreated
	b.Handle(OnManagedBot, func(c Context) error {
		upd = c.ManagedBot()
		assert.Equal(t, int64(7), c.Sender().ID)
		return nil
	})
	b.Handle(OnManagedBotCreated, func(c Context) error { created = c.Message().ManagedBotCreated; return nil })

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"managed_bot":{"user":{"id":7},"bot":{"id":99,"is_bot":true,"username":"helper_bot"}}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, upd)
	assert.Equal(t, "helper_bot", upd.Bot.Username)

	u = Update{}
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":2,"message":{"message_id":1,"date":1,"chat":{"id":7,"type":"private"},
		"managed_bot_created":{"bot":{"id":99,"is_bot":true}}}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, created)
	assert.Equal(t, int64(99), created.Bot.ID)
	assert.Contains(t, AllowedUpdates, "managed_bot")
}

func TestManagedBotMethods(t *testing.T) {
	b, srv := newFakeBot(t)
	bot := &User{ID: 99}

	srv.Reply("getManagedBotToken", "99:old")
	tok, err := b.ManagedBotToken(bot)
	require.NoError(t, err)
	assert.Equal(t, "99:old", tok)
	assert.Equal(t, map[string]string{"user_id": "99"}, srv.Last().Params)

	srv.Reply("replaceManagedBotToken", "99:new")
	tok, err = b.ReplaceManagedBotToken(bot)
	require.NoError(t, err)
	assert.Equal(t, "99:new", tok)

	srv.Reply("getManagedBotAccessSettings", json.RawMessage(`{"is_access_restricted":true,"added_users":[{"id":5}]}`))
	settings, err := b.ManagedBotAccessSettings(bot)
	require.NoError(t, err)
	assert.True(t, settings.IsAccessRestricted)
	assert.Equal(t, int64(5), settings.AddedUsers[0].ID)

	require.NoError(t, b.SetManagedBotAccessSettings(bot, BotAccessSettings{IsAccessRestricted: true, AddedUsers: []User{{ID: 5}, {ID: 6}}}))
	assert.Equal(t, map[string]string{"user_id": "99", "is_access_restricted": "true", "added_user_ids": "[5,6]"}, srv.Last().Params)

	srv.Reply("savePreparedKeyboardButton", json.RawMessage(`{"id":"pk1"}`))
	pk, err := b.SavePreparedKeyboardButton(&User{ID: 7}, ReplyButton{
		Text:       "Create a bot",
		ManagedBot: &ReplyManagedBot{ID: 1, SuggestedName: "Helper", SuggestedUsername: "helper_ai"},
	})
	require.NoError(t, err)
	assert.Equal(t, "pk1", pk.ID)
	assert.JSONEq(t, `{"text":"Create a bot","request_managed_bot":{"request_id":1,"suggested_name":"Helper","suggested_username":"helper_ai"}}`,
		srv.Last().Params["button"])
}
