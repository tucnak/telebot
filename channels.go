package telebot

import (
	"encoding/json"
	"fmt"
)

// ChatInfo struct is the main chat information from telegram
type ChatInfo struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Username string `json:"username"`
	Type     string `json:"type"`
}

// For each member of a chat , we have some information like usertype (admin , etc..) 
// and the main information about user like id and etc...
type ChannelUser struct {
	User   User   `json:"user"`
	Status string `json:"status"`
}

// GetChat allows to get main information of chat like channel username and etc...
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

// GetChatAdministrators return an array of ChannelUser struct which have two types, administrator and creator
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

// GetChatMembersCount... are you kidding me ? this method return number of users in a channel or group and etc...
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

// The bot will leave the chat if this method has been called.
func (b *Bot) LeaveChat(recipient Recipient) error {
	params := map[string]string{
		"chat_id": recipient.Destination(),
	}

	responseJSON, err := sendCommand("getChatMembersCount", b.Token, params)
	if err != nil {
		return err
	}

	var responseRecieved struct {
		Ok     bool `json:"ok"`
		Result bool `json:"result"`
	}

	err = json.Unmarshal(responseJSON, &responseRecieved)
	if err != nil {
		return err
	}

	if !responseRecieved.Ok {
		return fmt.Errorf("telebot: %s", responseRecieved.Result)
	}

	return nil
}

