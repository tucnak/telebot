package telebot

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/armon/go-radix"
	"github.com/pkg/errors"
)

// Bot represents a separate Telegram bot instance.
type Bot struct {
	Token     string
	Identity  User
	Messages  chan Message
	Queries   chan Query
	Callbacks chan Callback

	// Telebot debugging channel. If present, Telebot
	// will use it to report all occuring errors.
	Errors chan error

	tree *radix.Tree
}

// NewBot does try to build a Bot with token `token`, which
// is a secret API key assigned to particular bot.
func NewBot(token string) (*Bot, error) {
	bot := &Bot{
		Token: token,
		tree:  radix.New(),
	}

	user, err := bot.getMe()
	if err != nil {
		return nil, err
	}

	bot.Identity = user
	return bot, nil
}

// Listen starts a new polling goroutine, one that periodically looks for
// updates and delivers new messages to the subscription channel.
func (b *Bot) Listen(subscription chan Message, timeout time.Duration) {
	go b.poll(subscription, nil, nil, timeout, nil)
}

// ListenWithShutdown starts a new polling goroutine, one that periodically looks for
// updates and delivers new messages to the subscription channel.
// If shutdown channel received any value, close messages channels
func (b *Bot) ListenWithShutdown(subscription chan Message, timeout time.Duration, shutdown <-chan bool) {
	go b.poll(subscription, nil, nil, timeout, shutdown)
}

// Start periodically polls messages, updates and callbacks into their
// corresponding channels of the bot object.
//
// NOTE: It's a blocking method!
func (b *Bot) Start(timeout time.Duration) {
	b.poll(b.Messages, b.Queries, b.Callbacks, timeout, nil)
}

func (b *Bot) debug(err error) {
	if b.Errors != nil {
		b.Errors <- errors.WithStack(err)
	}
}

func (b *Bot) poll(
	messages chan Message,
	queries chan Query,
	callbacks chan Callback,
	timeout time.Duration,
	shutdown <-chan bool,
) {
	var latestUpdate int64
	if shutdown != nil {
		for {
			select {
			case <-shutdown:
				if messages != nil {
					close(messages)
				}
				if queries != nil {
					close(queries)
				}
				if callbacks != nil {
					close(callbacks)
				}
				return
			default:
				b.handleUpdates(&latestUpdate, messages, queries, callbacks, timeout)
			}
		}
	} else {
		for {
			b.handleUpdates(&latestUpdate, messages, queries, callbacks, timeout)
		}
	}
}

func (b *Bot) handleUpdates(
	latestUpdate *int64,
	messages chan Message,
	queries chan Query,
	callbacks chan Callback,
	timeout time.Duration,
) {
	*latestUpdate = *latestUpdate + 1
	updates, err := b.getUpdates(*latestUpdate, timeout)

	if err != nil {
		b.debug(errors.Wrap(err, "getUpdates() failed"))
		return
	}

	for _, update := range updates {
		if update.Payload != nil /* if message */ {
			if messages == nil {
				continue
			}

			messages <- *update.Payload
		} else if update.Query != nil /* if query */ {
			if queries == nil {
				continue
			}

			queries <- *update.Query
		} else if update.Callback != nil {
			if callbacks == nil {
				continue
			}

			callbacks <- *update.Callback
		}

		*latestUpdate = update.ID
	}
}

// SendMessage sends a text message to recipient.
func (b *Bot) SendMessage(recipient Recipient, message string, options *SendOptions) error {
	params := map[string]string{
		"chat_id": recipient.Destination(),
		"text":    message,
	}

	if options != nil {
		embedSendOptions(params, options)
	}

	responseJSON, err := b.sendCommand("sendMessage", params)
	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	return nil
}

// ForwardMessage forwards a message to recipient.
func (b *Bot) ForwardMessage(recipient Recipient, message Message) error {
	params := map[string]string{
		"chat_id":      recipient.Destination(),
		"from_chat_id": strconv.Itoa(message.Origin().ID),
		"message_id":   strconv.Itoa(message.ID),
	}

	responseJSON, err := b.sendCommand("forwardMessage", params)
	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	return nil
}

// SendPhoto sends a photo object to recipient.
//
// On success, photo object would be aliased to its copy on
// the Telegram servers, so sending the same photo object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b *Bot) SendPhoto(recipient Recipient, photo *Photo, options *SendOptions) error {
	params := map[string]string{
		"chat_id": recipient.Destination(),
		"caption": photo.Caption,
	}

	if options != nil {
		embedSendOptions(params, options)
	}

	var responseJSON []byte
	var err error

	if photo.Exists() {
		params["photo"] = photo.FileID
		responseJSON, err = b.sendCommand("sendPhoto", params)
	} else {
		responseJSON, err = b.sendFile("sendPhoto", "photo", photo.filename, params)
	}

	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	thumbnails := &responseReceived.Result.Photo
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
func (b *Bot) SendAudio(recipient Recipient, audio *Audio, options *SendOptions) error {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}

	if options != nil {
		embedSendOptions(params, options)
	}

	var responseJSON []byte
	var err error

	if audio.Exists() {
		params["audio"] = audio.FileID
		responseJSON, err = b.sendCommand("sendAudio", params)
	} else {
		responseJSON, err = b.sendFile("sendAudio", "audio", audio.filename, params)
	}

	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	filename := audio.filename
	*audio = responseReceived.Result.Audio
	audio.filename = filename

	return nil
}

