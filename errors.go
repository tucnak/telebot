package telebot

type AuthError struct {
	Payload string
}

func (e AuthError) Error() string {
	return "AuthError: " + e.Payload
}
