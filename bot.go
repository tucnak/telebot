package telebot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		stop:     make(chan struct{}),
		reporter: pref.Reporter,
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

	handlers map[string]interface{}
	reporter func(error)
	stop     chan struct{}
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

	// Reporter is a callback function that will get called
	// on any panics recovered from endpoint handlers.
	Reporter func(error)
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

	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
}

// ChosenInlineResult represents a result of an inline query that was chosen
// by the user and sent to their chat partner.
type ChosenInlineResult struct {
	ResultID string `json:"result_id"`
	Query    string `json:"query"`
	// Inline messages only!
	MessageID string `json:"inline_message_id"`

	From     User      `json:"from"`
	Location *Location `json:"location,omitempty"`
}

// Handle lets you set the handler for some command name or
// one of the supported endpoints.
//
// Example:
//
//     b.handle("/help", func (m *tb.Message) {})
//     b.handle(tb.OnEdited, func (m *tb.Message) {})
//     b.handle(tb.OnQuery, func (q *tb.Query) {})
//
//     // make a hook for one of your preserved (by-pointer)
//     // inline buttons.
//     b.handle(&inlineButton, func (c *tb.Callback) {})
//
func (b *Bot) Handle(endpoint interface{}, handler interface{}) {
	switch end := endpoint.(type) {
	case string:
		b.handlers[end] = handler
	case CallbackEndpoint:
		b.handlers[end.CallbackUnique()] = handler
	default:
		panic("telebot: unsupported endpoint")
	}
}

var (
	cmdRx   = regexp.MustCompile(`^(\/\w+)(@(\w+))?(\s|$)(.+)?`)
	cbackRx = regexp.MustCompile(`^\f(\w+)(\|(.+))?$`)
)

func (b *Bot) handleCommand(m *Message, cmdName, cmdBot string) bool {

	return false
}

// Start brings bot into motion by consuming incoming
// updates (see Bot.Updates channel).
func (b *Bot) Start() {
	if b.Poller == nil {
		panic("telebot: can't start without a poller")
	}

	stopPoller := make(chan struct{})

	go b.Poller.Poll(b, b.Updates, stopPoller)

	for {
		select {
		// handle incoming updates
		case upd := <-b.Updates:
			b.incomingUpdate(&upd)

		// call to stop polling
		case <-b.stop:
			stopPoller <- struct{}{}

		// polling has stopped
		case <-stopPoller:
			return
		}
	}
}

