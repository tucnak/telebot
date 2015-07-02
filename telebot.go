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
//		bot, err := telebot.Create("SECRET_TOKEN")
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
//					"Hello, "+message.Sender.FirstName+"!")
//			}
//		}
//	}
//
package telebot

import (
	"fmt"
	"os"
)

// Create does try to build a Bot with token `token`, which
// is a secret API key assigned to particular bot.
func Create(token string) (Bot, error) {
	user, err := api_getMe(token)
	if err != nil {
		return Bot{}, err
	}

	return Bot{
		Token:    token,
		Identity: user,
	}, nil
}

// NewFile attempts to create a File object, leading to a real
// file on the file system, that could be uploaded later.
//
// Notice that NewFile doesn't upload file, but only creates
// a descriptor for it.
func NewFile(path string) (File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return File{}, FileError{
			fmt.Sprintf("'%s' does not exist!", path),
		}
	}

	return File{filename: path}, nil
}
