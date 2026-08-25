package telebot

// InputRichBlock is the interface for all block types that can appear in
// InputRichMessage.Blocks (Bot API 10.2).
type InputRichBlock interface {
	inputRichBlock()
}

// InputRichBlockParagraph is a paragraph block.
type InputRichBlockParagraph struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// InputRichBlockSectionHeading is a section heading block with a level 1-6.
type InputRichBlockSectionHeading struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
	Size int      `json:"size"`
}

// InputRichBlockPreformatted is a preformatted (code) block.
type InputRichBlockPreformatted struct {
	Type     string   `json:"type"`
	Text     RichText `json:"text"`
	Language string   `json:"language,omitempty"`
}

// InputRichBlockFooter is a footer block.
type InputRichBlockFooter struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// InputRichBlockDivider is a horizontal divider block.
type InputRichBlockDivider struct {
	Type string `json:"type"`
}

// InputRichBlockMathematicalExpression is a LaTeX math block.
type InputRichBlockMathematicalExpression struct {
	Type       string `json:"type"`
	Expression string `json:"expression"`
}

// InputRichBlockAnchor is a named anchor block for in-page linking.
type InputRichBlockAnchor struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// InputRichBlockList is an ordered or unordered list block.
type InputRichBlockList struct {
	Type      string                   `json:"type"`
	Items     []InputRichBlockListItem `json:"items"`
	IsOrdered bool                     `json:"is_ordered,omitempty"`
}

// InputRichBlockListItem is a single entry of an InputRichBlockList.
type InputRichBlockListItem struct {
	Blocks      []InputRichBlock `json:"blocks"`
	HasCheckbox bool             `json:"has_checkbox,omitempty"`
	IsChecked   bool             `json:"is_checked,omitempty"`
	Value       int              `json:"value,omitempty"`
	Type        string           `json:"type,omitempty"`
}

// InputRichBlockBlockQuotation is a block quotation.
type InputRichBlockBlockQuotation struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// InputRichBlockPullQuotation is a pull quotation.
type InputRichBlockPullQuotation struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// InputRichBlockCollage is a collage of media items.
type InputRichBlockCollage struct {
	Type  string       `json:"type"`
	Media []Inputtable `json:"media"`
}

// InputRichBlockSlideshow is a slideshow of media items.
type InputRichBlockSlideshow struct {
	Type  string       `json:"type"`
	Media []Inputtable `json:"media"`
}

// InputRichBlockTable is a table block.
type InputRichBlockTable struct {
	Type       string                 `json:"type"`
	Cells      [][]RichBlockTableCell `json:"cells"`
	IsBordered bool                   `json:"is_bordered,omitempty"`
	IsStriped  bool                   `json:"is_striped,omitempty"`
	IsCompact  bool                   `json:"is_compact,omitempty"`
	Caption    *RichText              `json:"caption,omitempty"`
}

// InputRichBlockDetails is a collapsible details block.
type InputRichBlockDetails struct {
	Type    string           `json:"type"`
	Summary RichText        `json:"summary"`
	Blocks  []InputRichBlock `json:"blocks"`
	IsOpen  bool             `json:"is_open,omitempty"`
}

// InputRichBlockMap is a map block showing a geographic location.
type InputRichBlockMap struct {
	Type     string    `json:"type"`
	Location *Location `json:"location"`
	Zoom     int       `json:"zoom,omitempty"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
}

// InputRichBlockAnimation is an animation (GIF) block.
type InputRichBlockAnimation struct {
	Type      string     `json:"type"`
	Animation Inputtable `json:"-"`
}

// InputRichBlockAudio is an audio block.
type InputRichBlockAudio struct {
	Type  string     `json:"type"`
	Audio Inputtable `json:"-"`
}

// InputRichBlockPhoto is a photo block.
type InputRichBlockPhoto struct {
	Type  string     `json:"type"`
	Photo Inputtable `json:"-"`
}

// InputRichBlockVideo is a video block.
type InputRichBlockVideo struct {
	Type  string     `json:"type"`
	Video Inputtable `json:"-"`
}

// InputRichBlockVoiceNote is a voice note block.
type InputRichBlockVoiceNote struct {
	Type      string     `json:"type"`
	VoiceNote Inputtable `json:"-"`
}

// InputRichBlockThinking is a thinking/reasoning block.
type InputRichBlockThinking struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// InputRichBlockExpandableBlockQuotation is an expandable block quotation
// (Bot API 10.3).
type InputRichBlockExpandableBlockQuotation struct {
	Type   string    `json:"type"`
	Text   RichText  `json:"text"`
	Credit *RichText `json:"credit,omitempty"`
}

// InputRichBlockDocument is a document block (Bot API 10.3).
type InputRichBlockDocument struct {
	Type     string     `json:"type"`
	Document Inputtable `json:"-"`
}

// Interface compliance — every concrete block type implements InputRichBlock.
func (InputRichBlockParagraph) inputRichBlock()                     {}
func (InputRichBlockSectionHeading) inputRichBlock()                {}
func (InputRichBlockPreformatted) inputRichBlock()                  {}
func (InputRichBlockFooter) inputRichBlock()                        {}
func (InputRichBlockDivider) inputRichBlock()                       {}
func (InputRichBlockMathematicalExpression) inputRichBlock()         {}
func (InputRichBlockAnchor) inputRichBlock()                        {}
func (InputRichBlockList) inputRichBlock()                          {}
func (InputRichBlockBlockQuotation) inputRichBlock()                {}
func (InputRichBlockPullQuotation) inputRichBlock()                 {}
func (InputRichBlockCollage) inputRichBlock()                       {}
func (InputRichBlockSlideshow) inputRichBlock()                     {}
func (InputRichBlockTable) inputRichBlock()                         {}
func (InputRichBlockDetails) inputRichBlock()                       {}
func (InputRichBlockMap) inputRichBlock()                           {}
func (InputRichBlockAnimation) inputRichBlock()                     {}
func (InputRichBlockAudio) inputRichBlock()                         {}
func (InputRichBlockPhoto) inputRichBlock()                         {}
func (InputRichBlockVideo) inputRichBlock()                         {}
func (InputRichBlockVoiceNote) inputRichBlock()                     {}
func (InputRichBlockThinking) inputRichBlock()                      {}
func (InputRichBlockExpandableBlockQuotation) inputRichBlock()      {}
func (InputRichBlockDocument) inputRichBlock()                      {}
