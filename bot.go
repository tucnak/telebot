package telebot

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

// Bot represents a separate Telegram bot instance.
type Bot struct {
	Token     string
	Identity  User
	Messages  chan Message
	Queries   chan Query
	Callbacks chan Callback

	handlers map[*regexp.Regexp]Handler
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
		handlers: map[*regexp.Regexp]Handler{},
	}, nil
}

// Listen periodically looks for updates and delivers new messages
// to the subscription channel.
func (b *Bot) Listen(subscription chan Message, timeout time.Duration) {
	go b.poll(subscription, nil, nil, timeout)
}

// Start periodically polls messages and/or updates to corresponding
// channels from the bot object.
func (b *Bot) Start(timeout time.Duration) {
	b.poll(b.Messages, b.Queries, b.Callbacks, timeout)
}

// Handle registers a route handler for messages, matching the
// regular expression provided.
//
// It panics if command regex appears malformed.
func (b *Bot) Handle(command string, handler Handler) {
	reg := regexp.MustCompile(command)
	b.handlers[reg] = handler
}

// Serve listens for messages and routes them to appropiate handlers.
//
// Buffer is a size of buffered message channel.
//
// Freq is a poll frequency, e.g. 2*time.Second (once / 2 sec).
func (b *Bot) Serve(buffer int, freq time.Duration) {
	messages := make(chan Message, buffer)
	b.Listen(messages, freq)

	for message := range messages {
		if handler, args := b.route(&message); handler != nil {
			handler(Context{Message: &message, Args: args})
		}
	}
}

func (b *Bot) poll(
	messages chan Message,
	queries chan Query,
	callbacks chan Callback,
	timeout time.Duration,
) {
	var latestUpdate int64 = 0

	for {
		updates, err := getUpdates(b.Token,
			latestUpdate+1,
			int64(timeout/time.Second),
		)

		if err != nil {
			log.Println("failed to get updates:", err)
			continue
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

			latestUpdate = update.ID
		}
	}

}

