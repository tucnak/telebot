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
// client-side. Of the send options only ThreadID is honored.
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

	if sendOpts := b.extractOptions(opts); sendOpts != nil && sendOpts.ThreadID != 0 {
		params["message_thread_id"] = strconv.Itoa(sendOpts.ThreadID)
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
