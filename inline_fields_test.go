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
	loc := &LocationResult{Location: Location{Lat: 1, Lng: 2}, Title: "L", ThumbWidth: 64, ThumbHeight: 48}
	venue := &VenueResult{Location: Location{Lat: 1, Lng: 2}, Title: "V", Address: "A", GooglePlaceID: "gp"}
	text := &ArticleResult{Title: "T"}
	text.Content = &InputTextMessageContent{Text: "x", Entities: []MessageEntity{{Type: EntityItalic, Length: 1}}}
	invoice := &ArticleResult{Title: "Buy"}
	invoice.Content = &InputInvoiceMessageContent{Title: "Item", Description: "D", Payload: "p", Currency: "XTR",
		Prices: []Price{{Label: "Item", Amount: 5}}, PhotoURL: "https://example.com/i.jpg"}

	require.NoError(t, b.Answer(&Query{ID: "q"}, &QueryResponse{Results: Results{photo, loc, venue, text, invoice}}))

	var results []map[string]interface{}
	srv.Last().JSON(t, "results", &results)
	require.Len(t, results, 5)
	assert.Equal(t, []interface{}{map[string]interface{}{"type": "bold", "offset": 0.0, "length": 2.0}}, results[0]["caption_entities"])
	assert.NotContains(t, results[0], "parse_mode")
	assert.Equal(t, "HTML", results[1]["parse_mode"])
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
}

func TestSavePreparedInlineMessage(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("savePreparedInlineMessage", json.RawMessage(`{"id":"prep","expiration_date":9}`))

	result := &ArticleResult{Title: "T", Text: "x"}
	result.SetResultID("1")
	msg, err := b.SavePreparedInlineMessage(&User{ID: 7}, result, PreparedInlineChats{Bots: true, Groups: true, Channels: true})
	require.NoError(t, err)
	assert.Equal(t, "prep", msg.ID)

	req := srv.Last()
	assert.Equal(t, "true", req.Params["allow_bot_chats"])
	assert.Equal(t, "true", req.Params["allow_group_chats"])
	assert.Equal(t, "true", req.Params["allow_channel_chats"])
	assert.NotContains(t, req.Params, "allow_user_chats")
	assert.Contains(t, req.Params["result"], `"type":"article"`)
}
