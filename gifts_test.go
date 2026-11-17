package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGiftUnmarshal(t *testing.T) {
	jsonData := `{
		"id": "gift_abc123",
		"sticker": {
			"file_id": "sticker_xyz",
			"file_unique_id": "unique_sticker",
			"type": "regular",
			"width": 512,
			"height": 512,
			"is_animated": false,
			"is_video": false
		},
		"star_count": 150,
		"upgrade_star_count": 1000,
		"total_count": 5000,
		"remaining_count": 3500
	}`

	var gift Gift
	err := json.Unmarshal([]byte(jsonData), &gift)
	assert.NoError(t, err)
	assert.Equal(t, "gift_abc123", gift.ID)
	assert.Equal(t, 150, gift.StarCount)
	assert.Equal(t, 1000, gift.UpgradeStarCount)
	assert.Equal(t, 5000, gift.TotalCount)
	assert.Equal(t, 3500, gift.RemainingCount)
	assert.NotNil(t, gift.Sticker)
	assert.Equal(t, "sticker_xyz", gift.Sticker.FileID)
}

func TestGiftMinimal(t *testing.T) {
	jsonData := `{
		"id": "gift_minimal",
		"sticker": {
			"file_id": "sticker_123",
			"file_unique_id": "unique_123",
			"type": "regular",
			"width": 512,
			"height": 512,
			"is_animated": false,
			"is_video": false
		},
		"star_count": 50
	}`

	var gift Gift
	err := json.Unmarshal([]byte(jsonData), &gift)
	assert.NoError(t, err)
	assert.Equal(t, "gift_minimal", gift.ID)
	assert.Equal(t, 50, gift.StarCount)
	assert.Equal(t, 0, gift.UpgradeStarCount)
	assert.Equal(t, 0, gift.TotalCount)
	assert.Equal(t, 0, gift.RemainingCount)
}

func TestGiftsUnmarshal(t *testing.T) {
	jsonData := `{
		"gifts": [
			{
				"id": "gift_1",
				"sticker": {
					"file_id": "sticker_1",
					"file_unique_id": "unique_1",
					"type": "regular",
					"width": 512,
					"height": 512,
					"is_animated": false,
					"is_video": false
				},
				"star_count": 100
			},
			{
				"id": "gift_2",
				"sticker": {
					"file_id": "sticker_2",
					"file_unique_id": "unique_2",
					"type": "regular",
					"width": 512,
					"height": 512,
					"is_animated": false,
					"is_video": false
				},
				"star_count": 200,
				"upgrade_star_count": 500
			}
		]
	}`

	var gifts Gifts
	err := json.Unmarshal([]byte(jsonData), &gifts)
	assert.NoError(t, err)
	assert.Len(t, gifts.Gifts, 2)
	assert.Equal(t, "gift_1", gifts.Gifts[0].ID)
	assert.Equal(t, 100, gifts.Gifts[0].StarCount)
	assert.Equal(t, "gift_2", gifts.Gifts[1].ID)
	assert.Equal(t, 200, gifts.Gifts[1].StarCount)
	assert.Equal(t, 500, gifts.Gifts[1].UpgradeStarCount)
}

func TestGiftOmitEmpty(t *testing.T) {
	gift := Gift{
		ID:        "gift_test",
		StarCount: 75,
		// Sticker intentionally nil for this test
	}

	data, err := json.Marshal(gift)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, "gift_test", result["id"])
	assert.Equal(t, float64(75), result["star_count"])
	assert.NotContains(t, result, "upgrade_star_count")
	assert.NotContains(t, result, "total_count")
	assert.NotContains(t, result, "remaining_count")
}

func TestSendGift(t *testing.T) {
	b, srv := newFakeBot(t)

	require.NoError(t, b.SendGift(&User{ID: 7}, "g", "<b>Happy</b> birthday!", ModeHTML, true))
	assert.Equal(t, map[string]string{"user_id": "7", "gift_id": "g", "text": "<b>Happy</b> birthday!",
		"text_parse_mode": "HTML", "pay_for_upgrade": "true"}, srv.Last().Params)

	require.NoError(t, b.SendGift(&Chat{ID: -100}, "g", "Hi", &SendOptions{Entities: Entities{{Type: EntityBold, Length: 2}}}))
	req := srv.Last()
	assert.Equal(t, "-100", req.Params["chat_id"])
	assert.NotContains(t, req.Params, "user_id")
	assert.JSONEq(t, `[{"type":"bold","offset":0,"length":2}]`, req.Params["text_entities"])

	// A chat given by ID or username is a chat too.
	require.NoError(t, b.SendGift(ChatID(-100), "g"))
	assert.Equal(t, map[string]string{"chat_id": "-100", "gift_id": "g"}, srv.Last().Params)
	require.NoError(t, b.SendGift(ChatID(7), "g"))
	assert.Equal(t, map[string]string{"user_id": "7", "gift_id": "g"}, srv.Last().Params)
}
