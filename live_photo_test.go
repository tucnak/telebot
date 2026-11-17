package telebot

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const livePhotoJSON = `{"file_id":"lv","file_unique_id":"lu","width":1080,"height":1920,"duration":3,
	"photo":[{"file_id":"ps","file_unique_id":"psu","width":90,"height":160},{"file_id":"pl","file_unique_id":"plu","width":1080,"height":1920}]}`

func TestSendLivePhoto(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendLivePhoto", json.RawMessage(`{"message_id":1,"chat":{"id":7},"caption":"c","live_photo":`+livePhotoJSON+`}`))

	lp := &LivePhoto{
		File:       FromReader(strings.NewReader("mov")),
		Still:      &Photo{File: FromReader(strings.NewReader("jpg"))},
		Caption:    "c",
		HasSpoiler: true,
	}
	msg, err := b.Send(ChatID(7), lp)
	require.NoError(t, err)

	req := srv.Last()
	assert.Equal(t, "sendLivePhoto", req.Method)
	assert.Equal(t, "mov", string(req.Files["live_photo"].Data))
	assert.Equal(t, "jpg", string(req.Files["photo"].Data))
	assert.Equal(t, "true", req.Params["has_spoiler"])
	assert.Equal(t, "c", req.Params["caption"])

	assert.Equal(t, "lv", msg.LivePhoto.FileID)
	assert.Equal(t, "lv", lp.FileID, "the sent live photo takes the file ID")
	assert.Equal(t, "c", lp.Caption)

	// A received live photo is sent again with its largest static photo.
	srv.Reply("sendLivePhoto", json.RawMessage(`{"message_id":2,"chat":{"id":7},"live_photo":`+livePhotoJSON+`}`))
	_, err = b.Send(ChatID(8), msg.LivePhoto)
	require.NoError(t, err)
	assert.Equal(t, "lv", srv.Last().Params["live_photo"])
	assert.Equal(t, "pl", srv.Last().Params["photo"])

	_, err = b.Send(ChatID(8), &LivePhoto{File: File{FileID: "x"}})
	assert.Error(t, err, "a live photo needs a static photo")
}

func TestLivePhotoInAlbumAndEdit(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("sendMediaGroup", json.RawMessage(`[{"message_id":1,"chat":{"id":7},"live_photo":`+livePhotoJSON+`},{"message_id":2,"chat":{"id":7},"photo":[{"file_id":"p2"}]}]`))

	lp := &LivePhoto{File: FromReader(strings.NewReader("mov")), Still: &Photo{File: FromReader(strings.NewReader("jpg"))}}
	_, err := b.SendAlbum(ChatID(7), Album{lp, &Photo{File: File{FileID: "p2"}}})
	require.NoError(t, err)

	req := srv.Last()
	var media []map[string]interface{}
	req.JSON(t, "media", &media)
	assert.Equal(t, "live_photo", media[0]["type"])
	assert.Equal(t, "attach://0", media[0]["media"])
	assert.Equal(t, "attach://0s", media[0]["photo"])
	assert.Equal(t, "jpg", string(req.Files["0s"].Data))
	assert.Equal(t, "lv", lp.FileID)

	srv.Reply("editMessageMedia", json.RawMessage(`{"message_id":1,"chat":{"id":7}}`))
	_, err = b.EditMedia(&Message{ID: 1, Chat: &Chat{ID: 7}}, &LivePhoto{File: File{FileID: "lv"}, Still: &Photo{File: File{FileID: "pl"}}})
	require.NoError(t, err)
	assert.JSONEq(t, `{"type":"live_photo","media":"lv","photo":"pl","caption":""}`, srv.Last().Params["media"])

	srv.Reply("sendPaidMedia", json.RawMessage(`{"message_id":3,"chat":{"id":7}}`))
	_, err = b.SendPaidMedia(ChatID(7), 5, PaidAlbum{&LivePhoto{File: File{FileID: "lv"}, Still: &Photo{File: File{FileID: "pl"}}}})
	require.NoError(t, err)
	assert.JSONEq(t, `[{"type":"live_photo","media":"lv","photo":"pl","caption":""}]`, srv.Last().Params["media"])
}

func TestReceiveLivePhoto(t *testing.T) {
	b, _ := newFakeBot(t)

	var got *LivePhoto
	b.Handle(OnLivePhoto, func(c Context) error { got = c.Message().LivePhoto; return nil })

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":1,"date":1,"chat":{"id":7,"type":"private"},"live_photo":`+livePhotoJSON+`}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, got)
	assert.Equal(t, 3, got.Duration)
	assert.Len(t, got.Photo, 2)

	var pm PaidMedia
	require.NoError(t, json.Unmarshal([]byte(`{"type":"live_photo","live_photo":`+livePhotoJSON+`}`), &pm))
	assert.Equal(t, "lv", pm.LivePhoto.FileID)
}