func (b *Bot) incomingUpdate(upd *Update) {
	if upd.Message != nil {
		m := upd.Message

		if m.PinnedMessage != nil {
			b.handle(OnPinned, m)
			return
		}

		// Commands
		if m.Text != "" {
			// Filtering malicious messsages
			if m.Text[0] == '\a' {
				return
			}

			match := cmdRx.FindAllStringSubmatch(m.Text, -1)

			// Command found - handle and return
			if match != nil {
				// Syntax: "</command>@<bot> <payload>"
				command, botName := match[0][1], match[0][3]
				m.Payload = match[0][5]

				if botName != "" && !strings.EqualFold(b.Me.Username, botName) {
					return
				}

				if b.handle(command, m) {
					return
				}
			}

			// 1:1 satisfaction
			if b.handle(m.Text, m) {
				return
			}

			// OnText
			b.handle(OnText, m)
			return
		}

		// on media
		if b.handleMedia(m) {
			return
		}

		// OnAddedToGroup
		wasAdded := (m.UserJoined != nil && m.UserJoined.ID == b.Me.ID) ||
			(m.UsersJoined != nil && isUserInList(b.Me, m.UsersJoined))
		if m.GroupCreated || m.SuperGroupCreated || wasAdded {
			b.handle(OnAddedToGroup, m)
			return
		}

		if m.UserJoined != nil {
			b.handle(OnUserJoined, m)
			return
		}

		if m.UsersJoined != nil {
			for _, user := range m.UsersJoined {
				m.UserJoined = &user
				b.handle(OnUserJoined, m)
			}

			return
		}

		if m.UserLeft != nil {
			b.handle(OnUserLeft, m)
			return
		}

		if m.NewGroupTitle != "" {
			b.handle(OnNewGroupTitle, m)
			return
		}

		if m.NewGroupPhoto != nil {
			b.handle(OnNewGroupPhoto, m)
			return
		}

		if m.GroupPhotoDeleted {
			b.handle(OnGroupPhotoDeleted, m)
			return
		}

		if m.MigrateTo != 0 {
			if handler, ok := b.handlers[OnMigration]; ok {
				if handler, ok := handler.(func(int64, int64)); ok {
					// i'm not 100% sure that any of the values
					// won't be cached, so I pass them all in:
					go func(b *Bot, handler func(int64, int64), from, to int64) {
						defer b.deferDebug()
						handler(from, to)
					}(b, handler, m.MigrateFrom, m.MigrateTo)

				} else {
					panic("telebot: migration handler is bad")
				}
			}

			return
		}

		return
	}

	if upd.EditedMessage != nil {
		b.handle(OnEdited, upd.EditedMessage)
		return
	}

	if upd.ChannelPost != nil {
		b.handle(OnChannelPost, upd.ChannelPost)
		return
	}

	if upd.EditedChannelPost != nil {
		b.handle(OnEditedChannelPost, upd.EditedChannelPost)
		return
	}

	if upd.Callback != nil {
		if upd.Callback.Data != "" {
			data := upd.Callback.Data

			if data[0] == '\f' {
				match := cbackRx.FindAllStringSubmatch(data, -1)

				if match != nil {
					unique, payload := match[0][1], match[0][3]

					if handler, ok := b.handlers["\f"+unique]; ok {
						if handler, ok := handler.(func(*Callback)); ok {
							upd.Callback.Data = payload
							// i'm not 100% sure that any of the values
							// won't be cached, so I pass them all in:
							go func(b *Bot, handler func(*Callback), c *Callback) {
								defer b.deferDebug()
								handler(c)
							}(b, handler, upd.Callback)

							return
						}
					}

				}
			}
		}

		if handler, ok := b.handlers[OnCallback]; ok {
			if handler, ok := handler.(func(*Callback)); ok {
				// i'm not 100% sure that any of the values
				// won't be cached, so I pass them all in:
				go func(b *Bot, handler func(*Callback), c *Callback) {
					defer b.deferDebug()
					handler(c)
				}(b, handler, upd.Callback)

			} else {
				panic("telebot: callback handler is bad")
			}
		}
		return
	}

	if upd.Query != nil {
		if handler, ok := b.handlers[OnQuery]; ok {
			if handler, ok := handler.(func(*Query)); ok {
				// i'm not 100% sure that any of the values
				// won't be cached, so I pass them all in:
				go func(b *Bot, handler func(*Query), q *Query) {
					defer b.deferDebug()
					handler(q)
				}(b, handler, upd.Query)

			} else {
				panic("telebot: query handler is bad")
			}
		}
		return
	}

	if upd.ChosenInlineResult != nil {
		if handler, ok := b.handlers[OnChosenInlineResult]; ok {
			if handler, ok := handler.(func(*ChosenInlineResult)); ok {
				// i'm not 100% sure that any of the values
				// won't be cached, so I pass them all in:
				go func(b *Bot, handler func(*ChosenInlineResult),
					r *ChosenInlineResult) {
					defer b.deferDebug()
					handler(r)
				}(b, handler, upd.ChosenInlineResult)

			} else {
				panic("telebot: chosen inline result handler is bad")
			}
		}
		return
	}
}

func (b *Bot) handle(end string, m *Message) bool {
	handler, ok := b.handlers[end]
	if !ok {
		return false
	}

	if handler, ok := handler.(func(*Message)); ok {
		// i'm not 100% sure that any of the values
		// won't be cached, so I pass them all in:
		go func(b *Bot, handler func(*Message), m *Message) {
			defer b.deferDebug()
			handler(m)
		}(b, handler, m)

		return true
	}

	return false
}

func (b *Bot) handleMedia(m *Message) bool {
	if m.Photo != nil {
		b.handle(OnPhoto, m)
		return true
	}

	if m.Audio != nil {
		b.handle(OnAudio, m)
		return true
	}

	if m.Document != nil {
		b.handle(OnDocument, m)
		return true
	}

	if m.Sticker != nil {
		b.handle(OnSticker, m)
		return true
	}

	if m.Video != nil {
		b.handle(OnVideo, m)
		return true
	}

	if m.VideoNote != nil {
		b.handle(OnVideoNote, m)
		return true
	}

	if m.Contact != nil {
		b.handle(OnContact, m)
		return true
	}

	if m.Location != nil {
		b.handle(OnLocation, m)
		return true
	}

	if m.Venue != nil {
		b.handle(OnVenue, m)
		return true
	}

	return false
}

// Stop gracefully shuts the poller down.
func (b *Bot) Stop() {
	b.stop <- struct{}{}
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
		panic("telebot: unsupported sendable")
	}
}

