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
	"strings"
	"time"

	"github.com/pkg/errors"
)

func (b *Bot) sendCommand(method string, payload interface{}) ([]byte, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.Token, method)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return []byte{}, wrapSystem(err)
	}

	resp, err := http.Post(url, "application/json", &buf)
	if err != nil {
		return []byte{}, errors.Wrap(err, "http.Post failed")
	}
	resp.Close = true
	defer resp.Body.Close()
	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, wrapSystem(err)
	}

	return json, nil
}

func (b *Bot) sendFile(method, name, path string, params map[string]string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, wrapSystem(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return []byte{}, wrapSystem(err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return []byte{}, wrapSystem(err)
	}

	for field, value := range params {
		writer.WriteField(field, value)
	}

	if err = writer.Close(); err != nil {
		return []byte{}, wrapSystem(err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.Token, method)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return []byte{}, wrapSystem(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, errors.Wrap(err, "http.Post failed")
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return []byte{}, errors.New("api error: internal server error")
	}

	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, wrapSystem(err)
	}

	return json, nil
}

func (b *Bot) sendObject(f *File, what string, params map[string]string) (*Message, error) {
	sendWhat := "send" + strings.Title(what)

	if what == "videoNote" {
		what = "video_note"
	}

	var respJSON []byte
	var err error

	if f.InCloud() {
		params[what] = f.FileID
		respJSON, err = b.sendCommand(sendWhat, params)
	} else if f.FileURL != nil {
		params[what] = f.FileURL.String()
		respJSON, err = b.sendCommand(sendWhat, params)
	} else {
		respJSON, err = b.sendFile(sendWhat, what, f.FileLocal, params)
	}

	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}

func (b *Bot) getMe() (*User, error) {
	meJSON, err := b.sendCommand("getMe", nil)
	if err != nil {
		return nil, err
	}

	var botInfo struct {
		Ok          bool
		Result      *User
		Description string
	}

	err = json.Unmarshal(meJSON, &botInfo)
	if err != nil {
		return nil, errors.Wrap(err, "bad response json")
	}

	if !botInfo.Ok {
		return nil, errors.Errorf("api error: %s", botInfo.Description)
	}

	return botInfo.Result, nil

}

func (b *Bot) getUpdates(offset int, timeout time.Duration) (upd []Update, err error) {
	params := map[string]string{
		"offset":  strconv.Itoa(offset),
		"timeout": strconv.Itoa(int(timeout / time.Second)),
	}
	updatesJSON, errCommand := b.sendCommand("getUpdates", params)
	if errCommand != nil {
		err = errCommand
		return

	}
	var updatesReceived struct {
		Ok          bool
		Result      []Update
		Description string
	}

	err = json.Unmarshal(updatesJSON, &updatesReceived)
	if err != nil {
		err = errors.Wrap(err, "bad response json")
		return
	}

	if !updatesReceived.Ok {
		err = errors.Errorf("api error: %s", updatesReceived.Description)
		return
	}

	return updatesReceived.Result, nil
}
