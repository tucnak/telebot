package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tele "github.com/maxbolgarin/telebot/v4"
)

var b, _ = tele.NewBot(context.Background(), tele.Settings{Offline: true})

func TestRecover(t *testing.T) {
	onError := func(err error, c tele.Context) {
		require.Error(t, err, "recover test")
	}

	h := func(c tele.Context) error {
		panic("recover test")
	}

	assert.Panics(t, func() {
		h(nil)
	})

	assert.NotPanics(t, func() {
		Recover(onError)(h)(nil)
	})
}
