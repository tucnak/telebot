package telebot

// Default represents the empty value
const Default string = ""

// Context is passed to message handlers
type Context struct {
	Message *Message
	Args    map[string]string
}

// Handler represents a message handler
type Handler func(Context)
