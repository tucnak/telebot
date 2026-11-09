package telebot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbedSendOptionsReplyParams(t *testing.T) {
	b := &Bot{}

	tests := []struct {
		name string
		rp   *ReplyParams
		want string
	}{
		{
			name: "same chat leaves out chat_id",
			rp:   &ReplyParams{MessageID: 42},
			want: `{"message_id":42}`,
		},
		{
			name: "cross chat",
			rp:   &ReplyParams{MessageID: 42, ChatID: -100123},
			want: `{"message_id":42,"chat_id":-100123}`,
		},
		{
			name: "quote",
			rp:   &ReplyParams{MessageID: 42, Quote: `say "hi"`, QuoteParseMode: ModeHTML, QuotePosition: 3},
			want: `{"message_id":42,"quote":"say \"hi\"","quote_parse_mode":"HTML","quote_position":3}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := map[string]string{}
			b.embedSendOptions(params, &SendOptions{ReplyParams: tt.rp})
			assert.Equal(t, tt.want, params["reply_parameters"])
		})
	}

	params := map[string]string{}
	b.embedSendOptions(params, &SendOptions{})
	_, ok := params["reply_parameters"]
	assert.False(t, ok)
}

func TestSendWithReplyParams(t *testing.T) {
	var got map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":1,"chat":{"id":1}}}`))
	}))
	defer srv.Close()

	b, err := NewBot(Settings{URL: srv.URL, Token: "t", Offline: true, Synchronous: true})
	require.NoError(t, err)

	_, err = b.Send(ChatID(1), "text", &ReplyParams{MessageID: 42, ChatID: -100123})
	require.NoError(t, err)
	assert.Equal(t, `{"message_id":42,"chat_id":-100123}`, got["reply_parameters"])
}
