package telebot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMediaSendsCarryTheirSpoilerAndCaptionAbove(t *testing.T) {
	sent := map[string]map[string]string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseMultipartForm(1<<20))
		fields := map[string]string{}
		for k, v := range r.MultipartForm.Value {
			fields[k] = v[0]
		}
		sent[path.Base(r.URL.Path)] = fields
		w.Write([]byte(`{"ok":true,"result":{"message_id":1,"chat":{"id":2},` +
			`"photo":[{"file_id":"p"}],"video":{"file_id":"v"},"animation":{"file_id":"a"}}}`))
	}))
	defer srv.Close()

	b, err := NewBot(Settings{Offline: true, Token: "token", URL: srv.URL, Client: srv.Client()})
	require.NoError(t, err)

	file := func() File { return FromReader(bytes.NewReader([]byte("media"))) }
	for method, media := range map[string]Sendable{
		"sendPhoto":     &Photo{File: file(), Caption: "c", HasSpoiler: true, CaptionAbove: true},
		"sendVideo":     &Video{File: file(), Caption: "c", HasSpoiler: true, CaptionAbove: true},
		"sendAnimation": &Animation{File: file(), FileName: "a.mp4", Caption: "c", HasSpoiler: true, CaptionAbove: true},
	} {
		_, err := b.Send(ChatID(2), media)
		require.NoError(t, err, method)
		assert.Equal(t, "true", sent[method]["has_spoiler"], method)
		assert.Equal(t, "true", sent[method]["show_caption_above_media"], method)
	}

	_, err = b.Send(ChatID(2), &Video{File: file(), Caption: "c"})
	require.NoError(t, err)
	assert.NotContains(t, sent["sendVideo"], "has_spoiler")
	assert.NotContains(t, sent["sendVideo"], "show_caption_above_media")
}

func TestSendLocationAlertRadius(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendLocation", json.RawMessage(`{"message_id":1,"chat":{"id":7}}`))

	_, err := b.Send(ChatID(7), &Location{Lat: 1, Lng: 2, Heading: 90, AlertRadius: 300, LivePeriod: 60})
	require.NoError(t, err)
	assert.Equal(t, "90", srv.Last().Params["heading"])
	assert.Equal(t, "300", srv.Last().Params["proximity_alert_radius"])
}

func TestSendContact(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendContact", json.RawMessage(`{"message_id":1,"chat":{"id":7},"contact":{"phone_number":"+100","first_name":"Ada"}}`))

	msg, err := b.Send(ChatID(7), &Contact{PhoneNumber: "+100", FirstName: "Ada", VCard: "BEGIN:VCARD"}, Silent)
	require.NoError(t, err)
	assert.Equal(t, "Ada", msg.Contact.FirstName)
	assert.Equal(t, map[string]string{"chat_id": "7", "phone_number": "+100", "first_name": "Ada", "vcard": "BEGIN:VCARD",
		"disable_notification": "true"}, srv.Last().Params)
}

func TestSendPreviewOptions(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendMessage", json.RawMessage(`{"message_id":1,"chat":{"id":7}}`))
	srv.Reply("editMessageText", json.RawMessage(`{"message_id":1,"chat":{"id":7}}`))

	msg, err := b.Send(ChatID(7), "see https://example.com", &PreviewOptions{URL: "https://example.org", AboveText: true})
	require.NoError(t, err)
	var opts PreviewOptions
	srv.Last().JSON(t, "link_preview_options", &opts)
	assert.Equal(t, PreviewOptions{URL: "https://example.org", AboveText: true}, opts)

	_, err = b.Edit(msg, "edited", &SendOptions{PreviewOptions: &PreviewOptions{Disabled: true}})
	require.NoError(t, err)
	srv.Last().JSON(t, "link_preview_options", &opts)
	assert.True(t, opts.Disabled)
}
