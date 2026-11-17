package telebot

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v4/teletest"
)

// newFakeBot returns an offline bot wired to a fake Bot API server.
func newFakeBot(t *testing.T) (*Bot, *teletest.Server) {
	t.Helper()
	srv := teletest.NewServer(t)
	b, err := NewBot(Settings{
		URL:         srv.URL,
		Token:       srv.Token,
		Client:      srv.Client(),
		Offline:     true,
		Synchronous: true,
	})
	require.NoError(t, err)
	return b, srv
}

// Bot implements every method of API.
var _ API = (*Bot)(nil)
