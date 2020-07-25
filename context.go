package telebot

import "github.com/pkg/errors"

// HandlerFunc represents a handler function type
// which is used to handle endpoints.
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

	// ChosenInlineResult returns stored inline result if such presented.
	ChosenInlineResult() *ChosenInlineResult

	// ShippingQuery returns stored shipping query if such presented.
	ShippingQuery() *ShippingQuery

	// PreCheckoutQuery returns stored pre checkout query if such presented.
	PreCheckoutQuery() *PreCheckoutQuery

	// Poll returns stored poll if such presented.
	Poll() *Poll

	// PollAnswer returns stored poll answer if such presented.
	PollAnswer() *PollAnswer

	// Migration returns both migration from and to chat IDs.
	Migration() (int64, int64)

	// Sender returns a current recipient, depending on the context type.
	// Returns nil if user is not presented.
	Sender() Recipient

	// Text returns a current text, depending on the context type.
	// If the context contains payment, returns its payload.
	// In the case when no related data presented, returns an empty string.
	Text() string

	// Send sends a message to the current recipient.
	// See Send from bot.go.
	Send(what interface{}, opts ...interface{}) error

	// SendAlbum sends an album to the current recipient.
	// See SendAlbum from bot.go.
	SendAlbum(a Album, opts ...interface{}) error

	// Reply replies to the current message.
	// See Reply from bot.go.
	Reply(what interface{}, opts ...interface{}) error

	// Forward forwards a given message to the current recipient.
	// See Forward from bot.go.
	Forward(msg Editable, opts ...interface{}) error

	// Edit edits a current message.
	// See Edit from bot.go.
	Edit(what interface{}, opts ...interface{}) error

	// EditCaption edits a caption of the current message.
	// See EditCaption from bot.go.
	EditCaption(caption string, opts ...interface{}) error

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
}

// nativeContext is a native implementation of the Context interface.
// "context" is taken by context package, maybe there is a better name.
type nativeContext struct {
	b *Bot

	message            *Message
	callback           *Callback
	query              *Query
	chosenInlineResult *ChosenInlineResult
	shippingQuery      *ShippingQuery
	preCheckoutQuery   *PreCheckoutQuery
	poll               *Poll
	pollAnswer         *PollAnswer
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

func (c *nativeContext) ChosenInlineResult() *ChosenInlineResult {
	return c.chosenInlineResult
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

func (c *nativeContext) Sender() Recipient {
	switch {
	case c.message != nil:
		return c.message.Sender
	case c.callback != nil:
		return c.callback.Sender
	case c.query != nil:
		return c.query.Sender
	case c.chosenInlineResult != nil:
		return c.chosenInlineResult.Sender
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

func (c *nativeContext) Text() string {
	switch {
	case c.message != nil:
		return c.message.Text
	case c.callback != nil:
		return c.callback.Message.Text
	case c.query != nil:
		return c.query.Text
	case c.chosenInlineResult != nil:
		return c.chosenInlineResult.Query
	case c.shippingQuery != nil:
		return c.shippingQuery.Payload
	case c.preCheckoutQuery != nil:
		return c.preCheckoutQuery.Payload
	default:
		return ""
	}
}

func (c *nativeContext) Send(what interface{}, opts ...interface{}) error {
	_, err := c.b.Send(c.Sender(), what, opts...)
	return err
}

func (c *nativeContext) SendAlbum(a Album, opts ...interface{}) error {
	_, err := c.b.SendAlbum(c.Sender(), a, opts...)
	return err
}

func (c *nativeContext) Reply(what interface{}, opts ...interface{}) error {
	if c.message == nil {
		return ErrBadContext
	}
	_, err := c.b.Reply(c.message, what, opts...)
	return err
}

func (c *nativeContext) Forward(msg Editable, opts ...interface{}) error {
	_, err := c.b.Forward(c.Sender(), msg, opts...)
	return err
}

func (c *nativeContext) Edit(what interface{}, opts ...interface{}) error {
	if c.message == nil {
		return ErrBadContext
	}
	_, err := c.b.Edit(c.message, what, opts...)
	return err
}

func (c *nativeContext) EditCaption(caption string, opts ...interface{}) error {
	if c.message == nil {
		return ErrBadContext
	}
	_, err := c.b.EditCaption(c.message, caption, opts...)
	return err
}

func (c *nativeContext) Delete() error {
	if c.message == nil {
		return ErrBadContext
	}
	return c.b.Delete(c.message)
}

func (c *nativeContext) Notify(action ChatAction) error {
	return c.b.Notify(c.Sender(), action)
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
