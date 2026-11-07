package telebot

import (
	"bytes"
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
