package telebot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuestMessageUpdate(t *testing.T) {
	b, err := NewBot(Settings{Synchronous: true, Offline: true})
	require.NoError(t, err)

	b.Handle(OnText, func(c Context) error {
		t.Fatal("guest message should not trigger text handlers")
		return nil
	})
	b.Handle("/start", func(c Context) error {
		t.Fatal("guest message should not trigger command handlers")
		return nil
	})

	var handled bool
	b.Handle(OnGuestMessage, func(c Context) error {
		msg := c.Message()

		assert.Equal(t, "/start", c.Text())
		assert.Empty(t, c.Data())
		assert.Nil(t, c.Args())
		assert.Equal(t, int64(1), c.Sender().ID)
		assert.Equal(t, int64(2), c.Chat().ID)
		assert.Equal(t, int64(3), msg.GuestUser.ID)
		assert.Equal(t, int64(4), msg.GuestChat.ID)
		assert.Equal(t, "guest-query-id", msg.GuestQueryID)

		handled = true
		return nil
	})

	b.ProcessUpdate(Update{GuestMessage: &Message{
		Sender:       &User{ID: 1},
		Chat:         &Chat{ID: 2},
		GuestUser:    &User{ID: 3},
		GuestChat:    &Chat{ID: 4},
		Text:         "/start",
		GuestQueryID: "guest-query-id",
	}})

	assert.True(t, handled)
}

func TestGuestModeJSON(t *testing.T) {
	data := []byte(`{
		"guest_message": {
			"from": {"id": 1, "first_name": "Sender"},
			"chat": {"id": 2, "type": "private"},
			"guest_bot_caller_user": {"id": 3, "first_name": "Caller"},
			"guest_bot_caller_chat": {"id": 4, "type": "supergroup"},
			"guest_query_id": "guest-query-id"
		}
	}`)

	var update Update
	require.NoError(t, json.Unmarshal(data, &update))
	require.NotNil(t, update.GuestMessage)
	assert.Equal(t, int64(3), update.GuestMessage.GuestUser.ID)
	assert.Equal(t, int64(4), update.GuestMessage.GuestChat.ID)
	assert.Equal(t, "guest-query-id", update.GuestMessage.GuestQueryID)

	var user User
	require.NoError(t, json.Unmarshal([]byte(`{
		"id": 5,
		"supports_guest_queries": true
	}`), &user))
	assert.True(t, user.SupportsGuest)
}

func TestAnswerGuest(t *testing.T) {
	var body map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/answerGuestQuery", r.URL.Path)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		w.Write([]byte(`{"ok":true,"result":{"inline_message_id":"inline-id"}}`))
	}))
	defer srv.Close()

	b, err := NewBot(Settings{
		Offline:   true,
		Token:     "token",
		URL:       srv.URL,
		Client:    srv.Client(),
		ParseMode: ModeHTML,
	})
	require.NoError(t, err)

	sent, err := b.AnswerGuest(
		&Message{GuestQueryID: "guest-query-id"},
		&ArticleResult{Title: "Title", Text: "Text"},
	)
	require.NoError(t, err)
	require.NotNil(t, sent)

	assert.Equal(t, "inline-id", sent.InlineMessageID)
	assert.Equal(t, "guest-query-id", body["guest_query_id"])

	result, ok := body["result"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "article", result["type"])
	assert.Equal(t, ModeHTML, result["parse_mode"])
	assert.Equal(t, "Title", result["title"])
	assert.Equal(t, "Text", result["message_text"])
}

func TestContextAnswerGuest(t *testing.T) {
	c := &nativeContext{}
	assert.EqualError(t, c.AnswerGuest(&ArticleResult{}), "telebot: context guest message is nil")
}

func TestAllowedUpdatesIncludesGuestMessage(t *testing.T) {
	assert.Contains(t, AllowedUpdates, "guest_message")
}
