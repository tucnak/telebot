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

// Raw lets you call any method of Bot API manually.
func (b *Bot) Raw(method string, payload interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/bot%s/%s", b.URL, b.Token, method)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return []byte{}, wrapSystem(err)
	}

	resp, err := b.client.Post(url, "application/json", &buf)
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

func addFileToWriter(writer *multipart.Writer, fieldName string, file interface{}) error {
	var reader io.Reader
	var part io.Writer
	var err error
	var filename string

	if r, ok := file.(io.Reader); ok {
		// Telegram requires fields to have a filename, otherwise
		// `Bad Request: wrong URL host` would be returned
		filename = "empty"
		reader = r
	} else if path, ok := file.(string); ok {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		reader = f
		filename = filepath.Base(path)
	} else {
		return errors.Errorf("File for field `%v` should be an io.ReadCloser or string", fieldName)
	}

	part, err = writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, reader)
	return err
}

func (b *Bot) sendFiles(
	method string,
	files map[string]File,
	params map[string]string) ([]byte, error) {
	// ---
	body := &bytes.Buffer{}
	rawFiles := map[string]interface{}{}

	for name, f := range files {
		switch {
		case f.InCloud():
			params[name] = f.FileID
		case f.FileURL != "":
			params[name] = f.FileURL
		case f.OnDisk():
			rawFiles[name] = f.FileLocal
		case f.FileReader != nil:
			rawFiles[name] = f.FileReader
		default:
			return nil, errors.Errorf("sendFiles: File for field %s doesn't exist", name)
		}
	}

	if len(rawFiles) == 0 {
		return b.Raw(method, params)
	}

	writer := multipart.NewWriter(body)

	for field, file := range rawFiles {
		if err := addFileToWriter(writer, field, file); err != nil {
			return nil, wrapSystem(err)
		}
	}

	for field, value := range params {
		if err := writer.WriteField(field, value); err != nil {
			return nil, wrapSystem(err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, wrapSystem(err)
	}

	url := fmt.Sprintf("%s/bot%s/%s", b.URL, b.Token, method)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, wrapSystem(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http.Post failed")
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, errors.New("api error: internal server error")
	}

	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, wrapSystem(err)
	}

	return json, nil
}

func (b *Bot) sendObject(f *File, what string, params map[string]string) (*Message, error) {
	sendWhat := "send" + strings.Title(what)

	if what == "videoNote" {
		what = "video_note"
	}

	respJSON, err := b.sendFiles(sendWhat, map[string]File{what: *f}, params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}

func (b *Bot) getMe() (*User, error) {
	meJSON, err := b.Raw("getMe", nil)
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
	updatesJSON, errCommand := b.Raw("getUpdates", params)
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
