package telebot

import (
	"encoding/json"
	"errors"
	"strconv"
)

// InputRichBlock is a block of a rich message to send, set in
// InputRichMessage.Blocks. It is implemented by the InputRichBlock* types of
// this package; the "type" field Telegram expects is filled in for you.
//
//	&tele.InputRichMessage{Blocks: []tele.InputRichBlock{
//		tele.InputRichBlockSectionHeading{Text: tele.RichText{Plain: "Report"}, Size: 1},
//		tele.InputRichBlockParagraph{Text: tele.RichText{Plain: "All systems go."}},
//		tele.InputRichBlockPhoto{Photo: &tele.Photo{File: tele.FromDisk("chart.png")}},
//	}}
//
// Media in blocks may be file IDs, URLs or local files; local files are
// uploaded with the message.
//
// Bot API: InputRichBlock.
type InputRichBlock interface {
	encodeRichBlock(e *richEncoder) (richWire, error)
}

// InputRichBlockParagraph is a text paragraph, like the HTML tag <p>.
//
// Bot API: InputRichBlockParagraph.
type InputRichBlockParagraph struct {
	Text RichText `json:"text"`
}

// InputRichBlockSectionHeading is a section heading, like the HTML tags <h1>
// to <h6>. Size is 1 (largest) to 6 (smallest).
//
// Bot API: InputRichBlockSectionHeading.
type InputRichBlockSectionHeading struct {
	Text RichText `json:"text"`
	Size int      `json:"size"`
}

// InputRichBlockPreformatted is a preformatted block, like the nested HTML
// tags <pre> and <code>. Language is optional.
//
// Bot API: InputRichBlockPreformatted.
type InputRichBlockPreformatted struct {
	Text     RichText `json:"text"`
	Language string   `json:"language,omitempty"`
}

// InputRichBlockFooter is a footer, like the HTML tag <footer>.
//
// Bot API: InputRichBlockFooter.
type InputRichBlockFooter struct {
	Text RichText `json:"text"`
}

// InputRichBlockDivider is a divider, like the HTML tag <hr/>.
//
// Bot API: InputRichBlockDivider.
type InputRichBlockDivider struct{}

// InputRichBlockMathematicalExpression is a LaTeX expression shown as a block.
//
// Bot API: InputRichBlockMathematicalExpression.
type InputRichBlockMathematicalExpression struct {
	Expression string `json:"expression"`
}

// InputRichBlockAnchor is an anchor, like the HTML tag <a> with a name.
//
// Bot API: InputRichBlockAnchor.
type InputRichBlockAnchor struct {
	Name string `json:"name"`
}

// InputRichBlockList is a list, like the HTML tags <ul> or <ol>.
//
// Bot API: InputRichBlockList.
type InputRichBlockList struct {
	Items []InputRichBlockListItem `json:"items"`
}

// InputRichBlockListItem is an item of an InputRichBlockList.
//
// Bot API: InputRichBlockListItem.
type InputRichBlockListItem struct {
	// Content of the item.
	Blocks []InputRichBlock `json:"blocks"`

	// (Optional) The item has a checkbox, checked if IsChecked.
	HasCheckbox bool `json:"has_checkbox,omitempty"`
	IsChecked   bool `json:"is_checked,omitempty"`

	// (Optional) For ordered lists, the numeric value of the item label.
	Value int `json:"value,omitempty"`

	// (Optional) For ordered lists, the label style: "a", "A", "i", "I" or "1".
	Type string `json:"type,omitempty"`
}

// InputRichBlockBlockQuotation is a block quotation, like the HTML tag
// <blockquote>. Credit is optional.
//
// Bot API: InputRichBlockBlockQuotation.
type InputRichBlockBlockQuotation struct {
	Blocks []InputRichBlock `json:"blocks"`
	Credit *RichText        `json:"credit,omitempty"`
}

// InputRichBlockPullQuotation is a quotation with centered text, like the
// HTML tag <aside>. Credit is optional.
//
// Bot API: InputRichBlockPullQuotation.
type InputRichBlockPullQuotation struct {
	Text   RichText  `json:"text"`
	Credit *RichText `json:"credit,omitempty"`
}