func (b *Bot) route(message *Message) (Handler, map[string]string) {
	for reg, handler := range b.handlers {
		if matches := reg.FindStringSubmatch(message.Text); len(matches) > 0 {
			args := map[string]string{}

			for x, name := range reg.SubexpNames() {
				if x != 0 {
					args[name] = matches[x]
				}
			}

			return handler, args
		}
	}

	return nil, nil
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
func (b *Bot) ForwardMessage(recipient Recipient, message Message) error {
	params := map[string]string{
		"chat_id":      recipient.Destination(),
		"from_chat_id": strconv.Itoa(message.Origin().ID),
		"message_id":   strconv.Itoa(message.ID),
	}

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
func (b *Bot) SendLocation(recipient Recipient, geo *Location, options *SendOptions) error {
	params := map[string]string{
		"chat_id":   recipient.Destination(),
		"latitude":  fmt.Sprintf("%f", geo.Latitude),
		"longitude": fmt.Sprintf("%f", geo.Longitude),
	}

	if options != nil {
		embedSendOptions(params, options)
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

// SendVenue sends a venue object to recipient.
func (b *Bot) SendVenue(recipient Recipient, venue *Venue, options *SendOptions) error {
	params := map[string]string{
		"chat_id":   recipient.Destination(),
		"latitude":  fmt.Sprintf("%f", venue.Location.Latitude),
		"longitude": fmt.Sprintf("%f", venue.Location.Longitude),
		"title":     venue.Title,
		"address":   venue.Address}
	if venue.Foursquare_id != "" {
		params["foursquare_id"] = venue.Foursquare_id
	}

	if options != nil {
		embedSendOptions(params, options)
	}

	responseJSON, err := sendCommand("sendVenue", b.Token, params)
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
func (b *Bot) SendChatAction(recipient Recipient, action string) error {
	params := map[string]string{
		"chat_id": recipient.Destination(),
		"action":  action,
	}

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

// Respond publishes a set of responses for an inline query.
// This function is deprecated in favor of AnswerInlineQuery.
func (b *Bot) Respond(query Query, results []Result) error {
	params := map[string]string{
		"inline_query_id": query.ID,
	}

	if res, err := json.Marshal(results); err == nil {
		params["results"] = string(res)
	} else {
		return err
	}

	responseJSON, err := sendCommand("answerInlineQuery", b.Token, params)
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

// AnswerInlineQuery sends a response for a given inline query. A query can
// only be responded to once, subsequent attempts to respond to the same query
// will result in an error.
func (b *Bot) AnswerInlineQuery(query *Query, response *QueryResponse) error {
	response.QueryID = query.ID

	responseJSON, err := sendCommand("answerInlineQuery", b.Token, response)
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

// AnswerCallbackQuery sends a response for a given callback query. A callback can
// only be responded to once, subsequent attempts to respond to the same callback
// will result in an error.
func (b *Bot) AnswerCallbackQuery(callback *Callback, response *CallbackResponse) error {
	response.CallbackID = callback.ID

	responseJSON, err := sendCommand("answerCallbackQuery", b.Token, response)
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

// GetFile returns full file object including File.FilePath, which allow you to load file from Telegram
//
// Usually File objects does not contain any FilePath so you need to perform additional request
func (b *Bot) GetFile(fileID string) (File, error) {
	params := map[string]string{
		"file_id": fileID,
	}
	responseJSON, err := sendCommand("getFile", b.Token, params)
	if err != nil {
		return File{}, err
	}

	var responseRecieved struct {
		Ok          bool
		Description string
		Result      File
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return File{}, err
	}

	if !responseRecieved.Ok {
		return File{}, fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	return responseRecieved.Result, nil
}

// LeaveChat makes bot leave a group, supergroup or channel.
func (b *Bot) LeaveChat(recipient Recipient) error {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}
	responseJSON, err := sendCommand("leaveChat", b.Token, params)
	if err != nil {
		return err
	}

	var responseRecieved struct {
		Ok          bool
		Description string
		Result      bool
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return err
	}

	if !responseRecieved.Ok {
		return fmt.Errorf("telebot: leaveChat failure %s",
			responseRecieved.Description)
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
	responseJSON, err := sendCommand("getChat", b.Token, params)
	if err != nil {
		return Chat{}, err
	}

	var responseRecieved struct {
		Ok          bool
		Description string
		Result      Chat
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return Chat{}, err
	}

	if !responseRecieved.Ok {
		return Chat{}, fmt.Errorf("telebot: getChat failure %s",
			responseRecieved.Description)
	}

	return responseRecieved.Result, nil
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
	responseJSON, err := sendCommand("getChatAdministrators", b.Token, params)
	if err != nil {
		return []ChatMember{}, err
	}

	var responseRecieved struct {
		Ok          bool
		Result      []ChatMember
		Description string `json:"description",omitempty`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return []ChatMember{}, err
	}

	if !responseRecieved.Ok {
		return []ChatMember{}, fmt.Errorf("telebot: getChatAdministrators failure %s", responseRecieved.Description)
	}

	return responseRecieved.Result, nil
}

// GetChatMembersCount return the number of members in a chat.
//
// Returns Int on success.
func (b *Bot) GetChatMembersCount(recipient Recipient) (int, error) {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}
	responseJSON, err := sendCommand("getChatMembersCount", b.Token, params)
	if err != nil {
		return 0, err
	}

	var responseRecieved struct {
		Ok          bool
		Result      int
		Description string `json:"description",omitempty`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return 0, err
	}

	if !responseRecieved.Ok {
		return 0, fmt.Errorf("telebot: getChatMembersCount failure %s", responseRecieved.Description)
	}

	return responseRecieved.Result, nil
}

// GetUserProfilePhotos return list of profile pictures for a user.
//
// Returns a UserProfilePhotos object.
func (b *Bot) GetUserProfilePhotos(recipient Recipient) (UserProfilePhotos, error) {
	params := map[string]string{
		"user_id": recipient.Destination(),
	}
	responseJSON, err := sendCommand("getUserProfilePhotos", b.Token, params)
	if err != nil {
		return UserProfilePhotos{}, err
	}

	var responseRecieved struct {
		Ok          bool
		Result      UserProfilePhotos
		Description string `json:"description",omitempty`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return UserProfilePhotos{}, err
	}

	if !responseRecieved.Ok {
		return UserProfilePhotos{}, fmt.Errorf("telebot: getUserProfilePhotos failure %s", responseRecieved.Description)
	}

	return responseRecieved.Result, nil
}

// GetChatMember return information about a member of a chat.
//
// Returns a ChatMember object on success.
func (b *Bot) GetChatMember(recipient Recipient, user User) (ChatMember, error) {
	params := map[string]string{
		"chat_id": recipient.Destination(),
		"user_id": user.Destination(),
	}
	responseJSON, err := sendCommand("getChatMember", b.Token, params)
	if err != nil {
		return ChatMember{}, err
	}

	var responseRecieved struct {
		Ok          bool
		Result      ChatMember
		Description string `json:"description",omitempty`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return ChatMember{}, err
	}

	if !responseRecieved.Ok {
		return ChatMember{}, fmt.Errorf("telebot: getChatMember failure %s", responseRecieved.Description)
	}

	return responseRecieved.Result, nil
}

// GetFileDirectURL returns direct url for files using FileId which you can get from File object
func (b *Bot) GetFileDirectURL(fileID string) (string, error) {
	f, err := b.GetFile(fileID)
	if err != nil {
		return "", err
	}
	return "https://api.telegram.org/file/bot" + b.Token + "/" + f.FilePath, nil
}
