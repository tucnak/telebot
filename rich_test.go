package telebot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// richBot returns an offline bot whose requests hit srv and records the decoded
// request path and JSON body for assertions.
func richBot(t *testing.T, result string, path *string, body *map[string]interface{}) *Bot {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*path = r.URL.Path
		require.NoError(t, json.NewDecoder(r.Body).Decode(body))
		w.Write([]byte(result))
	}))
	t.Cleanup(srv.Close)

	b, err := NewBot(Settings{
		Offline:   true,
		Token:     "token",
		URL:       srv.URL,
		Client:    srv.Client(),
		ParseMode: ModeHTML, // proves parse_mode is never leaked onto rich calls
	})
	require.NoError(t, err)
	return b
}

func TestSendRich(t *testing.T) {
	var path string
	var body map[string]interface{}
	b := richBot(t, `{"ok":true,"result":{"message_id":42,"chat":{"id":2,"type":"private"}}}`, &path, &body)

	msg, err := b.Send(ChatID(2), &InputRichMessage{
		Markdown:            "# Title\n\n- a\n- b",
		IsRTL:               true,
		SkipEntityDetection: true,
	}, Silent)
	require.NoError(t, err)
	require.NotNil(t, msg)

	assert.Equal(t, "/bottoken/sendRichMessage", path)
	assert.Equal(t, 42, msg.ID)
	assert.Equal(t, "2", body["chat_id"])
	assert.Equal(t, "true", body["disable_notification"])
	assert.NotContains(t, body, "parse_mode")

	// Object params ride as JSON-encoded strings, as with reply_markup/media.
	var rich map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(body["rich_message"].(string)), &rich))
	assert.Equal(t, "# Title\n\n- a\n- b", rich["markdown"])
	assert.Equal(t, true, rich["is_rtl"])
	assert.Equal(t, true, rich["skip_entity_detection"])
	assert.NotContains(t, rich, "html")
}

func TestEditRich(t *testing.T) {
	var path string
	var body map[string]interface{}
	b := richBot(t, `{"ok":true,"result":{"message_id":7,"chat":{"id":2,"type":"private"}}}`, &path, &body)

	_, err := b.Edit(&Message{ID: 7, Chat: &Chat{ID: 2}}, &InputRichMessage{HTML: "<b>done</b>"})
	require.NoError(t, err)

	assert.Equal(t, "/bottoken/editMessageText", path)
	assert.Equal(t, "7", body["message_id"])
	assert.Equal(t, "2", body["chat_id"])
	assert.NotContains(t, body, "text")
	assert.NotContains(t, body, "parse_mode")

	var rich map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(body["rich_message"].(string)), &rich))
	assert.Equal(t, "<b>done</b>", rich["html"])
}

func TestSendRichDraft(t *testing.T) {
	var path string
	var body map[string]interface{}
	b := richBot(t, `{"ok":true,"result":true}`, &path, &body)

	err := b.SendRichDraft(ChatID(2), 99, &InputRichMessage{Markdown: "thinking…"})
	require.NoError(t, err)

	assert.Equal(t, "/bottoken/sendRichMessageDraft", path)
	assert.Equal(t, "99", body["draft_id"])
	assert.Equal(t, "2", body["chat_id"])
	assert.Contains(t, body, "rich_message")
}

func TestRichBadInput(t *testing.T) {
	b, err := NewBot(Settings{Offline: true, Token: "token"})
	require.NoError(t, err)

	_, err = b.Send(nil, &InputRichMessage{Markdown: "x"})
	assert.ErrorIs(t, err, ErrBadRecipient)

	assert.ErrorIs(t, b.SendRichDraft(nil, 1, &InputRichMessage{Markdown: "x"}), ErrBadRecipient)
	assert.ErrorIs(t, b.SendRichDraft(ChatID(2), 1, nil), ErrUnsupportedWhat)
}

func TestInputRichMessageContent(t *testing.T) {
	var _ InputMessageContent = (*InputRichMessageContent)(nil)

	data, err := json.Marshal(&InputRichMessageContent{
		RichMessage: &InputRichMessage{Markdown: "hi"},
	})
	require.NoError(t, err)
	assert.JSONEq(t, `{"rich_message":{"markdown":"hi"}}`, string(data))
}

func TestRichMessageUnmarshal(t *testing.T) {
	data := []byte(`{
		"message_id": 7,
		"chat": {"id": 2, "type": "private"},
		"rich_message": {
			"is_rtl": false,
			"blocks": [
				{"type": "heading", "text": "Quarterly Report", "size": 1},
				{"type": "paragraph", "text": ["Revenue is ", {"type": "bold", "text": "up 20%"}, "."]},
				{"type": "list", "items": [
					{"label": "1", "blocks": [{"type": "paragraph", "text": "First"}], "has_checkbox": true, "is_checked": true}
				]},
				{"type": "table", "is_bordered": true, "caption": "Results", "cells": [
					[{"text": "Q1", "is_header": true}, {"text": "Q2", "is_header": true}],
					[{"text": "100"}, {"text": "120"}]
				]},
				{"type": "photo", "photo": [{"file_id": "abc", "width": 100, "height": 100}], "caption": {"text": "A chart", "credit": "Acme"}}
			]
		}
	}`)

	var msg Message
	require.NoError(t, json.Unmarshal(data, &msg))
	require.NotNil(t, msg.RichMessage)
	rm := msg.RichMessage
	require.Len(t, rm.Blocks, 5)

	// heading: tagged-string text
	assert.Equal(t, RichBlockHeading, rm.Blocks[0].Type)
	assert.Equal(t, 1, rm.Blocks[0].Size)
	assert.Equal(t, "Quarterly Report", rm.Blocks[0].Text.String())

	// paragraph: array text mixing plain and a bold entity
	para := rm.Blocks[1].Text
	require.Equal(t, RichTextArray, para.Kind)
	assert.Equal(t, "Revenue is up 20%.", para.String())

	// list: task item
	require.Len(t, rm.Blocks[2].Items, 1)
	item := rm.Blocks[2].Items[0]
	assert.True(t, item.HasCheckbox)
	assert.True(t, item.IsChecked)
	assert.Equal(t, "First", item.Blocks[0].Text.String())

	// table: bare-RichText caption + header cells
	tbl := rm.Blocks[3]
	assert.True(t, tbl.IsBordered)
	require.NotNil(t, tbl.TableCaption)
	assert.Equal(t, "Results", tbl.TableCaption.String())
	assert.Nil(t, tbl.Caption)
	require.Len(t, tbl.Cells, 2)
	assert.True(t, tbl.Cells[0][0].IsHeader)
	assert.Equal(t, "Q1", tbl.Cells[0][0].Text.String())

	// photo: RichBlockCaption caption + media
	photo := rm.Blocks[4]
	require.Len(t, photo.Photo, 1)
	require.NotNil(t, photo.Caption)
	assert.Equal(t, "A chart", photo.Caption.Text.String())
	require.NotNil(t, photo.Caption.Credit)
	assert.Equal(t, "Acme", photo.Caption.Credit.String())
}
