package telebot

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// NewBot does try to build a Bot with token `token`, which
// is a secret API key assigned to particular bot.
func NewBot(pref Settings) (*Bot, error) {
	if pref.Updates == 0 {
		pref.Updates = 100
	}

	bot := &Bot{
		Token:   pref.Token,
		Updates: make(chan Update, pref.Updates),
		Poller:  pref.Poller,

		handlers: make(map[string]interface{}),
	}

	user, err := bot.getMe()
	if err != nil {
		return nil, err
	}

	bot.Me = user
	return bot, nil
}

// Bot represents a separate Telegram bot instance.
type Bot struct {
	Me      *User
	Token   string
	Updates chan Update
	Poller  Poller
	Errors  chan error

	handlers map[string]interface{}
}

// Settings represents a utility struct for passing certain
// properties of a bot around and is required to make bots.
type Settings struct {
	// Telegram token
	Token string

	// Updates channel capacity
	Updates int // Default: 100

	// Poller is the provider of Updates.
	Poller Poller
}

// Update object represents an incoming update.
type Update struct {
	ID int `json:"update_id"`

	Message           *Message  `json:"message,omitempty"`
	EditedMessage     *Message  `json:"edited_message,omitempty"`
	ChannelPost       *Message  `json:"channel_post,omitempty"`
	EditedChannelPost *Message  `json:"edited_channel_post,omitempty"`
	Callback          *Callback `json:"callback_query,omitempty"`
	Query             *Query    `json:"inline_query,omitempty"`
}

// Handle lets you set the handler for some command name or
// one of the supported endpoints.
//
// Example:
//
//     tb.Handle("/help", func (m *tb.Message) {})
//     tb.Handle(tb.OnEditedMessage, func (m *tb.Message) {})
//     tb.Handle(tb.OnQuery, func (q *tb.Query) {})
//
func (b *Bot) Handle(endpoint string, handler interface{}) {
	b.handlers[endpoint] = handler
}

var cmdRx = regexp.MustCompile(`^(\/\w+)(@(\w+))?(\s|$)`)

func (b *Bot) handleCommand(m *Message, cmdName, cmdBot string) bool {
	// Group-syntax: "/cmd@bot"
	if cmdBot != "" && !strings.EqualFold(b.Me.Username, cmdBot) {
		return false
	}

	if handler, ok := b.handlers[cmdName]; ok {
		if handler, ok := handler.(func(*Message)); ok {
			go handler(m)
			return true
		}
	}

	return false
}

// Start brings bot into motion by consuming incoming
// updates (see Bot.Updates channel).
func (b *Bot) Start() {
	if b.Poller == nil {
		panic("telebot: can't start without a poller")
	}

	go b.Poller.Poll(b, b.Updates)

	for upd := range b.Updates {
		if upd.Message != nil {
			m := upd.Message

			// Text message
			if m.Text != "" {
				if m.Text[0] == '\a' {
					continue
				}

				match := cmdRx.FindAllStringSubmatch(m.Text, -1)

				// Command found
				if match != nil {
					if b.handleCommand(m, match[0][1], match[0][3]) {
						continue
					}
				}
			}

			wasAdded := m.NewChatMembers != nil &&
				isUserInList(b.Me, m.NewChatMembers)

			if m.ChatCreated || wasAdded {
				if handler, ok := b.handlers[string(OnAddedToGroup)]; ok {
					if handler, ok := handler.(func(*Message)); ok {
						go handler(m)
						continue
					}
				}

				continue
			}

			// OnMessage
			if handler, ok := b.handlers[string(OnMessage)]; ok {
				if handler, ok := handler.(func(*Message)); ok {
					go handler(m)
					continue
				}
			}
			continue
		}

		if upd.EditedMessage != nil {
			if handler, ok := b.handlers[OnEditedMessage]; ok {
				if handler, ok := handler.(func(*Message)); ok {
					handler(upd.EditedMessage)
				}
			}
			continue
		}

		if upd.ChannelPost != nil {
			if handler, ok := b.handlers[OnChannelPost]; ok {
				if handler, ok := handler.(func(*Message)); ok {
					handler(upd.ChannelPost)
				}
			}
			continue
		}

		if upd.EditedChannelPost != nil {
			if handler, ok := b.handlers[OnEditedChannelPost]; ok {
				if handler, ok := handler.(func(*Message)); ok {
					handler(upd.EditedChannelPost)
				}
			}
			continue
		}

		if upd.Callback != nil {
			if handler, ok := b.handlers[OnCallback]; ok {
				if handler, ok := handler.(func(*Callback)); ok {
					handler(upd.Callback)
				}
			}
			continue
		}

		if upd.Query != nil {
			if handler, ok := b.handlers[OnQuery]; ok {
				if handler, ok := handler.(func(*Query)); ok {
					handler(upd.Query)
				}
			}
			continue
		}
	}
}

