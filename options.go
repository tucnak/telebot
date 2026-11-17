package telebot

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Option is a shortcut flag type for certain message features
// (so-called options). It means that instead of passing
// fully-fledged SendOptions* to Send(), you can use these
// flags instead.
//
// Supported options are defined as iota-constants.
type Option int

const (
	// NoPreview = SendOptions.DisableWebPagePreview
	NoPreview Option = iota

	// Silent = SendOptions.DisableNotification
	Silent

	// AllowWithoutReply = SendOptions.AllowWithoutReply
	AllowWithoutReply

	// Protected = SendOptions.Protected
	Protected

	// ForceReply = ReplyMarkup.ForceReply
	ForceReply

	// OneTimeKeyboard = ReplyMarkup.OneTimeKeyboard
	OneTimeKeyboard

	// RemoveKeyboard = ReplyMarkup.RemoveKeyboard
	RemoveKeyboard

	// IgnoreThread is used to ignore the thread when responding to a message via context.
	IgnoreThread
)

// Placeholder is used to set input field placeholder as a send option.
func Placeholder(text string) *SendOptions {
	return &SendOptions{
		ReplyMarkup: &ReplyMarkup{
			ForceReply:  true,
			Placeholder: text,
		},
	}
}

// SendOptions has most complete control over in what way the message
// must be sent, providing an API-complete set of custom properties
// and options.
//
// Despite its power, SendOptions is rather inconvenient to use all
// the way through bot logic, so you might want to consider storing
// and re-using it somewhere or be using Option flags instead.
type SendOptions struct {
	// If the message is a reply, original message.
	ReplyTo *Message

	// See ReplyMarkup struct definition.
	ReplyMarkup *ReplyMarkup

	// For text messages, disables previews for links in this message.
	DisableWebPagePreview bool

	// Sends the message silently. iOS users will not receive a notification, Android users will receive a notification with no sound.
	DisableNotification bool

	// ParseMode controls how client apps render your message.
	ParseMode ParseMode

	// Entities is a list of special entities that appear in message text, which can be specified instead of parse_mode.
	Entities Entities

	// AllowWithoutReply allows sending messages not a as reply if the replied-to message has already been deleted.
	AllowWithoutReply bool

	// Protected protects the contents of sent message from forwarding and saving.
	Protected bool

	// ThreadID supports sending messages to a thread.
	ThreadID int

	// HasSpoiler marks the message as containing a spoiler.
	HasSpoiler bool

	// ReplyParams Describes the message to reply to
	ReplyParams *ReplyParams

	// Unique identifier of the business connection
	BusinessConnectionID string

	// Unique identifier of the message effect to be added to the message; for private chats only
	EffectID string

	// Payload is a bot-defined paid media payload, returned in the PaidMediaPurchased update.
	// 0-128 characters. Used only for sendPaidMedia.
	Payload string

	// AllowPaidBroadcast allows the message to be sent to users who have not yet
	// unlocked access to paid channels. Only for messages sent to channels.
	AllowPaidBroadcast bool

	// Ephemeral sends the message as an ephemeral message, visible only to one
	// user and the bot (Bot API 10.3). It can also be passed to Send as an
	// option on its own.
	Ephemeral *EphemeralMessageParameters

	// CanStop shows the user a button to stop a streamed draft; pressing it
	// delivers OnGenerationStopped. KeepOnStop keeps the draft in the chat
	// for a short time after it is stopped. Used only by SendDraft and
	// SendRichDraft (Bot API 10.3).
	CanStop    bool
	KeepOnStop bool

	// DirectMessagesTopicID sends the message to a topic of a channel's direct
	// messages chat (Bot API 9.2).
	DirectMessagesTopicID int64

	// SuggestedPost sends the message as a post suggested to the channel of
	// a direct messages chat (Bot API 9.2). It can also be passed to Send as
	// an option on its own.
	SuggestedPost *SuggestedPostParameters

	// PreviewOptions sets how the link preview of a text message is
	// generated, when sending or editing it (Bot API 7.0). It can also be
	// passed to Send as an option on its own.
	PreviewOptions *PreviewOptions

	// RemoveCaption copies messages without their captions. Used only by
	// CopyMany.
	RemoveCaption bool

	// CaptionAbove shows an edited caption above the media. Used only when
	// editing a caption or copying with a new one; when sending media, set
	// the media's CaptionAbove.
	CaptionAbove bool

	// Caption replaces the caption of a copied message, with ParseMode or
	// Entities for its formatting. Used only by Copy.
	Caption string
}

