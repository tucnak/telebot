package teletest

import (
	"time"

	tele "gopkg.in/telebot.v3"
)

type Context struct {
	nc tele.Context

	send  Response
	reply Response
}

func (c *Context) Bot() *tele.Bot {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Update() tele.Update {
	return c.nc.Update()
}

func (c *Context) Message() *tele.Message {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Callback() *tele.Callback {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Query() *tele.Query {
	//TODO implement me
	panic("implement me")
}

func (c *Context) InlineResult() *tele.InlineResult {
	//TODO implement me
	panic("implement me")
}

func (c *Context) ShippingQuery() *tele.ShippingQuery {
	//TODO implement me
	panic("implement me")
}

func (c *Context) PreCheckoutQuery() *tele.PreCheckoutQuery {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Poll() *tele.Poll {
	//TODO implement me
	panic("implement me")
}

func (c *Context) PollAnswer() *tele.PollAnswer {
	//TODO implement me
	panic("implement me")
}

func (c *Context) ChatMember() *tele.ChatMemberUpdate {
	//TODO implement me
	panic("implement me")
}

func (c *Context) ChatJoinRequest() *tele.ChatJoinRequest {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Migration() (int64, int64) {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Topic() *tele.Topic {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Boost() *tele.BoostUpdated {
	//TODO implement me
	panic("implement me")
}

func (c *Context) BoostRemoved() *tele.BoostRemoved {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Sender() *tele.User {
	return c.nc.Sender()
}

func (c *Context) Chat() *tele.Chat {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Recipient() tele.Recipient {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Text() string {
	return c.nc.Text()
}

func (c *Context) Entities() tele.Entities {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Data() string {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Args() []string {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Send(what interface{}, opts ...interface{}) error {
	c.send.what = what
	c.send.opts = opts
	return nil
}

func (c *Context) SendAlbum(a tele.Album, opts ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Reply(what interface{}, opts ...interface{}) error {
	c.reply = Response{what, opts}
	return nil
}

func (c *Context) Forward(msg tele.Editable, opts ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) ForwardTo(to tele.Recipient, opts ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Edit(what interface{}, opts ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) EditCaption(caption string, opts ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) EditOrSend(what interface{}, opts ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) EditOrReply(what interface{}, opts ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Delete() error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) DeleteAfter(d time.Duration) *time.Timer {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Notify(action tele.ChatAction) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Ship(what ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Accept(errorMessage ...string) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Answer(resp *tele.QueryResponse) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Respond(resp ...*tele.CallbackResponse) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) RespondText(text string) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) RespondAlert(text string) error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Get(key string) interface{} {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Set(key string, val interface{}) {
	//TODO implement me
	panic("implement me")
}
