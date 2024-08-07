package teletest

import (
	"testing"

	tele "gopkg.in/telebot.v3"
)

type Recorder struct {
	b *tele.Bot
	e *Expect
}

func New(b *tele.Bot) (*Recorder, *Expect) {
	r := &Recorder{b: b, e: &Expect{}}
	return r, r.e
}

func (r *Recorder) Trigger(u tele.Update) {
	c := &Context{nc: r.b.NewContext(u)}
	r.b.ProcessContext(c)
	*r.e = Expect{c: c}
}

type Expect struct {
	c *Context
}

func (e Expect) Reply(t *testing.T, what interface{}, opts ...interface{}) {
	e.c.reply.Expect(t, what, opts)
}
