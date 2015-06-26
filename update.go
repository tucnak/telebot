package telebot

type Update struct {
	Id      int     `json:"update_id"`
	Payload Message `json:"message"`
}
