package telebot

type Message struct {
	Id       int    `json:"message_id"`
	Sender   User   `json:"from"`
	Unixtime int    `json:"date"`
	Text     string `json:"text"`
	// TBA: `chat`
}
