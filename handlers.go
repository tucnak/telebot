package telebot

// Default handler prefix.
const Default string = ""

type Context struct {
	Bot     *Bot
	Message Message
}

type Handler func(Context)

func (b *Bot) Handle(prefix string, handler Handler) {
	b.tree.Insert(prefix, handler)
}

func (b *Bot) Serve(msg Message) (ok bool) {
	request := msg.Text

	_, value, _ := b.tree.LongestPrefix(request)
	endpoint, ok := value.(Handler)

	if ok {
		go endpoint(Context{b, msg})
	}
	return
}
