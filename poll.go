package telebot

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// PollType defines poll types. As ReplyButton.Poll, it is sent as a
// KeyboardButtonPollType.
//
// Bot API: KeyboardButtonPollType.
type PollType string

const (
	// NOTE:
	// Despite "any" type isn't described in documentation,
	// it needed for proper KeyboardButtonPollType marshaling.
	PollAny PollType = "any"

	PollQuiz    PollType = "quiz"
	PollRegular PollType = "regular"
)

// Poll contains information about a poll.
//
// Bot API: Poll.
type Poll struct {
	ID         string       `json:"id"`
	Type       PollType     `json:"type"`
	Question   string       `json:"question"`
	Options    []PollOption `json:"options"`
	VoterCount int          `json:"total_voter_count"`

	// (Optional)
	Closed            bool            `json:"is_closed,omitempty"`
	CorrectOption     int             `json:"correct_option_id,omitempty"`
	MultipleAnswers   bool            `json:"allows_multiple_answers,omitempty"`
	Explanation       string          `json:"explanation,omitempty"`
	ParseMode         ParseMode       `json:"explanation_parse_mode,omitempty"`
	Entities          []MessageEntity `json:"explanation_entities,omitempty"`
	QuestionParseMode string          `json:"question_parse_mode,omitempty"`
	QuestionEntities  []MessageEntity `json:"question_entities,omitempty"`

	// True by default, shouldn't be omitted.
	Anonymous bool `json:"is_anonymous"`

	// (Mutually exclusive)
	OpenPeriod    int   `json:"open_period,omitempty"`
	CloseUnixdate int64 `json:"close_date,omitempty"`

	// Bot API 9.6: the correct options of a quiz, which may be several.
	// CorrectOption is still sent when this is empty, and set from it on
	// received polls with one correct option.
	CorrectOptions []int `json:"correct_option_ids,omitempty"`

	// Bot API 9.6: voters can change their answers. Unless told otherwise,
	// Telegram allows it in regular polls and forbids it in quizzes: when
	// sending, AllowsRevoting allows it and DisallowRevoting forbids it.
	AllowsRevoting   bool `json:"allows_revoting,omitempty"`
	DisallowRevoting bool `json:"-"`

	// Bot API 9.6: description of the poll, 0-1024 characters, with its
	// parse mode (when sending) or entities.
	Description          string          `json:"description,omitempty"`
	DescriptionParseMode ParseMode       `json:"-"`
	DescriptionEntities  []MessageEntity `json:"description_entities,omitempty"`

	// Bot API 10.0: only users who have been members of the channel for more
	// than 24 hours, and from the given countries ("FT" for Fragment
	// numbers), can vote.
	MembersOnly  bool     `json:"members_only,omitempty"`
	CountryCodes []string `json:"country_codes,omitempty"`

	// Bot API 10.0: media of the poll description and of the quiz
	// explanation, as received.
	Media            *PollMedia `json:"media,omitempty"`
	ExplanationMedia *PollMedia `json:"explanation_media,omitempty"`

	// Bot API 10.0: media of the poll description and of the quiz
	// explanation, to send.
	InputMedia            InputPollMedia `json:"-"`
	InputExplanationMedia InputPollMedia `json:"-"`

	// Bot API 9.6, sending only: show the options in random order, let users
	// add options (not for anonymous polls and quizzes), and show results only
	// once the poll closes.
	ShuffleOptions         bool `json:"-"`
	AllowAddingOptions     bool `json:"-"`
	HideResultsUntilCloses bool `json:"-"`
}

// UnmarshalJSON decodes a poll, setting CorrectOption from the only correct
// option of a quiz.
func (p *Poll) UnmarshalJSON(data []byte) error {
	type alias Poll
	if err := json.Unmarshal(data, (*alias)(p)); err != nil {
		return err
	}
	if len(p.CorrectOptions) == 1 && p.CorrectOption == 0 {
		p.CorrectOption = p.CorrectOptions[0]
	}
	return nil
}

// PollMedia is media in a poll description, quiz explanation or poll option
// (Bot API 10.0). At most one field is set.
//
// Bot API: PollMedia.
type PollMedia struct {
	Animation *Animation `json:"animation,omitempty"`
	Audio     *Audio     `json:"audio,omitempty"`
	Document  *Document  `json:"document,omitempty"`
	Link      *Link      `json:"link,omitempty"`
	LivePhoto *LivePhoto `json:"live_photo,omitempty"`
	Location  *Location  `json:"location,omitempty"`
	Photo     []Photo    `json:"photo,omitempty"`
	Sticker   *Sticker   `json:"sticker,omitempty"`
	Venue     *Venue     `json:"venue,omitempty"`
	Video     *Video     `json:"video,omitempty"`
}

// Link is an HTTP link, attached to a poll option (Bot API 10.1). Send it as
// an InputMediaLink by setting it as a PollOption's InputMedia.
//
// Bot API: Link, InputMediaLink.
type Link struct {
	URL string `json:"url"`
}

// InputPollMedia is media to send with a poll, a quiz explanation or a poll
// option (Bot API 10.0): a *Photo, *Video, *Animation, *LivePhoto or a
// *Location or *Venue; *Audio and *Document for the poll and explanation
// only; *Sticker and *Link for options only.
//
// Bot API: InputPollMedia, InputPollOptionMedia.
type InputPollMedia interface{}

