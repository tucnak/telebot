package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const uniqueGiftJSON = `{
	"gift_id": "g1", "base_name": "Cake", "name": "Cake-7", "number": 7,
	"model": {"name": "M", "sticker": {"file_id": "s1", "file_unique_id": "u1", "type": "regular", "width": 1, "height": 1}, "rarity_per_mille": 5, "rarity": "epic"},
	"symbol": {"name": "S", "sticker": {"file_id": "s2", "file_unique_id": "u2", "type": "regular", "width": 1, "height": 1}, "rarity_per_mille": 3},
	"backdrop": {"name": "B", "colors": {"center_color": 1, "edge_color": 2, "symbol_color": 3, "text_color": 4}, "rarity_per_mille": 9},
	"is_from_blockchain": true,
	"colors": {"model_custom_emoji_id": "e1", "symbol_custom_emoji_id": "e2", "light_theme_main_color": 5,
		"light_theme_other_colors": [6], "dark_theme_main_color": 7, "dark_theme_other_colors": [8, 9]}
}`

func TestGiftServiceMessages(t *testing.T) {
	b, _ := newFakeBot(t)

	var gift, unique *Message
	b.Handle(OnGift, func(c Context) error { gift = c.Message(); return nil })
	b.Handle(OnUniqueGift, func(c Context) error { unique = c.Message(); return nil })

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":1,"date":1,"chat":{"id":2,"type":"private"},
		"gift":{"gift":{"id":"g1","sticker":{"file_id":"s"},"star_count":50,"is_premium":true,"has_colors":true,
			"personal_total_count":10,"personal_remaining_count":3,"unique_gift_variant_count":40,
			"background":{"center_color":1,"edge_color":2,"text_color":3},"publisher_chat":{"id":-100}},
		"convert_star_count":40,"can_be_upgraded":true,"text":"hi","is_private":true,"unique_gift_number":12}}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, gift)
	g := gift.Gift
	assert.Equal(t, 40, g.ConvertStarCount)
	assert.True(t, g.CanBeUpgraded)
	assert.Equal(t, "hi", g.Text)
	assert.Equal(t, 12, g.UniqueGiftNumber)
	assert.True(t, g.Gift.IsPremium)
	assert.Equal(t, 3, g.Gift.PersonalRemainingCount)
	assert.Equal(t, &GiftBackground{CenterColor: 1, EdgeColor: 2, TextColor: 3}, g.Gift.Background)
	assert.Equal(t, int64(-100), g.Gift.PublisherChat.ID)

	u = Update{}
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":2,"message":{"message_id":2,"date":1,"chat":{"id":2,"type":"private"},
		"unique_gift":{"gift":`+uniqueGiftJSON+`,"origin":"resale","last_resale_currency":"TON","last_resale_amount":5000000000,
		"text":"enjoy","is_private":true,"next_transfer_date":1700000000}}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, unique)
	ug := unique.UniqueGift
	assert.Equal(t, "resale", ug.Origin)
	assert.Equal(t, int64(5000000000), ug.LastResaleAmount)
	assert.Equal(t, "enjoy", ug.Text)
	assert.True(t, ug.IsPrivate)
	assert.Equal(t, "Cake-7", ug.Gift.Name)
	assert.Equal(t, "epic", ug.Gift.Model.Rarity)
	assert.Equal(t, "s2", ug.Gift.Symbol.Sticker.FileID)
	assert.Equal(t, 3, ug.Gift.Backdrop.Colors.SymbolColor)
	assert.True(t, ug.Gift.IsFromBlockchain)
	assert.Equal(t, []int{8, 9}, ug.Gift.Colors.DarkThemeOtherColors)
}
