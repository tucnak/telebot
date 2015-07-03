package telebot

import (
	"time"
)

// Bot represents a separate Telegram bot instance.
type Bot struct {
	Token string

	// Bot as `User` on API level.
	Identity User
}

// Listen periodically looks for updates and delivers new messages
// to subscription channel.
func (b *Bot) Listen(subscription chan<- Message, interval time.Duration) {
	updates := make(chan Update)
	pulse := time.NewTicker(interval)
	latest_update := 0

	go func() {
		for range pulse.C {
			go api_getUpdates(b.Token,
				latest_update+1,
				updates)
		}
	}()

	go func() {
		for update := range updates {
			if update.Id > latest_update {
				latest_update = update.Id
			}

			subscription <- update.Payload
		}
	}()
}

// SendMessage sends a text message to recipient.
func (b *Bot) SendMessage(recipient User, message string) error {
	return api_sendMessage(b.Token, recipient, message)
}

// ForwardMessage forwards a message to recipient.
func (b *Bot) ForwardMessage(recipient User, message Message) error {
	return api_forwardMessage(b.Token, recipient, message)
}

// SendPhoto sends a photo object to recipient.
//
// On success, photo object would be aliased to its copy on
// the Telegram servers, so sending the same photo object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b *Bot) SendPhoto(recipient User, photo *Photo) error {
	return api_sendPhoto(b.Token, recipient, photo)
}

// SendAudio sends an audio object to recipient.
//
// On success, audio object would be aliased to its copy on
// the Telegram servers, so sending the same audio object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b *Bot) SendAudio(recipient User, audio *Audio) error {
	return api_sendAudio(b.Token, recipient, audio)
}

// SendDocument sends a general document object to recipient.
//
// On success, document object would be aliased to its copy on
// the Telegram servers, so sending the same document object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b *Bot) SendDocument(recipient User, doc *Document) error {
	return api_sendDocument(b.Token, recipient, doc)
}

// SendSticker sends a general document object to recipient.
//
// On success, sticker object would be aliased to its copy on
// the Telegram servers, so sending the same sticker object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b *Bot) SendSticker(recipient User, sticker *Sticker) error {
	return api_sendSticker(b.Token, recipient, sticker)
}

// SendVideo sends a general document object to recipient.
//
// On success, video object would be aliased to its copy on
// the Telegram servers, so sending the same video object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b *Bot) SendVideo(recipient User, video *Video) error {
	return api_sendVideo(b.Token, recipient, video)
}

// SendLocation sends a general document object to recipient.
//
// On success, video object would be aliased to its copy on
// the Telegram servers, so sending the same video object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b *Bot) SendLocation(recipient User, geo *Location) error {
	return api_sendLocation(b.Token, recipient, geo)
}
