package telebot

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDraftStop(t *testing.T) {
	b, srv := newFakeBot(t)
	stop := &SendOptions{CanStop: true, KeepOnStop: true, ThreadID: 3}

	require.NoError(t, b.SendDraft(ChatID(2), 9, "partial", stop))
	assert.Equal(t, map[string]string{
		"chat_id": "2", "draft_id": "9", "text": "partial", "message_thread_id": "3",
		"can_stop": "true", "keep_on_stop": "true",
	}, srv.Last().Params)

	require.NoError(t, b.SendRichDraft(ChatID(2), 9, &InputRichMessage{Markdown: "x"}, stop))
	req := srv.Last()
	assert.Equal(t, "sendRichMessageDraft", req.Method)
	assert.Equal(t, "true", req.Params["can_stop"])
	assert.Equal(t, "true", req.Params["keep_on_stop"])

	require.NoError(t, b.SendDraft(ChatID(2), 9, "partial"))
	assert.NotContains(t, srv.Last().Params, "can_stop")
}

func TestGenerationStopped(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendMessage", json.RawMessage(`{"message_id":1,"chat":{"id":2}}`))

	var got *MessageGenerationStopped
	b.Handle(OnGenerationStopped, func(c Context) error {
		got = c.GenerationStopped()
		return c.Send("stopped")
	})

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"stopped_message_generation":{
		"chat":{"id":2,"type":"private"},"message_thread_id":3,"draft_id":9}}`), &u))
	b.ProcessUpdate(u)

	require.NotNil(t, got)
	assert.Equal(t, 9, got.DraftID)
	assert.Equal(t, 3, got.ThreadID)
	assert.Equal(t, "2", srv.Last().Params["chat_id"])
	assert.Contains(t, AllowedUpdates, "stopped_message_generation")
}

func TestDraftsLive(t *testing.T) {
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}
	if userID == 0 {
		t.Skip("USER_ID is required for draft live tests")
	}

	stop := &SendOptions{CanStop: true, KeepOnStop: true}
	require.NoError(t, b.SendDraft(user, 1, "", stop))
	require.NoError(t, b.SendRichDraft(user, 1, &InputRichMessage{Blocks: []InputRichBlock{
		InputRichBlockThinking{Text: RichText{Plain: "Thinking…"}},
	}}, stop))

	msg, err := b.Send(user, &InputRichMessage{Blocks: []InputRichBlock{
		InputRichBlockExpandableBlockQuotation{Text: RichText{Plain: "A long quote"}, Credit: &RichText{Plain: "telebot"}},
		InputRichBlockDocument{Document: &Document{File: FromReader(strings.NewReader("telebot")), FileName: "telebot.txt"}},
		InputRichBlockTable{IsCompact: true, Cells: [][]RichBlockTableCell{{{Text: &RichText{Plain: "x"}, Align: "center", VAlign: "middle"}}}},
		InputRichBlockButtons{Align: "center", Buttons: []RichMessageButton{
			{Text: RichText{Plain: "Site"}, URL: "https://github.com/tucnak/telebot"},
			{Text: RichText{Plain: "Tap"}, Unique: "tap", Data: "1", Style: "primary"},
		}},
	}})
	require.NoError(t, err)
	require.NotNil(t, msg.RichMessage)

	var types []string
	for _, bl := range msg.RichMessage.Blocks {
		types = append(types, bl.Type)
	}
	assert.Equal(t, []string{"expandable_blockquote", "document", "table", "buttons"}, types)
	assert.Equal(t, "\ftap|1", msg.RichMessage.Blocks[3].Buttons[1].Data)
}