// InputRichBlockCollage is a collage of blocks, like the custom HTML tag
// <tg-collage>.
//
// Bot API: InputRichBlockCollage.
type InputRichBlockCollage struct {
	Blocks  []InputRichBlock  `json:"blocks"`
	Caption *RichBlockCaption `json:"caption,omitempty"`
}

// InputRichBlockSlideshow is a slideshow of blocks, like the custom HTML tag
// <tg-slideshow>.
//
// Bot API: InputRichBlockSlideshow.
type InputRichBlockSlideshow struct {
	Blocks  []InputRichBlock  `json:"blocks"`
	Caption *RichBlockCaption `json:"caption,omitempty"`
}

// InputRichBlockTable is a table, like the HTML tag <table>. Every cell needs
// Align and VAlign.
//
// Bot API: InputRichBlockTable.
type InputRichBlockTable struct {
	Cells      [][]RichBlockTableCell `json:"cells"`
	IsBordered bool                   `json:"is_bordered,omitempty"`
	IsStriped  bool                   `json:"is_striped,omitempty"`
	IsCompact  bool                   `json:"is_compact,omitempty"` // Bot API 10.3
	Caption    *RichText              `json:"caption,omitempty"`
}

// InputRichBlockDetails is an expandable block, like the HTML tag <details>.
//
// Bot API: InputRichBlockDetails.
type InputRichBlockDetails struct {
	Summary RichText         `json:"summary"`
	Blocks  []InputRichBlock `json:"blocks"`
	IsOpen  bool             `json:"is_open,omitempty"`
}

// InputRichBlockMap is a map centered on Location, like the custom HTML tag
// <tg-map>. Zoom is 0-24; Width and Height are optional.
//
// Bot API: InputRichBlockMap.
type InputRichBlockMap struct {
	Location Location          `json:"location"`
	Zoom     int               `json:"zoom,omitempty"`
	Width    int               `json:"width,omitempty"`
	Height   int               `json:"height,omitempty"`
	Caption  *RichBlockCaption `json:"caption,omitempty"`
}

// InputRichBlockAnimation is an animation block. The animation's own caption
// is ignored; use Caption.
//
// Bot API: InputRichBlockAnimation.
type InputRichBlockAnimation struct {
	Animation *Animation        `json:"animation"`
	Caption   *RichBlockCaption `json:"caption,omitempty"`
}

// InputRichBlockAudio is a music block. The audio's own caption is ignored;
// use Caption.
//
// Bot API: InputRichBlockAudio.
type InputRichBlockAudio struct {
	Audio   *Audio            `json:"audio"`
	Caption *RichBlockCaption `json:"caption,omitempty"`
}

// InputRichBlockPhoto is a photo block, like the HTML tag <img>. The photo's
// own caption is ignored; use Caption.
//
// Bot API: InputRichBlockPhoto.
type InputRichBlockPhoto struct {
	Photo   *Photo            `json:"photo"`
	Caption *RichBlockCaption `json:"caption,omitempty"`
}

// InputRichBlockVideo is a video block, like the HTML tag <video>. The
// video's own caption is ignored; use Caption.
//
// Bot API: InputRichBlockVideo.
type InputRichBlockVideo struct {
	Video   *Video            `json:"video"`
	Caption *RichBlockCaption `json:"caption,omitempty"`
}

// InputRichBlockVoiceNote is a voice note block. The voice note's own caption
// is ignored; use Caption.
//
// Bot API: InputRichBlockVoiceNote.
type InputRichBlockVoiceNote struct {
	VoiceNote *Voice            `json:"voice_note"`
	Caption   *RichBlockCaption `json:"caption,omitempty"`
}

// InputRichBlockThinking is a "Thinking…" placeholder, like the custom HTML
// tag <tg-thinking>. It can be used only in SendRichDraft.
//
// Bot API: InputRichBlockThinking.
type InputRichBlockThinking struct {
	Text RichText `json:"text"`
}

