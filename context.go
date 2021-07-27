package telebot

import (
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// HandlerFunc represents a handler function, which is
// used to handle actual endpoints.
type HandlerFunc func(Context) error

// Context represents a context of the current event. It stores data
// depending on its type, whether it's a message, callback or whatever.
type Context interface {
	// Message returns stored message if such presented.
	Message() *Message

	// Callback returns stored callback if such presented.
	Callback() *Callback

	// Query returns stored query if such presented.
	Query() *Query

	// InlineResult returns stored inline result if such presented.
	InlineResult() *InlineResult

	// ShippingQuery returns stored shipping query if such presented.
	ShippingQuery() *ShippingQuery

	// PreCheckoutQuery returns stored pre checkout query if such presented.
	PreCheckoutQuery() *PreCheckoutQuery

	// Poll returns stored poll if such presented.
	Poll() *Poll

	// PollAnswer returns stored poll answer if such presented.
	PollAnswer() *PollAnswer

	// ChatMember returns bot's chat member changes.
	ChatMember() *ChatMemberUpdated

	// Migration returns both migration from and to chat IDs.
	Migration() (int64, int64)

	// Sender returns the current recipient, depending on the context type.
	// Returns nil if user is not presented.
	Sender() *User

	// Chat returns the current chat, depending on the context type.
	// Returns nil if chat is not presented.
	Chat() *Chat

	// Recipient combines both Sender and Chat functions. If there is no user
	// the chat will be returned. The native context cannot be without sender,
	// but it is useful in the case when the context created intentionally
	// by the NewContext constructor and have only Chat field inside.
	Recipient() Recipient

	// Text returns the message text, depending on the context type.
	// In the case when no related data presented, returns an empty string.
	Text() string

	// Data returns the current data, depending on the context type.
	// If the context contains command, returns its arguments string.
	// If the context contains payment, returns its payload.
	// In the case when no related data presented, returns an empty string.
	Data() string

	// Args returns a raw slice of command or callback arguments as strings.
	// The message arguments split by space, while the callback's ones by a "|" symbol.
	Args() []string

	// Send sends a message to the current recipient.
	// See Send from bot.go.
	Send(what interface{}, opts ...interface{}) error

	// SendAlbum sends an album to the current recipient.
	// See SendAlbum from bot.go.
	SendAlbum(a Album, opts ...interface{}) error

	// Reply replies to the current message.
	// See Reply from bot.go.
	Reply(what interface{}, opts ...interface{}) error

	// Forward forwards the given message to the current recipient.
	// See Forward from bot.go.
	Forward(msg Editable, opts ...interface{}) error

	// ForwardTo forwards the current message to the given recipient.
	// See Forward from bot.go
	ForwardTo(to Recipient, opts ...interface{}) error

	// Edit edits the current message.
	// See Edit from bot.go.
	Edit(what interface{}, opts ...interface{}) error

	// EditCaption edits the caption of the current message.
	// See EditCaption from bot.go.
	EditCaption(caption string, opts ...interface{}) error

	// EditOrSend edits the current message if the update is callback,
	// otherwise the content is sent to the chat as a separate message.
	EditOrSend(what interface{}, opts ...interface{}) error

	// EditOrReply edits the current message if the update is callback,
	// otherwise the content is replied as a separate message.
	EditOrReply(what interface{}, opts ...interface{}) error

	// Delete removes the current message.
	// See Delete from bot.go.
	Delete() error

	// Notify updates the chat action for the current recipient.
	// See Notify from bot.go.
	Notify(action ChatAction) error

	// Ship replies to the current shipping query.
	// See Ship from bot.go.
	Ship(what ...interface{}) error

	// Accept finalizes the current deal.
	// See Accept from bot.go.
	Accept(errorMessage ...string) error

	// Answer sends a response to the current inline query.
	// See Answer from bot.go.
	Answer(resp *QueryResponse) error

	// Respond sends a response for the current callback query.
	// See Respond from bot.go.
	Respond(resp ...*CallbackResponse) error

	// Get retrieves data from the context.
	Get(key string) interface{}

	// Set saves data in the context.
	Set(key string, val interface{})
}

// nativeContext is a native implementation of the Context interface.
// "context" is taken by context package, maybe there is a better name.
type nativeContext struct {
	b *Bot

	message          *Message
	callback         *Callback
	query            *Query
	inlineResult     *InlineResult
	shippingQuery    *ShippingQuery
	preCheckoutQuery *PreCheckoutQuery
	poll             *Poll
	pollAnswer       *PollAnswer
	myChatMember     *ChatMemberUpdated
	chatMember       *ChatMemberUpdated

	lock  sync.RWMutex
	store map[string]interface{}
}

func (c *nativeContext) ChatMemberUpdated() *ChatMemberUpdated {
	switch {
	case c.chatMember != nil:
		return c.chatMember
	case c.myChatMember != nil:
		return c.myChatMember
	default:
		return nil
	}
}

func (c *nativeContext) Message() *Message {
	switch {
	case c.message != nil:
		return c.message
	case c.callback != nil:
		return c.callback.Message
	default:
		return nil
	}
}

func (c *nativeContext) Callback() *Callback {
	return c.callback
}

func (c *nativeContext) Query() *Query {
	return c.query
}

func (c *nativeContext) InlineResult() *InlineResult {
	return c.inlineResult
}

func (c *nativeContext) ShippingQuery() *ShippingQuery {
	return c.shippingQuery
}

func (c *nativeContext) PreCheckoutQuery() *PreCheckoutQuery {
	return c.preCheckoutQuery
}

func (c *nativeContext) Poll() *Poll {
	return c.poll
}

func (c *nativeContext) PollAnswer() *PollAnswer {
	return c.pollAnswer
}

func (c *nativeContext) Migration() (int64, int64) {
	return c.message.MigrateFrom, c.message.MigrateTo
}

func (c *nativeContext) Sender() *User {
	switch {
	case c.message != nil:
		return c.message.Sender
	case c.callback != nil:
		return c.callback.Sender
	case c.query != nil:
		return c.query.Sender
	case c.inlineResult != nil:
		return c.inlineResult.Sender
	case c.shippingQuery != nil:
		return c.shippingQuery.Sender
	case c.preCheckoutQuery != nil:
		return c.preCheckoutQuery.Sender
	case c.pollAnswer != nil:
		return c.pollAnswer.Sender
	default:
		return nil
	}
}

func (c *nativeContext) Chat() *Chat {
	switch {
	case c.message != nil:
		return c.message.Chat
	case c.callback != nil:
		return c.callback.Message.Chat
	default:
		return nil
	}
}

func (c *nativeContext) Recipient() Recipient {
	chat := c.Chat()
	if chat != nil {
		return chat
	}
	return c.Sender()
}

func (c *nativeContext) Text() string {
	switch {
	case c.message != nil:
		return c.message.Text
	case c.callback != nil:
		return c.callback.Message.Text
	default:
		return ""
	}
}

func (c *nativeContext) Data() string {
	switch {
	case c.message != nil:
		return c.message.Payload
	case c.callback != nil:
		return c.callback.Data
	case c.query != nil:
		return c.query.Text
	case c.inlineResult != nil:
		return c.inlineResult.Query
	case c.shippingQuery != nil:
		return c.shippingQuery.Payload
	case c.preCheckoutQuery != nil:
		return c.preCheckoutQuery.Payload
	default:
		return ""
	}
}

func (c *nativeContext) Args() []string {
	if c.message != nil {
		payload := strings.Trim(c.message.Payload, " ")
		if payload == "" {
			return nil
		}
		return strings.Split(payload, " ")
	}
	if c.callback != nil {
		return strings.Split(c.callback.Data, "|")
	}
	return nil
}

func (c *nativeContext) Send(what interface{}, opts ...interface{}) error {
	_, err := c.b.Send(c.Recipient(), what, opts...)
	return err
}

func (c *nativeContext) SendAlbum(a Album, opts ...interface{}) error {
	_, err := c.b.SendAlbum(c.Recipient(), a, opts...)
	return err
}

func (c *nativeContext) Reply(what interface{}, opts ...interface{}) error {
	msg := c.Message()
	if msg == nil {
		return ErrBadContext
	}
	_, err := c.b.Reply(msg, what, opts...)
	return err
}

func (c *nativeContext) Forward(msg Editable, opts ...interface{}) error {
	_, err := c.b.Forward(c.Recipient(), msg, opts...)
	return err
}

func (c *nativeContext) ForwardTo(to Recipient, opts ...interface{}) error {
	msg := c.Message()
	if msg == nil {
		return ErrBadContext
	}
	_, err := c.b.Forward(to, msg, opts...)
	return err
}

func (c *nativeContext) Edit(what interface{}, opts ...interface{}) error {
	if c.inlineResult != nil {
		_, err := c.b.Edit(c.inlineResult, what, opts...)
		return err
	}
	clb := c.callback
	if clb == nil || clb.Message == nil {
		return ErrBadContext
	}
	_, err := c.b.Edit(clb.Message, what, opts...)
	return err
}

func (c *nativeContext) EditCaption(caption string, opts ...interface{}) error {
	clb := c.callback
	if clb == nil || clb.Message == nil {
		return ErrBadContext
	}
	_, err := c.b.EditCaption(clb.Message, caption, opts...)
	return err
}

func (c *nativeContext) EditOrSend(what interface{}, opts ...interface{}) error {
	if c.callback != nil {
		return c.Edit(what, opts...)
	}
	return c.Send(what, opts...)
}

func (c *nativeContext) EditOrReply(what interface{}, opts ...interface{}) error {
	if c.callback != nil {
		return c.Edit(what, opts...)
	}
	return c.Reply(what, opts...)
}

func (c *nativeContext) Delete() error {
	msg := c.Message()
	if msg == nil {
		return ErrBadContext
	}
	return c.b.Delete(msg)
}

func (c *nativeContext) Notify(action ChatAction) error {
	return c.b.Notify(c.Recipient(), action)
}

func (c *nativeContext) Ship(what ...interface{}) error {
	if c.shippingQuery == nil {
		return errors.New("telebot: context shipping query is nil")
	}
	return c.b.Ship(c.shippingQuery, what...)
}

func (c *nativeContext) Accept(errorMessage ...string) error {
	if c.preCheckoutQuery == nil {
		return errors.New("telebot: context pre checkout query is nil")
	}
	return c.b.Accept(c.preCheckoutQuery, errorMessage...)
}

func (c *nativeContext) Answer(resp *QueryResponse) error {
	if c.query == nil {
		return errors.New("telebot: context inline query is nil")
	}
	return c.b.Answer(c.query, resp)
}

func (c *nativeContext) Respond(resp ...*CallbackResponse) error {
	if c.callback == nil {
		return errors.New("telebot: context callback is nil")
	}
	return c.b.Respond(c.callback, resp...)
}

func (c *nativeContext) Set(key string, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.store == nil {
		c.store = make(map[string]interface{})
	}
	c.store[key] = value
}

func (c *nativeContext) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.store[key]
}