// PollOption contains information about one answer option in a poll.
//
// Bot API: PollOption, InputPollOption.
type PollOption struct {
	Text       string          `json:"text"`
	VoterCount int             `json:"voter_count"`
	ParseMode  ParseMode       `json:"text_parse_mode,omitempty"`
	Entities   []MessageEntity `json:"text_entities,omitempty"`

	// Bot API 9.6: identifier of the option, kept when options are added and
	// deleted.
	PersistentID string `json:"persistent_id,omitempty"`

	// Bot API 9.6: who added the option after the poll was created, and when.
	AddedByUser  *User `json:"added_by_user,omitempty"`
	AddedByChat  *Chat `json:"added_by_chat,omitempty"`
	AdditionDate int64 `json:"addition_date,omitempty"`

	// Bot API 10.0: media of the option, as received.
	Media *PollMedia `json:"media,omitempty"`

	// Bot API 10.0: media of the option, to send.
	InputMedia InputPollMedia `json:"-"`
}

// PollAnswer represents an answer of a user in a non-anonymous poll.
//
// Bot API: PollAnswer.
type PollAnswer struct {
	PollID  string `json:"poll_id"`
	Sender  *User  `json:"user"`
	Chat    *Chat  `json:"voter_chat"`
	Options []int  `json:"option_ids"`

	// Bot API 9.6: persistent identifiers of the chosen options.
	OptionPersistentIDs []string `json:"option_persistent_ids,omitempty"`
}

// IsRegular says whether poll is a regular.
func (p *Poll) IsRegular() bool {
	return p.Type == PollRegular
}

// IsQuiz says whether poll is a quiz.
func (p *Poll) IsQuiz() bool {
	return p.Type == PollQuiz
}

// CloseDate returns the close date of poll in local time.
func (p *Poll) CloseDate() time.Time {
	return time.Unix(p.CloseUnixdate, 0)
}

// AddOptions adds text options to the poll.
func (p *Poll) AddOptions(opts ...string) {
	for _, t := range opts {
		p.Options = append(p.Options, PollOption{Text: t})
	}
}

// PollOptionAdded is a service message about an option added to a poll
// (Bot API 9.6).
//
// Bot API: PollOptionAdded.
type PollOptionAdded struct {
	// (Optional) The message with the poll, if known.
	PollMessage *Message `json:"poll_message,omitempty"`

	OptionPersistentID string          `json:"option_persistent_id"`
	OptionText         string          `json:"option_text"`
	OptionTextEntities []MessageEntity `json:"option_text_entities,omitempty"`
}

// PollOptionDeleted is a service message about an option deleted from a poll
// (Bot API 9.6).
//
// Bot API: PollOptionDeleted.
type PollOptionDeleted struct {
	// (Optional) The message with the poll, if known.
	PollMessage *Message `json:"poll_message,omitempty"`

	OptionPersistentID string          `json:"option_persistent_id"`
	OptionText         string          `json:"option_text"`
	OptionTextEntities []MessageEntity `json:"option_text_entities,omitempty"`
}

// encodePollMedia returns m as an InputMedia* object, adding the files it
// uploads to files under names starting with name.
func encodePollMedia(m InputPollMedia, name string, files map[string]File) (interface{}, error) {
	switch v := m.(type) {
	case *Location:
		w := richWire{"type": "location", "latitude": v.Lat, "longitude": v.Lng}
		if v.HorizontalAccuracy != nil {
			w["horizontal_accuracy"] = *v.HorizontalAccuracy
		}
		return w, nil
	case *Venue:
		w := richWire{
			"type": "venue", "latitude": v.Location.Lat, "longitude": v.Location.Lng,
			"title": v.Title, "address": v.Address,
		}
		w.set("foursquare_id", v.FoursquareID)
		w.set("foursquare_type", v.FoursquareType)
		w.set("google_place_id", v.GooglePlaceID)
		w.set("google_place_type", v.GooglePlaceType)
		return w, nil
	case *Link:
		return richWire{"type": "link", "url": v.URL}, nil
	case *Sticker:
		repr := v.File.process(name, files)
		if repr == "" {
			return nil, errors.New("telebot: poll sticker does not exist")
		}
		w := richWire{"type": "sticker", "media": repr}
		w.set("emoji", v.Emoji)
		return w, nil
	case Inputtable:
		repr := v.MediaFile().process(name, files)
		if repr == "" {
			return nil, errors.New("telebot: poll media does not exist")
		}
		im := v.InputMedia()
		im.Media = repr
		if err := embedStill(v, &im, name+"s", files); err != nil {
			return nil, err
		}
		return im, nil
	default:
		return nil, errors.New("telebot: unsupported poll media")
	}
}

// encodePollOptions returns the options of p as InputPollOption objects.
func encodePollOptions(options []PollOption, files map[string]File) (string, error) {
	out := make([]richWire, len(options))
	for i, o := range options {
		w := richWire{"text": o.Text}
		w.set("text_parse_mode", o.ParseMode)
		if len(o.Entities) > 0 {
			w["text_entities"] = o.Entities
		}
		if o.InputMedia != nil {
			media, err := encodePollMedia(o.InputMedia, "option"+strconv.Itoa(i), files)
			if err != nil {
				return "", err
			}
			w["media"] = media
		}
		out[i] = w
	}
	data, err := json.Marshal(out)
	return string(data), err
}