func (og *SendOptions) copy() *SendOptions {
	cp := *og
	if cp.ReplyMarkup != nil {
		cp.ReplyMarkup = cp.ReplyMarkup.copy()
	}
	return &cp
}

func (b *Bot) extractOptions(how []interface{}) *SendOptions {
	opts := &SendOptions{
		ParseMode: b.parseMode,
	}

	for _, prop := range how {
		switch opt := prop.(type) {
		case *SendOptions:
			opts = opt.copy()
		case *ReplyMarkup:
			if opt != nil {
				opts.ReplyMarkup = opt.copy()
			}
		case *ReplyParams:
			opts.ReplyParams = opt
		case *EphemeralMessageParameters:
			opts.Ephemeral = opt
		case *SuggestedPostParameters:
			opts.SuggestedPost = opt
		case *PreviewOptions:
			opts.PreviewOptions = opt
		case *Topic:
			opts.ThreadID = opt.ThreadID
		case Option:
			switch opt {
			case NoPreview:
				opts.DisableWebPagePreview = true
			case Silent:
				opts.DisableNotification = true
			case AllowWithoutReply:
				opts.AllowWithoutReply = true
			case ForceReply:
				if opts.ReplyMarkup == nil {
					opts.ReplyMarkup = &ReplyMarkup{}
				}
				opts.ReplyMarkup.ForceReply = true
			case OneTimeKeyboard:
				if opts.ReplyMarkup == nil {
					opts.ReplyMarkup = &ReplyMarkup{}
				}
				opts.ReplyMarkup.OneTimeKeyboard = true
			case RemoveKeyboard:
				if opts.ReplyMarkup == nil {
					opts.ReplyMarkup = &ReplyMarkup{}
				}
				opts.ReplyMarkup.RemoveKeyboard = true
			case Protected:
				opts.Protected = true
			case IgnoreThread:
				// Only Context reads it, to leave its thread out.
			default:
				panic("telebot: unsupported flag-option")
			}
		case ParseMode:
			opts.ParseMode = opt
		case Entities:
			opts.Entities = opt
		default:
			panic("telebot: unsupported send-option")
		}
	}

	return opts
}

