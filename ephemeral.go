package telebot

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Ephemeral messages (Bot API 10.2) are messages in a group that only one
// user and the bot can see. They are identified by their chat, the user who
// received them and Message.EphemeralID rather than Message.ID, so they are
// edited and deleted with their own methods, which take the ephemeral message
// as received or as returned when it was sent.

// EphemeralMessageParameters makes a sent message ephemeral (Bot API 10.3).
// Pass it to Send, or set it as SendOptions.Ephemeral:
//
//	b.Send(group, "only you can see this", &tele.EphemeralMessageParameters{
//		ReceiverUserID: user.ID,
//	})
//
// To reply ephemerally to a callback query in any chat, set CallbackQueryID
// and leave ReceiverUserID zero; the chat passed to Send is then ignored.
//
// Bot API: EphemeralMessageParameters.
type EphemeralMessageParameters struct {
	// (Optional) The user who will receive the message; required unless the
	// message replies to a callback query.
	ReceiverUserID int64 `json:"receiver_user_id,omitempty"`

	// (Optional) The callback query that triggered the message.
	CallbackQueryID string `json:"callback_query_id,omitempty"`

	// (Optional) Show the message in place of the callback query's message.
	// Must be false for callback queries from ephemeral messages.
	ReplaceCallbackQueryMessage bool `json:"replace_callback_query_message,omitempty"`
}

// ephemeralParams returns the parameters identifying an ephemeral message.
func ephemeralParams(msg *Message) (map[string]string, error) {
	if msg == nil || msg.Chat == nil || msg.ReceiverUser == nil || msg.EphemeralID == 0 {
		return nil, errors.New("telebot: not an ephemeral message")
	}
	return map[string]string{
		"chat_id":              msg.Chat.Recipient(),
		"receiver_user_id":     msg.ReceiverUser.Recipient(),
		"ephemeral_message_id": strconv.Itoa(msg.EphemeralID),
	}, nil
}

// EditEphemeral edits the text of an ephemeral message, or replaces it with
// rich content when what is an *InputRichMessage. It also takes an
// Inputtable, which it passes to EditEphemeralMedia, and a *ReplyMarkup,
// which it passes to EditEphemeralReplyMarkup. Of the options, the parse
// mode, entities, link preview options and the inline keyboard apply.
//
// Telegram does not guarantee that the user sees the edit.
//
// Bot API: editEphemeralMessageText.
func (b *Bot) EditEphemeral(msg *Message, what interface{}, opts ...interface{}) error {
	switch v := what.(type) {
	case *ReplyMarkup:
		return b.EditEphemeralReplyMarkup(msg, v)
	case Inputtable:
		return b.EditEphemeralMedia(msg, v, opts...)
	}

	params, err := ephemeralParams(msg)
	if err != nil {
		return err
	}

	var files map[string]File
	switch v := what.(type) {
	case string:
		params["text"] = v
	case *InputRichMessage:
		rich, f, err := encodeRich(v)
		if err != nil {
			return err
		}
		params["rich_message"] = rich
		files = f
	default:
		return ErrUnsupportedWhat
	}

	sendOpts := b.extractOptions(opts)
	b.embedEphemeralOptions(params, sendOpts)

	preview := sendOpts.PreviewOptions
	if preview == nil && sendOpts.DisableWebPagePreview {
		preview = &PreviewOptions{Disabled: true}
	}
	if preview != nil {
		data, _ := json.Marshal(preview)
		params["link_preview_options"] = string(data)
	}
	if _, ok := params["rich_message"]; ok {
		delete(params, "parse_mode")
		delete(params, "entities")
	}

	_, err = b.sendFiles("editEphemeralMessageText", files, params)
	return err
}

// EditEphemeralCaption edits the caption of an ephemeral message. Of the
// options, the parse mode, entities, CaptionAbove and the inline keyboard
// apply.
//
// Bot API: editEphemeralMessageCaption.
func (b *Bot) EditEphemeralCaption(msg *Message, caption string, opts ...interface{}) error {
	params, err := ephemeralParams(msg)
	if err != nil {
		return err
	}
	params["caption"] = caption

	b.embedEphemeralOptions(params, b.extractOptions(opts))

	_, err = b.Raw("editEphemeralMessageCaption", params)
	return err
}

// EditEphemeralMedia replaces the media of an ephemeral message. Local files
// are uploaded.
//
// Bot API: editEphemeralMessageMedia.
func (b *Bot) EditEphemeralMedia(msg *Message, media Inputtable, opts ...interface{}) error {
	params, err := ephemeralParams(msg)
	if err != nil {
		return err
	}

	sendOpts := b.extractOptions(opts)
	repr, files, err := encodeEditMedia(media, sendOpts)
	if err != nil {
		return err
	}
	params["media"] = repr
	b.embedEphemeralOptions(params, sendOpts)
	delete(params, "parse_mode")
	delete(params, "entities")

	_, err = b.sendFiles("editEphemeralMessageMedia", files, params)
	return err
}

// EditEphemeralReplyMarkup replaces the inline keyboard of an ephemeral
// message; a nil markup removes it.
//
// Bot API: editEphemeralMessageReplyMarkup.
func (b *Bot) EditEphemeralReplyMarkup(msg *Message, markup *ReplyMarkup) error {
	params, err := ephemeralParams(msg)
	if err != nil {
		return err
	}

	if markup == nil {
		markup = &ReplyMarkup{}
	}
	processButtons(markup.InlineKeyboard)
	data, _ := json.Marshal(markup)
	params["reply_markup"] = string(data)

	_, err = b.Raw("editEphemeralMessageReplyMarkup", params)
	return err
}

// DeleteEphemeral deletes an ephemeral message. Telegram does not guarantee
// that the user sees the deletion.
//
// Bot API: deleteEphemeralMessage.
func (b *Bot) DeleteEphemeral(msg *Message) error {
	params, err := ephemeralParams(msg)
	if err != nil {
		return err
	}

	_, err = b.Raw("deleteEphemeralMessage", params)
	return err
}

// embedEphemeralOptions adds the options the ephemeral edit methods accept.
func (b *Bot) embedEphemeralOptions(params map[string]string, opt *SendOptions) {
	if opt == nil {
		return
	}

	if opt.ParseMode != ModeDefault {
		params["parse_mode"] = opt.ParseMode
	}
	if len(opt.Entities) > 0 {
		delete(params, "parse_mode")
		entities, _ := json.Marshal(opt.Entities)
		if _, ok := params["caption"]; ok {
			params["caption_entities"] = string(entities)
		} else {
			params["entities"] = string(entities)
		}
	}
	if opt.CaptionAbove {
		params["show_caption_above_media"] = "true"
	}
	if opt.ReplyMarkup != nil {
		processButtons(opt.ReplyMarkup.InlineKeyboard)
		data, _ := json.Marshal(opt.ReplyMarkup)
		params["reply_markup"] = string(data)
	}
}
