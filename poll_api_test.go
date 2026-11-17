package telebot

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendPollWithMedia(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendPoll", json.RawMessage(`{"message_id":1,"chat":{"id":-100}}`))

	acc := float32(5)
	poll := &Poll{
		Type:                   PollQuiz,
		Question:               "Which?",
		QuestionEntities:       []MessageEntity{{Type: EntityCustomEmoji, Length: 2, CustomEmojiID: "1"}},
		CorrectOptions:         []int{0, 2},
		MultipleAnswers:        true,
		Explanation:            "because",
		Entities:               []MessageEntity{{Type: EntityBold, Length: 7}},
		InputExplanationMedia:  &Location{Lat: 1, Lng: 2, HorizontalAccuracy: &acc},
		Description:            "<i>pick</i>",
		DescriptionParseMode:   ModeHTML,
		InputMedia:             &Photo{File: FromReader(strings.NewReader("img"))},
		ShuffleOptions:         true,
		HideResultsUntilCloses: true,
		MembersOnly:            true,
		CountryCodes:           []string{"DE", "FT"},
		Options: []PollOption{
			{Text: "a", InputMedia: &Link{URL: "https://example.com"}},
			{Text: "b", InputMedia: &Sticker{File: File{FileID: "st"}}},
			{Text: "c", InputMedia: &Venue{Location: Location{Lat: 3, Lng: 4}, Title: "T", Address: "A"}},
		},
	}
	_, err := b.Send(ChatID(-100), poll)
	require.NoError(t, err)

	req := srv.Last()
	assert.Equal(t, "sendPoll", req.Method)
	assert.Equal(t, "[0,2]", req.Params["correct_option_ids"])
	assert.NotContains(t, req.Params, "correct_option_id")
	assert.JSONEq(t, `[{"type":"custom_emoji","offset":0,"length":2,"custom_emoji_id":"1"}]`, req.Params["question_entities"])
	assert.JSONEq(t, `[{"type":"bold","offset":0,"length":7}]`, req.Params["explanation_entities"])
	assert.NotContains(t, req.Params, "explanation_parse_mode")
	assert.JSONEq(t, `{"type":"location","latitude":1,"longitude":2,"horizontal_accuracy":5}`, req.Params["explanation_media"])
	assert.Equal(t, "<i>pick</i>", req.Params["description"])
	assert.Equal(t, "HTML", req.Params["description_parse_mode"])
	assert.JSONEq(t, `{"type":"photo","media":"attach://media","caption":""}`, req.Params["media"])
	assert.Equal(t, "img", string(req.Files["media"].Data))
	assert.Equal(t, "true", req.Params["shuffle_options"])
	assert.Equal(t, "true", req.Params["hide_results_until_closes"])
	assert.Equal(t, "true", req.Params["members_only"])
	assert.NotContains(t, req.Params, "allows_revoting")
	assert.Equal(t, `["DE","FT"]`, req.Params["country_codes"])
	assert.JSONEq(t, `[
		{"text":"a","media":{"type":"link","url":"https://example.com"}},
		{"text":"b","media":{"type":"sticker","media":"st"}},
		{"text":"c","media":{"type":"venue","latitude":3,"longitude":4,"title":"T","address":"A"}}
	]`, req.Params["options"])
}

func TestSendQuizSingleCorrectOption(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendPoll", json.RawMessage(`{"message_id":1,"chat":{"id":-100}}`))

	poll := &Poll{Type: PollQuiz, Question: "q", CorrectOption: 1}
	poll.AddOptions("x", "y")
	_, err := b.Send(ChatID(-100), poll)
	require.NoError(t, err)
	assert.Equal(t, "[1]", srv.Last().Params["correct_option_ids"])
	assert.JSONEq(t, `[{"text":"x"},{"text":"y"}]`, srv.Last().Params["options"])

	regular := &Poll{Question: "q", Type: PollRegular}
	regular.AddOptions("x")
	_, err = b.Send(ChatID(-100), regular)
	require.NoError(t, err)
	assert.NotContains(t, srv.Last().Params, "correct_option_ids")
}

