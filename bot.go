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

	client := pref.Client
	if client == nil {
		client = http.DefaultClient
	}

	if pref.URL == "" {
		pref.URL = DefaultApiURL
	}

	bot := &Bot{
		Token:   pref.Token,
		URL:     pref.URL,
		Updates: make(chan Update, pref.Updates),
		Poller:  pref.Poller,

		handlers:    make(map[string]interface{}),
		synchronous: pref.Synchronous,
		stop:        make(chan struct{}),
		reporter:    pref.Reporter,
		client:      client,
	}

	if pref.offline {
		bot.Me = &User{}
	} else {
		user, err := bot.getMe()
		if err != nil {
			return nil, err
		}
		bot.Me = user
	}

	return bot, nil
}

// Bot represents a separate Telegram bot instance.
type Bot struct {
	Me      *User
	Token   string
	URL     string
	Updates chan Update
	Poller  Poller

	handlers    map[string]interface{}
	synchronous bool
	reporter    func(error)
	stop        chan struct{}
	client      *http.Client
}

// Settings represents a utility struct for passing certain
// properties of a bot around and is required to make bots.
type Settings struct {
	// Telegram API Url
	URL string

	// Telegram token
	Token string

	// Updates channel capacity
	Updates int // Default: 100

	// Poller is the provider of Updates.
	Poller Poller

	// Synchronous prevents handlers from running in parallel.
	// It makes ProcessUpdate return after the handler is finished.
	Synchronous bool

	// Reporter is a callback function that will get called
	// on any panics recovered from endpoint handlers.
	Reporter func(error)

	// HTTP Client used to make requests to telegram api
	Client *http.Client

	// offline allows to create a bot without network for testing purposes.
	offline bool
}

// Update object represents an incoming update.
type Update struct {
	ID int `json:"update_id"`

	Message            *Message            `json:"message,omitempty"`
	EditedMessage      *Message            `json:"edited_message,omitempty"`
	ChannelPost        *Message            `json:"channel_post,omitempty"`
	EditedChannelPost  *Message            `json:"edited_channel_post,omitempty"`
	Callback           *Callback           `json:"callback_query,omitempty"`
	Query              *Query              `json:"inline_query,omitempty"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query,omitempty"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`
	Poll               *Poll               `json:"poll,omitempty"`
	PollAnswer         *PollAnswer         `json:"poll_answer,omitempty"`
}

// Command represents a bot command.
type Command struct {
	// Text is a text of the command, 1-32 characters.
	// Can contain only lowercase English letters, digits and underscores.
	Text string `json:"command"`

	// Description of the command, 3-256 characters.
	Description string `json:"description"`
}

// Handle lets you set the handler for some command name or
// one of the supported endpoints.
//
// Example:
//
//     b.Handle("/help", func (m *tb.Message) {})
//     b.Handle(tb.OnText, func (m *tb.Message) {})
//     b.Handle(tb.OnQuery, func (q *tb.Query) {})
//
//     // make a hook for one of your preserved (by-pointer) inline buttons.
//     b.Handle(&inlineButton, func (c *tb.Callback) {})
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
	cmdRx   = regexp.MustCompile(`^(/\w+)(@(\w+))?(\s|$)(.+)?`)
	cbackRx = regexp.MustCompile(`^\f(\w+)(\|(.+))?$`)
)

// Start brings bot into motion by consuming incoming
// updates (see Bot.Updates channel).
func (b *Bot) Start() {
	if b.Poller == nil {
		panic("telebot: can't start without a poller")
	}

	stop := make(chan struct{})
	go b.Poller.Poll(b, b.Updates, stop)

	for {
		select {
		// handle incoming updates
		case upd := <-b.Updates:
			b.ProcessUpdate(upd)
		// call to stop polling
		case <-b.stop:
			close(stop)
			return
		}
	}
}

// Stop gracefully shuts the poller down.
func (b *Bot) Stop() {
	b.stop <- struct{}{}
}

