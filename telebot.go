// Package telebot provides a handy wrapper for interactions
// with Telegram bots.
//
// Here is an example of helloworld bot implementation:
//
//	import (
//		"log"
//		"time"
//		"os'
//		"github.com/tucnak/telebot"
//	)
//
//	func main() {
//		bot, err := telebot.NewBot(os.Getenv("BOT_TOKEN"))
//		if err != nil {
//			log.Fatalln(err)
//		}
//
//		messages := make(chan telebot.Message, 100)
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

// A bunch of available chat actions.
const (
	Typing            = "typing"
	UploadingPhoto    = "upload_photo"
	UploadingVideo    = "upload_video"
	UploadingAudio    = "upload_audio"
	UploadingDocument = "upload_document"
	RecordingVideo    = "record_video"
	RecordingAudio    = "record_audio"
	FindingLocation   = "find_location"
)
