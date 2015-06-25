package telebot

// Login error, which basically occurs on broken token.
type AuthError struct {
	Payload string
}

func (e AuthError) Error() string {
	return "AuthError: " + e.Payload
}
