package telebot

import (
	"encoding/json"
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
