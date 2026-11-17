package telebot

import "time"

type (
	// VideoChatStarted represents a service message about a video chat
	// started in the chat.
	//
	// Bot API: VideoChatStarted.
	VideoChatStarted struct{}

	// VideoChatEnded represents a service message about a video chat
	// ended in the chat.
	//
	// Bot API: VideoChatEnded.
	VideoChatEnded struct {
		Duration int `json:"duration"` // in seconds
	}

	// VideoChatParticipants represents a service message about new
	// members invited to a video chat
	//
	// Bot API: VideoChatParticipantsInvited.
	VideoChatParticipants struct {
		Users []User `json:"users"`
	}

	// VideoChatScheduled represents a service message about a video chat scheduled in the chat.
	//
	// Bot API: VideoChatScheduled.
	VideoChatScheduled struct {
		Unixtime int64 `json:"start_date"`
	}
)

// StartsAt returns the point when the video chat is supposed to be started by a chat administrator.
func (v *VideoChatScheduled) StartsAt() time.Time {
	return time.Unix(v.Unixtime, 0)
}
