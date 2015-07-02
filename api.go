package telebot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func performApiCall(method string, token string, params url.Values) ([]byte, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s?%s",
		token, method, params.Encode())

	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return json, nil
}

func api_getMe(token string) (User, error) {
	me_json, err := performApiCall("getMe", token, url.Values{})
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
	params := url.Values{}
	params.Set("offset", strconv.FormatInt(int64(offset), 10))
	updates_json, err := performApiCall("getUpdates", token, params)
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
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(int64(recipient.Id), 10))
	params.Set("text", text)
	response_json, err := performApiCall("sendMessage", token, params)
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

func api_forwardMessage(token string, recipient User, message Message) error {
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(int64(recipient.Id), 10))
	params.Set("from_chat_id",
		strconv.FormatInt(int64(message.Origin().Id), 10))
	params.Set("message_id", strconv.FormatInt(int64(message.Id), 10))

	response_json, err := performApiCall("forwardMessage", token, params)
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
