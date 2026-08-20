package telebot

import "encoding/json"

// RichBlock type discriminators, as carried in the "type" field of a received
// rich block (Bot API 10.1).
const (
	RichBlockParagraph     = "paragraph"
	RichBlockHeading       = "heading"
	RichBlockPre           = "pre"
	RichBlockFooter        = "footer"
	RichBlockDivider       = "divider"
	RichBlockMath          = "mathematical_expression"
	RichBlockAnchor        = "anchor"
	RichBlockList          = "list"
	RichBlockBlockquote    = "blockquote"
	RichBlockPullquote     = "pullquote"
	RichBlockCollage       = "collage"
	RichBlockSlideshow     = "slideshow"
	RichBlockTable         = "table"
	RichBlockDetails       = "details"
	RichBlockMap           = "map"
	RichBlockAnimationType = "animation"
	RichBlockAudioType     = "audio"
	RichBlockPhotoType     = "photo"
	RichBlockVideoType     = "video"
	RichBlockVoiceNote     = "voice_note"
	RichBlockThinking      = "thinking"
)

// RichBlock is a single block of a received RichMessage (Bot API 10.1). It is a
// polymorphic type discriminated by Type; only the fields relevant to a given
// Type are populated. RichBlock is a received-only type.
type RichBlock struct {
	// Type is the block discriminator (see the RichBlock* constants).
	Type string `json:"type"`

	// Text is the block's formatted text (paragraph, heading, pre, footer,
	// pullquote, thinking).
	Text *RichText `json:"text,omitempty"`

	// Size is the heading level, 1 (largest) to 6 (smallest).
	Size int `json:"size,omitempty"`

	// Language is the optional source language of a "pre" block.
	Language string `json:"language,omitempty"`

	// Expression is the LaTeX of a "mathematical_expression" block.
	Expression string `json:"expression,omitempty"`

	// Name is the identifier of an "anchor" block.
	Name string `json:"name,omitempty"`

	// Items are the entries of a "list" block.
	Items []RichBlockListItem `json:"items,omitempty"`

	// Blocks are the nested blocks of a blockquote, collage, slideshow or
	// details block.
	Blocks []RichBlock `json:"blocks,omitempty"`

	// Summary is the disclosure summary of a "details" block.
	Summary *RichText `json:"summary,omitempty"`

	// IsOpen reports whether a "details" block is expanded by default.
	IsOpen bool `json:"is_open,omitempty"`

	// Credit is the attribution of a blockquote or pullquote block.
	Credit *RichText `json:"credit,omitempty"`

	// Cells are the rows-by-columns of a "table" block.
	Cells [][]RichBlockTableCell `json:"cells,omitempty"`

	// IsBordered and IsStriped style a "table" block.
	IsBordered bool `json:"is_bordered,omitempty"`
	IsStriped  bool `json:"is_striped,omitempty"`

	// Location, Zoom, Width and Height describe a "map" block.
	Location *Location `json:"location,omitempty"`
	Zoom     int       `json:"zoom,omitempty"`
	Width    int       `json:"width,omitempty"`
	Height   int       `json:"height,omitempty"`

	// Media of the corresponding block type.
	Animation *Animation `json:"animation,omitempty"`
	Audio     *Audio     `json:"audio,omitempty"`
	Photo     []Photo    `json:"photo,omitempty"`
	Video     *Video     `json:"video,omitempty"`
	VoiceNote *Voice     `json:"voice_note,omitempty"`

	// HasSpoiler marks an animation, photo or video block as a spoiler.
	HasSpoiler bool `json:"has_spoiler,omitempty"`

	// Caption is the caption of a collage, slideshow, map, animation, audio,
	// photo, video or voice-note block.
	Caption *RichBlockCaption `json:"-"`

	// TableCaption is the caption of a "table" block (a bare RichText rather
	// than a RichBlockCaption).
	TableCaption *RichText `json:"-"`
}

// UnmarshalJSON decodes a RichBlock, resolving the polymorphic "caption" field:
// a table carries a bare RichText caption, every other block a RichBlockCaption.
func (rb *RichBlock) UnmarshalJSON(data []byte) error {
	type alias RichBlock
	aux := struct {
		*alias
		Caption json.RawMessage `json:"caption,omitempty"`
	}{alias: (*alias)(rb)}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.Caption) == 0 || string(aux.Caption) == "null" {
		return nil
	}

	if rb.Type == RichBlockTable {
		rb.TableCaption = new(RichText)
		return json.Unmarshal(aux.Caption, rb.TableCaption)
	}

	rb.Caption = new(RichBlockCaption)
	return json.Unmarshal(aux.Caption, rb.Caption)
}

// RichBlockCaption is the caption of a media or container rich block.
type RichBlockCaption struct {
	// Text of the caption.
	Text RichText `json:"text"`

	// (Optional) Attribution shown alongside the caption.
	Credit *RichText `json:"credit,omitempty"`
}

// RichBlockTableCell is a single cell of a "table" rich block. An omitted Text
// denotes an invisible cell spanned into by a neighbour.
type RichBlockTableCell struct {
	Text     *RichText `json:"text,omitempty"`
	IsHeader bool      `json:"is_header,omitempty"`
	Colspan  int       `json:"colspan,omitempty"`
	Rowspan  int       `json:"rowspan,omitempty"`
	Align    string    `json:"align,omitempty"`  // "left", "center", "right"
	VAlign   string    `json:"valign,omitempty"` // "top", "middle", "bottom"
}

// RichBlockListItem is a single entry of a "list" rich block.
type RichBlockListItem struct {
	// Label is the rendered bullet or number of the item.
	Label string `json:"label,omitempty"`

	// Blocks is the content of the item.
	Blocks []RichBlock `json:"blocks,omitempty"`

	// HasCheckbox and IsChecked render a task-list item.
	HasCheckbox bool `json:"has_checkbox,omitempty"`
	IsChecked   bool `json:"is_checked,omitempty"`

	// Value is the ordinal of an ordered-list item.
	Value int `json:"value,omitempty"`

	// Type is the label style of an ordered item ("a", "A", "i", "I", "1").
	Type string `json:"type,omitempty"`
}
