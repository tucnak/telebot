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
	go b.poll(subscription, nil, nil, timeout)
}

// Start periodically polls messages, updates and callbacks into their
// corresponding channels of the bot object.
//
// NOTE: It's a blocking method!
func (b *Bot) Start(timeout time.Duration) {
	b.poll(b.Messages, b.Queries, b.Callbacks, timeout)
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
) {
	var latestUpdate int64

	for {
		updates, err := b.getUpdates(latestUpdate+1, timeout)

		if err != nil {
			b.debug(errors.Wrap(err, "getUpdates() failed"))
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

func (b *Bot) Send(to Recipient, what interface{}, how ...interface{}) (*Message, error) {
	options := extractOptions(how)

	switch object := what.(type) {
	case string:
		return b.sendText(to, object, options)
	case Sendable:
		return object.Send(b, to, options)
	default:
		panic(fmt.Sprintf("telebot: object %v is not Sendable", object))
	}
}

func (b *Bot) Reply(to *Message, what interface{}, how ...interface{}) (*Message, error) {
	options := extractOptions(how)
	if options == nil {
		options = &SendOptions{}
	}

	options.ReplyTo = to

	return b.Send(to.Chat, what, options)
}

func (b *Bot) Forward(to Recipient, what *Message, how ...interface{}) (*Message, error) {
	params := map[string]string{
		"chat_id":      to.Destination(),
		"from_chat_id": strconv.Itoa(what.Origin().ID),
		"message_id":   strconv.Itoa(what.ID),
	}

	options := extractOptions(how)
	if options == nil {
		options = &SendOptions{}
	}
	embedSendOptions(params, options)

	respJSON, err := b.sendCommand("forwardMessage", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}

func (b *Bot) Delete(what *Message) error {
	params := map[string]string{
		"chat_id":    what.Chat.Destination(),
		"message_id": strconv.Itoa(what.ID),
	}

	respJSON, err := b.sendCommand("deleteMessage", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
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

	var responseRecieved struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseRecieved.Ok {
		return errors.Errorf("api error: %s", responseRecieved.Description)
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

	var responseRecieved struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseRecieved.Ok {
		return errors.Errorf("api error: %s", responseRecieved.Description)
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

	var responseRecieved struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseRecieved.Ok {
		return errors.Errorf("api error: %s", responseRecieved.Description)
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

	var responseRecieved struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseRecieved.Ok {
		return errors.Errorf("api error: %s", responseRecieved.Description)
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

	var responseRecieved struct {
		Ok          bool
		Description string
		Result      File
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return File{}, errors.Wrap(err, "bad response json")
	}

	if !responseRecieved.Ok {
		return File{}, errors.Errorf("api error: %s", responseRecieved.Description)

	}

	return responseRecieved.Result, nil
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

	var responseRecieved struct {
		Ok          bool
		Description string
		Result      bool
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !responseRecieved.Ok {
		return errors.Errorf("api error: %s", responseRecieved.Description)
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

	var responseRecieved struct {
		Ok          bool
		Description string
		Result      Chat
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return Chat{}, errors.Wrap(err, "bad response json")
	}

	if !responseRecieved.Ok {
		return Chat{}, errors.Errorf("api error: %s", responseRecieved.Description)
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
	responseJSON, err := b.sendCommand("getChatAdministrators", params)
	if err != nil {
		return []ChatMember{}, err
	}

	var responseRecieved struct {
		Ok          bool
		Result      []ChatMember
		Description string `json:"description"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return []ChatMember{}, errors.Wrap(err, "bad response json")
	}

	if !responseRecieved.Ok {
		return []ChatMember{}, errors.Errorf("api error: %s", responseRecieved.Description)
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
	responseJSON, err := b.sendCommand("getChatMembersCount", params)
	if err != nil {
		return 0, err
	}

	var responseRecieved struct {
		Ok          bool
		Result      int
		Description string `json:"description"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return 0, errors.Wrap(err, "bad response json")
	}

	if !responseRecieved.Ok {
		return 0, errors.Errorf("api error: %s", responseRecieved.Description)
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
	responseJSON, err := b.sendCommand("getUserProfilePhotos", params)
	if err != nil {
		return UserProfilePhotos{}, err
	}

	var responseRecieved struct {
		Ok          bool
		Result      UserProfilePhotos
		Description string `json:"description"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return UserProfilePhotos{}, errors.Wrap(err, "bad response json")
	}

	if !responseRecieved.Ok {
		return UserProfilePhotos{}, errors.Errorf("api error: %s", responseRecieved.Description)
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
	responseJSON, err := b.sendCommand("getChatMember", params)
	if err != nil {
		return ChatMember{}, err
	}

	var responseRecieved struct {
		Ok          bool
		Result      ChatMember
		Description string `json:"description"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return ChatMember{}, errors.Wrap(err, "bad response json")
	}

	if !responseRecieved.Ok {
		return ChatMember{}, errors.Errorf("api error: %s", responseRecieved.Description)
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

// EditMessageText used to edit already sent message with known recepient and message id.
//
// On success, returns edited message object
func (b *Bot) EditMessageText(recipient Recipient, messageID int, message string, sendOptions *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":    recipient.Destination(),
		"message_id": strconv.Itoa(messageID),
		"text":       message,
	}

	if sendOptions != nil {
		embedSendOptions(params, sendOptions)
	}

	responseJSON, err := b.sendCommand("editMessageText", params)
	if err != nil {
		return nil, err
	}

	var responseRecieved struct {
		Ok          bool
		Description string
		Message     Message `json:"result"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return nil, err
	}

	if !responseRecieved.Ok {
		return nil, fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	return &responseRecieved.Message, err

}

// EditInlineMessageText used to edit already sent inline message with known inline message id.
//
// On success, returns edited message object
func (b *Bot) EditInlineMessageText(messageID string, message string, sendOptions *SendOptions) (*Message, error) {
	params := map[string]string{
		"inline_message_id": messageID,
		"text":              message,
	}

	if sendOptions != nil {
		embedSendOptions(params, sendOptions)
	}

	responseJSON, err := b.sendCommand("editMessageText", params)
	if err != nil {
		return nil, err
	}

	var responseRecieved struct {
		Ok          bool
		Description string
		Message     Message `json:"result"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return nil, err
	}

	if !responseRecieved.Ok {
		return nil, fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	return &responseRecieved.Message, err

}

// EditMessageCaption used to edit already sent photo caption with known recepient and message id.
//
// On success, returns edited message object
func (b *Bot) EditMessageCaption(recipient Recipient, messageID int, caption string, inlineKeyboard *InlineKeyboardMarkup) (*Message, error) {
	params := map[string]string{
		"chat_id":    recipient.Destination(),
		"message_id": strconv.Itoa(messageID),
		"caption":    caption,
	}

	if inlineKeyboard != nil {
		embedSendOptions(params, &SendOptions{
			ReplyMarkup: &ReplyMarkup{
				InlineKeyboard: inlineKeyboard.InlineKeyboard,
			},
		})
	}

	responseJSON, err := b.sendCommand("editMessageCaption", params)
	if err != nil {
		return nil, err
	}

	var responseRecieved struct {
		Ok          bool
		Description string
		Message     Message `json:"result"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return nil, err
	}

	if !responseRecieved.Ok {
		return nil, fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	return &responseRecieved.Message, err

}

// EditInlineMessageCaption used to edit already sent photo caption with known inline message id.
//
// On success, returns edited message object
func (b *Bot) EditInlineMessageCaption(messageID string, caption string, inlineKeyboard *InlineKeyboardMarkup) (*Message, error) {
	params := map[string]string{
		"inline_message_id": messageID,
		"caption":           caption,
	}

	if inlineKeyboard != nil {
		embedSendOptions(params, &SendOptions{
			ReplyMarkup: &ReplyMarkup{
				InlineKeyboard: inlineKeyboard.InlineKeyboard,
			},
		})
	}

	responseJSON, err := b.sendCommand("editMessageCaption", params)
	if err != nil {
		return nil, err
	}

	var responseRecieved struct {
		Ok          bool
		Description string
		Message     Message `json:"result"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return nil, err
	}

	if !responseRecieved.Ok {
		return nil, fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	return &responseRecieved.Message, err

}

// EditMessageReplyMarkup used to edit already sent message inline keyboard markup with known recepient and message id.
//
// On success, returns edited message object
func (b *Bot) EditMessageReplyMarkup(recipient Recipient, messageID int, inlineKeyboard *InlineKeyboardMarkup) (*Message, error) {
	params := map[string]string{
		"chat_id":    recipient.Destination(),
		"message_id": strconv.Itoa(messageID),
	}

	if inlineKeyboard != nil {
		embedSendOptions(params, &SendOptions{
			ReplyMarkup: &ReplyMarkup{
				InlineKeyboard: inlineKeyboard.InlineKeyboard,
			},
		})
	}

	responseJSON, err := b.sendCommand("editMessageReplyMarkup", params)
	if err != nil {
		return nil, err
	}

	var responseRecieved struct {
		Ok          bool
		Description string
		Message     Message `json:"result"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return nil, err
	}

	if !responseRecieved.Ok {
		return nil, fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	return &responseRecieved.Message, err

}

// EditInlineMessageReplyMarkup used to edit already sent message inline keyboard markup with known inline message id.
//
// On success, returns edited message object
func (b *Bot) EditInlineMessageReplyMarkup(messageID string, caption string, inlineKeyboard *InlineKeyboardMarkup) (*Message, error) {
	params := map[string]string{
		"inline_message_id": messageID,
	}

	if inlineKeyboard != nil {
		embedSendOptions(params, &SendOptions{
			ReplyMarkup: &ReplyMarkup{
				InlineKeyboard: inlineKeyboard.InlineKeyboard,
			},
		})
	}

	responseJSON, err := b.sendCommand("editMessageReplyMarkup", params)
	if err != nil {
		return nil, err
	}

	var responseRecieved struct {
		Ok          bool
		Description string
		Message     Message `json:"result"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return nil, err
	}

	if !responseRecieved.Ok {
		return nil, fmt.Errorf("telebot: %s", responseRecieved.Description)
	}

	return &responseRecieved.Message, err

}
