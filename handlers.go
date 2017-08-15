package telebot

// Default handler prefix.
const Default string = ""

type Context struct {
	Bot *Bot
	Msg *Message
}

type Handler func(Context)

func (b *Bot) Handle(prefix string, handler Handler) {
	b.tree.Insert(prefix, handler)
}

func (b *Bot) Serve(msg *Message) {
	request := msg.Text

	_, value, _ := b.tree.LongestPrefix(request)

	if endpoint, ok := value.(Handler); ok {
		endpoint(Context{b, msg})
	} else {
		panic("telebot: couldn't find a route")
	}
}