// SendDocument sends a general document object to recipient.
//
// On success, document object would be aliased to its copy on
// the Telegram servers, so sending the same document object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b *Bot) SendDocument(recipient Recipient, doc *Document, options *SendOptions) error {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}

	if options != nil {
		embedSendOptions(params, options)
	}

	var responseJSON []byte
	var err error

	if doc.Exists() {
		params["document"] = doc.FileID
		responseJSON, err = b.sendCommand("sendDocument", params)
	} else {
		responseJSON, err = b.sendFile("sendDocument", "document", doc.filename, params)
	}

	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	filename := doc.filename
	*doc = responseReceived.Result.Document
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
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}

	if options != nil {
		embedSendOptions(params, options)
	}

	var responseJSON []byte
	var err error

	if sticker.Exists() {
		params["sticker"] = sticker.FileID
		responseJSON, err = b.sendCommand("sendSticker", params)
	} else {
		responseJSON, err = b.sendFile("sendSticker", "sticker", sticker.filename, params)
	}

	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	filename := sticker.filename
	*sticker = responseReceived.Result.Sticker
	sticker.filename = filename

	return nil
}

// SendVideo sends a general document object to recipient.
//
// On success, video object would be aliased to its copy on
// the Telegram servers, so sending the same video object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b *Bot) SendVideo(recipient Recipient, video *Video, options *SendOptions) error {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}

	if options != nil {
		embedSendOptions(params, options)
	}

	var responseJSON []byte
	var err error

	if video.Exists() {
		params["video"] = video.FileID
		responseJSON, err = b.sendCommand("sendVideo", params)
	} else {
		responseJSON, err = b.sendFile("sendVideo", "video", video.filename, params)
	}

	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	filename := video.filename
	*video = responseReceived.Result.Video
	video.filename = filename

	return nil
}

// SendLocation sends a general document object to recipient.
//
// On success, video object would be aliased to its copy on
// the Telegram servers, so sending the same video object
// again, won't issue a new upload, but would make a use
// of existing file on Telegram servers.
func (b *Bot) SendLocation(recipient Recipient, geo *Location, options *SendOptions) error {
	params := map[string]string{
		"chat_id":   recipient.Destination(),
		"latitude":  fmt.Sprintf("%f", geo.Latitude),
		"longitude": fmt.Sprintf("%f", geo.Longitude),
	}

	if options != nil {
		embedSendOptions(params, options)
	}

	responseJSON, err := b.sendCommand("sendLocation", params)
	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	return nil
}

// SendVenue sends a venue object to recipient.
func (b *Bot) SendVenue(recipient Recipient, venue *Venue, options *SendOptions) error {
	params := map[string]string{
		"chat_id":   recipient.Destination(),
		"latitude":  fmt.Sprintf("%f", venue.Location.Latitude),
		"longitude": fmt.Sprintf("%f", venue.Location.Longitude),
		"title":     venue.Title,
		"address":   venue.Address}
	if venue.FoursquareID != "" {
		params["foursquare_id"] = venue.FoursquareID
	}

	if options != nil {
		embedSendOptions(params, options)
	}

	responseJSON, err := b.sendCommand("sendVenue", params)
	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
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
func (b *Bot) SendChatAction(recipient Recipient, action ChatAction) error {
	params := map[string]string{
		"chat_id": recipient.Destination(),
		"action":  string(action),
	}

	responseJSON, err := b.sendCommand("sendChatAction", params)
	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	return nil
}

// Respond publishes a set of responses for an inline query.
// This function is deprecated in favor of AnswerInlineQuery.
func (b *Bot) Respond(query Query, results []Result) error {
	params := map[string]string{
		"inline_query_id": query.ID,
	}

	if res, err := json.Marshal(results); err == nil {
		params["results"] = string(res)
	} else {
		b.debug(errors.Wrapf(err, "failed to respond to \"%s\"", query.Text))
		return err
	}

	responseJSON, err := b.sendCommand("answerInlineQuery", params)
	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	return nil
}

// AnswerInlineQuery sends a response for a given inline query. A query can
// only be responded to once, subsequent attempts to respond to the same query
// will result in an error.
func (b *Bot) AnswerInlineQuery(query *Query, response *QueryResponse) error {
	response.QueryID = query.ID

	responseJSON, err := b.sendCommand("answerInlineQuery", response)
	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	return nil
}

// AnswerCallbackQuery sends a response for a given callback query. A callback can
// only be responded to once, subsequent attempts to respond to the same callback
// will result in an error.
func (b *Bot) AnswerCallbackQuery(callback *Callback, response *CallbackResponse) error {
	response.CallbackID = callback.ID

	responseJSON, err := b.sendCommand("answerCallbackQuery", response)
	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	return nil
}

// GetFile returns full file object including File.FilePath, which allow you to load file from Telegram
//
// Usually File objects does not contain any FilePath so you need to perform additional request
func (b *Bot) GetFile(fileID string) (File, error) {
	params := map[string]string{
		"file_id": fileID,
	}
	responseJSON, err := b.sendCommand("getFile", params)
	if err != nil {
		return File{}, err
	}

	var responseReceived struct {
		Ok          bool
		Description string
		Result      File
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return File{}, errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return File{}, errors.Errorf("api error: %s", responseReceived.Description)

	}

	return responseReceived.Result, nil
}

// LeaveChat makes bot leave a group, supergroup or channel.
func (b *Bot) LeaveChat(recipient Recipient) error {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}
	responseJSON, err := b.sendCommand("leaveChat", params)
	if err != nil {
		return err
	}

	var responseReceived struct {
		Ok          bool
		Description string
		Result      bool
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return errors.Errorf("api error: %s", responseReceived.Description)
	}

	return nil
}