// InputRichBlockExpandableBlockQuotation is a block quotation that can be
// expanded and collapsed, like the HTML tag <blockquote expandable>.
// Credit is optional (Bot API 10.3).
//
// Bot API: InputRichBlockExpandableBlockQuotation.
type InputRichBlockExpandableBlockQuotation struct {
	Text   RichText  `json:"text"`
	Credit *RichText `json:"credit,omitempty"`
}

// InputRichBlockButtons is a row of 1-8 buttons, like the custom HTML tag
// <tg-button-row>. Align is "left", "center" or "right" (Bot API 10.3).
//
// Bot API: InputRichBlockButtons.
type InputRichBlockButtons struct {
	Buttons []RichMessageButton `json:"buttons"`
	Align   string              `json:"align,omitempty"`
}

// InputRichBlockDocument is a block with a general file, like the custom
// HTML tag <tg-document>. The document's own caption is ignored; use Caption
// (Bot API 10.3).
//
// Bot API: InputRichBlockDocument.
type InputRichBlockDocument struct {
	Document *Document         `json:"document"`
	Caption  *RichBlockCaption `json:"caption,omitempty"`
}

// InputRichMessageMedia is a media file referenced from InputRichMessage.HTML
// or Markdown by a tg://photo?id=, tg://video?id=, tg://document?id= or
// tg://audio?id= link with the same ID. Media may be an *Animation, *Audio,
// *Document, *Photo, *Video or *Voice; only the file and its properties are
// used.
//
// Bot API: InputRichMessageMedia.
type InputRichMessageMedia struct {
	ID    string `json:"id"`
	Media Media  `json:"media"`
}

// richWire is a block or media object ready for encoding/json.
type richWire map[string]interface{}

// set adds key unless v is the zero value of its type.
func (w richWire) set(key string, v interface{}) {
	switch x := v.(type) {
	case string:
		if x == "" {
			return
		}
	case int:
		if x == 0 {
			return
		}
	case bool:
		if !x {
			return
		}
	case *RichText:
		if x == nil {
			return
		}
	case *RichBlockCaption:
		if x == nil {
			return
		}
	}
	w[key] = v
}

// richEncoder turns rich blocks into their wire form, collecting the local
// files they upload.
type richEncoder struct {
	files map[string]File
}

