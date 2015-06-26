package telebot

type User struct {
	Id        int
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string

	// In case of group chat, Title will indicate
	// whether it's a chat or user: if Title is empty
	// it's a user, otherwise it's not.
	Title string
}