func TestReceivePoll96(t *testing.T) {
	var p Poll
	require.NoError(t, json.Unmarshal([]byte(`{"id":"p","question":"q","type":"quiz","total_voter_count":3,
		"is_closed":true,"is_anonymous":false,"allows_multiple_answers":false,"allows_revoting":true,"members_only":true,
		"country_codes":["DE"],"correct_option_ids":[1],"description":"d","media":{"link":{"url":"https://example.com"}},
		"explanation_media":{"photo":[{"file_id":"ph"}]},
		"options":[{"persistent_id":"o1","text":"a","voter_count":1,"added_by_user":{"id":7},"addition_date":5,
			"media":{"live_photo":`+livePhotoJSON+`}}]}`), &p))

	assert.Equal(t, 1, p.CorrectOption)
	assert.Equal(t, []int{1}, p.CorrectOptions)
	assert.True(t, p.AllowsRevoting)
	assert.True(t, p.MembersOnly)
	assert.Equal(t, "https://example.com", p.Media.Link.URL)
	assert.Equal(t, "ph", p.ExplanationMedia.Photo[0].FileID)
	o := p.Options[0]
	assert.Equal(t, "o1", o.PersistentID)
	assert.Equal(t, int64(7), o.AddedByUser.ID)
	assert.Equal(t, int64(5), o.AdditionDate)
	assert.Equal(t, "lv", o.Media.LivePhoto.FileID)

	var a PollAnswer
	require.NoError(t, json.Unmarshal([]byte(`{"poll_id":"p","option_ids":[0],"option_persistent_ids":["o1"]}`), &a))
	assert.Equal(t, []string{"o1"}, a.OptionPersistentIDs)
}

func TestPollOptionServiceMessages(t *testing.T) {
	b, srv := newFakeBot(t)

	var added *PollOptionAdded
	var deleted *PollOptionDeleted
	b.Handle(OnPollOptionAdded, func(c Context) error { added = c.Message().PollOptionAdded; return nil })
	b.Handle(OnPollOptionDeleted, func(c Context) error { deleted = c.Message().PollOptionDeleted; return nil })

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":2,"date":1,"chat":{"id":-100},
		"poll_option_added":{"poll_message":{"message_id":1,"date":1,"chat":{"id":-100}},"option_persistent_id":"o9","option_text":"new"}}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, added)
	assert.Equal(t, "o9", added.OptionPersistentID)
	assert.Equal(t, 1, added.PollMessage.ID)

	u = Update{}
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":2,"message":{"message_id":3,"date":1,"chat":{"id":-100},
		"reply_to_poll_option_id":"o9","poll_option_deleted":{"option_persistent_id":"o9","option_text":"new"}}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, deleted)
	assert.Equal(t, "new", deleted.OptionText)

	srv.Reply("sendMessage", json.RawMessage(`{"message_id":4,"chat":{"id":-100}}`))
	_, err := b.Send(ChatID(-100), "on option", &ReplyParams{MessageID: 1, PollOptionID: "o9"})
	require.NoError(t, err)
	assert.JSONEq(t, `{"message_id":1,"poll_option_id":"o9"}`, srv.Last().Params["reply_parameters"])
}

func TestPollMediaLive(t *testing.T) {
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}
	if userID == 0 {
		t.Skip("USER_ID is required for poll live tests")
	}

	poll := &Poll{
		Type:            PollQuiz,
		Question:        "telebot quiz",
		Description:     "Pick both",
		InputMedia:      &Location{Lat: 52.52, Lng: 13.40},
		CorrectOptions:  []int{0, 1},
		MultipleAnswers: true,
		Options: []PollOption{
			{Text: "one", InputMedia: &Link{URL: "https://github.com/tucnak/telebot"}},
			{Text: "two", InputMedia: &Location{Lat: 1, Lng: 2}},
			{Text: "three"},
		},
	}
	msg, err := b.Send(user, poll)
	require.NoError(t, err)
	require.NotNil(t, msg.Poll)
	assert.Equal(t, []int{0, 1}, msg.Poll.CorrectOptions)
	assert.Equal(t, "Pick both", msg.Poll.Description)
	require.NotNil(t, msg.Poll.Media)
	assert.NotNil(t, msg.Poll.Media.Location)
	require.Len(t, msg.Poll.Options, 3)
	assert.NotEmpty(t, msg.Poll.Options[0].PersistentID)
	require.NotNil(t, msg.Poll.Options[0].Media)
	assert.Equal(t, "https://github.com/tucnak/telebot", msg.Poll.Options[0].Media.Link.URL)
}
