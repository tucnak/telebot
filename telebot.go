// Package telebot is a framework for Telegram bots.
//
// Example:
//
//		import (
//			"time"
//			tb "gopkg.in/tucnak/telebot.v2"
//		)
//
//		func main() {
//			b, err := tb.NewBot(tb.Settings{
//				Token: "TOKEN_HERE",
//				Poller: &tb.LongPoller{10 * time.Second},
//			})
//
//			if err != nil {
//				return
//			}
//
//			b.Handle(tb.OnMessage, func(m *tb.Message) {
//				b.Send(m.From, "hello world")
//			}
//
//			b.Start()
//		}
//
package telebot

// These are one of the possible events Handle() can deal with.
//
// For convenience, all Telebot-provided endpoints start with
// an "alert" character \a.
const (
	// Basic generic message handlers.
	//
	// Handler: func(*Message)
	OnMessage           = "\amessage"
	OnEditedMessage     = "\aedited_msg"
	OnChannelPost       = "\achan_post"
	OnEditedChannelPost = "\achan_post"

	// Will fire on callback requests.
	//
	// Handler: func(*Callback)
	OnCallback = "\acallback"

	// Will fire on incoming inline queries.
	//
	// Handler: func(*Query)
	OnQuery = "\aquery"

	// Will fire when bot gets added to some chat,
	//
	// Handler: func(*Message)
	OnAddedToGroup = "\aadded_to_group"
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
	FindingLocation   ChatAction = "find_location"
)

// ParseMode determines the way client applications treat the text of the message
type ParseMode string

const (
	ModeDefault  ParseMode = ""
	ModeMarkdown ParseMode = "Markdown"
	ModeHTML     ParseMode = "HTML"
)

// EntityType is a MessageEntity type.
type EntityType string

const (
	EntityMention   EntityType = "mention"
	EntityTMention  EntityType = "text_mention"
	EntityHashtag   EntityType = "hashtag"
	EntityCommand   EntityType = "bot_command"
	EntityURL       EntityType = "url"
	EntityEmail     EntityType = "email"
	EntityBold      EntityType = "bold"
	EntityItalic    EntityType = "italic"
	EntityCode      EntityType = "code"
	EntityCodeBlock EntityType = "pre"
	EntityTextLink  EntityType = "text_link"
)

// ChatType represents one of the possible chat types.
type ChatType string

const (
	ChatPrivate    ChatType = "private"
	ChatGroup      ChatType = "group"
	ChatSuperGroup ChatType = "supergroup"
	ChatChannel    ChatType = "channel"
)

// MemberStatus is one's chat status
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
