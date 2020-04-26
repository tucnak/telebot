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

		handlers: make(map[string]interface{}),
		stop:     make(chan struct{}),
		reporter: pref.Reporter,
		client:   client,
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
	URL     string
	Updates chan Update
	Poller  Poller

	handlers map[string]interface{}
	reporter func(error)
	stop     chan struct{}
	client   *http.Client
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

	// Reporter is a callback function that will get called
	// on any panics recovered from endpoint handlers.
	Reporter func(error)

	// HTTP Client used to make requests to telegram api
	Client *http.Client
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
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`
	Poll               *Poll               `json:"poll,omitempty"`
	PollAnswer         *PollAnswer         `json:"poll_answer,omitempty"`
}

// Command represents a bot command.
type Command struct {
	// Text is a aext of the command, 1-32 characters.
	// Can contain only lowercase English letters, digits and underscores.
	Text string `json:"command"`

	// Description of the command, 3-256 characters.
	Description string `json:"description"`
}

// ChosenInlineResult represents a result of an inline query that was chosen
// by the user and sent to their chat partner.
type ChosenInlineResult struct {
	From      User      `json:"from"`
	Location  *Location `json:"location,omitempty"`
	ResultID  string    `json:"result_id"`
	Query     string    `json:"query"`
	MessageID string    `json:"inline_message_id"` // inline messages only!
}

type PreCheckoutQuery struct {
	Sender   *User  `json:"from"`
	ID       string `json:"id"`
	Currency string `json:"currency"`
	Payload  string `json:"invoice_payload"`
	Total    int    `json:"total_amount"`
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
			b.incomingUpdate(&upd)

		// call to stop polling
		case <-b.stop:
			stop <- struct{}{}

		// polling has stopped
		case <-stop:
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

				go func(b *Bot, handler func(int64, int64), from, to int64) {
					if b.reporter == nil {
						defer b.deferDebug()
					}
					handler(from, to)
				}(b, handler, m.Chat.ID, m.MigrateTo)
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
						go func(b *Bot, handler func(*Callback), c *Callback) {
							if b.reporter == nil {
								defer b.deferDebug()
							}
							handler(c)
						}(b, handler, upd.Callback)

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

			go func(b *Bot, handler func(*Callback), c *Callback) {
				if b.reporter == nil {
					defer b.deferDebug()
				}
				handler(c)
			}(b, handler, upd.Callback)
		}

		return
	}

	if upd.Query != nil {
		if handler, ok := b.handlers[OnQuery]; ok {
			handler, ok := handler.(func(*Query))
			if !ok {
				panic("telebot: query handler is bad")
			}

			go func(b *Bot, handler func(*Query), q *Query) {
				if b.reporter == nil {
					defer b.deferDebug()
				}
				handler(q)
			}(b, handler, upd.Query)
		}

		return
	}

	if upd.ChosenInlineResult != nil {
		if handler, ok := b.handlers[OnChosenInlineResult]; ok {
			handler, ok := handler.(func(*ChosenInlineResult))
			if !ok {
				panic("telebot: chosen inline result handler is bad")
			}

			go func(b *Bot, handler func(*ChosenInlineResult), r *ChosenInlineResult) {
				if b.reporter == nil {
					defer b.deferDebug()
				}
				handler(r)
			}(b, handler, upd.ChosenInlineResult)
		}

		return
	}

	if upd.PreCheckoutQuery != nil {
		if handler, ok := b.handlers[OnCheckout]; ok {
			handler, ok := handler.(func(*PreCheckoutQuery))
			if !ok {
				panic("telebot: pre checkout query handler is bad")
			}

			go func(b *Bot, handler func(*PreCheckoutQuery), pre *PreCheckoutQuery) {
				if b.reporter == nil {
					defer b.deferDebug()
				}
				handler(pre)
			}(b, handler, upd.PreCheckoutQuery)
		}

		return
	}

	if upd.Poll != nil {
		if handler, ok := b.handlers[OnPoll]; ok {
			handler, ok := handler.(func(*Poll))
			if !ok {
				panic("telebot: poll handler is bad")
			}

			go func(b *Bot, handler func(*Poll), p *Poll) {
				if b.reporter == nil {
					defer b.deferDebug()
				}
				handler(p)
			}(b, handler, upd.Poll)
		}

		return
	}

	if upd.PollAnswer != nil {
		if handler, ok := b.handlers[OnPollAnswer]; ok {
			handler, ok := handler.(func(*PollAnswer))
			if !ok {
				panic("telebot: poll answer handler is bad")
			}

			go func(b *Bot, handler func(*PollAnswer), pa *PollAnswer) {
				if b.reporter == nil {
					defer b.deferDebug()
				}
				handler(pa)
			}(b, handler, upd.PollAnswer)
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

		go func(b *Bot, handler func(*Message), m *Message) {
			if b.reporter == nil {
				defer b.deferDebug()
			}
			handler(m)
		}(b, handler, m)

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

// Stop gracefully shuts the poller down.
func (b *Bot) Stop() {
	b.stop <- struct{}{}
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
//     - Option (a shorcut flag for popular options)
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

// SendAlbum is used when sending multiple instances of media as a single
// message (so-called album).
//
// From all existing options, it only supports telebot.Silent.
func (b *Bot) SendAlbum(to Recipient, a Album, options ...interface{}) ([]Message, error) {
	if to == nil {
		return nil, ErrBadRecipient
	}

	media := make([]string, len(a))
	files := make(map[string]File)

	for i, x := range a {
		var mediaRepr string
		var jsonRepr []byte

		f := x.MediaFile()

		switch {
		case f.InCloud():
			mediaRepr = f.FileID
		case f.FileURL != "":
			mediaRepr = f.FileURL
		case f.OnDisk() || f.FileReader != nil:
			mediaRepr = "attach://" + strconv.Itoa(i)
			files[strconv.Itoa(i)] = *f
		default:
			return nil, errors.Errorf(
				"telebot: album entry #%d doesn't exist anywhere", i)
		}

		switch y := x.(type) {
		case *Photo:
			jsonRepr, _ = json.Marshal(struct {
				Type      string    `json:"type"`
				Media     string    `json:"media"`
				Caption   string    `json:"caption,omitempty"`
				ParseMode ParseMode `json:"parse_mode,omitempty"`
			}{
				"photo",
				mediaRepr,
				y.Caption,
				y.ParseMode,
			})
		case *Video:
			jsonRepr, _ = json.Marshal(struct {
				Type              string `json:"type"`
				Caption           string `json:"caption"`
				Media             string `json:"media"`
				Width             int    `json:"width,omitempty"`
				Height            int    `json:"height,omitempty"`
				Duration          int    `json:"duration,omitempty"`
				SupportsStreaming bool   `json:"supports_streaming,omitempty"`
			}{
				"video",
				y.Caption,
				mediaRepr,
				y.Width,
				y.Height,
				y.Duration,
				y.SupportsStreaming,
			})
		default:
			return nil, errors.Errorf("telebot: album entry #%d is not valid", i)
		}

		media[i] = string(jsonRepr)
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
		Ok          bool
		Result      []Message
		Description string
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return nil, errors.Errorf("api error: %s", resp.Description)
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
// This function will panic upon nil Message.
func (b *Bot) Forward(to Recipient, what *Message, options ...interface{}) (*Message, error) {
	if to == nil {
		return nil, ErrBadRecipient
	}

	params := map[string]string{
		"chat_id":      to.Recipient(),
		"from_chat_id": what.Chat.Recipient(),
		"message_id":   strconv.Itoa(what.ID),
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
		params = map[string]string{}
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

// EditReplyMarkup used to edit reply markup of already sent message.
// Pass nil or empty ReplyMarkup to delete it from the message.
//
// On success, returns edited message object.
// This function will panic upon nil Editable.
func (b *Bot) EditReplyMarkup(message Editable, markup *ReplyMarkup) (*Message, error) {
	messageID, chatID := message.MessageSig()
	params := map[string]string{}

	if chatID == 0 { // if inline message
		params["inline_message_id"] = messageID
	} else {
		params["chat_id"] = strconv.FormatInt(chatID, 10)
		params["message_id"] = messageID
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

// EditCaption used to edit already sent photo caption with known recipient and message id.
//
// On success, returns edited message object
func (b *Bot) EditCaption(message Editable, caption string, options ...interface{}) (*Message, error) {
	messageID, chatID := message.MessageSig()

	params := map[string]string{"caption": caption}

	// if inline message
	if chatID == 0 {
		params["inline_message_id"] = messageID
	} else {
		params["chat_id"] = strconv.FormatInt(chatID, 10)
		params["message_id"] = messageID
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	data, err := b.Raw("editMessageCaption", params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
}

// EditMedia used to edit already sent media with known recipient and message id.
//
// Use cases:
//
//     bot.EditMedia(msg, &tb.Photo{File: tb.FromDisk("chicken.jpg")});
//     bot.EditMedia(msg, &tb.Video{File: tb.FromURL("http://video.mp4")});
//
func (b *Bot) EditMedia(message Editable, inputMedia InputMedia, options ...interface{}) (*Message, error) {
	var mediaRepr string
	var jsonRepr []byte
	var thumb *Photo

	file := make(map[string]File)

	f := inputMedia.MediaFile()
	thumbAttachName := "thumb"

	switch {
	case f.InCloud():
		mediaRepr = f.FileID
	case f.FileURL != "":
		mediaRepr = f.FileURL
	case f.OnDisk() || f.FileReader != nil:
		s := f.FileLocal
		if f.FileReader != nil {
			s = "0"
		}
		if s == thumbAttachName {
			thumbAttachName = "thumb2"
		}
		mediaRepr = "attach://" + s
		file[s] = *f
	default:
		return nil, errors.Errorf(
			"telebot: can't edit media, it doesn't exist anywhere")
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

	resultMedia := &FileJSON{Media: mediaRepr}

	sendOpts := extractOptions(options)
	if sendOpts != nil {
		resultMedia.ParseMode = sendOpts.ParseMode
	}

	switch y := inputMedia.(type) {
	case *Photo:
		resultMedia.Type = "photo"
		resultMedia.Caption = y.Caption
	case *Video:
		resultMedia.Type = "video"
		resultMedia.Caption = y.Caption
		resultMedia.Width = y.Width
		resultMedia.Height = y.Height
		resultMedia.Duration = y.Duration
		resultMedia.SupportsStreaming = y.SupportsStreaming
		resultMedia.MIME = y.MIME
		thumb = y.Thumbnail
		if thumb != nil {
			resultMedia.Thumbnail = "attach://" + thumbAttachName
		}
	case *Document:
		resultMedia.Type = "document"
		resultMedia.Caption = y.Caption
		resultMedia.FileName = y.FileName
		resultMedia.MIME = y.MIME
		thumb = y.Thumbnail
		if thumb != nil {
			resultMedia.Thumbnail = "attach://" + thumbAttachName
		}
	case *Audio:
		resultMedia.Type = "audio"
		resultMedia.Caption = y.Caption
		resultMedia.Duration = y.Duration
		resultMedia.MIME = y.MIME
		resultMedia.Title = y.Title
		resultMedia.Performer = y.Performer
		thumb = y.Thumbnail
		if thumb != nil {
			resultMedia.Thumbnail = "attach://" + thumbAttachName
		}
	default:
		return nil, errors.Errorf("telebot: inputMedia entry is not valid")
	}

	messageID, chatID := message.MessageSig()

	jsonRepr, _ = json.Marshal(resultMedia)
	params := map[string]string{}
	params["media"] = string(jsonRepr)

	// If inline message.
	if chatID == 0 {
		params["inline_message_id"] = messageID
	} else {
		params["chat_id"] = strconv.FormatInt(chatID, 10)
		params["message_id"] = messageID
	}

	if thumb != nil {
		file[thumbAttachName] = *thumb.MediaFile()
	}

	embedSendOptions(params, sendOpts)

	data, err := b.sendFiles("editMessageMedia", file, params)
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
func (b *Bot) Delete(message Editable) error {
	messageID, chatID := message.MessageSig()

	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": messageID,
	}

	data, err := b.Raw("deleteMessage", params)
	if err != nil {
		return err
	}

	return extractOk(data)
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

	data, err := b.Raw("sendChatAction", params)
	if err != nil {
		return err
	}

	return extractOk(data)
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

	data, err := b.Raw("answerPreCheckoutQuery", params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// Answer sends a response for a given inline query. A query can only
// be responded to once, subsequent attempts to respond to the same query
// will result in an error.
func (b *Bot) Answer(query *Query, response *QueryResponse) error {
	response.QueryID = query.ID

	for _, result := range response.Results {
		result.Process()
	}

	data, err := b.Raw("answerInlineQuery", response)
	if err != nil {
		return err
	}

	return extractOk(data)
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
	data, err := b.Raw("answerCallbackQuery", response)
	if err != nil {
		return err
	}

	return extractOk(data)
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
		Ok          bool
		Description string
		Result      File
	}

	err = json.Unmarshal(data, &resp)
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

// GetFile from Telegram servers
func (b *Bot) GetFile(file *File) (io.ReadCloser, error) {
	f, err := b.FileByID(file.FileID)
	if err != nil {
		return nil, err
	}
	// save FilePath
	file.FilePath = f.FilePath

	req, err := http.NewRequest("GET", b.URL+"/file/bot"+b.Token+"/"+f.FilePath, nil)
	if err != nil {
		return nil, wrapError(err)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "file http.GET failed")
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errors.Errorf("api error: expected 200 OK but got %s", resp.Status)
	}

	return resp.Body, nil
}

// StopLiveLocation should be called to stop broadcasting live message location
// before Location.LivePeriod expires.
//
// It supports telebot.ReplyMarkup.
func (b *Bot) StopLiveLocation(message Editable, options ...interface{}) (*Message, error) {
	messageID, chatID := message.MessageSig()

	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": messageID,
	}

	sendOpts := extractOptions(options)
	embedSendOptions(params, sendOpts)

	data, err := b.Raw("stopMessageLiveLocation", params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
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
		Ok          bool
		Description string
		Result      string
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return "", errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return "", errors.Errorf("api error: %s", resp.Description)
	}

	return resp.Result, nil
}

// SetGroupTitle should be used to update group title.
func (b *Bot) SetGroupTitle(chat *Chat, newTitle string) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"title":   newTitle,
	}

	data, err := b.Raw("setChatTitle", params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// SetGroupDescription should be used to update group title.
func (b *Bot) SetGroupDescription(chat *Chat, description string) error {
	params := map[string]string{
		"chat_id":     chat.Recipient(),
		"description": description,
	}

	data, err := b.Raw("setChatDescription", params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// SetGroupPhoto should be used to update group photo.
func (b *Bot) SetGroupPhoto(chat *Chat, p *Photo) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	data, err := b.sendFiles("setChatPhoto", map[string]File{"photo": p.File}, params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// SetGroupStickerSet should be used to update group's group sticker set.
func (b *Bot) SetGroupStickerSet(chat *Chat, setName string) error {
	params := map[string]string{
		"chat_id":          chat.Recipient(),
		"sticker_set_name": setName,
	}

	data, err := b.Raw("setChatStickerSet", params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// SetGroupPermissions sets default chat permissions for all members.
func (b *Bot) SetGroupPermissions(chat *Chat, perms Rights) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	embedRights(params, perms)

	data, err := b.Raw("setChatPermissions", params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// DeleteGroupPhoto should be used to just remove group photo.
func (b *Bot) DeleteGroupPhoto(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	data, err := b.Raw("deleteChatPhoto", params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// DeleteGroupStickerSet should be used to just remove group sticker set.
func (b *Bot) DeleteGroupStickerSet(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	data, err := b.Raw("deleteChatStickerSet", params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// Leave makes bot leave a group, supergroup or channel.
func (b *Bot) Leave(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	data, err := b.Raw("leaveChat", params)
	if err != nil {
		return err
	}

	return extractOk(data)
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

	data, err := b.Raw("pinChatMessage", params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// Use this method to unpin a message in a supergroup or a channel.
//
// It supports telebot.Silent option.
func (b *Bot) Unpin(chat *Chat) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
	}

	data, err := b.Raw("unpinChatMessage", params)
	if err != nil {
		return err
	}

	return extractOk(data)
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
		Ok          bool
		Description string
		Result      *Chat
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return nil, errors.Errorf("api error: %s", resp.Description)
	}

	if resp.Result.Type == ChatChannel && resp.Result.Username == "" {
		// Channel is Private
		resp.Result.Type = ChatChannelPrivate
	}

	return resp.Result, nil
}

// ProfilePhotosOf return list of profile pictures for a user.
func (b *Bot) ProfilePhotosOf(user *User) ([]Photo, error) {
	params := map[string]string{
		"user_id": user.Recipient(),
	}

	data, err := b.Raw("getUserProfilePhotos", params)
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

	err = json.Unmarshal(data, &resp)
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

	data, err := b.Raw("getChatMember", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Ok          bool
		Result      *ChatMember
		Description string `json:"description"`
	}

	err = json.Unmarshal(data, &resp)
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
	return b.URL + "/file/bot" + b.Token + "/" + f.FilePath, nil
}

// UploadStickerFile returns uploaded File on success.
func (b *Bot) UploadStickerFile(userID int, pngSticker *File) (*File, error) {
	files := map[string]File{
		"png_sticker": *pngSticker,
	}
	params := map[string]string{
		"user_id": strconv.Itoa(userID),
	}

	data, err := b.sendFiles("uploadStickerFile", files, params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Ok          bool
		Result      File
		Description string
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	if !resp.Ok {
		return nil, errors.Errorf("api error: %s", resp.Description)
	}

	return &resp.Result, nil
}

// GetStickerSet returns StickerSet on success.
func (b *Bot) GetStickerSet(name string) (*StickerSet, error) {
	data, err := b.Raw("getStickerSet", map[string]string{"name": name})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Ok          bool
		Description string
		Result      *StickerSet
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	if !resp.Ok {
		return nil, errors.Errorf("api error: %s", resp.Description)
	}

	return resp.Result, nil
}

// CreateNewStickerSet creates new sticker set.
func (b *Bot) CreateNewStickerSet(sp StickerSetParams, containsMasks bool, maskPosition MaskPosition) error {
	files := map[string]File{
		"png_sticker": *sp.PngSticker,
	}
	params := map[string]string{
		"user_id": strconv.Itoa(sp.UserID),
		"name":    sp.Name,
		"title":   sp.Title,
		"emojis":  sp.Emojis,
	}

	if containsMasks {
		mp, err := json.Marshal(&maskPosition)
		if err != nil {
			return err
		}
		params["mask_position"] = string(mp)
	}

	data, err := b.sendFiles("createNewStickerSet", files, params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// AddStickerToSet adds new sticker to existing sticker set.
func (b *Bot) AddStickerToSet(sp StickerSetParams, maskPosition MaskPosition) error {
	files := map[string]File{
		"png_sticker": *sp.PngSticker,
	}
	params := map[string]string{
		"user_id": strconv.Itoa(sp.UserID),
		"name":    sp.Name,
		"title":   sp.Title,
		"emojis":  sp.Emojis,
	}

	if maskPosition != (MaskPosition{}) {
		mp, err := json.Marshal(&maskPosition)
		if err != nil {
			return err
		}
		params["mask_position"] = string(mp)
	}

	data, err := b.sendFiles("addStickerToSet", files, params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// SetStickerPositionInSet moves a sticker in set to a specific position.
func (b *Bot) SetStickerPositionInSet(sticker string, position int) error {
	params := map[string]string{
		"sticker":  sticker,
		"position": strconv.Itoa(position),
	}
	data, err := b.Raw("setStickerPositionInSet", params)
	if err != nil {
		return err
	}

	return extractOk(data)
}

// DeleteStickerFromSet deletes sticker from set created by the bot.
func (b *Bot) DeleteStickerFromSet(sticker string) error {
	data, err := b.Raw("deleteStickerFromSet", map[string]string{"sticker": sticker})
	if err != nil {
		return err
	}

	return extractOk(data)
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

	data, err := b.Raw("setMyCommands", params)
	if err != nil {
		return err
	}

	return extractOk(data)
}
