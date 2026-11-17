package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInlineResultFields(t *testing.T) {
	b, srv := newFakeBot(t)
	b.parseMode = ModeHTML // not sent where entities are given

	photo := &PhotoResult{URL: "https://example.com/p.jpg", ThumbURL: "https://example.com/t.jpg", Caption: "hi"}
	photo.CaptionEntities = []MessageEntity{{Type: EntityBold, Length: 2}}
	plain := &PhotoResult{URL: "https://example.com/p.jpg", ThumbURL: "https://example.com/t.jpg", Caption: "<b>hi</b>"}
	loc := &LocationResult{Location: Location{Lat: 1, Lng: 2}, Title: "L", ThumbWidth: 64, ThumbHeight: 48}
	venue := &VenueResult{Location: Location{Lat: 1, Lng: 2}, Title: "V", Address: "A", GooglePlaceID: "gp"}
	text := &ArticleResult{Title: "T"}
	text.Content = &InputTextMessageContent{Text: "x", Entities: []MessageEntity{{Type: EntityItalic, Length: 1}}}
	invoice := &ArticleResult{Title: "Buy"}
	invoice.Content = &InputInvoiceMessageContent{Title: "Item", Description: "D", Payload: "p", Currency: "XTR",
		Prices: []Price{{Label: "Item", Amount: 5}}, PhotoURL: "https://example.com/i.jpg"}

	require.NoError(t, b.Answer(&Query{ID: "q"}, &QueryResponse{Results: Results{photo, loc, venue, text, invoice, plain}}))

	var results []map[string]interface{}
	srv.Last().JSON(t, "results", &results)
	require.Len(t, results, 6)
	assert.Equal(t, []interface{}{map[string]interface{}{"type": "bold", "offset": 0.0, "length": 2.0}}, results[0]["caption_entities"])
	assert.NotContains(t, results[0], "parse_mode")
	assert.Equal(t, "HTML", results[5]["parse_mode"])
	assert.Equal(t, 64.0, results[1]["thumbnail_width"])
	assert.Equal(t, "gp", results[2]["google_place_id"])
	assert.Equal(t, map[string]interface{}{"message_text": "x", "entities": []interface{}{map[string]interface{}{"type": "italic", "offset": 0.0, "length": 1.0}}},
		results[3]["input_message_content"])
	content := results[4]["input_message_content"].(map[string]interface{})
	assert.Equal(t, "XTR", content["currency"])
	assert.Equal(t, "https://example.com/i.jpg", content["photo_url"])
	assert.NotContains(t, content, "provider_token")
}

func TestInputMessageContentFields(t *testing.T) {
	data, err := json.Marshal(&InputLocationMessageContent{Lat: 1, Lng: 2, HorizontalAccuracy: 5, LivePeriod: 60, Heading: 90, AlertRadius: 100})
	require.NoError(t, err)
	assert.JSONEq(t, `{"latitude":1,"longitude":2,"horizontal_accuracy":5,"live_period":60,"heading":90,"proximity_alert_radius":100}`, string(data))

	data, err = json.Marshal(&InputContactMessageContent{PhoneNumber: "+1", FirstName: "A", VCard: "V"})
	require.NoError(t, err)
	assert.JSONEq(t, `{"phone_number":"+1","first_name":"A","vcard":"V"}`, string(data))
}

func TestForwardCopyDeleteMany(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("forwardMessages", json.RawMessage(`[{"message_id":8},{"message_id":9}]`))
	srv.Reply("copyMessages", json.RawMessage(`[{"message_id":10}]`))
	msgs := []Editable{&Message{ID: 1, Chat: &Chat{ID: 2}}, &Message{ID: 3, Chat: &Chat{ID: 2}}}

	sent, err := b.ForwardMany(ChatID(7), msgs)
	require.NoError(t, err)
	assert.Equal(t, 9, sent[1].ID)
	assert.Equal(t, map[string]string{"chat_id": "7", "from_chat_id": "2", "message_ids": "[1,3]"}, srv.Last().Params)

	_, err = b.CopyMany(ChatID(7), msgs, &SendOptions{RemoveCaption: true})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"chat_id": "7", "from_chat_id": "2", "message_ids": "[1,3]", "remove_caption": "true"}, srv.Last().Params)

	require.NoError(t, b.DeleteMany(msgs))
	assert.Equal(t, map[string]string{"chat_id": "2", "message_ids": "[1,3]"}, srv.Last().Params)

	assert.Error(t, b.DeleteMany(nil))
	assert.Error(t, b.DeleteMany([]Editable{StoredMessage{MessageID: "1,2", ChatID: 2}}), "not one message ID")
	_, err = b.ForwardMany(ChatID(7), nil)
	assert.Error(t, err)
}

