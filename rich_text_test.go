package telebot

import (
	"encoding/json"
	"testing"
)

func TestRichTextMarshalPlain(t *testing.T) {
	rt := RichText{Kind: RichTextPlain, Plain: "hello world"}
	b, err := json.Marshal(rt)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `"hello world"` {
		t.Fatalf("expected %q, got %s", "hello world", b)
	}
}

func TestRichTextMarshalPlainEmpty(t *testing.T) {
	rt := RichText{Kind: RichTextPlain, Plain: ""}
	b, err := json.Marshal(rt)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `""` {
		t.Fatalf("expected empty string, got %s", b)
	}
}

func TestRichTextMarshalArray(t *testing.T) {
	rt := RichText{
		Kind: RichTextArray,
		Parts: []RichText{
			{Kind: RichTextPlain, Plain: "hello "},
			{Kind: RichTextPlain, Plain: "world"},
		},
	}
	b, err := json.Marshal(rt)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `["hello ","world"]` {
		t.Fatalf("unexpected: %s", b)
	}
}

func TestRichTextMarshalEntity(t *testing.T) {
	inner := RichText{Kind: RichTextPlain, Plain: "bold text"}
	rt := RichText{
		Kind: RichTextEntity,
		Type: RichBold,
		Text: &inner,
	}
	b, err := json.Marshal(rt)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m["type"] != "bold" {
		t.Fatalf("expected type=bold, got %v", m["type"])
	}
	if m["text"] != "bold text" {
		t.Fatalf("expected text=bold text, got %v", m["text"])
	}
}

func TestRichTextRoundTrip(t *testing.T) {
	original := RichText{
		Kind: RichTextArray,
		Parts: []RichText{
			{Kind: RichTextPlain, Plain: "hello "},
			{
				Kind: RichTextEntity,
				Type: RichBold,
				Text: &RichText{Kind: RichTextPlain, Plain: "world"},
			},
		},
	}
	b, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}

	var decoded RichText
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Kind != RichTextArray {
		t.Fatalf("expected array, got %d", decoded.Kind)
	}
	if len(decoded.Parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(decoded.Parts))
	}
	if decoded.Parts[0].Plain != "hello " {
		t.Fatalf("part 0: %q", decoded.Parts[0].Plain)
	}
	if decoded.Parts[1].Type != RichBold {
		t.Fatalf("part 1 type: %q", decoded.Parts[1].Type)
	}
}
