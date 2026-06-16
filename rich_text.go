package telebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// RichTextKind tells which of the three wire forms a RichText holds.
type RichTextKind int

const (
	// RichTextPlain is a bare string of plain text.
	RichTextPlain RichTextKind = iota

	// RichTextArray is an ordered run of RichText values (mixed formatting).
	RichTextArray

	// RichTextEntity is a single tagged formatting entity (see RichText.Type).
	RichTextEntity
)

// RichText type discriminators, as carried in the "type" field of a tagged
// RichText entity (Bot API 10.1).
const (
	RichBold                   = "bold"
	RichItalic                 = "italic"
	RichUnderline              = "underline"
	RichStrikethrough          = "strikethrough"
	RichSpoiler                = "spoiler"
	RichDateTime               = "date_time"
	RichTextMentionType        = "text_mention"
	RichSubscript              = "subscript"
	RichSuperscript            = "superscript"
	RichMarked                 = "marked"
	RichCode                   = "code"
	RichCustomEmoji            = "custom_emoji"
	RichMathematicalExpression = "mathematical_expression"
	RichURL                    = "url"
	RichEmailAddress           = "email_address"
	RichPhoneNumber            = "phone_number"
	RichBankCardNumber         = "bank_card_number"
	RichMention                = "mention"
	RichHashtag                = "hashtag"
	RichCashtag                = "cashtag"
	RichBotCommand             = "bot_command"
	RichAnchor                 = "anchor"
	RichAnchorLink             = "anchor_link"
	RichReference              = "reference"
	RichReferenceLink          = "reference_link"
)

// RichText is a polymorphic piece of rich-formatted text (Bot API 10.1). On the
// wire it is one of: a plain string, an array of RichText, or a tagged object
// carrying a "type" discriminator and a (usually nested) text. Kind reports
// which form was decoded; the relevant fields are populated accordingly.
//
// RichText is a received-only type — rich content is sent as HTML or Markdown
// via InputRichMessage, never as a block tree.
type RichText struct {
	// Kind is the wire form this value was decoded from.
	Kind RichTextKind `json:"-"`

	// Plain holds the text when Kind is RichTextPlain.
	Plain string `json:"-"`

	// Parts holds the run when Kind is RichTextArray.
	Parts []RichText `json:"-"`

	// Type is the entity discriminator when Kind is RichTextEntity.
	Type string `json:"type,omitempty"`

	// Text is the nested, formatted text of the entity (most entity types).
	Text *RichText `json:"text,omitempty"`

	// URL is the target of a "url" entity.
	URL string `json:"url,omitempty"`

	// EmailAddress is the address of an "email_address" entity.
	EmailAddress string `json:"email_address,omitempty"`

	// PhoneNumber is the number of a "phone_number" entity.
	PhoneNumber string `json:"phone_number,omitempty"`

	// BankCardNumber is the number of a "bank_card_number" entity.
	BankCardNumber string `json:"bank_card_number,omitempty"`

	// Username is the @username of a "mention" entity.
	Username string `json:"username,omitempty"`

	// Hashtag is the tag of a "hashtag" entity.
	Hashtag string `json:"hashtag,omitempty"`

	// Cashtag is the tag of a "cashtag" entity.
	Cashtag string `json:"cashtag,omitempty"`

	// BotCommand is the command of a "bot_command" entity.
	BotCommand string `json:"bot_command,omitempty"`

	// User is the mentioned user of a "text_mention" entity.
	User *User `json:"user,omitempty"`

	// UnixTime is the moment of a "date_time" entity, in Unix time.
	UnixTime int64 `json:"unix_time,omitempty"`

	// DateTimeFormat is the rendering format of a "date_time" entity.
	DateTimeFormat string `json:"date_time_format,omitempty"`

	// CustomEmojiID is the identifier of a "custom_emoji" entity.
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`

	// AlternativeText is the fallback text of a "custom_emoji" entity.
	AlternativeText string `json:"alternative_text,omitempty"`

	// Expression is the LaTeX of a "mathematical_expression" entity.
	Expression string `json:"expression,omitempty"`

	// Name is the identifier of an "anchor" or "reference" entity.
	Name string `json:"name,omitempty"`

	// AnchorName is the target of an "anchor_link" entity (empty = top).
	AnchorName string `json:"anchor_name,omitempty"`

	// ReferenceName is the target of a "reference_link" entity.
	ReferenceName string `json:"reference_name,omitempty"`
}

// UnmarshalJSON decodes the string, array and tagged-object forms of RichText.
func (r *RichText) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	switch data[0] {
	case '"':
		r.Kind = RichTextPlain
		return json.Unmarshal(data, &r.Plain)
	case '[':
		r.Kind = RichTextArray
		return json.Unmarshal(data, &r.Parts)
	case '{':
		r.Kind = RichTextEntity
		type alias RichText
		return json.Unmarshal(data, (*alias)(r))
	default:
		return wrapError(fmt.Errorf("telebot: invalid RichText value: %s", data))
	}
}

// String renders the RichText as plain text, recursively flattening arrays and
// entities. Entities without textual content (anchors, dividers) contribute
// nothing; custom emoji contribute their alternative text.
func (r RichText) String() string {
	switch r.Kind {
	case RichTextPlain:
		return r.Plain
	case RichTextArray:
		var b strings.Builder
		for i := range r.Parts {
			b.WriteString(r.Parts[i].String())
		}
		return b.String()
	default: // RichTextEntity
		switch r.Type {
		case RichCustomEmoji:
			return r.AlternativeText
		case RichMathematicalExpression:
			return r.Expression
		case RichAnchor:
			return ""
		default:
			if r.Text != nil {
				return r.Text.String()
			}
			return ""
		}
	}
}
