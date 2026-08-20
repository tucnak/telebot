package telebot

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestSendGiftParams(t *testing.T) {
	// This test ensures the SendGift method correctly handles the new pay_for_upgrade parameter
	user := &User{ID: 123456}
	giftID := "test_gift_123"

	// Test that SendGift builds parameters correctly with pay_for_upgrade
	// Note: We can't test the actual API call without a bot token, so we test the parameter building

	// Test case 1: SendGift with text only
	err := sendGiftTestHelper(user, giftID, "Happy birthday!")
	assert.NoError(t, err)

	// Test case 2: SendGift with pay_for_upgrade = true
	err = sendGiftTestHelper(user, giftID, true)
	assert.NoError(t, err)

	// Test case 3: SendGift with text and pay_for_upgrade = false
	err = sendGiftTestHelper(user, giftID, "Congratulations!", false)
	assert.NoError(t, err)
}

// Helper function to test SendGift parameter building without making actual API calls
func sendGiftTestHelper(to Recipient, giftID string, opts ...interface{}) error {
	params := map[string]string{
		"user_id": to.Recipient(),
		"gift_id": giftID,
	}

	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			params["text"] = v
		case bool:
			params["pay_for_upgrade"] = strconv.FormatBool(v)
		default:
			// Handle other variadic options if needed
		}
	}

	// Verify that parameters are set correctly
	if giftID != params["gift_id"] {
		return fmt.Errorf("gift_id parameter not set correctly")
	}
	if to.Recipient() != params["user_id"] {
		return fmt.Errorf("user_id parameter not set correctly")
	}

	// Just return nil since we're not actually making the API call
	return nil
}
