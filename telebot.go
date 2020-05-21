// Package telebot is a framework for Telegram bots.
//
// Example:
//
//		package main
//
//		import (
//			"time"
//			tb "gopkg.in/tucnak/telebot.v2"
//		)
//
//		func main() {
//			b, err := tb.NewBot(tb.Settings{
//				Token: "TOKEN_HERE",
//				Poller: &tb.LongPoller{Timeout: 10 * time.Second},
//			})
//
//			if err != nil {
//				return
//			}
//
//			b.Handle(tb.OnText, func(m *tb.Message) {
//				b.Send(m.Sender, "hello world")
//			})
//
//			b.Start()
//		}
//
package telebot

import "github.com/pkg/errors"

var (
	ErrBadRecipient    = errors.New("telebot: recipient is nil")
	ErrUnsupportedWhat = errors.New("telebot: unsupported what argument")
	ErrCouldNotUpdate  = errors.New("telebot: could not fetch new updates")
	ErrNoGameMessage   = errors.New("telebot: no game message")
)

const DefaultApiURL = "https://api.telegram.org"

// These are one of the possible events Handle() can deal with.
//
// For convenience, all Telebot-provided endpoints start with
// an "alert" character \a.
const (
	// Basic message handlers.
	//
	// Handler: func(*Message)
	OnText              = "\atext"
	OnPhoto             = "\aphoto"
	OnAudio             = "\aaudio"
	OnAnimation         = "\aanimation"
	OnDocument          = "\adocument"
	OnSticker           = "\asticker"
	OnVideo             = "\avideo"
	OnVoice             = "\avoice"
	OnVideoNote         = "\avideo_note"
	OnContact           = "\acontact"
	OnLocation          = "\alocation"
	OnVenue             = "\avenue"
	OnEdited            = "\aedited"
	OnPinned            = "\apinned"
	OnChannelPost       = "\achan_post"
	OnEditedChannelPost = "\achan_edited_post"
	OnDice              = "\adice"
	OnInvoice           = "\ainvoice"
	OnPayment           = "\apayment"

	// Will fire when bot is added to a group.
	OnAddedToGroup = "\aadded_to_group"

	// Group events:
	OnUserJoined        = "\auser_joined"
	OnUserLeft          = "\auser_left"
	OnNewGroupTitle     = "\anew_chat_title"
	OnNewGroupPhoto     = "\anew_chat_photo"
	OnGroupPhotoDeleted = "\achat_photo_del"

	// Migration happens when group switches to
	// a super group. You might want to update
	// your internal references to this chat
	// upon switching as its ID will change.
	//
	// Handler: func(from, to int64)
	OnMigration = "\amigration"

	// Will fire on callback requests.
	//
	// Handler: func(*Callback)
	OnCallback = "\acallback"

	// Will fire on incoming inline queries.
	//
	// Handler: func(*Query)
	OnQuery = "\aquery"

	// Will fire on chosen inline results.
	//
	// Handler: func(*ChosenInlineResult)
	OnChosenInlineResult = "\achosen_inline_result"

	// Will fire on ShippingQuery.
	//
	// Handler: func(*ShippingQuery)
	OnShipping = "\ashipping_query"

	// Will fire on PreCheckoutQuery.
	//
	// Handler: func(*PreCheckoutQuery)
	OnCheckout = "\apre_checkout_query"

	// Will fire on Poll.
	//
	// Handler: func(*Poll)
	OnPoll = "\apoll"

	// Will fire on PollAnswer.
	//
	// Handler: func(*PollAnswer)
	OnPollAnswer = "\apoll_answer"
)

// ChatAction is a client-side status indicating bot activity.
type ChatAction string

const (
	Typing            ChatAction = "typing"
	UploadingPhoto    ChatAction = "upload_photo"
	UploadingVideo    ChatAction = "upload_video"
	UploadingAudio    ChatAction = "upload_audio"
	UploadingDocument ChatAction = "upload_document"
	UploadingVNote    ChatAction = "upload_video_note"
	RecordingVideo    ChatAction = "record_video"
	RecordingAudio    ChatAction = "record_audio"
	RecordingVNote    ChatAction = "record_video_note"
	FindingLocation   ChatAction = "find_location"
)

// ParseMode determines the way client applications treat the text of the message
type ParseMode = string

const (
	ModeDefault    ParseMode = ""
	ModeMarkdown   ParseMode = "Markdown"
	ModeMarkdownV2 ParseMode = "MarkdownV2"
	ModeHTML       ParseMode = "HTML"
)

// EntityType is a MessageEntity type.
type EntityType string

const (
	EntityMention       EntityType = "mention"
	EntityTMention      EntityType = "text_mention"
	EntityHashtag       EntityType = "hashtag"
	EntityCashtag       EntityType = "cashtag"
	EntityCommand       EntityType = "bot_command"
	EntityURL           EntityType = "url"
	EntityEmail         EntityType = "email"
	EntityPhone         EntityType = "phone_number"
	EntityBold          EntityType = "bold"
	EntityItalic        EntityType = "italic"
	EntityUnderline     EntityType = "underline"
	EntityStrikethrough EntityType = "strikethrough"
	EntityCode          EntityType = "code"
	EntityCodeBlock     EntityType = "pre"
	EntityTextLink      EntityType = "text_link"
)

// ChatType represents one of the possible chat types.
type ChatType string

const (
	ChatPrivate        ChatType = "private"
	ChatGroup          ChatType = "group"
	ChatSuperGroup     ChatType = "supergroup"
	ChatChannel        ChatType = "channel"
	ChatChannelPrivate ChatType = "privatechannel"
)

// MemberStatus is one's chat status.
type MemberStatus string

const (
	Creator       MemberStatus = "creator"
	Administrator MemberStatus = "administrator"
	Member        MemberStatus = "member"
	Restricted    MemberStatus = "restricted"
	Left          MemberStatus = "left"
	Kicked        MemberStatus = "kicked"
)

// MaskFeature defines sticker mask position.
type MaskFeature string

const (
	FeatureForehead MaskFeature = "forehead"
	FeatureEyes     MaskFeature = "eyes"
	FeatureMouth    MaskFeature = "mouth"
	FeatureChin     MaskFeature = "chin"
)

// PollType defines poll types.
type PollType string

const (
	// Despite "any" type isn't described in documentation,
	// it needed for proper KeyboardButtonPollType marshaling.
	PollAny PollType = "any"

	PollQuiz    PollType = "quiz"
	PollRegular PollType = "regular"
)

type DiceType string

var (
	Cube = &Dice{Type: "🎲"}
	Dart = &Dice{Type: "🎯"}
)
