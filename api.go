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

func api_GET(method string, token string, params url.Values) ([]byte, error) {
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

func api_POST(method, token, name, path string, params url.Values) ([]byte, error) {
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

func api_getMe(token string) (User, error) {
	me_json, err := api_GET("getMe", token, url.Values{})
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
	params.Set("offset", strconv.Itoa(offset))
	updates_json, err := api_GET("getUpdates", token, params)
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
	params.Set("chat_id", strconv.Itoa(recipient.Id))
	params.Set("text", text)
	response_json, err := api_GET("sendMessage", token, params)
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
	params.Set("chat_id", strconv.Itoa(recipient.Id))
	params.Set("from_chat_id", strconv.Itoa(message.Origin().Id))
	params.Set("message_id", strconv.Itoa(message.Id))

	response_json, err := api_GET("forwardMessage", token, params)
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

func api_sendPhoto(token string, recipient User, photo *Photo) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Id))
	params.Set("caption", photo.Caption)

	var response_json []byte
	var err error

	if photo.Exists() {
		params.Set("photo", photo.FileId)
		response_json, err = api_GET("sendPhoto", token, params)
	} else {
		response_json, err = api_POST("sendPhoto", token, "photo",
			photo.filename, params)
	}

	if err != nil {
		return err
	}

	var response_recieved struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(response_json, &response_recieved)
	if err != nil {
		return err
	}

	if !response_recieved.Ok {
		return SendError{response_recieved.Description}
	}

	thumbnails := &response_recieved.Result.Photo
	photo.File = (*thumbnails)[len(*thumbnails)-1].File

	return nil
}

func api_sendAudio(token string, recipient User, audio *Audio) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Id))

	var response_json []byte
	var err error

	if audio.Exists() {
		params.Set("audio", audio.FileId)
		response_json, err = api_GET("sendAudio", token, params)
	} else {
		response_json, err = api_POST("sendAudio", token, "audio",
			audio.filename, params)
	}

	if err != nil {
		return err
	}

	var response_recieved struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(response_json, &response_recieved)
	if err != nil {
		return err
	}

	if !response_recieved.Ok {
		return SendError{response_recieved.Description}
	}

	*audio = response_recieved.Result.Audio

	return nil
}

func api_sendDocument(token string, recipient User, doc *Document) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Id))

	var response_json []byte
	var err error

	if doc.Exists() {
		params.Set("document", doc.FileId)
		response_json, err = api_GET("sendDocument", token, params)
	} else {
		response_json, err = api_POST("sendDocument", token, "document",
			doc.filename, params)
	}

	if err != nil {
		return err
	}

	var response_recieved struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(response_json, &response_recieved)
	if err != nil {
		return err
	}

	if !response_recieved.Ok {
		return SendError{response_recieved.Description}
	}

	*doc = response_recieved.Result.Document

	return nil
}

func api_sendSticker(token string, recipient User, sticker *Sticker) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Id))

	var response_json []byte
	var err error

	if sticker.Exists() {
		params.Set("sticker", sticker.FileId)
		response_json, err = api_GET("sendSticker", token, params)
	} else {
		response_json, err = api_POST("sendSticker", token, "sticker",
			sticker.filename, params)
	}

	if err != nil {
		return err
	}

	var response_recieved struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(response_json, &response_recieved)
	if err != nil {
		return err
	}

	if !response_recieved.Ok {
		return SendError{response_recieved.Description}
	}

	*sticker = response_recieved.Result.Sticker

	return nil
}

func api_sendVideo(token string, recipient User, video *Video) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Id))

	var response_json []byte
	var err error

	if video.Exists() {
		params.Set("video", video.FileId)
		response_json, err = api_GET("sendVideo", token, params)
	} else {
		response_json, err = api_POST("sendVideo", token, "video",
			video.filename, params)
	}

	if err != nil {
		return err
	}

	var response_recieved struct {
		Ok          bool
		Result      Message
		Description string
	}

	err = json.Unmarshal(response_json, &response_recieved)
	if err != nil {
		return err
	}

	if !response_recieved.Ok {
		return SendError{response_recieved.Description}
	}

	*video = response_recieved.Result.Video

	return nil
}

func api_sendLocation(token string, recipient User, geo *Location) error {
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(recipient.Id))
	params.Set("latitude", fmt.Sprintf("%f", geo.Latitude))
	params.Set("longitude", fmt.Sprintf("%f", geo.Longitude))

	response_json, err := api_GET("sendLocation", token, params)

	if err != nil {
		return err
	}

	var response_recieved struct {
		Ok          bool
		Result      Message
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
