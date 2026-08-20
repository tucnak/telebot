package telebot

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlbumSetCaption(t *testing.T) {
	tests := []struct {
		name  string
		media Inputtable
	}{
		{
			name:  "photo",
			media: &Photo{Caption: "wrong_caption"},
		},
		{
			name:  "animation",
			media: &Animation{Caption: "wrong_caption"},
		},
		{
			name:  "video",
			media: &Video{Caption: "wrong_caption"},
		},
		{
			name:  "audio",
			media: &Audio{Caption: "wrong_caption"},
		},
		{
			name:  "document",
			media: &Document{Caption: "wrong_caption"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a Album
			a = append(a, tt.media)
			a = append(a, &Photo{Caption: "random_caption"})
			a.SetCaption("correct_caption")
			assert.Equal(t, "correct_caption", a[0].InputMedia().Caption)
			assert.Equal(t, "random_caption", a[1].InputMedia().Caption)
		})
	}
}

func TestPaidMediaPurchased(t *testing.T) {
	jsonData := `{
		"from": {
			"id": 123,
			"is_bot": false,
			"first_name": "Test"
		},
		"paid_media_payload": "test_payload_123"
	}`

	var pmp PaidMediaPurchased
	err := json.Unmarshal([]byte(jsonData), &pmp)
	assert.NoError(t, err)
	assert.NotNil(t, pmp.From)
	assert.Equal(t, int64(123), pmp.From.ID)
	assert.Equal(t, "test_payload_123", pmp.Payload)
}

func TestPaidMediaPayload(t *testing.T) {
	jsonData := `{
		"type": "photo",
		"payload": "custom_payload"
	}`

	var pm PaidMedia
	err := json.Unmarshal([]byte(jsonData), &pm)
	assert.NoError(t, err)
	assert.Equal(t, "photo", pm.Type)
	assert.Equal(t, "custom_payload", pm.Payload)
}
