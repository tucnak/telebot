package telebot

// Login error, which basically occurs on broken token.
type AuthError struct {
	Payload string
}

type FetchError struct {
	Payload string
}

type SendError struct {
	Payload string
}

func (e AuthError) Error() string {
	return "AuthError: " + e.Payload
}

func (e FetchError) Error() string {
	return "FetchError: " + e.Payload
}

func (e SendError) Error() string {
	return "SendError: " + e.Payload
}
