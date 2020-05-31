package telebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractOk(t *testing.T) {
	data := []byte(`{"ok":true,"result":{"foo":"bar"}}`)
	assert.NoError(t, extractOk(data))

	data = []byte(`{"ok":false,"error_code":429,"description":"Too Many Requests: retry after 8","parameters":{"retry_after":8}}`)
	assert.Error(t, extractOk(data))

	data = []byte(`{"ok":false,"error_code":400,"description":"Bad Request: reply message not found"}`)
	assert.EqualError(t, extractOk(data), ErrToReplyNotFound.Error())
}

func TestExtractMessage(t *testing.T) {
	data := []byte(`{"ok":true,"result":true}`)
	_, err := extractMessage(data)
	assert.Equal(t, ErrTrueResult, err)

	data = []byte(`{"ok":true,"result":{"foo":"bar"}}`)
	_, err = extractMessage(data)
	assert.NoError(t, err)
}
