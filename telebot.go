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
