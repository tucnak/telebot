package telebot

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const richReply = `{"message_id":5,"chat":{"id":2,"type":"private"}}`

func plain(s string) RichText { return RichText{Plain: s} }

func TestSendRichBlocks(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendRichMessage", json.RawMessage(richReply))

	msg, err := b.Send(ChatID(2), &InputRichMessage{
		IsRTL: true,
		Blocks: []InputRichBlock{
			InputRichBlockSectionHeading{Text: plain("Title"), Size: 1},
			InputRichBlockParagraph{Text: RichText{Kind: RichTextArray, Parts: []RichText{
				plain("a "),
				{Kind: RichTextEntity, Type: RichBold, Text: &RichText{Plain: "b"}},
			}}},
			InputRichBlockDivider{},
			InputRichBlockList{Items: []InputRichBlockListItem{
				{Blocks: []InputRichBlock{InputRichBlockParagraph{Text: plain("one")}}, HasCheckbox: true, IsChecked: true},
				{Blocks: []InputRichBlock{InputRichBlockParagraph{Text: plain("two")}}, Value: 2, Type: "1"},
			}},
			InputRichBlockTable{
				Cells:      [][]RichBlockTableCell{{{Text: &RichText{Plain: "x"}, Align: "left", VAlign: "top"}}},
				IsBordered: true,
				Caption:    &RichText{Plain: "t"},
			},
			InputRichBlockDetails{Summary: plain("more"), Blocks: []InputRichBlock{InputRichBlockFooter{Text: plain("f")}}},
			InputRichBlockMap{Location: Location{Lat: 1, Lng: 2}, Zoom: 10},
			InputRichBlockCollage{
				Blocks: []InputRichBlock{
					InputRichBlockPhoto{Photo: &Photo{File: File{FileID: "photo-id"}}},
					InputRichBlockVideo{Video: &Video{File: FromReader(strings.NewReader("vid")), FileName: "v.mp4"}},
				},
				Caption: &RichBlockCaption{Text: plain("cap"), Credit: &RichText{Plain: "me"}},
			},
			InputRichBlockVoiceNote{VoiceNote: &Voice{File: FromURL("https://example.com/v.ogg"), Duration: 3}},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 5, msg.ID)

	req := srv.Last()
	assert.Equal(t, "sendRichMessage", req.Method)
	require.Contains(t, req.Files, "rich0")
	assert.Equal(t, "vid", string(req.Files["rich0"].Data))

	var rich struct {
		IsRTL  bool                     `json:"is_rtl"`
		Blocks []map[string]interface{} `json:"blocks"`
	}
	req.JSON(t, "rich_message", &rich)
	assert.True(t, rich.IsRTL)
	require.Len(t, rich.Blocks, 9)

	var types []string
	for _, bl := range rich.Blocks {
		types = append(types, bl["type"].(string))
	}
	assert.Equal(t, []string{"heading", "paragraph", "divider", "list", "table", "details", "map", "collage", "voice_note"}, types)

	assert.Equal(t, map[string]interface{}{"type": "heading", "text": "Title", "size": 1.0}, rich.Blocks[0])
	assert.Equal(t, []interface{}{"a ", map[string]interface{}{"type": "bold", "text": "b"}}, rich.Blocks[1]["text"])
	assert.Equal(t, map[string]interface{}{"type": "divider"}, rich.Blocks[2])

	items := rich.Blocks[3]["items"].([]interface{})
	assert.Equal(t, map[string]interface{}{
		"blocks":       []interface{}{map[string]interface{}{"type": "paragraph", "text": "one"}},
		"has_checkbox": true,
		"is_checked":   true,
	}, items[0])
	assert.Equal(t, 2.0, items[1].(map[string]interface{})["value"])

	assert.Equal(t, "t", rich.Blocks[4]["caption"])
	assert.Equal(t, true, rich.Blocks[4]["is_bordered"])
	assert.Equal(t, []interface{}{[]interface{}{map[string]interface{}{"text": "x", "align": "left", "valign": "top"}}}, rich.Blocks[4]["cells"])

	assert.Equal(t, "more", rich.Blocks[5]["summary"])
	assert.Equal(t, map[string]interface{}{"latitude": 1.0, "longitude": 2.0}, rich.Blocks[6]["location"])

	collage := rich.Blocks[7]
	assert.Equal(t, map[string]interface{}{"text": "cap", "credit": "me"}, collage["caption"])
	inner := collage["blocks"].([]interface{})
	photo := inner[0].(map[string]interface{})["photo"].(map[string]interface{})
	assert.Equal(t, "photo", photo["type"])
	assert.Equal(t, "photo-id", photo["media"])
	video := inner[1].(map[string]interface{})["video"].(map[string]interface{})
	assert.Equal(t, "attach://rich0", video["media"])

	voice := rich.Blocks[8]["voice_note"].(map[string]interface{})
	assert.Equal(t, "voice_note", voice["type"])
	assert.Equal(t, "https://example.com/v.ogg", voice["media"])
	assert.Equal(t, 3.0, voice["duration"])
}

func TestSendRichMedia(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendRichMessage", json.RawMessage(richReply))

	_, err := b.Send(ChatID(2), &InputRichMessage{
		HTML: `<img src="tg://photo?id=p1">`,
		Media: []InputRichMessageMedia{
			{ID: "p1", Media: &Photo{File: FromReader(strings.NewReader("img"))}},
			{ID: "a1", Media: &Voice{File: File{FileID: "v"}, Duration: 3}},
		},
	})
	require.NoError(t, err)

	req := srv.Last()
	var rich struct {
		HTML  string `json:"html"`
		Media []struct {
			ID    string                 `json:"id"`
			Media map[string]interface{} `json:"media"`
		} `json:"media"`
	}
	req.JSON(t, "rich_message", &rich)
	assert.Equal(t, `<img src="tg://photo?id=p1">`, rich.HTML)
	require.Len(t, rich.Media, 2)
	assert.Equal(t, "p1", rich.Media[0].ID)
	assert.Equal(t, "photo", rich.Media[0].Media["type"])
	assert.Equal(t, "attach://rich0", rich.Media[0].Media["media"])
	assert.Equal(t, "img", string(req.Files["rich0"].Data))
	assert.Equal(t, "voice_note", rich.Media[1].Media["type"])
	assert.Equal(t, "v", rich.Media[1].Media["media"])

	// A voice note is input media in rich messages only.
	_, err = b.Edit(&Message{ID: 1, Chat: &Chat{ID: 2}}, &Voice{File: File{FileID: "v"}})
	assert.Equal(t, ErrUnsupportedWhat, err)
	_, err = b.Send(ChatID(2), &InputRichMessage{Media: []InputRichMessageMedia{{ID: "s", Media: &Sticker{File: File{FileID: "s"}}}}})
	assert.Error(t, err)
}

func TestSendRichDraftBlocks(t *testing.T) {
	b, srv := newFakeBot(t)

	err := b.SendRichDraft(ChatID(2), 9, &InputRichMessage{Blocks: []InputRichBlock{
		InputRichBlockThinking{Text: plain("Thinking…")},
	}})
	require.NoError(t, err)

	req := srv.Last()
	assert.Equal(t, "sendRichMessageDraft", req.Method)
	assert.Equal(t, "9", req.Params["draft_id"])
	assert.JSONEq(t, `{"blocks":[{"type":"thinking","text":"Thinking…"}]}`, req.Params["rich_message"])
}

func TestEditRichBlocks(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("editMessageText", json.RawMessage(richReply))

	_, err := b.Edit(&Message{ID: 5, Chat: &Chat{ID: 2}}, &InputRichMessage{Blocks: []InputRichBlock{
		InputRichBlockPullQuotation{Text: plain("q"), Credit: &RichText{Plain: "c"}},
		InputRichBlockBlockQuotation{Blocks: []InputRichBlock{InputRichBlockAnchor{Name: "a"}}},
		InputRichBlockPreformatted{Text: plain("x := 1"), Language: "go"},
		InputRichBlockMathematicalExpression{Expression: "e^{i\\pi}"},
		InputRichBlockSlideshow{Blocks: []InputRichBlock{
			InputRichBlockAnimation{Animation: &Animation{File: File{FileID: "gif"}}},
			InputRichBlockAudio{Audio: &Audio{File: File{FileID: "mp3"}}},
		}},
	}})
	require.NoError(t, err)

	req := srv.Last()
	assert.Equal(t, "editMessageText", req.Method)
	assert.Equal(t, "5", req.Params["message_id"])

	var rich struct {
		Blocks []json.RawMessage `json:"blocks"`
	}
	req.JSON(t, "rich_message", &rich)
	require.Len(t, rich.Blocks, 5)
	assert.JSONEq(t, `{"type":"pullquote","text":"q","credit":"c"}`, string(rich.Blocks[0]))
	assert.JSONEq(t, `{"type":"blockquote","blocks":[{"type":"anchor","name":"a"}]}`, string(rich.Blocks[1]))
	assert.JSONEq(t, `{"type":"pre","text":"x := 1","language":"go"}`, string(rich.Blocks[2]))
	assert.JSONEq(t, `{"type":"mathematical_expression","expression":"e^{i\\pi}"}`, string(rich.Blocks[3]))
	assert.Contains(t, string(rich.Blocks[4]), `"media":"gif"`)
	assert.Contains(t, string(rich.Blocks[4]), `"media":"mp3"`)
}

func TestInputRichMessageMarshalJSON(t *testing.T) {
	data, err := json.Marshal(&InputRichMessageContent{RichMessage: &InputRichMessage{
		Blocks: []InputRichBlock{InputRichBlockPhoto{Photo: &Photo{File: File{FileID: "id"}}}},
	}})
	require.NoError(t, err)
	assert.Contains(t, string(data), `"media":"id"`)
	assert.Contains(t, string(data), `"type":"photo"`)

	_, err = json.Marshal(&InputRichMessage{
		Blocks: []InputRichBlock{InputRichBlockPhoto{Photo: &Photo{File: FromDisk("a.png")}}},
	})
	assert.Error(t, err, "local files cannot be uploaded in inline results")

	_, err = json.Marshal(&InputRichMessage{Blocks: []InputRichBlock{InputRichBlockPhoto{}}})
	assert.Error(t, err)
}

func TestRichTextMarshalRoundTrip(t *testing.T) {
	for _, in := range []string{
		`"plain"`,
		`["a ",{"type":"italic","text":"b"}]`,
		`{"type":"url","text":"site","url":"https://example.com"}`,
		`{"type":"custom_emoji","custom_emoji_id":"1","alternative_text":":)"}`,
		`{"type":"anchor_link","text":"up","anchor_name":"a"}`,
		`{"type":"anchor_link","text":"top","anchor_name":""}`, // back to the top
	} {
		var r RichText
		require.NoError(t, json.Unmarshal([]byte(in), &r))
		out, err := json.Marshal(r)
		require.NoError(t, err)
		assert.JSONEq(t, in, string(out))
	}
}

func TestRichTextWithoutKind(t *testing.T) {
	bold := RichText{Type: RichBold, Text: &RichText{Plain: "b"}}
	run := RichText{Parts: []RichText{{Plain: "a "}, bold}}

	data, err := json.Marshal(run)
	require.NoError(t, err)
	assert.JSONEq(t, `["a ",{"type":"bold","text":"b"}]`, string(data))
	assert.Equal(t, "a b", run.String())
}

func TestRichBlocksLive(t *testing.T) {
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}
	if userID == 0 {
		t.Skip("USER_ID is required for rich message live tests")
	}

	msg, err := b.Send(user, &InputRichMessage{Blocks: []InputRichBlock{
		InputRichBlockSectionHeading{Text: plain("telebot"), Size: 2},
		InputRichBlockParagraph{Text: RichText{Kind: RichTextArray, Parts: []RichText{
			plain("Rich blocks "),
			{Kind: RichTextEntity, Type: RichBold, Text: &RichText{Plain: "work"}},
		}}},
		InputRichBlockList{Items: []InputRichBlockListItem{
			{Blocks: []InputRichBlock{InputRichBlockParagraph{Text: plain("one")}}},
			{Blocks: []InputRichBlock{InputRichBlockParagraph{Text: plain("two")}}, HasCheckbox: true, IsChecked: true},
		}},
		InputRichBlockTable{Cells: [][]RichBlockTableCell{{
			{Text: &RichText{Plain: "a"}, IsHeader: true, Align: "left", VAlign: "top"},
			{Text: &RichText{Plain: "b"}, Align: "right", VAlign: "middle"},
		}}, IsBordered: true},
		InputRichBlockPhoto{Photo: &Photo{File: logo}, Caption: &RichBlockCaption{Text: plain("by URL")}},
		InputRichBlockAudio{Audio: &Audio{File: FromReader(strings.NewReader("ID3")), FileName: "a.mp3"}},
	}})
	require.NoError(t, err)
	require.NotNil(t, msg.RichMessage)

	var types []string
	for _, bl := range msg.RichMessage.Blocks {
		types = append(types, bl.Type)
	}
	assert.Equal(t, []string{"heading", "paragraph", "list", "table", "photo", "audio"}, types)
}
