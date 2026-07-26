package telebot

import (
	"encoding/json"
	"net/http"
)

// Telegram lets a bot answer the webhook request itself: instead of issuing a
// separate HTTPS call, the method and its parameters are written into the body
// of the reply to Telegram's POST. That saves one round trip per update.
//
// The mechanism is deliberately limited:
//
//   - Exactly one method fits in a webhook reply. Anything a handler sends
//     afterwards goes out over the regular API, transparently.
//   - The result is not readable. Telegram's documentation is explicit: "It's
//     not possible to know that such a request was successful or get its
//     result." Send therefore reports no *Message for the reply it produced.
//   - Files cannot be uploaded this way, so media with a local file always
//     takes the regular path.
//
// Consuming the reply is up to the HTTP handler that received the update:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//		body, _ := io.ReadAll(r.Body)
//		var u telebot.Update
//		json.Unmarshal(body, &u)
//		reply := bot.ProcessUpdateReply(u, body)
//		reply.WriteTo(w)
//	}
//
// The bundled Webhook poller does not use it: it forwards updates over a
// channel and answers Telegram before handlers have run.

// WebhookReply is the Bot API call a handler produced to be answered inside
// the webhook response. A nil reply means there is nothing to write, which is
// the normal outcome when the handler sent nothing or the bot is polling.
type WebhookReply struct {
	params map[string]string
	method string
}

// Method reports the Bot API method to be invoked, or "" for a nil reply.
func (r *WebhookReply) Method() string {
	if r == nil {
		return ""
	}
	return r.method
}

// Payload renders the reply as the JSON object Telegram expects. It returns
// nil for a nil reply.
func (r *WebhookReply) Payload() map[string]interface{} {
	if r == nil {
		return nil
	}

	payload := make(map[string]interface{}, len(r.params)+1)
	for k, v := range r.params {
		// Nested parameters travel as JSON strings, but the webhook reply is
		// itself JSON, so they are unwrapped back into real objects.
		if len(v) > 0 && (v[0] == '{' || v[0] == '[') {
			var nested interface{}
			if json.Unmarshal([]byte(v), &nested) == nil {
				payload[k] = nested
				continue
			}
		}
		payload[k] = v
	}
	payload["method"] = r.method

	return payload
}

// WriteTo answers the webhook request with the reply. A nil reply writes an
// empty 200, which tells Telegram the update was handled and nothing is to
// be sent. Any non-2xx response would make Telegram redeliver the update.
func (r *WebhookReply) WriteTo(w http.ResponseWriter) error {
	if r == nil {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(r.Payload())
}

// takeWebhookReply claims the single webhook-reply slot for the given call.
// It reports false once the slot is taken, so every later call falls through
// to a regular API request.
func (c *nativeContext) takeWebhookReply(method string, params map[string]string) bool {
	if !c.webhookReplies {
		return false
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.reply != nil {
		return false
	}

	// The params map is reused by the caller, so it is copied.
	cp := make(map[string]string, len(params))
	for k, v := range params {
		cp[k] = v
	}
	c.reply = &WebhookReply{method: method, params: cp}

	return true
}

// webhookReply returns the claimed reply, if any.
func (c *nativeContext) webhookReply() *WebhookReply {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.reply
}