func TestSavePreparedInlineMessage(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("savePreparedInlineMessage", json.RawMessage(`{"id":"prep","expiration_date":9}`))

	result := &ArticleResult{Title: "T", Text: "x"}
	result.SetResultID("1")
	msg, err := b.SavePreparedInlineMessage(&User{ID: 7}, result, &PreparedInlineChats{Bots: true, Groups: true, Channels: true})
	require.NoError(t, err)
	assert.Equal(t, "prep", msg.ID)

	req := srv.Last()
	assert.Equal(t, "true", req.Params["allow_bot_chats"])
	assert.Equal(t, "true", req.Params["allow_group_chats"])
	assert.Equal(t, "true", req.Params["allow_channel_chats"])
	assert.NotContains(t, req.Params, "allow_user_chats")
	assert.Contains(t, req.Params["result"], `"type":"article"`)
}

func TestForwardCopyVideoStartTimestamp(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("forwardMessage", json.RawMessage(`{"message_id":8,"chat":{"id":7}}`))
	srv.Reply("copyMessage", json.RawMessage(`{"message_id":9}`))
	msg := &Message{ID: 1, Chat: &Chat{ID: 2}}

	_, err := b.Forward(ChatID(7), msg, Silent, int64(30))
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"chat_id": "7", "from_chat_id": "2", "message_id": "1",
		"disable_notification": "true", "video_start_timestamp": "30"}, srv.Last().Params)

	_, err = b.Copy(ChatID(7), msg, int64(30))
	require.NoError(t, err)
	assert.Equal(t, "30", srv.Last().Params["video_start_timestamp"])

	_, err = b.Copy(ChatID(7), msg, &SendOptions{Caption: "new", CaptionAbove: true, Entities: Entities{{Type: EntityBold, Length: 3}}})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"chat_id": "7", "from_chat_id": "2", "message_id": "1", "caption": "new",
		"show_caption_above_media": "true", "caption_entities": `[{"type":"bold","offset":0,"length":3}]`}, srv.Last().Params)
}

func TestEditReplyMarkupBusinessConnection(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("editMessageReplyMarkup", json.RawMessage(`{"message_id":1,"chat":{"id":2}}`))
	msg := &Message{ID: 1, Chat: &Chat{ID: 2}}

	_, err := b.Edit(msg, &ReplyMarkup{}, &SendOptions{BusinessConnectionID: "bc"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"chat_id": "2", "message_id": "1", "business_connection_id": "bc", "reply_markup": "{}"}, srv.Last().Params)

	_, err = b.EditReplyMarkup(msg, nil)
	require.NoError(t, err)
	assert.NotContains(t, srv.Last().Params, "business_connection_id")
}

func TestMessageIsForwardedAndInaccessible(t *testing.T) {
	var fwd, post, gone Message
	require.NoError(t, json.Unmarshal([]byte(`{"message_id":1,"date":5,"chat":{"id":2},"forward_origin":{"type":"hidden_user","date":4,"sender_user_name":"x"}}`), &fwd))
	require.NoError(t, json.Unmarshal([]byte(`{"message_id":1,"date":5,"chat":{"id":-100,"type":"channel"}}`), &post))
	require.NoError(t, json.Unmarshal([]byte(`{"message_id":1,"date":0,"chat":{"id":2}}`), &gone))

	assert.True(t, fwd.IsForwarded())
	assert.False(t, post.IsForwarded())
	assert.False(t, post.Inaccessible(), "a channel post has no sender")
	assert.True(t, gone.Inaccessible())
}