// Send accepts 2+ arguments, starting with destination chat, followed by
// some Sendable (or string!) and optional send options.
//
// Note: since most arguments are of type interface{}, make sure to pass
//       them by-pointer, NOT by-value, which will result in a panic.
//
// What is a send option exactly? It can be one of the following types:
//
//     - *SendOptions (the actual object accepted by Telegram API)
//     - *ReplyMarkup (a component of SendOptions)
//     - Option (a shorcut flag for popular options)
//     - ParseMode (HTML, Markdown, etc)
//
// This function will panic upon unsupported payloads and options!
func (b *Bot) Send(to Recipient, what interface{}, options ...interface{}) (*Message, error) {
	sendOpts := extractOptions(options)

	switch object := what.(type) {
	case string:
		return b.sendText(to, object, sendOpts)
	case Sendable:
		return object.Send(b, to, sendOpts)
	default:
		panic(fmt.Sprintf("telebot: object %v is not Sendable", object))
	}
}

// Reply behaves just like Send() with an exception of "reply-to" indicator.
//
func (b *Bot) Reply(to *Message, what interface{}, options ...interface{}) (*Message, error) {
	// This function will panic upon unsupported payloads and options!
	sendOpts := extractOptions(options)
	if sendOpts == nil {
		sendOpts = &SendOptions{}
	}

	sendOpts.ReplyTo = to

	return b.Send(to.Chat, what, sendOpts)
}

