package telebot

import (
	"strconv"
)

// InputRichMessage describes a rich message to send with the Bot API 10.1
// rich-message methods (sendRichMessage, sendRichMessageDraft, and the
// rich_message parameter of editMessageText). Exactly one of HTML, Markdown or
// Blocks must be set. Telegram parses HTML and Markdown server-side into its
// rich block tree — headings, tables, ordered and nested lists, blockquotes,
// dividers, media blocks and so on — and renders it as a single message
// regardless of length; Blocks gives that tree directly.
//
// InputRichMessage implements Sendable, so the idiomatic way to send one is
// through Send/Reply:
//
//	b.Send(to, &telebot.InputRichMessage{Markdown: "# Title\n\nBody"})
//
// and to edit a message into rich content through Edit:
//
//	b.Edit(msg, &telebot.InputRichMessage{HTML: "<b>done</b>"})
//
// Bot API: InputRichMessage.
type InputRichMessage struct {
	// (Optional) Content described using Telegram's extended HTML formatting.
	HTML string `json:"html,omitempty"`

	// (Optional) Content described using Telegram's extended Markdown formatting.
	Markdown string `json:"markdown,omitempty"`

	// (Optional) Pass true to render the message right-to-left.
	IsRTL bool `json:"is_rtl,omitempty"`

	// (Optional) Pass true to skip automatic detection of URLs, email addresses,
	// mentions, hashtags, cashtags, bot commands and phone numbers in the text.
	SkipEntityDetection bool `json:"skip_entity_detection,omitempty"`

	// (Optional) Content described as a list of blocks (Bot API 10.2).
	Blocks []InputRichBlock `json:"blocks,omitempty"`

	// (Optional) Media referenced from HTML or Markdown by tg://photo?id=,
	// tg://video?id=, tg://document?id= and tg://audio?id= links (Bot API 10.2).
	Media []InputRichMessageMedia `json:"media,omitempty"`
}

// Send delivers the rich message through bot b to recipient via the
// sendRichMessage method, satisfying the Sendable interface.
func (i *InputRichMessage) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	rich, files, err := encodeRich(i)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"chat_id":      to.Recipient(),
		"rich_message": rich,
	}
	b.embedSendOptions(params, opt)

	// Parse mode and entities apply to plain text, never to rich content; the
	// markup is carried inside the InputRichMessage itself.
	delete(params, "parse_mode")
	delete(params, "entities")

	raw, err := b.sendFiles("sendRichMessage", files, params)
	if err != nil {
		return nil, err
	}

	return extractMessage(raw)
}

// SendRichDraft streams a partial rich message to a private chat while it is
// being generated (the sendRichMessageDraft method). The streamed draft is
// ephemeral and acts as a temporary 30-second preview; once the output is
// finalized you must Send the complete InputRichMessage to persist it. draftID
// must be non-zero — updates sharing the same identifier are animated
// client-side. Of the send options only ThreadID, CanStop and KeepOnStop are
// honored.
//
// Bot API: sendRichMessageDraft.
func (b *Bot) SendRichDraft(to Recipient, draftID int, rich *InputRichMessage, opts ...interface{}) error {
	if to == nil {
		return ErrBadRecipient
	}
	if rich == nil {
		return ErrUnsupportedWhat
	}

	data, files, err := encodeRich(rich)
	if err != nil {
		return err
	}

	params := map[string]string{
		"chat_id":      to.Recipient(),
		"draft_id":     strconv.Itoa(draftID),
		"rich_message": data,
	}

	if sendOpts := b.extractOptions(opts); sendOpts != nil {
		if sendOpts.ThreadID != 0 {
			params["message_thread_id"] = strconv.Itoa(sendOpts.ThreadID)
		}
		embedDraftStop(params, sendOpts)
	}

	_, err = b.sendFiles("sendRichMessageDraft", files, params)
	return err
}

// InputRichMessageContent represents the content of a rich message to be sent
// as the result of an inline, guest or Web App query (Bot API 10.1).
//
// Bot API: InputRichMessageContent.
type InputRichMessageContent struct {
	// The rich message to be sent.
	RichMessage *InputRichMessage `json:"rich_message"`
}

// IsInputMessageContent marks InputRichMessageContent as InputMessageContent.
func (i *InputRichMessageContent) IsInputMessageContent() bool {
	return true
}

// RichMessage represents a received rich formatted message (Bot API 10.1),
// available on Message.RichMessage. Its content is a tree of RichBlock values.
//
// Bot API: RichMessage.
type RichMessage struct {
	// Content of the message, as an ordered list of blocks.
	Blocks []RichBlock `json:"blocks"`

	// (Optional) True if the message must be shown right-to-left.
	IsRTL bool `json:"is_rtl,omitempty"`
}

// RichMessageButton is a button in a rich message (Bot API 10.3), shown in a
// buttons block or inline in text. Exactly one of the fields other than
// Text, Style and Unique must be set.
//
// Like InlineButton, a button with Unique set reaches the handler registered
// with b.Handle(&button), its Data becoming the callback's data.
//
// Bot API: RichMessageButton.
type RichMessageButton struct {
	// Unique names the callback endpoint, as in InlineButton.
	Unique string `json:"-"`

	// Text of the button: plain text, custom emoji and date-time entities.
	Text RichText `json:"text"`

	// (Optional) "danger", "success", "primary" or "link".
	Style ButtonStyle `json:"style,omitempty"`

	URL                   string             `json:"url,omitempty"`
	Data                  string             `json:"callback_data,omitempty"`
	WebApp                *WebApp            `json:"web_app,omitempty"`
	Login                 *Login             `json:"login_url,omitempty"`
	InlineQuery           *string            `json:"switch_inline_query,omitempty"`
	InlineQueryChat       *string            `json:"switch_inline_query_current_chat,omitempty"`
	InlineQueryChosenChat *SwitchInlineQuery `json:"switch_inline_query_chosen_chat,omitempty"`
	CopyText              *CopyTextButton    `json:"copy_text,omitempty"`
	Disabled              *DisabledButton    `json:"disabled,omitempty"`
}

// CallbackUnique returns the endpoint of the button for b.Handle.
func (t *RichMessageButton) CallbackUnique() string {
	return "\f" + t.Unique
}

// wire returns the button as Telegram expects it, its Unique folded into
// the callback data.
func (t RichMessageButton) wire() RichMessageButton {
	if t.Unique != "" {
		if t.Data == "" {
			t.Data = "\f" + t.Unique
		} else {
			t.Data = "\f" + t.Unique + "|" + t.Data
		}
	}
	return t
}

// MessageGenerationStopped is received when a user presses the stop button of
// a draft sent with the CanStop option (Bot API 10.3).
//
// Bot API: MessageGenerationStopped.
type MessageGenerationStopped struct {
	// Chat in which the message is generated.
	Chat *Chat `json:"chat"`

	// (Optional) Message thread in which the message is generated.
	ThreadID int `json:"message_thread_id,omitempty"`

	// The draft that was stopped, as passed to SendDraft or SendRichDraft.
	DraftID int `json:"draft_id"`
}
