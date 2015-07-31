package telebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

func sendCommand(method string, token string, params url.Values) ([]byte, error) {
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

func sendFile(method, token, name, path string, params url.Values) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return []byte{}, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return []byte{}, err
	}

	for field, values := range params {
		if len(values) > 0 {
			writer.WriteField(field, values[0])
		}
	}

	if err = writer.Close(); err != nil {
		return []byte{}, err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, method)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return json, nil
}

func embedSendOptions(params *url.Values, options *SendOptions) {
	if params == nil || options == nil {
		return
	}

	if options.ReplyTo.ID != 0 {
		params.Set("reply_to_message_id", strconv.Itoa(options.ReplyTo.ID))
	}

	if options.DisableWebPagePreview {
		params.Set("disable_web_page_preview", "true")
	}

	// process reply_markup
	if options.ReplyMarkup.ForceReply || options.ReplyMarkup.CustomKeyboard != nil || options.ReplyMarkup.HideCustomKeyboard {
		replyMarkup, _ := json.Marshal(options.ReplyMarkup)
		params.Set("reply_markup", string(replyMarkup))
	}
}

func getMe(token string) (User, error) {
	meJSON, err := sendCommand("getMe", token, url.Values{})
	if err != nil {
		return User{}, err
	}

	var botInfo struct {
		Ok          bool
		Result      User
		Description string
	}

	err = json.Unmarshal(meJSON, &botInfo)
	if err != nil {
		return User{}, AuthError{"invalid token"}
	}

	if botInfo.Ok {
		return botInfo.Result, nil
	}

	return User{}, AuthError{botInfo.Description}
}

func getUpdates(token string, offset int, updates chan<- Update) error {
	params := url.Values{}
	params.Set("offset", strconv.Itoa(offset))
	updatesJSON, err := sendCommand("getUpdates", token, params)
	if err != nil {
		return err
	}

	var updatesRecieved struct {
		Ok          bool
		Result      []Update
		Description string
	}

	err = json.Unmarshal(updatesJSON, &updatesRecieved)
	if err != nil {
		return err
	}

	if !updatesRecieved.Ok {
		return FetchError{updatesRecieved.Description}
	}

	for _, update := range updatesRecieved.Result {
		updates <- update
	}

	return nil
}
