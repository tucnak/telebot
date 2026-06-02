package telebot

import "encoding/json"

// SentGuestMessage describes an inline message sent by a guest bot.
type SentGuestMessage struct {
	InlineMessageID string `json:"inline_message_id"`
}

// AnswerGuest sends a response for a given guest message.
func (b *Bot) AnswerGuest(msg *Message, result Result) (*SentGuestMessage, error) {
	result.Process(b)
	if err := inferIQR(result); err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"guest_query_id": msg.GuestQueryID,
		"result":         result,
	}

	data, err := b.Raw("answerGuestQuery", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *SentGuestMessage
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}
