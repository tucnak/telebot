package telebot

import (
	"encoding/json"
	"testing"
)

func TestInputRichBlockParagraphMarshal(t *testing.T) {
	block := InputRichBlockParagraph{
		Type: "paragraph",
		Text: RichText{Kind: RichTextPlain, Plain: "hello"},
	}
	b, err := json.Marshal(block)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	if m["type"] != "paragraph" {
		t.Fatalf("type: %v", m["type"])
	}
	if m["text"] != "hello" {
		t.Fatalf("text: %v", m["text"])
	}
}

func TestInputRichBlockDetailsMarshal(t *testing.T) {
	block := InputRichBlockDetails{
		Type:    "details",
		Summary: RichText{Kind: RichTextPlain, Plain: "Thinking"},
		Blocks: []InputRichBlock{
			InputRichBlockParagraph{
				Type: "paragraph",
				Text: RichText{Kind: RichTextPlain, Plain: "inner content"},
			},
		},
		IsOpen: false,
	}
	b, err := json.Marshal(block)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	if m["type"] != "details" {
		t.Fatalf("type: %v", m["type"])
	}
	if m["summary"] != "Thinking" {
		t.Fatalf("summary: %v", m["summary"])
	}
	blocks := m["blocks"].([]interface{})
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
	inner := blocks[0].(map[string]interface{})
	if inner["type"] != "paragraph" {
		t.Fatalf("inner type: %v", inner["type"])
	}
}

func TestInputRichBlockOmitEmptyRichText(t *testing.T) {
	block := InputRichBlockExpandableBlockQuotation{
		Type: "expandable_blockquote",
		Text: RichText{Kind: RichTextPlain, Plain: "quote"},
	}
	b, err := json.Marshal(block)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	if _, exists := m["credit"]; exists {
		t.Fatal("credit should be omitted when nil")
	}
}

func TestInputRichBlockTableCaption(t *testing.T) {
	caption := RichText{Kind: RichTextPlain, Plain: "Table 1"}
	block := InputRichBlockTable{
		Type: "table",
		Cells: [][]RichBlockTableCell{
			{{Text: &RichText{Kind: RichTextPlain, Plain: "A"}}},
		},
		Caption: &caption,
	}
	b, err := json.Marshal(block)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	if m["caption"] != "Table 1" {
		t.Fatalf("caption: %v", m["caption"])
	}
}

func TestInputRichMessageBlocksMarshal(t *testing.T) {
	msg := InputRichMessage{
		Blocks: []InputRichBlock{
			InputRichBlockThinking{
				Type: "thinking",
				Text: RichText{Kind: RichTextPlain, Plain: "Thinking..."},
			},
			InputRichBlockParagraph{
				Type: "paragraph",
				Text: RichText{Kind: RichTextPlain, Plain: "The answer"},
			},
		},
	}
	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	blocks := m["blocks"].([]interface{})
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}
	first := blocks[0].(map[string]interface{})
	if first["type"] != "thinking" {
		t.Fatalf("first block type: %v", first["type"])
	}
	second := blocks[1].(map[string]interface{})
	if second["type"] != "paragraph" {
		t.Fatalf("second block type: %v", second["type"])
	}
}
