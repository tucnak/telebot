package telebot

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const widgetReply = `{"message_id":3,"date":1,"chat":{"id":7,"type":"private"},
	"widget":{"id":"w1","height":200,"html":"<h1>Hi</h1>","public_data":"p","chat_data":"c"}}`

func TestSendWidget(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendWidget", json.RawMessage(widgetReply))

	msg, err := b.Send(ChatID(7), &Widget{
		HTML: "<h1>Hi</h1>", Width: 300, Height: 200, PublicData: "p", ChatData: "c",
		Caption: "<b>look</b>", CaptionAbove: true, AllowQueries: true, IsShared: true,
	}, ModeHTML, Protected)
	require.NoError(t, err)
	assert.Equal(t, &Widget{ID: "w1", Height: 200, HTML: "<h1>Hi</h1>", PublicData: "p", ChatData: "c"}, msg.Widget)

	assert.Equal(t, map[string]string{
		"chat_id": "7", "html": "<h1>Hi</h1>", "width": "300", "height": "200",
		"public_data": "p", "chat_data": "c", "caption": "<b>look</b>", "parse_mode": "HTML",
		"show_caption_above_media": "true", "allow_queries": "true", "is_shared": "true",
		"protect_content": "true",
	}, srv.Last().Params)

	// A received shared widget is sent again by its ID.
	srv.Reply("sendSharedWidget", json.RawMessage(widgetReply))
	_, err = b.Send(ChatID(8), &Widget{ID: msg.Widget.ID, ChatData: "other"})
	require.NoError(t, err)
	assert.Equal(t, "sendSharedWidget", srv.Last().Method)
	assert.Equal(t, map[string]string{"chat_id": "8", "widget_id": "w1", "chat_data": "other"}, srv.Last().Params)
}

func TestWidgetQueries(t *testing.T) {
	b, srv := newFakeBot(t)

	var q *WidgetQuery
	var sq *SharedWidgetQuery
	b.Handle(OnWidgetQuery, func(c Context) error {
		q = c.WidgetQuery()
		assert.Equal(t, int64(5), c.Sender().ID)
		assert.Equal(t, int64(7), c.Chat().ID)
		return c.Bot().AnswerWidgetQuery(q.ID, "pong")
	})
	b.Handle(OnSharedWidgetQuery, func(c Context) error {
		sq = c.SharedWidgetQuery()
		assert.Equal(t, int64(5), c.Sender().ID)
		return nil
	})

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"widget_query":{"id":"q1","from":{"id":5},
		"chat":{"id":7,"type":"private"},"message_id":3,"chat_instance":"ci","data":"ping",
		"type":{"type":"sign_data","wallet":{"address":"EQW","amount":"10","public_key":"ab","state_init":"te6"},"sign_data":"{}"}}}`), &u))
	b.ProcessUpdate(u)

	require.NotNil(t, q)
	assert.Equal(t, WidgetQuerySignData, q.Type.Type)
	assert.Equal(t, &WidgetQueryWallet{Address: "EQW", Amount: "10", PublicKey: "ab", StateInit: "te6"}, q.Type.Wallet)
	assert.Equal(t, "answerWidgetQuery", srv.Last().Method)
	assert.Equal(t, map[string]string{"widget_query_id": "q1", "data": "pong"}, srv.Last().Params)

	// The query identifies its message for EditWidgetData.
	pub, chat := "", "new"
	require.NoError(t, b.EditWidgetData(q, &pub, &chat))
	assert.Equal(t, "editMessageWidgetData", srv.Last().Method)
	assert.Equal(t, map[string]string{"chat_id": "7", "message_id": "3", "public_data": "", "chat_data": "new"}, srv.Last().Params)

	require.NoError(t, b.EditWidgetData(&WidgetQuery{InlineMessageID: "im"}, nil, &chat))
	assert.Equal(t, map[string]string{"inline_message_id": "im", "chat_data": "new"}, srv.Last().Params)

	u = Update{}
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":2,"shared_widget_query":{"id":"q2","from":{"id":5},
		"data":"x","widget_id":"w1","type":{"type":"sent_transaction","wallet":{"address":"EQW","amount":"1"},"msg_hash":"h"}}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, sq)
	assert.Equal(t, "h", sq.Type.MsgHash)
	assert.Equal(t, "w1", sq.WidgetID)

	require.NoError(t, b.EditSharedWidgetData("w1", "pub", 60))
	assert.Equal(t, map[string]string{"widget_id": "w1", "public_data": "pub", "cache_time": "60"}, srv.Last().Params)

	assert.Contains(t, AllowedUpdates, "widget_query")
	assert.Contains(t, AllowedUpdates, "shared_widget_query")
}

func TestWidgetQueryWalletCard(t *testing.T) {
	var q WidgetQuery
	require.NoError(t, json.Unmarshal([]byte(`{"id":"q","from":{"id":1},"chat_instance":"c","data":"d",
		"type":{"type":"data","wallet_card":{"address":"EQA","amount":"0","stickers":[]}}}`), &q))
	assert.Equal(t, WidgetQueryData, q.Type.Type)
	assert.Equal(t, "EQA", q.Type.WalletCard.Address)
}

func TestWidgetInlineContent(t *testing.T) {
	data, err := json.Marshal(&InputWidgetMessageContent{HTML: "<p>x</p>", AllowQueries: true})
	require.NoError(t, err)
	assert.JSONEq(t, `{"html":"<p>x</p>","allow_queries":true}`, string(data))

	data, err = json.Marshal(&InputSharedWidgetMessageContent{WidgetID: "w1", Caption: "c"})
	require.NoError(t, err)
	assert.JSONEq(t, `{"widget_id":"w1","caption":"c"}`, string(data))

	var _ InputMessageContent = &InputWidgetMessageContent{}
	var _ InputMessageContent = &InputSharedWidgetMessageContent{}
}

func TestWidgetLive(t *testing.T) {
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}
	if userID == 0 {
		t.Skip("USER_ID is required for widget live tests")
	}

	msg, err := b.Send(user, &Widget{
		HTML:       "<!doctype html><p>telebot widget</p>",
		Height:     120,
		PublicData: "v1",
		ChatData:   "c1",
		IsShared:   true,
		Caption:    "telebot widget",
	})
	if err != nil && strings.Contains(err.Error(), "WIDGET_NOT_ALLOWED") {
		t.Skip("the test bot is not allowed to send widgets")
	}
	require.NoError(t, err)
	require.NotNil(t, msg.Widget)
	assert.NotEmpty(t, msg.Widget.ID)

	pub := "v2"
	require.NoError(t, b.EditWidgetData(msg, &pub, nil))
	require.NoError(t, b.EditSharedWidgetData(msg.Widget.ID, "v3", 60))

	again, err := b.Send(user, &Widget{ID: msg.Widget.ID})
	require.NoError(t, err)
	assert.Equal(t, msg.Widget.ID, again.Widget.ID)
}
