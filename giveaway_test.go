package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGiveawayPrizeStarCount(t *testing.T) {
	jsonData := `{
		"chats": [],
		"winners_selection_date": 1234567890,
		"winner_count": 5,
		"prize_star_count": 100
	}`

	var g Giveaway
	err := json.Unmarshal([]byte(jsonData), &g)
	assert.NoError(t, err)
	assert.Equal(t, 5, g.WinnerCount)
	assert.Equal(t, 100, g.PrizeStarCount)
}

func TestGiveawayCreatedPrizeStarCount(t *testing.T) {
	jsonData := `{
		"prize_star_count": 50
	}`

	var gc GiveawayCreated
	err := json.Unmarshal([]byte(jsonData), &gc)
	assert.NoError(t, err)
	assert.Equal(t, 50, gc.PrizeStarCount)
}

func TestGiveawayWinnersPrizeStarCount(t *testing.T) {
	jsonData := `{
		"chat": {
			"id": 123,
			"type": "private"
		},
		"giveaway_message_id": 456,
		"winners_selection_date": 1234567890,
		"winner_count": 3,
		"winners": [],
		"prize_star_count": 150
	}`

	var gw GiveawayWinners
	err := json.Unmarshal([]byte(jsonData), &gw)
	assert.NoError(t, err)
	assert.Equal(t, 3, gw.WinnerCount)
	assert.Equal(t, 150, gw.PrizeStarCount)
}

func TestGiveawayCompletedIsStarGiveaway(t *testing.T) {
	jsonData := `{
		"winner_count": 2,
		"is_star_giveaway": true
	}`

	var gc GiveawayCompleted
	err := json.Unmarshal([]byte(jsonData), &gc)
	assert.NoError(t, err)
	assert.Equal(t, 2, gc.WinnerCount)
	assert.True(t, gc.IsStarGiveaway)
}

func TestBoostSourcePrizeStarCount(t *testing.T) {
	jsonData := `{
		"source": "giveaway",
		"giveaway_message_id": 789,
		"prize_star_count": 200
	}`

	var bs BoostSource
	err := json.Unmarshal([]byte(jsonData), &bs)
	assert.NoError(t, err)
	assert.Equal(t, BoostGiveaway, bs.Source)
	assert.Equal(t, 789, bs.GiveawayMessageID)
	assert.Equal(t, 200, bs.PrizeStarCount)
}
