package telebot

import "time"

// VoiceChatStarted represents a service message about a voice chat
// started in the chat.
type VoiceChatStarted struct{}

// VoiceChatEnded represents a service message about a voice chat
// ended in the chat.
type VoiceChatEnded struct {
	Duration int `json:"duration"`
}

// VoiceChatPartecipantsInvited represents a service message about new
// members invited to a voice chat
type VoiceChatPartecipantsInvited struct {
	Users []User `json:"users"`
}

// VoiceChatSchedule represents a service message about a voice chat scheduled in the chat.
type VoiceChatSchedule struct {
	StartUnixTime int64 `json:"start_date"`
}

// ExpireDate returns the point when the voice chat is supposed to be started by a chat administrator.
func (v *VoiceChatSchedule) ExpireDate() time.Time {
	return time.Unix(v.StartUnixTime, 0)
}
