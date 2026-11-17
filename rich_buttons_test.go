package telebot

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendRichBlocks103(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendRichMessage", json.RawMessage(richReply))

	empty := ""
	_, err := b.Send(ChatID(2), &InputRichMessage{Blocks: []InputRichBlock{
		InputRichBlockButtons{Align: "center", Buttons: []RichMessageButton{
			{Text: plain("Go"), Unique: "go", Data: "1", Style: "primary"},
			{Text: plain("Site"), URL: "https://example.com"},
			{Text: plain("Here"), InlineQueryChat: &empty},
			{Text: plain("Off"), Disabled: &DisabledButton{}},
		}},
		InputRichBlockExpandableBlockQuotation{Text: plain("long"), Credit: &RichText{Plain: "c"}},
		InputRichBlockDocument{Document: &Document{File: FromReader(strings.NewReader("doc")), FileName: "a.txt"}},
		InputRichBlockTable{Cells: [][]RichBlockTableCell{}, IsCompact: true},
		InputRichBlockParagraph{Text: RichText{Kind: RichTextEntity, Type: RichButton,
			Button: &RichMessageButton{Text: plain("inline"), Unique: "in"}}},
	}})
	require.NoError(t, err)

	req := srv.Last()
	var rich struct {
		Blocks []json.RawMessage `json:"blocks"`
	}
	req.JSON(t, "rich_message", &rich)
	require.Len(t, rich.Blocks, 5)

	assert.JSONEq(t, `{"type":"buttons","align":"center","buttons":[
		{"text":"Go","style":"primary","callback_data":"\fgo|1"},
		{"text":"Site","url":"https://example.com"},
		{"text":"Here","switch_inline_query_current_chat":""},
		{"text":"Off","disabled":{}}]}`, string(rich.Blocks[0]))
	assert.JSONEq(t, `{"type":"expandable_blockquote","text":"long","credit":"c"}`, string(rich.Blocks[1]))
	assert.JSONEq(t, `{"type":"document","document":{"type":"document","media":"attach://rich0","caption":""}}`, string(rich.Blocks[2]))
	assert.Equal(t, "doc", string(req.Files["rich0"].Data))
	assert.JSONEq(t, `{"type":"table","cells":[],"is_compact":true}`, string(rich.Blocks[3]))
	assert.JSONEq(t, `{"type":"paragraph","text":{"type":"button","button":{"text":"inline","callback_data":"\fin"}}}`, string(rich.Blocks[4]))
}

func TestReceiveRichBlocks103(t *testing.T) {
	var msg RichMessage
	require.NoError(t, json.Unmarshal([]byte(`{"blocks":[
		{"type":"buttons","align":"right","buttons":[{"text":"Pay","callback_data":"x","style":"success"}]},
		{"type":"expandable_blockquote","text":"t","credit":"c"},
		{"type":"document","document":{"file_id":"f","file_unique_id":"u","file_name":"a.pdf"},"caption":{"text":"cap"}},
		{"type":"table","cells":[],"is_compact":true,"caption":"tc"},
		{"type":"paragraph","text":{"type":"button","button":{"text":"b","url":"https://example.com"}}}
	]}`), &msg))
	require.Len(t, msg.Blocks, 5)

	bt := msg.Blocks[0]
	assert.Equal(t, RichBlockButtons, bt.Type)
	assert.Equal(t, "right", bt.Align)
	require.Len(t, bt.Buttons, 1)
	assert.Equal(t, "Pay", bt.Buttons[0].Text.String())
	assert.Equal(t, "x", bt.Buttons[0].Data)
	assert.Equal(t, ButtonStyle("success"), bt.Buttons[0].Style)

	assert.Equal(t, RichBlockExpandableBlockquote, msg.Blocks[1].Type)
	assert.Equal(t, "t", msg.Blocks[1].Text.String())
	assert.Equal(t, "c", msg.Blocks[1].Credit.String())

	doc := msg.Blocks[2]
	require.NotNil(t, doc.Document)
	assert.Equal(t, "a.pdf", doc.Document.FileName)
	assert.Equal(t, "cap", doc.Caption.Text.String())

	assert.True(t, msg.Blocks[3].IsCompact)
	assert.Equal(t, "tc", msg.Blocks[3].TableCaption.String())

	text := msg.Blocks[4].Text
	assert.Equal(t, RichButton, text.Type)
	assert.Equal(t, "https://example.com", text.Button.URL)
}

func TestRichButtonHandler(t *testing.T) {
	b, _ := newFakeBot(t)

	btn := RichMessageButton{Unique: "go"}
	var data string
	b.Handle(&btn, func(c Context) error { data = c.Data(); return nil })

	b.ProcessUpdate(Update{Callback: &Callback{ID: "1", Data: "\fgo|42", Sender: &User{ID: 1}}})
	assert.Equal(t, "42", data)
}
