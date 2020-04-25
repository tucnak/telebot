package telebot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPoll(t *testing.T) {
	assert.True(t, (&Poll{Type: PollRegular}).IsRegular())
	assert.True(t, (&Poll{Type: PollQuiz}).IsQuiz())

	p := &Poll{}
	opts := []PollOption{{Text: "Option 1"}, {Text: "Option 2"}}
	p.AddOptions(opts[0].Text, opts[1].Text)
	assert.Equal(t, opts, p.Options)
}

func TestPollSend(t *testing.T) {
	if token == "" {
		t.Skip("TELEBOT_SECRET is required")
	}
	if userID == 0 {
		t.Skip("USER_ID is required for Poll methods test")
	}

	p := &Poll{
		Type:          PollQuiz,
		Question:      "Test Poll",
		CloseUnixdate: time.Now().Unix() + 60,
		Explanation:   "Explanation",
	}
	p.AddOptions("1", "2")

	markup := &ReplyMarkup{
		ReplyKeyboard: [][]ReplyButton{{{
			Text: "Poll",
			Poll: PollAny,
		}}},
	}

	msg, err := b.Send(user, p, markup)
	assert.NoError(t, err)
	assert.Equal(t, p.Question, msg.Poll.Question)
	assert.Equal(t, p.Options, msg.Poll.Options)
	assert.Equal(t, p.CloseUnixdate, msg.Poll.CloseUnixdate)
	assert.Equal(t, p.CloseDate(), msg.Poll.CloseDate())

	_, err = b.Send(user, &Poll{}) // empty poll
	assert.Equal(t, ErrBadPollOptions, err)
}
