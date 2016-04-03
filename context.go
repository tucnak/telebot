package telebot

const Default string = ""

type Context struct {
	Message *Message
	Args    map[string]string
}

type Handler func(Context)