func (b *Bot) embedSendOptions(params map[string]string, opt *SendOptions) {
	if opt == nil {
		return
	}

	if opt.ReplyTo != nil && opt.ReplyTo.ID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(opt.ReplyTo.ID)
	}

	if opt.ReplyParams != nil {
		replyParams, _ := json.Marshal(opt.ReplyParams)
		params["reply_parameters"] = string(replyParams)
	}

	if opt.DisableWebPagePreview {
		params["disable_web_page_preview"] = "true"
	}

	if opt.DisableNotification {
		params["disable_notification"] = "true"
	}

	if opt.ParseMode != ModeDefault {
		params["parse_mode"] = opt.ParseMode
	}

	if len(opt.Entities) > 0 {
		delete(params, "parse_mode")
		entities, _ := json.Marshal(opt.Entities)

		if params["caption"] != "" {
			params["caption_entities"] = string(entities)
		} else {
			params["entities"] = string(entities)
		}
	}

	if opt.AllowWithoutReply {
		params["allow_sending_without_reply"] = "true"
	}

	if opt.ReplyMarkup != nil {
		processButtons(opt.ReplyMarkup.InlineKeyboard)
		replyMarkup, _ := json.Marshal(opt.ReplyMarkup)
		params["reply_markup"] = string(replyMarkup)
	}

	if opt.Protected {
		params["protect_content"] = "true"
	}

	if opt.ThreadID != 0 {
		params["message_thread_id"] = strconv.Itoa(opt.ThreadID)
	}

	if opt.HasSpoiler {
		params["has_spoiler"] = "true"
	}

	if opt.BusinessConnectionID != "" {
		params["business_connection_id"] = opt.BusinessConnectionID
	}

	if opt.EffectID != "" {
		params["message_effect_id"] = opt.EffectID
	}

	if opt.AllowPaidBroadcast {
		params["allow_paid_broadcast"] = "true"
	}

	if opt.PreviewOptions != nil {
		data, _ := json.Marshal(opt.PreviewOptions)
		params["link_preview_options"] = string(data)
	}

	if opt.RemoveCaption {
		params["remove_caption"] = "true"
	}

	if opt.DirectMessagesTopicID != 0 {
		params["direct_messages_topic_id"] = strconv.FormatInt(opt.DirectMessagesTopicID, 10)
	}

	if opt.SuggestedPost != nil {
		data, _ := json.Marshal(opt.SuggestedPost)
		params["suggested_post_parameters"] = string(data)
	}

	if opt.Ephemeral != nil {
		data, _ := json.Marshal(opt.Ephemeral)
		params["ephemeral_message_parameters"] = string(data)
	}
}

// embedDraftStop adds the stop button options of the draft methods.
func embedDraftStop(params map[string]string, opt *SendOptions) {
	if opt.CanStop {
		params["can_stop"] = "true"
	}
	if opt.KeepOnStop {
		params["keep_on_stop"] = "true"
	}
}

func processButtons(keys [][]InlineButton) {
	if len(keys) < 1 || len(keys[0]) < 1 {
		return
	}

	for i := range keys {
		for j := range keys[i] {
			key := &keys[i][j]
			if key.Unique != "" {
				// Format: "\f<callback_name>|<data>"
				data := key.Data
				if data == "" {
					key.Data = "\f" + key.Unique
				} else {
					key.Data = "\f" + key.Unique + "|" + data
				}
			}
		}
	}
}

// PreviewOptions describes the options used for link preview generation.
//
// Bot API: LinkPreviewOptions.
type PreviewOptions struct {
	// (Optional) True, if the link preview is disabled.
	Disabled bool `json:"is_disabled"`

	// (Optional) URL to use for the link preview. If empty, then the first URL
	// found in the message text will be used.
	URL string `json:"url"`

	// (Optional) True, if the media in the link preview is supposed to be shrunk;
	// ignored if the URL isn't explicitly specified or media size change.
	// isn't supported for the preview.
	SmallMedia bool `json:"prefer_small_media"`

	// (Optional) True, if the media in the link preview is supposed to be enlarged;
	// ignored if the URL isn't explicitly specified or media size change.
	// isn't supported for the preview.
	LargeMedia bool `json:"prefer_large_media"`

	// (Optional) True, if the link preview must be shown above the message text;
	// otherwise, the link preview will be shown below the message text.
	AboveText bool `json:"show_above_text"`
}

// embedMessages sets message_ids to the IDs of msgs and chatKey to the chat
// they are in.
func embedMessages(params map[string]string, chatKey string, msgs []Editable) error {
	if len(msgs) == 0 {
		return errors.New("telebot: no messages given")
	}

	ids := make([]int, len(msgs))
	for i, msg := range msgs {
		msgID, _ := msg.MessageSig()
		id, err := strconv.Atoi(msgID)
		if err != nil {
			return fmt.Errorf("telebot: message ID %q is not a number", msgID)
		}
		ids[i] = id
	}
	data, _ := json.Marshal(ids)
	_, chatID := msgs[0].MessageSig()

	params["message_ids"] = string(data)
	params[chatKey] = strconv.FormatInt(chatID, 10)
	return nil
}
