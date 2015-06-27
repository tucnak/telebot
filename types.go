package telebot

// User object represents a Telegram user, bot or group chat.
type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`

	// Title differs a group chat apart from users and bots.
	Title string `json:"title"`
}

// Message object represents a message.
type Message struct {
	Id       int    `json:"message_id"`
	Sender   User   `json:"from"`
	Unixtime int    `json:"date"`
	Text     string `json:"text"`
	Chat     User   `json:"chat"`
}

// Update object represents an incoming update.
type Update struct {
	Id      int     `json:"update_id"`
	Payload Message `json:"message"`
}
