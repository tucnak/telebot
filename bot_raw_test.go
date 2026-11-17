package telebot

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v4/teletest"

	"github.com/stretchr/testify/assert"
)

// testPayload implements json.Marshaler
// to test json encoding error behaviour.
type testPayload struct{}

func (testPayload) MarshalJSON() ([]byte, error) {
	return nil, errors.New("test error")
}

func testRawServer(w http.ResponseWriter, r *http.Request) {
	switch {
	// causes EOF error on ioutil.ReadAll
	case strings.HasSuffix(r.URL.Path, "/testReadError"):
		// tells the body is 1 byte length but actually it's 0
		w.Header().Set("Content-Length", "1")

	// returns unknown telegram error
	case strings.HasSuffix(r.URL.Path, "/testUnknownError"):
		data, _ := json.Marshal(struct {
			Ok          bool   `json:"ok"`
			Code        int    `json:"error_code"`
			Description string `json:"description"`
		}{
			Ok:          false,
			Code:        400,
			Description: "unknown error",
		})

		w.WriteHeader(400)
		w.Write(data)
	}
}

func TestRaw(t *testing.T) {
	if token == "" {
		t.Skip("TELEBOT_SECRET is required")
	}

	b, err := newTestBot()
	if err != nil {
		t.Fatal(err)
	}

	_, err = b.Raw("BAD METHOD", nil)
	assert.EqualError(t, err, ErrNotFound.Error())

	_, err = b.Raw("", &testPayload{})
	assert.Error(t, err)

	srv := httptest.NewServer(http.HandlerFunc(testRawServer))
	defer srv.Close()

	b.URL = srv.URL
	b.client = srv.Client()

	_, err = b.Raw("testReadError", nil)
	assert.EqualError(t, err, "telebot: "+io.ErrUnexpectedEOF.Error())

	_, err = b.Raw("testUnknownError", nil)
	assert.EqualError(t, err, "telegram: unknown error (400)")
}

func TestExtractOk(t *testing.T) {
	data := []byte(`{"ok": true, "result": {}}`)
	require.NoError(t, extractOk(data))

	data = []byte(`{
		"ok": false,
		"error_code": 400,
		"description": "Bad Request: reply message not found"
	}`)
	assert.EqualError(t, extractOk(data), ErrNotFoundToReply.Error())

	data = []byte(`{
		"ok": false,
		"error_code": 429,
		"description": "Too Many Requests: retry after 8",
		"parameters": {"retry_after": 8}
	}`)
	assert.Equal(t, FloodError{
		err:        NewError(429, "Too Many Requests: retry after 8"),
		RetryAfter: 8,
	}, extractOk(data))

	data = []byte(`{
		"ok": false,
		"error_code": 400,
		"description": "Bad Request: group chat was upgraded to a supergroup chat",
		"parameters": {"migrate_to_chat_id": -100123456789}
	}`)
	assert.Equal(t, GroupError{
		err:        ErrGroupMigrated,
		MigratedTo: -100123456789,
	}, extractOk(data))
}

func TestExtractMessage(t *testing.T) {
	data := []byte(`{"ok":true,"result":true}`)
	_, err := extractMessage(data)
	assert.Equal(t, ErrTrueResult, err)

	data = []byte(`{"ok":true,"result":{"foo":"bar"}}`)
	_, err = extractMessage(data)
	require.NoError(t, err)
}

func TestRawFiles(t *testing.T) {
	b, srv := newFakeBot(t)

	_, err := b.RawFiles("setChatPhoto", map[string]string{"chat_id": "-100"},
		map[string]File{"photo": FromReader(strings.NewReader("jpg"))})
	require.NoError(t, err)
	req := srv.Last()
	assert.Equal(t, "setChatPhoto", req.Method)
	assert.Equal(t, "-100", req.Params["chat_id"])
	assert.Equal(t, "jpg", string(req.Files["photo"].Data))

	_, err = b.RawFiles("sendDocument", nil, map[string]File{"document": {FileID: "abc"}})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"document": "abc"}, srv.Last().Params)
	assert.Empty(t, srv.Last().Files)
}

func TestRawFilesErrorHidesToken(t *testing.T) {
	srv := teletest.NewServer(t)
	url := srv.URL
	b, err := NewBot(Settings{URL: url, Token: srv.Token, Offline: true})
	require.NoError(t, err)

	// A server that is gone makes the transport fail with the request URL.
	b.URL = "http://127.0.0.1:1"
	_, err = b.RawFiles("sendDocument", nil, map[string]File{"document": FromReader(strings.NewReader("x"))})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "127.0.0.1:1", "the error still describes the request")
	assert.NotContains(t, err.Error(), srv.Token)
}