func (e *richEncoder) blocks(bs []InputRichBlock) ([]richWire, error) {
	out := make([]richWire, 0, len(bs))
	for _, b := range bs {
		if b == nil {
			return nil, errors.New("telebot: nil rich block")
		}
		w, err := b.encodeRichBlock(e)
		if err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, nil
}

// media encodes m as an InputMedia object, uploading it if it is local.
func (e *richEncoder) media(m Media) (*InputMedia, error) {
	if m == nil || m.MediaFile() == nil {
		return nil, errors.New("telebot: rich message media is missing")
	}

	var im InputMedia
	switch v := m.(type) {
	case *Voice:
		// Only rich messages take a voice note as input media, so Voice is
		// not an Inputtable.
		im = InputMedia{Type: "voice_note", Caption: v.Caption, Duration: v.Duration}
	case *LivePhoto:
		return nil, errors.New("telebot: unsupported rich message media")
	case Inputtable:
		im = v.InputMedia()
	default:
		return nil, errors.New("telebot: unsupported rich message media")
	}

	im.Media = m.MediaFile().process("rich"+strconv.Itoa(len(e.files)), e.files)
	if im.Media == "" {
		return nil, errors.New("telebot: rich message media does not exist")
	}
	return &im, nil
}

func (e *richEncoder) mediaBlock(typ, key string, m Media, caption *RichBlockCaption) (richWire, error) {
	im, err := e.media(m)
	if err != nil {
		return nil, err
	}
	w := richWire{"type": typ, key: im}
	w.set("caption", caption)
	return w, nil
}

// encode returns the wire form of the message.
func (e *richEncoder) message(r *InputRichMessage) (richWire, error) {
	w := richWire{}
	w.set("html", r.HTML)
	w.set("markdown", r.Markdown)
	w.set("is_rtl", r.IsRTL)
	w.set("skip_entity_detection", r.SkipEntityDetection)
	if r.Blocks != nil {
		blocks, err := e.blocks(r.Blocks)
		if err != nil {
			return nil, err
		}
		w["blocks"] = blocks
	}
	if r.Media != nil {
		media := make([]richWire, 0, len(r.Media))
		for _, m := range r.Media {
			im, err := e.media(m.Media)
			if err != nil {
				return nil, err
			}
			media = append(media, richWire{"id": m.ID, "media": im})
		}
		w["media"] = media
	}
	return w, nil
}

// encodeRich returns the rich_message parameter and the files it uploads.
func encodeRich(r *InputRichMessage) (string, map[string]File, error) {
	e := &richEncoder{files: map[string]File{}}
	w, err := e.message(r)
	if err != nil {
		return "", nil, err
	}
	data, err := json.Marshal(w)
	return string(data), e.files, err
}

// MarshalJSON encodes the message as Telegram expects it in inline, guest
// and Web App results. Media must then be file IDs or URLs, as nothing can be
// uploaded alongside.
func (r InputRichMessage) MarshalJSON() ([]byte, error) {
	e := &richEncoder{files: map[string]File{}}
	w, err := e.message(&r)
	if err != nil {
		return nil, err
	}
	if len(e.files) > 0 {
		return nil, errors.New("telebot: rich message media must be file IDs or URLs here")
	}
	return json.Marshal(w)
}

func textBlock(typ string, text RichText) richWire {
	return richWire{"type": typ, "text": text}
}

func (b InputRichBlockParagraph) encodeRichBlock(*richEncoder) (richWire, error) {
	return textBlock(RichBlockParagraph, b.Text), nil
}

func (b InputRichBlockSectionHeading) encodeRichBlock(*richEncoder) (richWire, error) {
	return richWire{"type": RichBlockHeading, "text": b.Text, "size": b.Size}, nil
}

func (b InputRichBlockPreformatted) encodeRichBlock(*richEncoder) (richWire, error) {
	w := textBlock(RichBlockPre, b.Text)
	w.set("language", b.Language)
	return w, nil
}

func (b InputRichBlockFooter) encodeRichBlock(*richEncoder) (richWire, error) {
	return textBlock(RichBlockFooter, b.Text), nil
}

func (b InputRichBlockDivider) encodeRichBlock(*richEncoder) (richWire, error) {
	return richWire{"type": RichBlockDivider}, nil
}

func (b InputRichBlockMathematicalExpression) encodeRichBlock(*richEncoder) (richWire, error) {
	return richWire{"type": RichBlockMath, "expression": b.Expression}, nil
}

func (b InputRichBlockAnchor) encodeRichBlock(*richEncoder) (richWire, error) {
	return richWire{"type": RichBlockAnchor, "name": b.Name}, nil
}

func (b InputRichBlockList) encodeRichBlock(e *richEncoder) (richWire, error) {
	items := make([]richWire, 0, len(b.Items))
	for _, it := range b.Items {
		blocks, err := e.blocks(it.Blocks)
		if err != nil {
			return nil, err
		}
		w := richWire{"blocks": blocks}
		w.set("has_checkbox", it.HasCheckbox)
		w.set("is_checked", it.IsChecked)
		w.set("value", it.Value)
		w.set("type", it.Type)
		items = append(items, w)
	}
	return richWire{"type": RichBlockList, "items": items}, nil
}

func (b InputRichBlockBlockQuotation) encodeRichBlock(e *richEncoder) (richWire, error) {
	blocks, err := e.blocks(b.Blocks)
	if err != nil {
		return nil, err
	}
	w := richWire{"type": RichBlockBlockquote, "blocks": blocks}
	w.set("credit", b.Credit)
	return w, nil
}

func (b InputRichBlockPullQuotation) encodeRichBlock(*richEncoder) (richWire, error) {
	w := textBlock(RichBlockPullquote, b.Text)
	w.set("credit", b.Credit)
	return w, nil
}

func (b InputRichBlockCollage) encodeRichBlock(e *richEncoder) (richWire, error) {
	blocks, err := e.blocks(b.Blocks)
	if err != nil {
		return nil, err
	}
	w := richWire{"type": RichBlockCollage, "blocks": blocks}
	w.set("caption", b.Caption)
	return w, nil
}

func (b InputRichBlockSlideshow) encodeRichBlock(e *richEncoder) (richWire, error) {
	blocks, err := e.blocks(b.Blocks)
	if err != nil {
		return nil, err
	}
	w := richWire{"type": RichBlockSlideshow, "blocks": blocks}
	w.set("caption", b.Caption)
	return w, nil
}

func (b InputRichBlockTable) encodeRichBlock(*richEncoder) (richWire, error) {
	w := richWire{"type": RichBlockTable, "cells": b.Cells}
	w.set("is_bordered", b.IsBordered)
	w.set("is_striped", b.IsStriped)
	w.set("is_compact", b.IsCompact)
	w.set("caption", b.Caption)
	return w, nil
}

func (b InputRichBlockDetails) encodeRichBlock(e *richEncoder) (richWire, error) {
	blocks, err := e.blocks(b.Blocks)
	if err != nil {
		return nil, err
	}
	w := richWire{"type": RichBlockDetails, "summary": b.Summary, "blocks": blocks}
	w.set("is_open", b.IsOpen)
	return w, nil
}

func (b InputRichBlockMap) encodeRichBlock(*richEncoder) (richWire, error) {
	w := richWire{"type": RichBlockMap, "location": b.Location}
	w.set("zoom", b.Zoom)
	w.set("width", b.Width)
	w.set("height", b.Height)
	w.set("caption", b.Caption)
	return w, nil
}

func (b InputRichBlockAnimation) encodeRichBlock(e *richEncoder) (richWire, error) {
	if b.Animation == nil {
		return nil, errors.New("telebot: rich animation block has no animation")
	}
	return e.mediaBlock(RichBlockAnimationType, "animation", b.Animation, b.Caption)
}

func (b InputRichBlockAudio) encodeRichBlock(e *richEncoder) (richWire, error) {
	if b.Audio == nil {
		return nil, errors.New("telebot: rich audio block has no audio")
	}
	return e.mediaBlock(RichBlockAudioType, "audio", b.Audio, b.Caption)
}

func (b InputRichBlockPhoto) encodeRichBlock(e *richEncoder) (richWire, error) {
	if b.Photo == nil {
		return nil, errors.New("telebot: rich photo block has no photo")
	}
	return e.mediaBlock(RichBlockPhotoType, "photo", b.Photo, b.Caption)
}

func (b InputRichBlockVideo) encodeRichBlock(e *richEncoder) (richWire, error) {
	if b.Video == nil {
		return nil, errors.New("telebot: rich video block has no video")
	}
	return e.mediaBlock(RichBlockVideoType, "video", b.Video, b.Caption)
}

func (b InputRichBlockVoiceNote) encodeRichBlock(e *richEncoder) (richWire, error) {
	if b.VoiceNote == nil {
		return nil, errors.New("telebot: rich voice note block has no voice note")
	}
	return e.mediaBlock(RichBlockVoiceNote, "voice_note", b.VoiceNote, b.Caption)
}

func (b InputRichBlockThinking) encodeRichBlock(*richEncoder) (richWire, error) {
	return textBlock(RichBlockThinking, b.Text), nil
}

func (b InputRichBlockExpandableBlockQuotation) encodeRichBlock(*richEncoder) (richWire, error) {
	w := textBlock(RichBlockExpandableBlockquote, b.Text)
	w.set("credit", b.Credit)
	return w, nil
}

func (b InputRichBlockButtons) encodeRichBlock(*richEncoder) (richWire, error) {
	buttons := make([]RichMessageButton, len(b.Buttons))
	for i, btn := range b.Buttons {
		buttons[i] = btn.wire()
	}
	w := richWire{"type": RichBlockButtons, "buttons": buttons}
	w.set("align", b.Align)
	return w, nil
}

func (b InputRichBlockDocument) encodeRichBlock(e *richEncoder) (richWire, error) {
	if b.Document == nil {
		return nil, errors.New("telebot: rich document block has no document")
	}
	return e.mediaBlock(RichBlockDocumentType, "document", b.Document, b.Caption)
}
