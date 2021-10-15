package telebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ Context = (*nativeContext)(nil)

func TestContext(t *testing.T) {
	t.Run("Get,Set", func(t *testing.T) {
		var c Context
		c = new(nativeContext)
		c.Set("name", "Jon Snow")
		assert.Equal(t, "Jon Snow", c.Get("name"))
	})
}

func Test_nativeContext_Text(t *testing.T) {
	type fields struct {
		b                *Bot
		message          *Message
		callback         *Callback
		query            *Query
		shippingQuery    *ShippingQuery
		preCheckoutQuery *PreCheckoutQuery
		poll             *Poll
		pollAnswer       *PollAnswer
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"message",
			fields{
				message: &Message{Text: "message text"},
			},
			"message text",
		},
		{
			"callback",
			fields{
				callback: &Callback{Message: &Message{Text: "callback message text"}},
			},
			"callback message text",
		},
		{
			"query",
			fields{
				query: &Query{Text: "query text"},
			},
			"query text",
		},
		{
			"nothing",
			fields{},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &nativeContext{
				b:                tt.fields.b,
				message:          tt.fields.message,
				callback:         tt.fields.callback,
				query:            tt.fields.query,
				shippingQuery:    tt.fields.shippingQuery,
				preCheckoutQuery: tt.fields.preCheckoutQuery,
				poll:             tt.fields.poll,
				pollAnswer:       tt.fields.pollAnswer,
			}
			if got := c.Text(); got != tt.want {
				t.Errorf("nativeContext.Text() = %v, want %v", got, tt.want)
			}
		})
	}
}