// ProcessUpdate processes a single incoming update.
// A started bot calls this function automatically.
func (b *Bot) ProcessUpdate(upd Update) {
	if upd.Message != nil {
		m := upd.Message

		if m.PinnedMessage != nil {
			b.handle(OnPinned, m)
			return
		}

		// Commands
		if m.Text != "" {
			// Filtering malicious messages
			if m.Text[0] == '\a' {
				return
			}

			match := cmdRx.FindAllStringSubmatch(m.Text, -1)
			if match != nil {
				// Syntax: "</command>@<bot> <payload>"

				command, botName := match[0][1], match[0][3]
				if botName != "" && !strings.EqualFold(b.Me.Username, botName) {
					return
				}

				m.Payload = match[0][5]
				if b.handle(command, m) {
					return
				}
			}

			// 1:1 satisfaction
			if b.handle(m.Text, m) {
				return
			}

			b.handle(OnText, m)
			return
		}

		if b.handleMedia(m) {
			return
		}

		if m.Invoice != nil {
			b.handle(OnInvoice, m)
			return
		}

		if m.Payment != nil {
			b.handle(OnPayment, m)
			return
		}

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
				handler, ok := handler.(func(int64, int64))
				if !ok {
					panic("telebot: migration handler is bad")
				}

				b.runHandler(func() { handler(m.Chat.ID, m.MigrateTo) })
			}

			return
		}
	}

	if upd.EditedMessage != nil {
		b.handle(OnEdited, upd.EditedMessage)
		return
	}

	if upd.ChannelPost != nil {
		m := upd.ChannelPost

		if m.PinnedMessage != nil {
			b.handle(OnPinned, m)
			return
		}

		b.handle(OnChannelPost, upd.ChannelPost)
		return
	}

	if upd.EditedChannelPost != nil {
		b.handle(OnEditedChannelPost, upd.EditedChannelPost)
		return
	}

	if upd.Callback != nil {
		if upd.Callback.Data != "" {
			if upd.Callback.MessageID != "" {
				upd.Callback.Message = &Message{
					// InlineID indicates that message
					// is inline so we have only its id
					InlineID: upd.Callback.MessageID,
				}
			}

			data := upd.Callback.Data
			if data[0] == '\f' {
				match := cbackRx.FindAllStringSubmatch(data, -1)
				if match != nil {
					unique, payload := match[0][1], match[0][3]

					if handler, ok := b.handlers["\f"+unique]; ok {
						handler, ok := handler.(func(*Callback))
						if !ok {
							panic(fmt.Errorf("telebot: %s callback handler is bad", unique))
						}

						upd.Callback.Data = payload
						b.runHandler(func() { handler(upd.Callback) })

						return
					}
				}
			}
		}

		if handler, ok := b.handlers[OnCallback]; ok {
			handler, ok := handler.(func(*Callback))
			if !ok {
				panic("telebot: callback handler is bad")
			}

			b.runHandler(func() { handler(upd.Callback) })
		}

		return
	}

	if upd.Query != nil {
		if handler, ok := b.handlers[OnQuery]; ok {
			handler, ok := handler.(func(*Query))
			if !ok {
				panic("telebot: query handler is bad")
			}

			b.runHandler(func() { handler(upd.Query) })
		}

		return
	}

	if upd.ChosenInlineResult != nil {
		if handler, ok := b.handlers[OnChosenInlineResult]; ok {
			handler, ok := handler.(func(*ChosenInlineResult))
			if !ok {
				panic("telebot: chosen inline result handler is bad")
			}

			b.runHandler(func() { handler(upd.ChosenInlineResult) })
		}

		return
	}

	if upd.ShippingQuery != nil {
		if handler, ok := b.handlers[OnShipping]; ok {
			handler, ok := handler.(func(*ShippingQuery))
			if !ok {
				panic("telebot: shipping query handler is bad")
			}

			b.runHandler(func() { handler(upd.ShippingQuery) })
		}

		return
	}

	if upd.PreCheckoutQuery != nil {
		if handler, ok := b.handlers[OnCheckout]; ok {
			handler, ok := handler.(func(*PreCheckoutQuery))
			if !ok {
				panic("telebot: pre checkout query handler is bad")
			}

			b.runHandler(func() { handler(upd.PreCheckoutQuery) })
		}

		return
	}

	if upd.Poll != nil {
		if handler, ok := b.handlers[OnPoll]; ok {
			handler, ok := handler.(func(*Poll))
			if !ok {
				panic("telebot: poll handler is bad")
			}

			b.runHandler(func() { handler(upd.Poll) })
		}

		return
	}

	if upd.PollAnswer != nil {
		if handler, ok := b.handlers[OnPollAnswer]; ok {
			handler, ok := handler.(func(*PollAnswer))
			if !ok {
				panic("telebot: poll answer handler is bad")
			}

			b.runHandler(func() { handler(upd.PollAnswer) })
		}

		return
	}
}

