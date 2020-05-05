package telebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChat(t *testing.T) {
	user := &User{ID: 1}
	chat := &Chat{ID: 1}
	chatID := ChatID(1)

	assert.Implements(t, (*Recipient)(nil), user)
	assert.Implements(t, (*Recipient)(nil), chat)
	assert.Implements(t, (*Recipient)(nil), chatID)

	assert.Equal(t, "1", user.Recipient())
	assert.Equal(t, "1", chat.Recipient())
	assert.Equal(t, "1", chatID.Recipient())
}
