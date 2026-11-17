package telebot

import (
	"encoding/json"
	"strconv"
)

// Topic is a forum topic.
//
// Bot API: ForumTopic, ForumTopicCreated, ForumTopicEdited.
type Topic struct {
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
	ThreadID          int    `json:"message_thread_id"`

	// Bot API 9.3: (Optional) The creator didn't name the topic, so its name
	// likely needs changing.
	IsNameImplicit bool `json:"is_name_implicit,omitempty"`
}

// ForumTopicClosed is a service message about a forum topic closed in the
// chat. It holds no information.
//
// Bot API: ForumTopicClosed.
type ForumTopicClosed struct{}

// GeneralForumTopicHidden is a service message about the General forum topic
// hidden in the chat. It holds no information.
//
// Bot API: GeneralForumTopicHidden.
type GeneralForumTopicHidden struct{}

// GeneralForumTopicUnhidden is a service message about the General forum
// topic unhidden in the chat. It holds no information.
//
// Bot API: GeneralForumTopicUnhidden.
type GeneralForumTopicUnhidden struct{}

// CreateTopic creates a topic in a forum supergroup chat.
//
// Bot API: createForumTopic.
func (b *Bot) CreateTopic(chat *Chat, topic *Topic) (*Topic, error) {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"name":    topic.Name,
	}

	if topic.IconColor != 0 {
		params["icon_color"] = strconv.Itoa(topic.IconColor)
	}
	if topic.IconCustomEmojiID != "" {
		params["icon_custom_emoji_id"] = topic.IconCustomEmojiID
	}

	data, err := b.Raw("createForumTopic", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *Topic
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, err
}

// EditTopic edits name and icon of a topic in a forum supergroup chat.
//
// Bot API: editForumTopic.
func (b *Bot) EditTopic(chat *Chat, topic *Topic) error {
	params := map[string]interface{}{
		"chat_id":           chat.Recipient(),
		"message_thread_id": topic.ThreadID,
	}

	if topic.Name != "" {
		params["name"] = topic.Name
	}
	if topic.IconCustomEmojiID != "" {
		params["icon_custom_emoji_id"] = topic.IconCustomEmojiID
	}

	_, err := b.Raw("editForumTopic", params)
	return err
}

// CloseTopic closes an open topic in a forum supergroup chat.
//
// Bot API: closeForumTopic.
func (b *Bot) CloseTopic(chat *Chat, topic *Topic) error {
	params := map[string]interface{}{
		"chat_id":           chat.Recipient(),
		"message_thread_id": topic.ThreadID,
	}

	_, err := b.Raw("closeForumTopic", params)
	return err
}

// ReopenTopic reopens a closed topic in a forum supergroup chat.
//
// Bot API: reopenForumTopic.
func (b *Bot) ReopenTopic(chat *Chat, topic *Topic) error {
	params := map[string]interface{}{
		"chat_id":           chat.Recipient(),
		"message_thread_id": topic.ThreadID,
	}

	_, err := b.Raw("reopenForumTopic", params)
	return err
}

// DeleteTopic deletes a forum topic along with all its messages in a forum supergroup chat.
//
// Bot API: deleteForumTopic.
func (b *Bot) DeleteTopic(chat *Chat, topic *Topic) error {
	params := map[string]interface{}{
		"chat_id":           chat.Recipient(),
		"message_thread_id": topic.ThreadID,
	}

	_, err := b.Raw("deleteForumTopic", params)
	return err
}

// UnpinAllTopicMessages clears the list of pinned messages in a forum topic. The bot must be an administrator in the chat for this to work and must have the can_pin_messages administrator right in the supergroup.
//
// Bot API: unpinAllForumTopicMessages.
func (b *Bot) UnpinAllTopicMessages(chat *Chat, topic *Topic) error {
	params := map[string]interface{}{
		"chat_id":           chat.Recipient(),
		"message_thread_id": topic.ThreadID,
	}

	_, err := b.Raw("unpinAllForumTopicMessages", params)
	return err
}

// TopicIconStickers gets custom emoji stickers, which can be used as a forum topic icon by any user.
//
// Bot API: getForumTopicIconStickers.
func (b *Bot) TopicIconStickers() ([]Sticker, error) {
	params := map[string]string{}

	data, err := b.Raw("getForumTopicIconStickers", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result []Sticker
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// EditGeneralTopic edits name of the 'General' topic in a forum supergroup chat.
//
// Bot API: editGeneralForumTopic.
func (b *Bot) EditGeneralTopic(chat *Chat, topic *Topic) error {
	params := map[string]interface{}{
		"chat_id": chat.Recipient(),
		"name":    topic.Name,
	}

	_, err := b.Raw("editGeneralForumTopic", params)
	return err
}

// CloseGeneralTopic closes an open 'General' topic in a forum supergroup chat.
//
// Bot API: closeGeneralForumTopic.
func (b *Bot) CloseGeneralTopic(chat *Chat) error {
	params := map[string]interface{}{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("closeGeneralForumTopic", params)
	return err
}

// ReopenGeneralTopic reopens a closed 'General' topic in a forum supergroup chat.
//
// Bot API: reopenGeneralForumTopic.
func (b *Bot) ReopenGeneralTopic(chat *Chat) error {
	params := map[string]interface{}{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("reopenGeneralForumTopic", params)
	return err
}

// HideGeneralTopic hides the 'General' topic in a forum supergroup chat.
//
// Bot API: hideGeneralForumTopic.
func (b *Bot) HideGeneralTopic(chat *Chat) error {
	params := map[string]interface{}{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("hideGeneralForumTopic", params)
	return err
}

// UnhideGeneralTopic unhides the 'General' topic in a forum supergroup chat.
//
// Bot API: unhideGeneralForumTopic.
func (b *Bot) UnhideGeneralTopic(chat *Chat) error {
	params := map[string]interface{}{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("unhideGeneralForumTopic", params)
	return err
}

// UnpinAllGeneralTopicMessages clears the list of pinned messages in a General forum topic.
// The bot must be an administrator in the chat for this to work and must have the
// can_pin_messages administrator right in the supergroup.
//
// Bot API: unpinAllGeneralForumTopicMessages.
func (b *Bot) UnpinAllGeneralTopicMessages(chat *Chat) error {
	params := map[string]interface{}{
		"chat_id": chat.Recipient(),
	}

	_, err := b.Raw("unpinAllGeneralForumTopicMessages", params)
	return err
}
