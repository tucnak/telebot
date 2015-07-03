package telebot

// AuthError occurs if the token appears to be invalid.
type AuthError struct {
	Payload string
}

// FetchError occurs when something goes wrong
// while fetching updates.
type FetchError struct {
	Payload string
}

// SendError occurs when something goes wrong
// while posting images, documents, etc.
type SendError struct {
	Payload string
}

// FileError occurs when local file can't be read.
type FileError struct {
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

func (e FileError) Error() string {
	return "FileError: " + e.Payload
}
