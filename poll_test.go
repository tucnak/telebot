package telebot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}
	if userID == 0 {
		t.Skip("USER_ID is required for Poll methods test")
	}

	_, err := b.Send(user, &Poll{}) // empty poll
	assert.Equal(t, ErrBadPollOptions, err)

	poll := &Poll{
		Type:          PollQuiz,
		Question:      "Test Poll",
		CloseUnixdate: time.Now().Unix() + 60,
		Explanation:   "Explanation",
	}
	poll.AddOptions("1", "2")

	msg, err := b.Send(user, poll)
	require.NoError(t, err)
	assert.Equal(t, poll.Type, msg.Poll.Type)
	assert.Equal(t, poll.Question, msg.Poll.Question)
	assert.Equal(t, poll.Options, msg.Poll.Options)
	assert.Equal(t, poll.CloseUnixdate, msg.Poll.CloseUnixdate)
	assert.Equal(t, poll.CloseDate(), msg.Poll.CloseDate())

	p, err := b.StopPoll(msg)
	require.NoError(t, err)
	assert.Equal(t, poll.Options, p.Options)
	assert.Equal(t, 0, p.VoterCount)
}
