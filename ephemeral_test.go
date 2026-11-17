package telebot

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ephemeralMessage(t *testing.T) *Message {
	var msg Message
	require.NoError(t, json.Unmarshal([]byte(`{
		"message_id": 0,
		"ephemeral_message_id": 33,
		"date": 1,
		"chat": {"id": -100, "type": "supergroup"},
		"receiver_user": {"id": 7, "first_name": "R"},
		"text": "only you"
	}`), &msg))
	return &msg
}

func TestEphemeralMessageDecode(t *testing.T) {
	msg := ephemeralMessage(t)
	assert.Equal(t, 33, msg.EphemeralID)
	require.NotNil(t, msg.ReceiverUser)
	assert.Equal(t, int64(7), msg.ReceiverUser.ID)
}

func TestEditEphemeral(t *testing.T) {
	b, srv := newFakeBot(t)
	msg := ephemeralMessage(t)
	ids := map[string]string{"chat_id": "-100", "receiver_user_id": "7", "ephemeral_message_id": "33"}
	with := func(extra map[string]string) map[string]string {
		out := map[string]string{}
		for k, v := range ids {
			out[k] = v
		}
		for k, v := range extra {
			out[k] = v
		}
		return out
	}

	markup := &ReplyMarkup{InlineKeyboard: [][]InlineButton{{{Text: "ok", Data: "ok"}}}}
	require.NoError(t, b.EditEphemeral(msg, "new", ModeHTML, markup))
	req := srv.Last()
	assert.Equal(t, "editEphemeralMessageText", req.Method)
	assert.Equal(t, "new", req.Params["text"])
	assert.Equal(t, "HTML", req.Params["parse_mode"])
	assert.Contains(t, req.Params["reply_markup"], `"text":"ok"`)
	assert.Equal(t, "33", req.Params["ephemeral_message_id"])

	require.NoError(t, b.EditEphemeral(msg, &InputRichMessage{Markdown: "*x*"}, ModeHTML))
	req = srv.Last()
	assert.Equal(t, "editEphemeralMessageText", req.Method)
	assert.Equal(t, with(map[string]string{"rich_message": `{"markdown":"*x*"}`}), req.Params)

	require.NoError(t, b.EditEphemeralCaption(msg, "cap", &SendOptions{CaptionAbove: true, Entities: Entities{{Type: EntityBold, Length: 3}}}))
	req = srv.Last()
	assert.Equal(t, "editEphemeralMessageCaption", req.Method)
	assert.Equal(t, "cap", req.Params["caption"])
	assert.Equal(t, "true", req.Params["show_caption_above_media"])
	assert.JSONEq(t, `[{"type":"bold","offset":0,"length":3}]`, req.Params["caption_entities"])

	require.NoError(t, b.EditEphemeral(msg, &Photo{File: FromReader(strings.NewReader("png"))}))
	req = srv.Last()
	assert.Equal(t, "editEphemeralMessageMedia", req.Method)
	assert.Equal(t, "png", string(req.Files["0"].Data))
	assert.JSONEq(t, `{"type":"photo","media":"attach://0","caption":""}`, req.Params["media"])

	require.NoError(t, b.EditEphemeral(msg, (*ReplyMarkup)(nil)))
	req = srv.Last()
	assert.Equal(t, "editEphemeralMessageReplyMarkup", req.Method)
	assert.Equal(t, with(map[string]string{"reply_markup": "{}"}), req.Params)

	require.NoError(t, b.DeleteEphemeral(msg))
	req = srv.Last()
	assert.Equal(t, "deleteEphemeralMessage", req.Method)
	assert.Equal(t, ids, req.Params)

	assert.Error(t, b.DeleteEphemeral(&Message{ID: 1, Chat: &Chat{ID: 1}}))
	assert.Error(t, b.EditEphemeral(msg, 42))
}

func TestReplyToEphemeral(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendMessage", json.RawMessage(`{"message_id":1,"chat":{"id":-100}}`))

	_, err := b.Send(ChatID(-100), "hi", &ReplyParams{EphemeralMessageID: 33})
	require.NoError(t, err)
	assert.Equal(t, `{"ephemeral_message_id":33}`, srv.Last().Params["reply_parameters"])
}