// GetChat get up to date information about the chat.
//
// Including current name of the user for one-on-one conversations,
// current username of a user, group or channel, etc.
//
// Returns a Chat object on success.
func (b *Bot) GetChat(recipient Recipient) (Chat, error) {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}
	responseJSON, err := b.sendCommand("getChat", params)
	if err != nil {
		return Chat{}, err
	}

	var responseReceived struct {
		Ok          bool
		Description string
		Result      Chat
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return Chat{}, errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return Chat{}, errors.Errorf("api error: %s", responseReceived.Description)
	}

	return responseReceived.Result, nil
}

// GetChatAdministrators return list of administrators in a chat.
//
// On success, returns an Array of ChatMember objects that
// contains information about all chat administrators except other bots.
//
// If the chat is a group or a supergroup and
// no administrators were appointed, only the creator will be returned.
func (b *Bot) GetChatAdministrators(recipient Recipient) ([]ChatMember, error) {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}
	responseJSON, err := b.sendCommand("getChatAdministrators", params)
	if err != nil {
		return []ChatMember{}, err
	}

	var responseReceived struct {
		Ok          bool
		Result      []ChatMember
		Description string `json:"description"`
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return []ChatMember{}, errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return []ChatMember{}, errors.Errorf("api error: %s", responseReceived.Description)
	}

	return responseReceived.Result, nil
}

// GetChatMembersCount return the number of members in a chat.
//
// Returns Int on success.
func (b *Bot) GetChatMembersCount(recipient Recipient) (int, error) {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}
	responseJSON, err := b.sendCommand("getChatMembersCount", params)
	if err != nil {
		return 0, err
	}

	var responseReceived struct {
		Ok          bool
		Result      int
		Description string `json:"description"`
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return 0, errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return 0, errors.Errorf("api error: %s", responseReceived.Description)
	}

	return responseReceived.Result, nil
}

// GetUserProfilePhotos return list of profile pictures for a user.
//
// Returns a UserProfilePhotos object.
func (b *Bot) GetUserProfilePhotos(recipient Recipient) (UserProfilePhotos, error) {
	params := map[string]string{
		"user_id": recipient.Destination(),
	}
	responseJSON, err := b.sendCommand("getUserProfilePhotos", params)
	if err != nil {
		return UserProfilePhotos{}, err
	}

	var responseReceived struct {
		Ok          bool
		Result      UserProfilePhotos
		Description string `json:"description"`
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return UserProfilePhotos{}, errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return UserProfilePhotos{}, errors.Errorf("api error: %s", responseReceived.Description)
	}

	return responseReceived.Result, nil
}

// GetChatMember return information about a member of a chat.
//
// Returns a ChatMember object on success.
func (b *Bot) GetChatMember(recipient Recipient, user User) (ChatMember, error) {
	params := map[string]string{
		"chat_id": recipient.Destination(),
		"user_id": user.Destination(),
	}
	responseJSON, err := b.sendCommand("getChatMember", params)
	if err != nil {
		return ChatMember{}, err
	}

	var responseReceived struct {
		Ok          bool
		Result      ChatMember
		Description string `json:"description"`
	}

	err = json.Unmarshal(responseJSON, &responseReceived)
	if err != nil {
		return ChatMember{}, errors.Wrap(err, "bad response json")
	}

	if !responseReceived.Ok {
		return ChatMember{}, errors.Errorf("api error: %s", responseReceived.Description)
	}

	return responseReceived.Result, nil
}

// GetFileDirectURL returns direct url for files using FileId which you can get from File object
func (b *Bot) GetFileDirectURL(fileID string) (string, error) {
	f, err := b.GetFile(fileID)
	if err != nil {
		return "", err
	}
	return "https://api.telegram.org/file/bot" + b.Token + "/" + f.FilePath, nil
}