func (b *Bot) handle(end string, m *Message) bool {
	if handler, ok := b.handlers[end]; ok {
		handler, ok := handler.(func(*Message))
		if !ok {
			panic(fmt.Errorf("telebot: %s handler is bad", end))
		}

		b.runHandler(func() { handler(m) })

		return true
	}

	return false
}

func (b *Bot) handleMedia(m *Message) bool {
	switch {
	case m.Photo != nil:
		b.handle(OnPhoto, m)
	case m.Voice != nil:
		b.handle(OnVoice, m)
	case m.Audio != nil:
		b.handle(OnAudio, m)
	case m.Animation != nil:
		b.handle(OnAnimation, m)
	case m.Document != nil:
		b.handle(OnDocument, m)
	case m.Sticker != nil:
		b.handle(OnSticker, m)
	case m.Video != nil:
		b.handle(OnVideo, m)
	case m.VideoNote != nil:
		b.handle(OnVideoNote, m)
	case m.Contact != nil:
		b.handle(OnContact, m)
	case m.Location != nil:
		b.handle(OnLocation, m)
	case m.Venue != nil:
		b.handle(OnVenue, m)
	case m.Dice != nil:
		b.handle(OnDice, m)
	default:
		return false
	}
	return true
}

// Send accepts 2+ arguments, starting with destination chat, followed by
// some Sendable (or string!) and optional send options.
//
// Note: since most arguments are of type interface{}, but have pointer
// 		method receivers, make sure to pass them by-pointer, NOT by-value.
//
// What is a send option exactly? It can be one of the following types:
//
//     - *SendOptions (the actual object accepted by Telegram API)
//     - *ReplyMarkup (a component of SendOptions)
//     - Option (a shortcut flag for popular options)
//     - ParseMode (HTML, Markdown, etc)
//
func (b *Bot) Send(to Recipient, what interface{}, options ...interface{}) (*Message, error) {
	if to == nil {
		return nil, ErrBadRecipient
	}

	sendOpts := extractOptions(options)

	switch object := what.(type) {
	case string:
		return b.sendText(to, object, sendOpts)
	case Sendable:
		return object.Send(b, to, sendOpts)
	default:
		return nil, ErrUnsupportedWhat
	}
}

