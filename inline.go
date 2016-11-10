package telebot

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/mitchellh/hashstructure"
)

// inlineQueryHashOptions sets the HashOptions to be used when hashing
// an inline query result (used to generate IDs).
var inlineQueryHashOptions = &hashstructure.HashOptions{
	Hasher: fnv.New64(),
}

// Query is an incoming inline query. When the user sends
// an empty query, your bot could return some default or
// trending results.
type Query struct {
	// Unique identifier for this query.
	ID string `json:"id"`

	// Sender.
	From User `json:"from"`

	// (Optional) Sender location, only for bots that request user location.
	Location Location `json:"location"`

	// Text of the query (up to 512 characters).
	Text string `json:"query"`

	// Offset of the results to be returned, can be controlled by the bot.
	Offset string `json:"offset"`
}

// QueryResponse builds a response to an inline Query.
// See also: https://core.telegram.org/bots/api#answerinlinequery
type QueryResponse struct {
	// The ID of the query to which this is a response.
	// It is not necessary to specify this field manually.
	QueryID string `json:"inline_query_id"`

	// The results for the inline query.
	Results InlineQueryResults `json:"results"`

	// (Optional) The maximum amount of time in seconds that the result
	// of the inline query may be cached on the server.
	CacheTime int `json:"cache_time,omitempty"`

	// (Optional) Pass True, if results may be cached on the server side
	// only for the user that sent the query. By default, results may
	// be returned to any user who sends the same query.
	IsPersonal bool `json:"is_personal"`

	// (Optional) Pass the offset that a client should send in the next
	// query with the same text to receive more results. Pass an empty
	// string if there are no more results or if you don‘t support
	// pagination. Offset length can’t exceed 64 bytes.
	NextOffset string `json:"next_offset"`

	// (Optional) If passed, clients will display a button with specified
	// text that switches the user to a private chat with the bot and sends
	// the bot a start message with the parameter switch_pm_parameter.
	SwitchPMText string `json:"switch_pm_text,omitempty"`

	// (Optional) Parameter for the start message sent to the bot when user
	// presses the switch button.
	SwitchPMParameter string `json:"switch_pm_parameter,omitempty"`
}

// InlineQueryResult represents one result of an inline query.
type InlineQueryResult interface {
	GetID() string
	SetID(string)
}

// InlineQueryResults is a slice wrapper for convenient marshalling.
type InlineQueryResults []InlineQueryResult

// MarshalJSON makes sure IQRs have proper IDs and Type variables set.
//
// If ID of some result appears empty, it gets set to a new hash.
// JSON-specific Type gets infered from the actual (specific) IQR type.
func (results *InlineQueryResults) MarshalJSON() ([]byte, error) {
	for i, result := range *results {
		if result.GetID() == "" {
			hash, err := hashstructure.Hash(result, inlineQueryHashOptions)
			if err != nil {
				return nil, fmt.Errorf("telebot: can't hash IQR #%d: %s",
					i, err)
			}

			result.SetID(strconv.FormatUint(hash, 16))
		}

		if err := inferIQR(result); err != nil {
			return nil, fmt.Errorf("telebot: can't infer type of IQR #%d: %s",
				i, err)
		}
	}

	return json.Marshal([]InlineQueryResult(*results))
}

func inferIQR(result InlineQueryResult) error {
	switch r := result.(type) {
	case *InlineQueryResultArticle:
		r.Type = "article"
	case *InlineQueryResultAudio:
		r.Type = "audio"
	case *InlineQueryResultContact:
		r.Type = "contact"
	case *InlineQueryResultDocument:
		r.Type = "document"
	case *InlineQueryResultGif:
		r.Type = "gif"
	case *InlineQueryResultLocation:
		r.Type = "location"
	case *InlineQueryResultMpeg4Gif:
		r.Type = "mpeg4_gif"
	case *InlineQueryResultPhoto:
		r.Type = "photo"
	case *InlineQueryResultVenue:
		r.Type = "venue"
	case *InlineQueryResultVideo:
		r.Type = "video"
	case *InlineQueryResultVoice:
		r.Type = "voice"
	default:
		return fmt.Errorf("%T is not an IQR", result)
	}

	return nil
}

// Result is a deprecated type, superseded by InlineQueryResult.
type Result interface {
	MarshalJSON() ([]byte, error)
}