func TestEphemeralCommand(t *testing.T) {
	b, srv := newFakeBot(t)
	require.NoError(t, b.SetCommands([]Command{{Text: "secret", Description: "only you", IsEphemeral: true}}))
	assert.JSONEq(t, `[{"command":"secret","description":"only you","is_ephemeral":true}]`, srv.Last().Params["commands"])
}

func TestEditCaptionAbove(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("editMessageCaption", json.RawMessage(`{"message_id":1,"chat":{"id":2}}`))

	_, err := b.EditCaption(&Message{ID: 1, Chat: &Chat{ID: 2}}, "c", &SendOptions{CaptionAbove: true})
	require.NoError(t, err)
	assert.Equal(t, "true", srv.Last().Params["show_caption_above_media"])
}

func TestSendEphemeral(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendMessage", json.RawMessage(`{"message_id":0,"ephemeral_message_id":4,"date":1,
		"chat":{"id":-100,"type":"supergroup"},"receiver_user":{"id":7,"first_name":"R"},"text":"psst"}`))

	msg, err := b.Send(ChatID(-100), "psst", &EphemeralMessageParameters{ReceiverUserID: 7})
	require.NoError(t, err)
	assert.Equal(t, 4, msg.EphemeralID)
	assert.JSONEq(t, `{"receiver_user_id":7}`, srv.Last().Params["ephemeral_message_parameters"])

	// The returned message can be edited right away.
	require.NoError(t, b.EditEphemeral(msg, "edited"))
	assert.Equal(t, "4", srv.Last().Params["ephemeral_message_id"])

	srv.Reply("sendPhoto", json.RawMessage(`{"message_id":1,"chat":{"id":-100},"photo":[{"file_id":"p","width":1,"height":1}]}`))
	_, err = b.Send(ChatID(-100), &Photo{File: File{FileID: "p"}}, &SendOptions{
		Ephemeral: &EphemeralMessageParameters{CallbackQueryID: "cb", ReplaceCallbackQueryMessage: true},
	})
	require.NoError(t, err)
	req := srv.Last()
	assert.Equal(t, "sendPhoto", req.Method)
	assert.JSONEq(t, `{"callback_query_id":"cb","replace_callback_query_message":true}`, req.Params["ephemeral_message_parameters"])

	srv.Reply("sendRichMessage", json.RawMessage(richReply))
	_, err = b.Send(ChatID(-100), &InputRichMessage{HTML: "<b>x</b>"}, &EphemeralMessageParameters{ReceiverUserID: 7})
	require.NoError(t, err)
	assert.JSONEq(t, `{"receiver_user_id":7}`, srv.Last().Params["ephemeral_message_parameters"])
}

func TestEphemeralLive(t *testing.T) {
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}
	groupID, _ := strconv.ParseInt(os.Getenv("GROUP_ID"), 10, 64)
	if groupID == 0 || userID == 0 {
		t.Skip("GROUP_ID and USER_ID are required for ephemeral live tests")
	}
	group := &Chat{ID: groupID}

	msg, err := b.Send(group, "telebot ephemeral", &EphemeralMessageParameters{ReceiverUserID: userID})
	require.NoError(t, err)
	require.NotZero(t, msg.EphemeralID)
	require.NotNil(t, msg.ReceiverUser)
	assert.Equal(t, userID, msg.ReceiverUser.ID)

	markup := &ReplyMarkup{InlineKeyboard: [][]InlineButton{{{Text: "ok", Data: "ok"}}}}
	require.NoError(t, b.EditEphemeral(msg, "telebot ephemeral, edited", markup))
	require.NoError(t, b.EditEphemeral(msg, &InputRichMessage{Markdown: "**telebot** ephemeral, rich"}))
	require.NoError(t, b.EditEphemeralReplyMarkup(msg, nil))

	require.NoError(t, b.DeleteEphemeral(msg))
}