func TestRespondCacheTime(t *testing.T) {
	b, srv := newFakeBot(t)
	require.NoError(t, b.Respond(&Callback{ID: "c"}, &CallbackResponse{Text: "ok", CacheTime: 30}))
	assert.Equal(t, map[string]string{"callback_query_id": "c", "text": "ok", "cache_time": "30"}, srv.Last().Params)
}

func TestStopLiveLocation(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("stopMessageLiveLocation", json.RawMessage(`{"message_id":1,"chat":{"id":7}}`))
	srv.Reply("stopMessageLiveLocation", true)

	_, err := b.StopLiveLocation(&Message{ID: 1, Chat: &Chat{ID: 7}})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"chat_id": "7", "message_id": "1"}, srv.Last().Params)

	// A live location sent through inline mode has an inline message ID.
	_, err = b.StopLiveLocation(&Callback{MessageID: "inline"})
	assert.Equal(t, ErrTrueResult, err)
	assert.Equal(t, map[string]string{"inline_message_id": "inline"}, srv.Last().Params)
}

func TestNotifyUnpinProfilePhotosOptions(t *testing.T) {
	b, srv := newFakeBot(t)

	require.NoError(t, b.Notify(ChatID(7), Typing))
	assert.Equal(t, map[string]string{"chat_id": "7", "action": "typing"}, srv.Last().Params)
	require.NoError(t, b.Notify(ChatID(7), Typing, 5))
	assert.Equal(t, "5", srv.Last().Params["message_thread_id"])
	require.NoError(t, b.Notify(ChatID(7), Typing, &SendOptions{ThreadID: 5, BusinessConnectionID: "bc"}))
	assert.Equal(t, map[string]string{"chat_id": "7", "action": "typing", "message_thread_id": "5", "business_connection_id": "bc"}, srv.Last().Params)

	require.NoError(t, b.Unpin(ChatID(7)))
	assert.Equal(t, map[string]string{"chat_id": "7"}, srv.Last().Params)
	require.NoError(t, b.Unpin(ChatID(7), 3, &SendOptions{BusinessConnectionID: "bc"}))
	assert.Equal(t, map[string]string{"chat_id": "7", "message_id": "3", "business_connection_id": "bc"}, srv.Last().Params)
	require.NoError(t, b.Unpin(ChatID(7), &Message{ID: 4, Chat: &Chat{ID: 7}}))
	assert.Equal(t, "4", srv.Last().Params["message_id"])
	require.NoError(t, b.Notify(ChatID(7), Typing, &Topic{ThreadID: 6}))
	assert.Equal(t, "6", srv.Last().Params["message_thread_id"])

	// An option neither takes is a mistake, not something to skip.
	assert.Panics(t, func() { b.Notify(ChatID(7), Typing, int32(5)) })
	assert.Panics(t, func() { b.Unpin(ChatID(7), int64(3)) })

	srv.Reply("getUserProfilePhotos", json.RawMessage(`{"total_count":9,"photos":[[{"file_id":"s"},{"file_id":"l"}]]}`))
	photos, err := b.ProfilePhotosOf(&User{ID: 7})
	require.NoError(t, err)
	require.Len(t, photos, 1)
	assert.Equal(t, map[string]string{"user_id": "7"}, srv.Last().Params)
	_, err = b.ProfilePhotosOf(&User{ID: 7}, 5, 2)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"user_id": "7", "offset": "5", "limit": "2"}, srv.Last().Params)
}

func TestAnswerWebAppResultType(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("answerWebAppQuery", json.RawMessage(`{"inline_message_id":"im"}`))

	result := &ArticleResult{Title: "T", Text: "x"}
	msg, err := b.AnswerWebApp(&Query{ID: "q"}, result)
	require.NoError(t, err)
	assert.Equal(t, "im", msg.InlineMessageID)
	assert.Contains(t, srv.Last().Params["result"], `"type":"article"`)
	assert.Contains(t, srv.Last().Params["result"], `"id":"0"`, "a result sent on its own gets an ID")
}