// SendAlbum is used when sending multiple instances of media as a single
// message (so-called album).
//
// From all existing options, it only supports telebot.Silent.
func (b *Bot) SendAlbum(to Recipient, a Album, options ...interface{}) ([]Message, error) {
	media := make([]string, len(a))
	files := make(map[string]string)

	for i, x := range a {
		var (
			f                    *File
			caption              string
			mediaRepr, mediaType string
		)

		switch y := x.(type) {
		case *Photo:
			f = &y.File
			mediaType = "photo"
			caption = y.Caption
		case *Video:
			f = &y.File
			mediaType = "video"
			caption = y.Caption
		default:
			return nil, errors.Errorf("telebot: album entry #%d is not valid", i)
		}

		if f.InCloud() {
			mediaRepr = f.FileID
		} else if f.FileURL != "" {
			mediaRepr = f.FileURL
		} else if f.OnDisk() {
			mediaRepr = fmt.Sprintf("attach://%d", i)
			files[strconv.Itoa(i)] = f.FileLocal
		} else {
			return nil, errors.Errorf(
				"telebot: album entry #%d doesn't exist anywhere", i)
		}

		jsonRepr, _ := json.Marshal(map[string]string{
			"type":    mediaType,
			"media":   mediaRepr,
			"caption": caption,
		})

		media[i] = string(jsonRepr)
	}

	params := map[string]string{
		"chat_id": to.Recipient(),
		"media":   "[" + strings.Join(media, ",") + "]",
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	respJSON, err := b.sendFiles("sendMediaGroup", files, params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Ok          bool
		Result      []Message
		Description string
	}

	err = json.Unmarshal(respJSON, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return nil, errors.Errorf("api error: %s", resp.Description)
	}

	for attachName, _ := range files {
		i, _ := strconv.Atoi(attachName)

		var newID string
		if resp.Result[i].Photo != nil {
			newID = resp.Result[i].Photo.FileID
		} else {
			newID = resp.Result[i].Video.FileID
		}

		a[i].MediaFile().FileID = newID
	}

	return resp.Result, nil
}

// Reply behaves just like Send() with an exception of "reply-to" indicator.
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
	embedSendOptions(params, sendOpts)

	respJSON, err := b.Raw("forwardMessage", params)
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
//     // Edit live location:
//     b.Edit(liveMsg, tb.Location{42.1337, 69.4242})
//
func (b *Bot) Edit(message Editable, what interface{}, options ...interface{}) (*Message, error) {
	messageID, chatID := message.MessageSig()
	// TODO: add support for inline messages (chatID = 0)

	params := map[string]string{}

	switch v := what.(type) {
	case string:
		params["text"] = v
	case Location:
		params["latitude"] = fmt.Sprintf("%f", v.Lat)
		params["longitude"] = fmt.Sprintf("%f", v.Lng)
	default:
		panic("telebot: unsupported what argument")
	}

	// if inline message
	if chatID == 0 {
		params["inline_message_id"] = messageID
	} else {
		params["chat_id"] = strconv.FormatInt(chatID, 10)
		params["message_id"] = messageID
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	respJSON, err := b.Raw("editMessageText", params)
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
		params["inline_message_id"] = messageID
	} else {
		params["chat_id"] = strconv.FormatInt(chatID, 10)
		params["message_id"] = messageID
	}

	respJSON, err := b.Raw("editMessageCaption", params)
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
		"message_id": messageID,
	}

	respJSON, err := b.Raw("deleteMessage", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// Notify updates the chat action for recipient.
//
// Chat action is a status message that recipient would see where
// you typically see "Harry is typing" status message. The only
// difference is that bots' chat actions live only for 5 seconds
// and die just once the client recieves a message from the bot.
//
// Currently, Telegram supports only a narrow range of possible
// actions, these are aligned as constants of this package.
func (b *Bot) Notify(recipient Recipient, action ChatAction) error {
	params := map[string]string{
		"chat_id": recipient.Recipient(),
		"action":  string(action),
	}

	respJSON, err := b.Raw("sendChatAction", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// Answer sends a response for a given inline query. A query can only
// be responded to once, subsequent attempts to respond to the same query
// will result in an error.
func (b *Bot) Answer(query *Query, response *QueryResponse) error {
	response.QueryID = query.ID

	for _, result := range response.Results {
		result.Process()
	}

	respJSON, err := b.Raw("answerInlineQuery", response)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// Respond sends a response for a given callback query. A callback can
// only be responded to once, subsequent attempts to respond to the same callback
// will result in an error.
//
// Example:
//
//		bot.Respond(c)
//		bot.Respond(c, response)
//
func (b *Bot) Respond(callback *Callback, responseOptional ...*CallbackResponse) error {
	var response *CallbackResponse
	if responseOptional == nil {
		response = &CallbackResponse{}
	} else {
		response = responseOptional[0]
	}

	response.CallbackID = callback.ID
	respJSON, err := b.Raw("answerCallbackQuery", response)
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

	respJSON, err := b.Raw("getFile", params)
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

// Download saves the file from Telegram servers locally.
//
// Maximum file size to download is 20 MB.
func (b *Bot) Download(f *File, localFilename string) error {
	g, err := b.FileByID(f.FileID)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s",
		b.Token, g.FilePath)

	out, err := os.Create(localFilename)
	if err != nil {
		return wrapSystem(err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return wrapSystem(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return wrapSystem(err)
	}

	g.FileLocal = localFilename
	*f = g

	return nil
}

// StopLiveLocation should be called to stop broadcasting live message location
// before Location.LivePeriod expires.
//
// It supports telebot.ReplyMarkup.
func (b *Bot) StopLiveLocation(message Editable, options ...interface{}) (*Message, error) {
	messageID, chatID := message.MessageSig()

	params := map[string]string{
		"chat_id":    fmt.Sprintf("%d", chatID),
		"message_id": fmt.Sprintf("%d", messageID),
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	respJSON, err := b.Raw("stopMessageLiveLocation", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}

// GetInviteLink should be used to export chat's invite link.
func (b *Bot) GetInviteLink(chat *Chat) (string, error) {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	respJSON, err := b.Raw("exportChatInviteLink", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Ok          bool
		Description string
		Result      string
	}

	err = json.Unmarshal(respJSON, &resp)
	if err != nil {
		return "", errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return "", errors.Errorf("api error: %s", resp.Description)
	}

	return resp.Result, nil
}

// SetChatTitle should be used to update group title.
func (b *Bot) SetGroupTitle(chat *Chat, newTitle string) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"title":   newTitle,
	}

	respJSON, err := b.Raw("setChatTitle", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// SetGroupDescription should be used to update group title.
func (b *Bot) SetGroupDescription(chat *Chat, description string) error {
	params := map[string]string{
		"chat_id":     chat.Recipient(),
		"description": description,
	}

	respJSON, err := b.Raw("setChatDescription", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// SetGroupPhoto should be used to update group photo.
func (b *Bot) SetGroupPhoto(chat *Chat, p *Photo) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	respJSON, err := b.sendFiles("setChatPhoto",
		map[string]string{"photo": p.FileLocal}, params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// SetGroupStickerSet should be used to update group's group sticker set.
func (b *Bot) SetGroupStickerSet(chat *Chat, setName string) error {
	params := map[string]string{
		"chat_id":          chat.Recipient(),
		"sticker_set_name": setName,
	}

	respJSON, err := b.Raw("setChatStickerSet", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// DeleteGroupPhoto should be used to just remove group photo.
func (b *Bot) DeleteGroupPhoto(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	respJSON, err := b.Raw("deleteGroupPhoto", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// DeleteGroupStickerSet should be used to just remove group sticker set.
func (b *Bot) DeleteGroupStickerSet(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	respJSON, err := b.Raw("deleteChatStickerSet", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// Leave makes bot leave a group, supergroup or channel.
func (b *Bot) Leave(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	respJSON, err := b.Raw("leaveChat", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// Use this method to pin a message in a supergroup or a channel.
//
// It supports telebot.Silent option.
func (b *Bot) Pin(message Editable, options ...interface{}) error {
	messageID, chatID := message.MessageSig()

	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": messageID,
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	respJSON, err := b.Raw("pinChatMessage", params)
	if err != nil {
		return err
	}

	return extractOkResponse(respJSON)
}

// Use this method to unpin a message in a supergroup or a channel.
//
// It supports telebot.Silent option.
func (b *Bot) Unpin(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	respJSON, err := b.Raw("upinChatMessage", params)
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

	respJSON, err := b.Raw("getChat", params)
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

	if resp.Result.Type == ChatChannel && resp.Result.Username != "" {
		//Channel is Private
		resp.Result.Type = ChatChannelPrivate
	}

	return resp.Result, nil
}

// ProfilePhotosOf return list of profile pictures for a user.
func (b *Bot) ProfilePhotosOf(user *User) ([]Photo, error) {
	params := map[string]string{
		"user_id": user.Recipient(),
	}

	respJSON, err := b.Raw("getUserProfilePhotos", params)
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

	respJSON, err := b.Raw("getChatMember", params)
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
