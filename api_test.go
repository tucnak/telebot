package telebot

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	b, err := newTestBot()
	assert.NoError(t, err)

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
