package telebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRaw(t *testing.T) {
	b, err := newTestBot()
	assert.NoError(t, err)

	_, err = b.Raw("BAD METHOD", nil)
	assert.EqualError(t, err, ErrNotFound.Error())
}
