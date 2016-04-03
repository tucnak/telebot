package telebot

type Controller func(*Message, *map[string]string) error
