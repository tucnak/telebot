package telebot

// Default is a waste!
const Default string = ""

// Handler is smth.
type Handler func(Context)

// Context is smth.
type Context struct {
	Message *Message
}
