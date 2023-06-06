package telebot

import (
	"encoding/json"
	"strconv"
)

type Topic struct {
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
	MessageThreadID   int    `json:"message_thread_id"`
}

type TopicCreated struct {
	Topic
}

type TopicClosed struct{}

type TopicDeleted struct {
	Name              string `json:"name"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
}

type TopicReopened struct {
	Topic
}

// CreateTopic creates a topic in a forum supergroup chat.
func (b *Bot) CreateTopic(chat *Chat, forum *Topic) error {
	params := map[string]string{
		"chat_id": chat.Recipient(),
		"name":    forum.Name,
	}

	if forum.IconColor != 0 {
		params["icon_color"] = strconv.Itoa(forum.IconColor)
	}
	if forum.IconCustomEmojiID != "" {
		params["icon_custom_emoji_id"] = forum.IconCustomEmojiID
	}

	_, err := b.Raw("createForumTopic", params)
	return err
}

// EditTopic edits name and icon of a topic in a forum supergroup chat.
func (b *Bot) EditTopic(chat *Chat, forum *Topic) error {
	params := map[string]interface{}{
		"chat_id":           chat.Recipient(),
		"message_thread_id": forum.MessageThreadID,
	}

	if forum.Name != "" {
		params["name"] = forum.Name
	}
	if forum.IconCustomEmojiID != "" {
		params["icon_custom_emoji_id"] = forum.IconCustomEmojiID
	}

	_, err := b.Raw("editForumTopic", params)
	return err
}

// CloseTopic closes an open topic in a forum supergroup chat.
func (b *Bot) CloseTopic(chat *Chat, forum *Topic) error {
	params := map[string]interface{}{
		"chat_id":           chat.Recipient(),
		"message_thread_id": forum.MessageThreadID,
	}

	_, err := b.Raw("closeForumTopic", params)
	return err
}

// ReopenTopic reopens a closed topic in a forum supergroup chat.
func (b *Bot) ReopenTopic(chat *Chat, forum *Topic) error {
	params := map[string]interface{}{
		"chat_id":           chat.Recipient(),
		"message_thread_id": forum.MessageThreadID,
	}

	_, err := b.Raw("reopenForumTopic", params)
	return err
}

// DeleteTopic deletes a forum topic along with all its messages in a forum supergroup chat.
func (b *Bot) DeleteTopic(chat *Chat, forum *Topic) error {
	params := map[string]interface{}{
		"chat_id":           chat.Recipient(),
		"message_thread_id": forum.MessageThreadID,
	}

	_, err := b.Raw("deleteForumTopic", params)
	return err
}

// UnpinAllTopicMessages clears the list of pinned messages in a forum topic. The bot must be an administrator in the chat for this to work and must have the can_pin_messages administrator right in the supergroup.
func (b *Bot) UnpinAllTopicMessages(chat *Chat, forum *Topic) error {
	params := map[string]interface{}{
		"chat_id":           chat.Recipient(),
		"message_thread_id": forum.MessageThreadID,
	}

	_, err := b.Raw("unpinAllForumTopicMessages", params)
	return err
}

// TopicIconStickers gets custom emoji stickers, which can be used as a forum topic icon by any user
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
