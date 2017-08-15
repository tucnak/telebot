package telebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func sendCommand(method, token string, payload interface{}) ([]byte, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, method)

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(payload); err != nil {
		return []byte{}, err
	}

	resp, err := http.Post(url, "application/json", &b)
	if err != nil {
		return []byte{}, err
	}
	resp.Close = true
	defer resp.Body.Close()
	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return json, nil
}

func sendFile(method, token, name, path string, params map[string]string) ([]byte, error) {
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

	for field, value := range params {
		writer.WriteField(field, value)
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

	if resp.StatusCode == http.StatusInternalServerError {
		return []byte{}, fmt.Errorf("telegram: internal server error")
	}

	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return json, nil
}

func embedSendOptions(params map[string]string, options *SendOptions) {
	if options == nil {
		return
	}

	if options.ReplyTo.ID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(options.ReplyTo.ID)
	}

	if options.DisableWebPagePreview {
		params["disable_web_page_preview"] = "true"
	}

	if options.DisableNotification {
		params["disable_notification"] = "true"
	}

	if options.ParseMode != ModeDefault {
		params["parse_mode"] = string(options.ParseMode)
	}

	// Processing force_reply:
	{
		forceReply := options.ReplyMarkup.ForceReply
		customKeyboard := (options.ReplyMarkup.CustomKeyboard != nil)
		inlineKeyboard := options.ReplyMarkup.InlineKeyboard != nil
		hiddenKeyboard := options.ReplyMarkup.HideCustomKeyboard
		if forceReply || customKeyboard || hiddenKeyboard || inlineKeyboard {
			replyMarkup, _ := json.Marshal(options.ReplyMarkup)
			params["reply_markup"] = string(replyMarkup)
		}
	}
}

func getMe(token string) (User, error) {
	meJSON, err := sendCommand("getMe", token, nil)
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
		return User{}, fmt.Errorf("telebot: invalid token")
	}

	if botInfo.Ok {
		return botInfo.Result, nil
	}

	return User{}, fmt.Errorf("telebot: %s", botInfo.Description)
}

func getUpdates(token string, offset, timeout time.Duration) (upd []Update, err error) {
	params := map[string]string{
		"offset":  strconv.Itoa(offset),
		"timeout": strconv.Itoa(timeout / time.Second),
	}
	updatesJSON, err := sendCommand("getUpdates", token, params)
	if err != nil {
		return
	}

	var updatesRecieved struct {
		Ok          bool
		Result      []Update
		Description string
	}

	err = json.Unmarshal(updatesJSON, &updatesRecieved)
	if err != nil {
		return
	}

	if !updatesRecieved.Ok {
		err = fmt.Errorf("telebot: %s", updatesRecieved.Description)
		return
	}

	return updatesRecieved.Result, nil
}
