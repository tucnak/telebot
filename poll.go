package telebot

// Poll object represents a poll, it contains information about a poll.
type Poll struct {
	ID                    string       `json:"id"`
	Question              string       `json:"question"`
	Options               []PollOption `json:"options"`
	TotalVoterCount       int          `json:"total_voter_count"`
	IsClosed              bool         `json:"is_closed"`
	IsAnonymous           bool         `json:"is_anonymous"`
	Type                  string       `json:"type"`
	AllowsMultipleAnswers bool         `json:"allows_multiple_answers"`

	// Optional. 0-based identifier of the correct answer option.
	// Available only for polls in the quiz mode, which are closed,
	// or was sent (not forwarded) by the bot or to the private chat with the bot.
	CorrectOptionID int `json:"correct_option_id"`
}

// PollOption object represents a option of a poll
type PollOption struct {
	Text       string
	VoterCount int
}

// PollAnswer object represents an answer of a user in a non-anonymous poll.
type PollAnswer struct {
	PollID string `json:"poll_id"`
	User   User   `json:"user"`

	// 0-based identifiers of answer options, chosen by the user. May be empty if the user retracted their vote.
	OptionIDs []int `json:"option_ids"`
}

func (p *Poll) IsRegular() bool {
	return p.Type == "regular"
}

func (p *Poll) IsQuiz() bool {
	return p.Type == "quiz"
}
