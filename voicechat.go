package telebot

// VoiceChatStarted represents a service message about a voice chat
// started in the chat.
type VoiceChatStarted struct{}

// VoiceChatEnded represents a service message about a voice chat
// ended in the chat.
type VoiceChatEnded struct {
	Duration int `json:"duration"`
}

// VoiceChatParticipants represents a service message about new
// members invited to a voice chat
type VoiceChatParticipants struct {
	Users []User `json:"users"`
}
