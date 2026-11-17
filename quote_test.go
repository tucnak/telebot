package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRichQuotes(t *testing.T) {
	var msg Message
	require.NoError(t, json.Unmarshal([]byte(`{"message_id":2,"date":1,"chat":{"id":7},
		"quote":{"text":"","position":4,"rich_texts":["plain",{"type":"bold","text":"b"}]}}`), &msg))
	require.Len(t, msg.Quote.RichTexts, 2)
	assert.Equal(t, "plain", msg.Quote.RichTexts[0].String())
	assert.Equal(t, RichBold, msg.Quote.RichTexts[1].Type)

	b, srv := newFakeBot(t)
	srv.Reply("sendMessage", json.RawMessage(`{"message_id":3,"chat":{"id":7}}`))
	_, err := b.Send(ChatID(7), "reply", &ReplyParams{MessageID: 1, QuoteRichTexts: msg.Quote.RichTexts})
	require.NoError(t, err)
	assert.JSONEq(t, `{"message_id":1,"quote_rich_texts":["plain",{"type":"bold","text":"b"}]}`, srv.Last().Params["reply_parameters"])
}

func TestPaidMessagesWaived(t *testing.T) {
	b, _ := newFakeBot(t)

	fired := false
	b.Handle(OnPaidMessagesWaived, func(c Context) error { fired = c.Message().PaidMessagesWaived != nil; return nil })

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":1,"date":1,
		"chat":{"id":7,"type":"private"},"paid_messages_waived":{}}}`), &u))
	b.ProcessUpdate(u)
	assert.True(t, fired)
}
