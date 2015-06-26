package telebot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func api_getMe(token string) (User, error) {
	request := "https://api.telegram.org/bot" + token + "/getMe"

	resp, err := http.Get(request)
	if err != nil {
		return User{}, err
	}

	defer resp.Body.Close()
	me_json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return User{}, err
	}

	var bot_info struct {
		Ok          bool
		Result      User
		Description string
	}

	err = json.Unmarshal(me_json, &bot_info)
	if err != nil {
		return User{}, err
	}

	if bot_info.Ok {
		return bot_info.Result, nil
	} else {
		return User{}, AuthError{bot_info.Description}
	}
}

func api_getUpdates(token string, offset int, updates chan<- Update) error {
	command := fmt.Sprintf("getUpdates?offset=%d", offset)
	request := "https://api.telegram.org/bot" + token + "/" + command

	resp, err := http.Get(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	updates_json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var updates_recieved struct {
		Ok          bool
		Result      []Update
		Description string
	}

	err = json.Unmarshal(updates_json, &updates_recieved)
	if err != nil {
		return err
	}

	if !updates_recieved.Ok {
		return FetchError{updates_recieved.Description}
	}

	for _, update := range updates_recieved.Result {
		updates <- update
	}

	return nil
}

func api_sendMessage(token string, recipient User, text string) error {
	resource := "https://api.telegram.org/bot" + token + "/sendMessage"

	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(int64(recipient.Id), 10))
	params.Set("text", text)

	request := resource + "?" + params.Encode()

	resp, err := http.Get(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	response_json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response_recieved struct {
		Ok          bool
		Description string
	}

	err = json.Unmarshal(response_json, &response_recieved)
	if err != nil {
		return err
	}

	if !response_recieved.Ok {
		return SendError{response_recieved.Description}
	}

	return nil
}
