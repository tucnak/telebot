package telebot

type User struct {
	Id        int
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string
}
