package telebot

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Bot represents a separate Telegram bot instance.
type Bot struct {
	Token string

	// Bot as `User` on API level.
	Identity User
}

// NewBot does try to build a Bot with token `token`, which
// is a secret API key assigned to particular bot.
func NewBot(token string) (*Bot, error) {
	user, err := getMe(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Token:    token,
		Identity: user,
	}, nil
}

// Listen periodically looks for updates and delivers new messages
// to subscription channel.
func (b Bot) Listen(subscription chan<- Message, timeout time.Duration) {

	go func() {
		latestUpdate := 0
		for {
			if updates, err := getUpdates(b.Token, latestUpdate+1, int(timeout/time.Second)); err == nil {
				for _, update := range updates {
					latestUpdate = update.ID
					subscription <- update.Payload
				}
			}
		}
	}()
}

// SendMessage sends a text message to recipient.
func (b Bot) SendMessage(recipient Recipient, message string, options *SendOptions) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Destination()))
	params.Set("text", message)

	if options != nil {
		embedSendOptions(&params, options)
	}

	responseJSON, err := sendCommand("sendMessage", b.Token, params)
	if err != nil {
		return err
	}

	var responseRecieved struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return err
	}

	if !responseRecieved.Ok {
		return fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	return nil
}

// ForwardMessage forwards a message to recipient.
func (b Bot) ForwardMessage(recipient Recipient, message Message) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Destination()))
	params.Set("from_chat_id", strconv.Itoa(message.Origin().ID))
	params.Set("message_id", strconv.Itoa(message.ID))

	responseJSON, err := sendCommand("forwardMessage", b.Token, params)
	if err != nil {
		return err
	}

	var responseRecieved struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return err
	}

	if !responseRecieved.Ok {
		return fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	return nil
}

// SendPhoto sends a photo object to recipient.
//
// On success, photo object would be aliased to its copy on
// the Telegram servers, so sending the same photo object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b Bot) SendPhoto(recipient Recipient, photo *Photo, options *SendOptions) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Destination()))
	params.Set("caption", photo.Caption)

	if options != nil {
		embedSendOptions(&params, options)
	}

	var responseJSON []byte
	var err error

	if photo.Exists() {
		params.Set("photo", photo.FileID)
		responseJSON, err = sendCommand("sendPhoto", b.Token, params)
	} else {
		responseJSON, err = sendFile("sendPhoto", b.Token, "photo",
			photo.filename, params)
	}

	if err != nil {
		return err
	}

	var responseRecieved struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return err
	}

	if !responseRecieved.Ok {
		return fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	thumbnails := &responseRecieved.Result.Photo
	filename := photo.filename
	photo.File = (*thumbnails)[len(*thumbnails)-1].File
	photo.filename = filename

	return nil
}

// SendAudio sends an audio object to recipient.
//
// On success, audio object would be aliased to its copy on
// the Telegram servers, so sending the same audio object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b Bot) SendAudio(recipient Recipient, audio *Audio, options *SendOptions) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Destination()))

	if options != nil {
		embedSendOptions(&params, options)
	}

	var responseJSON []byte
	var err error

	if audio.Exists() {
		params.Set("audio", audio.FileID)
		responseJSON, err = sendCommand("sendAudio", b.Token, params)
	} else {
		responseJSON, err = sendFile("sendAudio", b.Token, "audio",
			audio.filename, params)
	}

	if err != nil {
		return err
	}

	var responseRecieved struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return err
	}

	if !responseRecieved.Ok {
		return fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	filename := audio.filename
	*audio = responseRecieved.Result.Audio
	audio.filename = filename

	return nil
}

// SendDocument sends a general document object to recipient.
//
// On success, document object would be aliased to its copy on
// the Telegram servers, so sending the same document object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b Bot) SendDocument(recipient Recipient, doc *Document, options *SendOptions) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Destination()))

	if options != nil {
		embedSendOptions(&params, options)
	}

	var responseJSON []byte
	var err error

	if doc.Exists() {
		params.Set("document", doc.FileID)
		responseJSON, err = sendCommand("sendDocument", b.Token, params)
	} else {
		responseJSON, err = sendFile("sendDocument", b.Token, "document",
			doc.filename, params)
	}

	if err != nil {
		return err
	}

	var responseRecieved struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return err
	}

	if !responseRecieved.Ok {
		return fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	filename := doc.filename
	*doc = responseRecieved.Result.Document
	doc.filename = filename

	return nil
}

// SendSticker sends a general document object to recipient.
//
// On success, sticker object would be aliased to its copy on
// the Telegram servers, so sending the same sticker object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b *Bot) SendSticker(recipient Recipient, sticker *Sticker, options *SendOptions) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Destination()))

	if options != nil {
		embedSendOptions(&params, options)
	}

	var responseJSON []byte
	var err error

	if sticker.Exists() {
		params.Set("sticker", sticker.FileID)
		responseJSON, err = sendCommand("sendSticker", b.Token, params)
	} else {
		responseJSON, err = sendFile("sendSticker", b.Token, "sticker",
			sticker.filename, params)
	}

	if err != nil {
		return err
	}

	var responseRecieved struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return err
	}

	if !responseRecieved.Ok {
		return fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	filename := sticker.filename
	*sticker = responseRecieved.Result.Sticker
	sticker.filename = filename

	return nil
}

// SendVideo sends a general document object to recipient.
//
// On success, video object would be aliased to its copy on
// the Telegram servers, so sending the same video object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b Bot) SendVideo(recipient Recipient, video *Video, options *SendOptions) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Destination()))

	if options != nil {
		embedSendOptions(&params, options)
	}

	var responseJSON []byte
	var err error

	if video.Exists() {
		params.Set("video", video.FileID)
		responseJSON, err = sendCommand("sendVideo", b.Token, params)
	} else {
		responseJSON, err = sendFile("sendVideo", b.Token, "video",
			video.filename, params)
	}

	if err != nil {
		return err
	}

	var responseRecieved struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return err
	}

	if !responseRecieved.Ok {
		return fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	filename := video.filename
	*video = responseRecieved.Result.Video
	video.filename = filename

	return nil
}

// SendLocation sends a general document object to recipient.
//
// On success, video object would be aliased to its copy on
// the Telegram servers, so sending the same video object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b Bot) SendLocation(recipient Recipient, geo *Location, options *SendOptions) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Destination()))
	params.Set("latitude", fmt.Sprintf("%f", geo.Latitude))
	params.Set("longitude", fmt.Sprintf("%f", geo.Longitude))

	if options != nil {
		embedSendOptions(&params, options)
	}

	responseJSON, err := sendCommand("sendLocation", b.Token, params)

	if err != nil {
		return err
	}

	var responseRecieved struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return err
	}

	if !responseRecieved.Ok {
		return fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	return nil
}

// SendChatAction updates a chat action for recipient.
//
// Chat action is a status message that recipient would see where
// you typically see "Harry is typing" status message. The only
// difference is that bots' chat actions live only for 5 seconds
// and die just once the client recieves a message from the bot.
//
// Currently, Telegram supports only a narrow range of possible
// actions, these are aligned as constants of this package.
func (b Bot) SendChatAction(recipient Recipient, action string) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Destination()))
	params.Set("action", action)

	responseJSON, err := sendCommand("sendChatAction", b.Token, params)

	if err != nil {
		return err
	}

	var responseRecieved struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return err
	}

	if !responseRecieved.Ok {
		return fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	return nil
}
