package telebot

import (
	"encoding/json"
	"fmt"
)

type ChatInfo struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Username string `json:"username"`
	Type     string `json:"type"`
}

type ChannelUser struct {
	User   User   `json:"user"`
	Status string `json:"status"`
}

func (b *Bot) GetChat(recipient Recipient) (ChatInfo, error) {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}

	responseJSON, err := sendCommand("getChat", b.Token, params)
	if err != nil {
		return ChatInfo{}, err
	}

	var responseRecieved struct {
		Ok     bool     `json:"ok"`
		Result ChatInfo `json:"result"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return ChatInfo{}, err
	}

	if !responseRecieved.Ok {
		return ChatInfo{}, fmt.Errorf("telebot: %s", responseRecieved.Result)
	}

	return responseRecieved.Result, nil
}

func (b *Bot) GetChatAdministrators(recipient Recipient) ([]ChannelUser, error) {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}

	responseJSON, err := sendCommand("getChatAdministrators", b.Token, params)
	if err != nil {
		return []ChannelUser{}, err
	}

	var responseRecieved struct {
		Ok     bool          `json:"ok"`
		Result []ChannelUser `json:"result"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return []ChannelUser{}, err
	}

	if !responseRecieved.Ok {
		return []ChannelUser{}, fmt.Errorf("telebot: %s", responseRecieved.Result)
	}

	return responseRecieved.Result, nil
}

func (b *Bot) GetChatMembersCount(recipient Recipient) (int, error) {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}

	responseJSON, err := sendCommand("getChatMembersCount", b.Token, params)
	if err != nil {
		return 0, err
	}

	var responseRecieved struct {
		Ok     bool `json:"ok"`
		Result int  `json:"result"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return 0, err
	}

	if !responseRecieved.Ok {
		return 0, fmt.Errorf("telebot: %s", responseRecieved.Result)
	}

	return responseRecieved.Result, nil
}