// SendAlbum sends multiple instances of media as a single message.
//
// From all existing options, it only supports tb.Silent option.
func (b *Bot) SendAlbum(to Recipient, a Album, options ...interface{}) ([]Message, error) {
	if to == nil {
		return nil, ErrBadRecipient
	}

	media := make([]string, len(a))
	files := make(map[string]File)

	for i, x := range a {
		var (
			repr string
			data []byte
			f    = x.MediaFile()
		)

		switch {
		case f.InCloud():
			repr = f.FileID
		case f.FileURL != "":
			repr = f.FileURL
		case f.OnDisk() || f.FileReader != nil:
			repr = "attach://" + strconv.Itoa(i)
			files[strconv.Itoa(i)] = *f
		default:
			return nil, errors.Errorf("telebot: album entry #%d does not exist", i)
		}

		switch y := x.(type) {
		case *Photo:
			data, _ = json.Marshal(struct {
				Type      string    `json:"type"`
				Media     string    `json:"media"`
				Caption   string    `json:"caption,omitempty"`
				ParseMode ParseMode `json:"parse_mode,omitempty"`
			}{
				Type:      "photo",
				Media:     repr,
				Caption:   y.Caption,
				ParseMode: y.ParseMode,
			})
		case *Video:
			data, _ = json.Marshal(struct {
				Type              string `json:"type"`
				Caption           string `json:"caption"`
				Media             string `json:"media"`
				Width             int    `json:"width,omitempty"`
				Height            int    `json:"height,omitempty"`
				Duration          int    `json:"duration,omitempty"`
				SupportsStreaming bool   `json:"supports_streaming,omitempty"`
			}{
				Type:              "video",
				Caption:           y.Caption,
				Media:             repr,
				Width:             y.Width,
				Height:            y.Height,
				Duration:          y.Duration,
				SupportsStreaming: y.SupportsStreaming,
			})
		default:
			return nil, errors.Errorf("telebot: album entry #%d is not valid", i)
		}

		media[i] = string(data)
	}

	params := map[string]string{
		"chat_id": to.Recipient(),
		"media":   "[" + strings.Join(media, ",") + "]",
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	data, err := b.sendFiles("sendMediaGroup", files, params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result []Message
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}

	for attachName := range files {
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
//
// This function will panic upon nil Message.
func (b *Bot) Reply(to *Message, what interface{}, options ...interface{}) (*Message, error) {
	sendOpts := extractOptions(options)
	if sendOpts == nil {
		sendOpts = &SendOptions{}
	}

	sendOpts.ReplyTo = to
	return b.Send(to.Chat, what, sendOpts)
}

// Forward behaves just like Send() but of all options it only supports Silent (see Bots API).
//
// This function will panic upon nil Editable.
func (b *Bot) Forward(to Recipient, msg Editable, options ...interface{}) (*Message, error) {
	if to == nil {
		return nil, ErrBadRecipient
	}
	msgID, chatID := msg.MessageSig()

	params := map[string]string{
		"chat_id":      to.Recipient(),
		"from_chat_id": strconv.FormatInt(chatID, 10),
		"message_id":   msgID,
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	data, err := b.Raw("forwardMessage", params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
}

// Edit is magic, it lets you change already sent message.
//
// Use cases:
//
//     b.Edit(msg, msg.Text, newMarkup)
//     b.Edit(msg, "new <b>text</b>", tb.ModeHTML)
//     b.Edit(msg, tb.Location{42.1337, 69.4242})
//
// This function will panic upon nil Editable.
func (b *Bot) Edit(msg Editable, what interface{}, options ...interface{}) (*Message, error) {
	var (
		method string
		params = make(map[string]string)
	)

	switch v := what.(type) {
	case string:
		method = "editMessageText"
		params["text"] = v
	case Location:
		method = "editMessageLiveLocation"
		params["latitude"] = fmt.Sprintf("%f", v.Lat)
		params["longitude"] = fmt.Sprintf("%f", v.Lng)
	default:
		return nil, ErrUnsupportedWhat
	}

	msgID, chatID := msg.MessageSig()

	if chatID == 0 { // if inline message
		params["inline_message_id"] = msgID
	} else {
		params["chat_id"] = strconv.FormatInt(chatID, 10)
		params["message_id"] = msgID
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	data, err := b.Raw(method, params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
}

// EditReplyMarkup edits reply markup of already sent message.
// Pass nil or empty ReplyMarkup to delete it from the message.
//
// On success, returns edited message object.
// This function will panic upon nil Editable.
func (b *Bot) EditReplyMarkup(msg Editable, markup *ReplyMarkup) (*Message, error) {
	msgID, chatID := msg.MessageSig()
	params := make(map[string]string)

	if chatID == 0 { // if inline message
		params["inline_message_id"] = msgID
	} else {
		params["chat_id"] = strconv.FormatInt(chatID, 10)
		params["message_id"] = msgID
	}

	if markup == nil {
		// will delete reply markup
		markup = &ReplyMarkup{}
	}

	processButtons(markup.InlineKeyboard)
	data, _ := json.Marshal(markup)
	params["reply_markup"] = string(data)

	data, err := b.Raw("editMessageReplyMarkup", params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
}

// EditCaption edits already sent photo caption with known recipient and message id.
//
// On success, returns edited message object.
// This function will panic upon nil Editable.
func (b *Bot) EditCaption(msg Editable, caption string, options ...interface{}) (*Message, error) {
	msgID, chatID := msg.MessageSig()

	params := map[string]string{
		"caption": caption,
	}

	if chatID == 0 { // if inline message
		params["inline_message_id"] = msgID
	} else {
		params["chat_id"] = strconv.FormatInt(chatID, 10)
		params["message_id"] = msgID
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	data, err := b.Raw("editMessageCaption", params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
}

// EditMedia edits already sent media with known recipient and message id.
//
// Use cases:
//
//     bot.EditMedia(msg, &tb.Photo{File: tb.FromDisk("chicken.jpg")})
//     bot.EditMedia(msg, &tb.Video{File: tb.FromURL("http://video.mp4")})
//
// This function will panic upon nil Editable.
func (b *Bot) EditMedia(msg Editable, media InputMedia, options ...interface{}) (*Message, error) {
	var (
		repr  string
		thumb *Photo

		thumbName = "thumb"
		file      = media.MediaFile()
		files     = make(map[string]File)
	)

	switch {
	case file.InCloud():
		repr = file.FileID
	case file.FileURL != "":
		repr = file.FileURL
	case file.OnDisk() || file.FileReader != nil:
		s := file.FileLocal
		if file.FileReader != nil {
			s = "0"
		} else if s == thumbName {
			thumbName = "thumb2"
		}

		repr = "attach://" + s
		files[s] = *file
	default:
		return nil, errors.Errorf("telebot: can't edit media, it does not exist")
	}

	type FileJSON struct {
		// All types.
		Type      string    `json:"type"`
		Caption   string    `json:"caption"`
		Media     string    `json:"media"`
		ParseMode ParseMode `json:"parse_mode,omitempty"`

		// Video.
		Width             int  `json:"width,omitempty"`
		Height            int  `json:"height,omitempty"`
		SupportsStreaming bool `json:"supports_streaming,omitempty"`

		// Video and audio.
		Duration int `json:"duration,omitempty"`

		// Document.
		FileName string `json:"file_name"`

		// Document, video and audio.
		Thumbnail string `json:"thumb,omitempty"`
		MIME      string `json:"mime_type,omitempty"`

		// Audio.
		Title     string `json:"title,omitempty"`
		Performer string `json:"performer,omitempty"`
	}

	result := &FileJSON{Media: repr}

	switch m := media.(type) {
	case *Photo:
		result.Type = "photo"
		result.Caption = m.Caption
	case *Video:
		result.Type = "video"
		result.Caption = m.Caption
		result.Width = m.Width
		result.Height = m.Height
		result.Duration = m.Duration
		result.SupportsStreaming = m.SupportsStreaming
		result.MIME = m.MIME
		thumb = m.Thumbnail
	case *Document:
		result.Type = "document"
		result.Caption = m.Caption
		result.FileName = m.FileName
		result.MIME = m.MIME
		thumb = m.Thumbnail
	case *Audio:
		result.Type = "audio"
		result.Caption = m.Caption
		result.Duration = m.Duration
		result.MIME = m.MIME
		result.Title = m.Title
		result.Performer = m.Performer
		thumb = m.Thumbnail
	default:
		return nil, errors.Errorf("telebot: media entry is not valid")
	}

	msgID, chatID := msg.MessageSig()
	params := make(map[string]string)

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	if sendOpts != nil {
		result.ParseMode = sendOpts.ParseMode
	}
	if thumb != nil {
		result.Thumbnail = "attach://" + thumbName
		files[thumbName] = *thumb.MediaFile()
	}

	data, _ := json.Marshal(result)
	params["media"] = string(data)

	if chatID == 0 { // If inline message.
		params["inline_message_id"] = msgID
	} else {
		params["chat_id"] = strconv.FormatInt(chatID, 10)
		params["message_id"] = msgID
	}

	data, err := b.sendFiles("editMessageMedia", files, params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
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
// This function will panic upon nil Editable.
func (b *Bot) Delete(msg Editable) error {
	msgID, chatID := msg.MessageSig()

	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": msgID,
	}

	_, err := b.Raw("deleteMessage", params)
	return err
}

// Notify updates the chat action for recipient.
//
// Chat action is a status message that recipient would see where
// you typically see "Harry is typing" status message. The only
// difference is that bots' chat actions live only for 5 seconds
// and die just once the client receives a message from the bot.
//
// Currently, Telegram supports only a narrow range of possible
// actions, these are aligned as constants of this package.
func (b *Bot) Notify(to Recipient, action ChatAction) error {
	if to == nil {
		return ErrBadRecipient
	}

	params := map[string]string{
		"chat_id": to.Recipient(),
		"action":  string(action),
	}

	_, err := b.Raw("sendChatAction", params)
	return err
}

// Ship replies to the shipping query, if you sent an invoice
// requesting an address and the parameter is_flexible was specified.
//
// Usage:
//
//		b.Ship(query)          // OK
//		b.Ship(query, opts...) // OK with options
//		b.Ship(query, "Oops!") // Error message
//
func (b *Bot) Ship(query *ShippingQuery, what ...interface{}) error {
	params := map[string]string{
		"shipping_query_id": query.ID,
	}

	if len(what) == 0 {
		params["ok"] = "True"
	} else if s, ok := what[0].(string); ok {
		params["ok"] = "False"
		params["error_message"] = s
	} else {
		var opts []ShippingOption
		for _, v := range what {
			opt, ok := v.(ShippingOption)
			if !ok {
				return ErrUnsupportedWhat
			}
			opts = append(opts, opt)
		}

		params["ok"] = "True"
		data, _ := json.Marshal(opts)
		params["shipping_options"] = string(data)
	}

	_, err := b.Raw("answerShippingQuery", params)
	return err
}

// Accept finalizes the deal.
func (b *Bot) Accept(query *PreCheckoutQuery, errorMessage ...string) error {
	params := map[string]string{
		"pre_checkout_query_id": query.ID,
	}

	if len(errorMessage) == 0 {
		params["ok"] = "True"
	} else {
		params["ok"] = "False"
		params["error_message"] = errorMessage[0]
	}

	_, err := b.Raw("answerPreCheckoutQuery", params)
	return err
}

// Answer sends a response for a given inline query. A query can only
// be responded to once, subsequent attempts to respond to the same query
// will result in an error.
func (b *Bot) Answer(query *Query, resp *QueryResponse) error {
	resp.QueryID = query.ID

	for _, result := range resp.Results {
		result.Process()
	}

	_, err := b.Raw("answerInlineQuery", resp)
	return err
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
func (b *Bot) Respond(c *Callback, response ...*CallbackResponse) error {
	var resp *CallbackResponse
	if response == nil {
		resp = &CallbackResponse{}
	} else {
		resp = response[0]
	}

	resp.CallbackID = c.ID
	_, err := b.Raw("answerCallbackQuery", resp)
	return err
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

	data, err := b.Raw("getFile", params)
	if err != nil {
		return File{}, err
	}

	var resp struct {
		Result File
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return File{}, wrapError(err)
	}
	return resp.Result, nil
}

// Download saves the file from Telegram servers locally.
//
// Maximum file size to download is 20 MB.
func (b *Bot) Download(file *File, localFilename string) error {
	reader, err := b.GetFile(file)
	if err != nil {
		return wrapError(err)
	}
	defer reader.Close()

	out, err := os.Create(localFilename)
	if err != nil {
		return wrapError(err)
	}
	defer out.Close()

	_, err = io.Copy(out, reader)
	if err != nil {
		return wrapError(err)
	}

	file.FileLocal = localFilename
	return nil
}

// GetFile gets a file from Telegram servers.
func (b *Bot) GetFile(file *File) (io.ReadCloser, error) {
	f, err := b.FileByID(file.FileID)
	if err != nil {
		return nil, err
	}

	url := b.URL + "/file/bot" + b.Token + "/" + f.FilePath
	file.FilePath = f.FilePath // saving file path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, wrapError(err)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, wrapError(err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errors.Errorf("telebot: expected status 200 but got %s", resp.Status)
	}

	return resp.Body, nil
}

// StopLiveLocation stops broadcasting live message location
// before Location.LivePeriod expires.
//
// It supports tb.ReplyMarkup.
// This function will panic upon nil Editable.
func (b *Bot) StopLiveLocation(msg Editable, options ...interface{}) (*Message, error) {
	msgID, chatID := msg.MessageSig()

	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": msgID,
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	data, err := b.Raw("stopMessageLiveLocation", params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
}

// StopPoll stops a poll which was sent by the bot and returns
// the stopped Poll object with the final results.
//
// It supports ReplyMarkup.
// This function will panic upon nil Editable.
func (b *Bot) StopPoll(msg Editable, options ...interface{}) (*Poll, error) {
	msgID, chatID := msg.MessageSig()

	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": msgID,
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	data, err := b.Raw("stopPoll", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *Poll
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// GetInviteLink should be used to export chat's invite link.
func (b *Bot) GetInviteLink(chat *Chat) (string, error) {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	data, err := b.Raw("exportChatInviteLink", params)
	if err != nil {
		return "", err
	}

	var resp struct {
		Result string
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", wrapError(err)
	}
	return resp.Result, nil
}

// SetGroupTitle should be used to update group title.
func (b *Bot) SetGroupTitle(chat *Chat, title string) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"title":   title,
	}

	_, err := b.Raw("setChatTitle", params)
	return err
}

// SetGroupDescription should be used to update group title.
func (b *Bot) SetGroupDescription(chat *Chat, description string) error {
	params := map[string]string{
		"chat_id":     chat.Recipient(),
		"description": description,
	}

	_, err := b.Raw("setChatDescription", params)
	return err
}

// SetGroupPhoto should be used to update group photo.
func (b *Bot) SetGroupPhoto(chat *Chat, p *Photo) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.sendFiles("setChatPhoto", map[string]File{"photo": p.File}, params)
	return err
}

// SetGroupStickerSet should be used to update group's group sticker set.
func (b *Bot) SetGroupStickerSet(chat *Chat, setName string) error {
	params := map[string]string{
		"chat_id":          chat.Recipient(),
		"sticker_set_name": setName,
	}

	_, err := b.Raw("setChatStickerSet", params)
	return err
}

// SetGroupPermissions sets default chat permissions for all members.
func (b *Bot) SetGroupPermissions(chat *Chat, perms Rights) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}
	embedRights(params, perms)

	_, err := b.Raw("setChatPermissions", params)
	return err
}

// DeleteGroupPhoto should be used to just remove group photo.
func (b *Bot) DeleteGroupPhoto(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("deleteChatPhoto", params)
	return err
}

// DeleteGroupStickerSet should be used to just remove group sticker set.
func (b *Bot) DeleteGroupStickerSet(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("deleteChatStickerSet", params)
	return err
}

// Leave makes bot leave a group, supergroup or channel.
func (b *Bot) Leave(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("leaveChat", params)
	return err
}

// Pin pins a message in a supergroup or a channel.
//
// It supports tb.Silent option.
// This function will panic upon nil Editable.
func (b *Bot) Pin(msg Editable, options ...interface{}) error {
	msgID, chatID := msg.MessageSig()

	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": msgID,
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	_, err := b.Raw("pinChatMessage", params)
	return err
}

// Unpin unpins a message in a supergroup or a channel.
//
// It supports tb.Silent option.
func (b *Bot) Unpin(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("unpinChatMessage", params)
	return err
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

	data, err := b.Raw("getChat", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *Chat
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	if resp.Result.Type == ChatChannel && resp.Result.Username == "" {
		resp.Result.Type = ChatChannelPrivate
	}
	return resp.Result, nil
}

// ProfilePhotosOf returns list of profile pictures for a user.
func (b *Bot) ProfilePhotosOf(user *User) ([]Photo, error) {
	params := map[string]string{
		"user_id": user.Recipient(),
	}

	data, err := b.Raw("getUserProfilePhotos", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result struct {
			Count  int     `json:"total_count"`
			Photos []Photo `json:"photos"`
		}
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result.Photos, nil
}

// ChatMemberOf returns information about a member of a chat.
//
// Returns a ChatMember object on success.
func (b *Bot) ChatMemberOf(chat *Chat, user *User) (*ChatMember, error) {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"user_id": user.Recipient(),
	}

	data, err := b.Raw("getChatMember", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *ChatMember
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// FileURLByID returns direct url for files using FileId which you can get from File object
func (b *Bot) FileURLByID(fileID string) (string, error) {
	f, err := b.FileByID(fileID)
	if err != nil {
		return "", err
	}

	return b.URL + "/file/bot" + b.Token + "/" + f.FilePath, nil
}

// GetCommands returns the current list of the bot's commands.
func (b *Bot) GetCommands() ([]Command, error) {
	data, err := b.Raw("getMyCommands", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result []Command
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// SetCommands changes the list of the bot's commands.
func (b *Bot) SetCommands(cmds []Command) error {
	data, _ := json.Marshal(cmds)

	params := map[string]string{
		"commands": string(data),
	}

	_, err := b.Raw("setMyCommands", params)
	return err
}

func (b *Bot) NewMarkup() *ReplyMarkup {
	return &ReplyMarkup{}
}
