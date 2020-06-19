package telebot

import "github.com/pkg/errors"

// Сontext represents a context of the current event. It stores data
// depending on its type, whether it is a message, callback or whatever.
type Context interface {
	// Sender returns the current recipient, depending on the context type.
	// Returns nil if user is not presented.
	Sender() Recipient

	// Text returns the current text, depending on the context type.
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

	// Forward forwards given message to the current recipient.
	// See Forward from bot.go.
	Forward(msg Editable, opts ...interface{}) error

	// Edit edits the current message.
	// See Edit from bot.go.
	Edit(what interface{}, opts ...interface{}) error

	// EditCaption edits the caption of the current message.
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

	*Message
	*Callback
	*Query
	*ChosenInlineResult
	*ShippingQuery
	*PreCheckoutQuery
	*Poll
	*PollAnswer
}

func (c *nativeContext) Sender() Recipient {
	switch {
	case c.Message != nil:
		return c.Message.Sender
	case c.Callback != nil:
		return c.Callback.Sender
	case c.Query != nil:
		return c.Query.Sender
	case c.ChosenInlineResult != nil:
		return c.ChosenInlineResult.Sender
	case c.ShippingQuery != nil:
		return c.ShippingQuery.Sender
	case c.PreCheckoutQuery != nil:
		return c.PreCheckoutQuery.Sender
	case c.PollAnswer != nil:
		return c.PollAnswer.Sender
	default:
		return nil
	}
}

func (c *nativeContext) Text() string {
	switch {
	case c.Message != nil:
		return c.Message.Text
	case c.Callback != nil:
		return c.Callback.Message.Text
	case c.Query != nil:
		return c.Query.Text
	case c.ChosenInlineResult != nil:
		return c.ChosenInlineResult.Query
	case c.ShippingQuery != nil:
		return c.ShippingQuery.Payload
	case c.PreCheckoutQuery != nil:
		return c.PreCheckoutQuery.Payload
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
	if c.Message == nil {
		return ErrBadContext
	}
	_, err := c.b.Reply(c.Message, what, opts...)
	return err
}

func (c *nativeContext) Forward(msg Editable, opts ...interface{}) error {
	_, err := c.b.Forward(c.Sender(), msg, opts...)
	return err
}

func (c *nativeContext) Edit(what interface{}, opts ...interface{}) error {
	if c.Message == nil {
		return ErrBadContext
	}
	_, err := c.b.Edit(c.Message, what, opts...)
	return err
}

func (c *nativeContext) EditCaption(caption string, opts ...interface{}) error {
	if c.Message == nil {
		return ErrBadContext
	}
	_, err := c.b.EditCaption(c.Message, caption, opts...)
	return err
}

func (c *nativeContext) Delete() error {
	if c.Message == nil {
		return ErrBadContext
	}
	return c.b.Delete(c.Message)
}

func (c *nativeContext) Notify(action ChatAction) error {
	return c.b.Notify(c.Sender(), action)
}

func (c *nativeContext) Ship(what ...interface{}) error {
	if c.ShippingQuery == nil {
		return errors.New("telebot: context shipping query is nil")
	}
	return c.b.Ship(c.ShippingQuery, what...)
}

func (c *nativeContext) Accept(errorMessage ...string) error {
	if c.PreCheckoutQuery == nil {
		return errors.New("telebot: context pre checkout query is nil")
	}
	return c.b.Accept(c.PreCheckoutQuery, errorMessage...)
}

func (c *nativeContext) Answer(resp *QueryResponse) error {
	if c.Query == nil {
		return errors.New("telebot: context inline query is nil")
	}
	return c.b.Answer(c.Query, resp)
}

func (c *nativeContext) Respond(resp ...*CallbackResponse) error {
	if c.Callback == nil {
		return errors.New("telebot: context callback is nil")
	}
	return c.b.Respond(c.Callback, resp...)
}
