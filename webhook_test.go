package telebot

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhook_Poll(t *testing.T) {
	t.Run("graceful shutdown", func(t *testing.T) {
		webhook := &Webhook{
			Listen:           "127.0.0.1:80",
			IgnoreSetWebhook: true,
		}

		pref := defaultSettings()
		pref.Offline = true
		pref.Poller = webhook

		b, err := NewBot(pref)
		require.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			b.Start()
		}()

		// Give bot time to start
		time.Sleep(100 * time.Millisecond)

		b.Stop()
		wg.Wait()
	})

	t.Run("no Listen", func(t *testing.T) {
		webhook := &Webhook{
			Listen:           "",
			IgnoreSetWebhook: true,
		}

		pref := defaultSettings()
		pref.Offline = true
		pref.Poller = webhook

		b, err := NewBot(pref)
		require.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			b.Start()
		}()

		// Give bot time to start
		time.Sleep(100 * time.Millisecond)

		b.Stop()
		wg.Wait()
	})
}

func TestWebhook_ServeHTTP(t *testing.T) {
	webhook := &Webhook{
		IgnoreSetWebhook: true,
	}

	pref := defaultSettings()
	pref.Offline = true
	b, err := NewBot(pref)
	require.NoError(t, err)

	dest := make(chan Update, 1)
	webhook.dest = dest
	webhook.bot = b

	t.Run("valid update", func(t *testing.T) {
		update := Update{
			ID: 123,
			Message: &Message{
				ID:   1,
				Text: "test message",
				Chat: &Chat{ID: 456},
			},
		}

		body, _ := json.Marshal(update)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		w := httptest.NewRecorder()

		webhook.ServeHTTP(w, req)

		select {
		case received := <-dest:
			assert.Equal(t, update.ID, received.ID)
			assert.Equal(t, update.Message.Text, received.Message.Text)
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for update")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		webhook.ServeHTTP(w, req)

		select {
		case <-dest:
			t.Fatal("should not receive update for invalid JSON")
		case <-time.After(100 * time.Millisecond):
			// Expected - no update should be sent
		}
	})

	t.Run("secret token validation", func(t *testing.T) {
		webhook.SecretToken = "secret123"

		update := Update{ID: 123}
		body, _ := json.Marshal(update)

		// Without secret token
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		w := httptest.NewRecorder()
		webhook.ServeHTTP(w, req)

		select {
		case <-dest:
			t.Fatal("should not receive update without secret token")
		case <-time.After(100 * time.Millisecond):
			// Expected
		}

		// With correct secret token
		req = httptest.NewRequest("POST", "/", bytes.NewReader(body))
		req.Header.Set("X-Telegram-Bot-Api-Secret-Token", "secret123")
		w = httptest.NewRecorder()
		webhook.ServeHTTP(w, req)

		select {
		case received := <-dest:
			assert.Equal(t, update.ID, received.ID)
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for update")
		}

		// Reset for other tests
		webhook.SecretToken = ""
	})
}
