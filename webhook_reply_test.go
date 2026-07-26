package telebot

import (
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func replyBot(t *testing.T) *Bot {
	t.Helper()

	b, err := NewBot(Settings{Token: "1:offline", Offline: true, Synchronous: true})
	require.NoError(t, err)
	// Nothing may reach the network; any call that is not answered inside the
	// webhook response has to fail loudly rather than silently succeed.
	b.URL = "http://127.0.0.1:1"

	return b
}

func textUpdate(chatID int64, text string) Update {
	return Update{Message: &Message{Chat: &Chat{ID: chatID}, Text: text}}
}

// The first call is answered in the response, and it carries the parameters
// the handler asked for.
func TestWebhookReplyCapturesFirstCall(t *testing.T) {
	b := replyBot(t)
	b.Handle(OnText, func(c Context) error {
		return c.Send("hello")
	})

	reply := b.ProcessUpdateReply(textUpdate(42, "hi"))
	require.NotNil(t, reply)
	assert.Equal(t, "sendMessage", reply.Method())

	payload := reply.Payload()
	assert.Equal(t, "sendMessage", payload["method"])
	assert.Equal(t, "42", payload["chat_id"])
	assert.Equal(t, "hello", payload["text"])
}

// Only one method fits in a webhook response, so later calls must go out as
// real requests instead of overwriting the first one.
func TestWebhookReplyFallsBackAfterFirstCall(t *testing.T) {
	b := replyBot(t)

	var second error
	b.Handle(OnText, func(c Context) error {
		if err := c.Send("first"); err != nil {
			return err
		}
		second = c.Send("second")
		return nil
	})

	reply := b.ProcessUpdateReply(textUpdate(42, "hi"))
	require.NotNil(t, reply)

	// The first call is the one that got the slot.
	assert.Equal(t, "first", reply.Payload()["text"])
	// The second attempted a request rather than vanishing.
	require.Error(t, second, "the second send must not be silently dropped")
	assert.NotErrorIs(t, second, ErrWebhookReply)
}

// A handler that reports an error still has its own reply preserved: the
// OnError callback sending something must not replace it.
func TestWebhookReplySurvivesOnError(t *testing.T) {
	b, err := NewBot(Settings{
		Token: "1:offline", Offline: true, Synchronous: true,
		OnError: func(err error, c Context) {
			_ = c.Send("error: " + err.Error())
		},
	})
	require.NoError(t, err)
	b.URL = "http://127.0.0.1:1"

	b.Handle(OnText, func(c Context) error {
		_ = c.Send("the real result")
		return fmt.Errorf("boom")
	})

	reply := b.ProcessUpdateReply(textUpdate(42, "hi"))
	require.NotNil(t, reply)
	assert.Equal(t, "the real result", reply.Payload()["text"])
}

// Each update owns its reply, so concurrent updates cannot leak answers to
// one another. This is the regression test for responses being kept on the
// shared Bot.
func TestWebhookReplyIsolatedPerUpdate(t *testing.T) {
	b := replyBot(t)
	b.Handle(OnText, func(c Context) error {
		return c.Send(fmt.Sprintf("reply-to-%d", c.Chat().ID))
	})

	const n = 200

	var wg sync.WaitGroup
	errs := make(chan string, n)

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()

			reply := b.ProcessUpdateReply(textUpdate(id, "hi"))
			if reply == nil {
				errs <- fmt.Sprintf("chat %d: no reply", id)
				return
			}

			payload := reply.Payload()
			want := fmt.Sprintf("reply-to-%d", id)
			if payload["text"] != want || payload["chat_id"] != fmt.Sprint(id) {
				errs <- fmt.Sprintf("chat %d: got %v/%v", id, payload["chat_id"], payload["text"])
			}
		}(int64(i))
	}

	wg.Wait()
	close(errs)

	var found []string
	for e := range errs {
		found = append(found, e)
	}
	assert.Empty(t, found, "replies must never cross between updates")
}

// Media backed by a local file cannot be uploaded in a webhook response, so
// it has to take the regular path.
func TestWebhookReplySkipsFileUploads(t *testing.T) {
	f, err := os.CreateTemp("", "telebot")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	require.NoError(t, f.Close())

	b := replyBot(t)

	var sendErr error
	b.Handle(OnText, func(c Context) error {
		sendErr = c.Send(&Document{File: FromDisk(f.Name()), FileName: "x.txt"})
		return nil
	})

	reply := b.ProcessUpdateReply(textUpdate(42, "hi"))
	assert.Nil(t, reply, "uploads cannot be answered in the response")
	require.Error(t, sendErr, "the upload must be attempted over HTTP")
}

// Without ProcessUpdateReply nothing is intercepted, so polling bots keep
// issuing ordinary requests.
func TestWebhookReplyOptIn(t *testing.T) {
	b := replyBot(t)

	var sendErr error
	b.Handle(OnText, func(c Context) error {
		sendErr = c.Send("hello")
		return nil
	})

	b.ProcessUpdate(textUpdate(42, "hi"))
	require.Error(t, sendErr, "ProcessUpdate must not capture replies")
}

func TestWebhookReplyWriteTo(t *testing.T) {
	b := replyBot(t)
	b.Handle(OnText, func(c Context) error {
		return c.Send("hello", &ReplyMarkup{InlineKeyboard: [][]InlineButton{{{Text: "b", Data: "d"}}}})
	})

	rec := httptest.NewRecorder()
	require.NoError(t, b.ProcessUpdateReply(textUpdate(42, "hi")).WriteTo(rec))

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	body := rec.Body.String()
	assert.Contains(t, body, `"method":"sendMessage"`)
	// reply_markup must be a real object, not a quoted JSON string.
	assert.Contains(t, body, `"reply_markup":{`)
	assert.False(t, strings.Contains(body, `"reply_markup":"`))

	// A nil reply is still a valid, empty acknowledgement.
	rec = httptest.NewRecorder()
	require.NoError(t, (*WebhookReply)(nil).WriteTo(rec))
	assert.Equal(t, 200, rec.Code)
	assert.Empty(t, rec.Body.String())
}