// Forward behaves just like Send() but of all options it
// only supports Silent (see Bots API).
//
// This function will panic upon unsupported payloads and options!
func (b *Bot) Forward(to Recipient, what *Message, options ...interface{}) (*Message, error) {
	params := map[string]string{
		"chat_id":      to.Recipient(),
		"from_chat_id": what.Chat.Recipient(),
		"message_id":   strconv.Itoa(what.ID),
	}

	sendOpts := extractOptions(options)
	if sendOpts == nil {
		sendOpts = &SendOptions{}
	}
	embedSendOptions(params, sendOpts)

	respJSON, err := b.sendCommand("forwardMessage", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}

// Edit is magic, it lets you change already sent message.
//
// Use cases:
//
//     b.Edit(msg, msg.Text, newMarkup)
//     b.Edit(msg, "new <b>text</b>", tb.ModeHTML)
//
func (b *Bot) Edit(originalMsg Editable, text string, options ...interface{}) (*Message, error) {
	messageID, chatID := originalMsg.MessageSig()
	// TODO: add support for inline messages (chatID = 0)

	params := map[string]string{"text": text}

	// if inline message
	if chatID == 0 {
		params["inline_message_id"] = strconv.Itoa(messageID)
	} else {
		params["chat_id"] = strconv.FormatInt(chatID, 10)
		params["message_id"] = strconv.Itoa(messageID)
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	respJSON, err := b.sendCommand("editMessageText", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}

// EditCaption used to edit already sent photo caption with known recepient and message id.
//
// On success, returns edited message object
func (b *Bot) EditCaption(originalMsg Editable, caption string) (*Message, error) {
	messageID, chatID := originalMsg.MessageSig()

	params := map[string]string{"caption": caption}

	// if inline message
	if chatID == 0 {
		params["inline_message_id"] = strconv.Itoa(messageID)
	} else {
		params["chat_id"] = strconv.FormatInt(chatID, 10)
		params["message_id"] = strconv.Itoa(messageID)
	}

	respJSON, err := b.sendCommand("editMessageCaption", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}

// Delete removes the message, including service messages,
// with the following limitations:
//
//     * A message can only be deleted if it was sent less than 48 hours ago.
//     * Bots can delete outgoing messages in groups and supergroups.
//     * Bots granted can_post_messages permissions can delete outgoing
//       messages in channels.
//     * If the bot is an administrator of a group, it can delete any message there.
//     * If the bot has can_delete_messages permission in a supergroup or a
//       channel, it can delete any message there.
//
func (b *Bot) Delete(message Editable) error {
	messageID, chatID := message.MessageSig()

	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": strconv.Itoa(messageID),
	}

	respJSON, err := b.sendCommand("deleteMessage", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// Action updates a chat action for recipient.
//
// Chat action is a status message that recipient would see where
// you typically see "Harry is typing" status message. The only
// difference is that bots' chat actions live only for 5 seconds
// and die just once the client recieves a message from the bot.
//
// Currently, Telegram supports only a narrow range of possible
// actions, these are aligned as constants of this package.
func (b *Bot) Action(recipient Recipient, action ChatAction) error {
	params := map[string]string{
		"chat_id": recipient.Recipient(),
		"action":  string(action),
	}

	respJSON, err := b.sendCommand("sendChatAction", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// AnswerInlineQuery sends a response for a given inline query. A query can
// only be responded to once, subsequent attempts to respond to the same query
// will result in an error.
func (b *Bot) AnswerInlineQuery(query *Query, response *QueryResponse) error {
	response.QueryID = query.ID

	respJSON, err := b.sendCommand("answerInlineQuery", response)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// AnswerCallbackQuery sends a response for a given callback query. A callback can
// only be responded to once, subsequent attempts to respond to the same callback
// will result in an error.
func (b *Bot) AnswerCallbackQuery(callback *Callback, response *CallbackResponse) error {
	response.CallbackID = callback.ID

	respJSON, err := b.sendCommand("answerCallbackQuery", response)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// FileByID returns full file object including File.FilePath, allowing you to
// download the file from the server.
//
// Usually, Telegram-provided File objects miss FilePath so you might need to
// perform an additional request to fetch them.
func (b *Bot) FileByID(fileID string) (File, error) {
	params := map[string]string{
		"file_id": fileID,
	}

	respJSON, err := b.sendCommand("getFile", params)
	if err != nil {
		return File{}, err
	}

	var resp struct {
		Ok          bool
		Description string
		Result      File
	}

	err = json.Unmarshal(respJSON, &resp)
	if err != nil {
		return File{}, errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return File{}, errors.Errorf("api error: %s", resp.Description)

	}

	return resp.Result, nil
}

// Leave makes bot leave a group, supergroup or channel.
func (b *Bot) Leave(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	respJSON, err := b.sendCommand("leaveChat", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// ChatByID fetches chat info of its ID.
//
// Including current name of the user for one-on-one conversations,
// current username of a user, group or channel, etc.
//
// Returns a Chat object on success.
func (b *Bot) ChatByID(id string) (*Chat, error) {
	params := map[string]string{
		"chat_id": id,
	}

	respJSON, err := b.sendCommand("getChat", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Ok          bool
		Description string
		Result      *Chat
	}

	err = json.Unmarshal(respJSON, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return nil, errors.Errorf("api error: %s", resp.Description)
	}

	return resp.Result, nil
}

// AdminsOf return a member list of chat admins.
//
// On success, returns an Array of ChatMember objects that
// contains information about all chat administrators except other bots.
// If the chat is a group or a supergroup and
// no administrators were appointed, only the creator will be returned.
func (b *Bot) AdminsOf(chat *Chat) ([]ChatMember, error) {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	respJSON, err := b.sendCommand("getChatAdministrators", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Ok          bool
		Result      []ChatMember
		Description string `json:"description"`
	}

	err = json.Unmarshal(respJSON, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return nil, errors.Errorf("api error: %s", resp.Description)
	}

	return resp.Result, nil
}

// Len return the number of members in a chat.
func (b *Bot) Len(chat *Chat) (int, error) {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	respJSON, err := b.sendCommand("getChatMembersCount", params)
	if err != nil {
		return 0, err
	}

	var resp struct {
		Ok          bool
		Result      int
		Description string `json:"description"`
	}

	err = json.Unmarshal(respJSON, &resp)
	if err != nil {
		return 0, errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return 0, errors.Errorf("api error: %s", resp.Description)
	}

	return resp.Result, nil
}

// ProfilePhotosOf return list of profile pictures for a user.
func (b *Bot) ProfilePhotosOf(user *User) ([]Photo, error) {
	params := map[string]string{
		"user_id": user.Recipient(),
	}

	respJSON, err := b.sendCommand("getUserProfilePhotos", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Ok     bool
		Result struct {
			Count  int     `json:"total_count"`
			Photos []Photo `json:"photos"`
		}

		Description string `json:"description"`
	}

	err = json.Unmarshal(respJSON, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return nil, errors.Errorf("api error: %s", resp.Description)
	}

	return resp.Result.Photos, nil
}

// ChatMemberOf return information about a member of a chat.
//
// Returns a ChatMember object on success.
func (b *Bot) ChatMemberOf(chat *Chat, user *User) (*ChatMember, error) {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"user_id": user.Recipient(),
	}

	respJSON, err := b.sendCommand("getChatMember", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Ok          bool
		Result      *ChatMember
		Description string `json:"description"`
	}

	err = json.Unmarshal(respJSON, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return nil, errors.Errorf("api error: %s", resp.Description)
	}

	return resp.Result, nil
}

// FileURLByID returns direct url for files using FileId which you can get from File object
func (b *Bot) FileURLByID(fileID string) (string, error) {
	f, err := b.FileByID(fileID)
	if err != nil {
		return "", err
	}
	return "https://api.telegram.org/file/bot" + b.Token + "/" + f.FilePath, nil
}
