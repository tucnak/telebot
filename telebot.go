// Package telebot provides a handy wrapper for interactions
// with Telegram bots.
//
// Here is an example of helloworld bot implementation:
//
//	import (
//		"time"
//		"github.com/tucnak/telebot"
//	)
//
//	func main() {
//		bot, err := telebot.NewBot("SECRET_TOKEN")
//		if err != nil {
//			return
//		}
//
//		messages := make(chan telebot.Message)
//		bot.Listen(messages, 1*time.Second)
//
//		for message := range messages {
//			if message.Text == "/hi" {
//				bot.SendMessage(message.Chat,
//					"Hello, "+message.Sender.FirstName+"!", nil)
//			}
//		}
//	}
//
package telebot

// ChatAction is a client-side status indicating bot activity.
type ChatAction string

const (
	Typing            ChatAction = "typing"
	UploadingPhoto    ChatAction = "upload_photo"
	UploadingVideo    ChatAction = "upload_video"
	UploadingAudio    ChatAction = "upload_audio"
	UploadingDocument ChatAction = "upload_document"
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
