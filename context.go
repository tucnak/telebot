package telebot

type Context struct {
	Update *Update
	abort  bool
}

func (c *Context) Abort() {
	c.abort = true
}
