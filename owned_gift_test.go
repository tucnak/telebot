package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOwnedGifts(t *testing.T) {
	b, srv := newFakeBot(t)
	page := json.RawMessage(`{"total_count":2,"next_offset":"n","gifts":[
		{"type":"regular","gift":{"id":"g1","sticker":{"file_id":"s"},"star_count":10},"send_date":1,"text":"hi",
			"can_be_upgraded":true,"convert_star_count":8,"owned_gift_id":"og1","sender_user":{"id":3}},
		{"type":"unique","gift":` + uniqueGiftJSON + `,"send_date":2,"can_be_transferred":true,"transfer_star_count":25}]}`)
	srv.Reply("getBusinessAccountGifts", page)
	srv.Reply("getUserGifts", page)
	srv.Reply("getChatGifts", page)

	gifts, err := b.BusinessGifts("bc", GiftQuery{ExcludeUnsaved: true, SortByPrice: true, Limit: 10})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"business_connection_id": "bc", "exclude_unsaved": "true", "sort_by_price": "true", "limit": "10"},
		srv.Last().Params)

	assert.Equal(t, 2, gifts.TotalCount)
	assert.Equal(t, "n", gifts.NextOffset)
	reg, uniq := gifts.Gifts[0], gifts.Gifts[1]
	assert.Equal(t, OwnedGiftRegular, reg.Type)
	assert.Equal(t, "g1", reg.Gift.ID)
	assert.Nil(t, reg.UniqueGift)
	assert.Equal(t, 8, reg.ConvertStarCount)
	assert.Equal(t, int64(3), reg.SenderUser.ID)
	assert.Equal(t, OwnedGiftUnique, uniq.Type)
	assert.Equal(t, "Cake-7", uniq.UniqueGift.Name)
	assert.Nil(t, uniq.Gift)
	assert.Equal(t, 25, uniq.TransferStarCount)

	_, err = b.UserGifts(&User{ID: 7}, GiftQuery{ExcludeFromBlockchain: true, Offset: "x"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"user_id": "7", "exclude_from_blockchain": "true", "offset": "x"}, srv.Last().Params)

	_, err = b.ChatGifts(ChatID(-100), GiftQuery{})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"chat_id": "-100"}, srv.Last().Params)
}

func TestGiftMethods(t *testing.T) {
	b, srv := newFakeBot(t)

	require.NoError(t, b.ConvertGiftToStars("bc", "og1"))
	assert.Equal(t, map[string]string{"business_connection_id": "bc", "owned_gift_id": "og1"}, srv.Last().Params)

	require.NoError(t, b.UpgradeGift("bc", "og1", true, 25))
	assert.Equal(t, "upgradeGift", srv.Last().Method)
	assert.Equal(t, map[string]string{"business_connection_id": "bc", "owned_gift_id": "og1", "keep_original_details": "true", "star_count": "25"}, srv.Last().Params)

	require.NoError(t, b.TransferGift("bc", "og2", ChatID(9), 0))
	assert.Equal(t, map[string]string{"business_connection_id": "bc", "owned_gift_id": "og2", "new_owner_chat_id": "9"}, srv.Last().Params)

	require.NoError(t, b.GiftPremiumSubscription(&User{ID: 7}, 3, 1000, "<b>enjoy</b>", ModeHTML))
	assert.Equal(t, "giftPremiumSubscription", srv.Last().Method)
	assert.Equal(t, map[string]string{"user_id": "7", "month_count": "3", "star_count": "1000", "text": "<b>enjoy</b>", "text_parse_mode": "HTML"}, srv.Last().Params)

	srv.Reply("getChat", json.RawMessage(`{"id":7,"type":"private","accepted_gift_types":{"unlimited_gifts":true,"limited_gifts":false,"unique_gifts":true,"premium_subscription":true,"gifts_from_channels":false}}`))
	chat, err := b.ChatByID(7)
	require.NoError(t, err)
	assert.Equal(t, &AcceptedGiftTypes{UnlimitedGifts: true, UniqueGifts: true, PremiumSubscription: true}, chat.AcceptedGiftTypes)

	fired := false
	b.Handle(OnGiftUpgradeSent, func(c Context) error { fired = c.Message().GiftUpgradeSent.Gift.ID == "g1"; return nil })
	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":1,"date":1,"chat":{"id":7},"gift_upgrade_sent":{"gift":{"id":"g1","star_count":1}}}}`), &u))
	b.ProcessUpdate(u)
	assert.True(t, fired)
}

func TestUserGiftsLive(t *testing.T) {
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}
	if userID == 0 {
		t.Skip("USER_ID is required for gift live tests")
	}
	gifts, err := b.UserGifts(user, GiftQuery{Limit: 5})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, gifts.TotalCount, len(gifts.Gifts))
}
