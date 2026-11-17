package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionPartnerSubscriptionPeriod(t *testing.T) {
	jsonData := `{
		"type": "user",
		"user": {
			"id": 123456,
			"first_name": "John",
			"is_bot": false
		},
		"invoice_payload": "monthly_sub",
		"subscription_period": 2592000
	}`

	var tp TransactionPartner
	err := json.Unmarshal([]byte(jsonData), &tp)
	assert.NoError(t, err)
	assert.Equal(t, TransactionTypeUser, tp.Type)
	assert.Equal(t, 2592000, tp.SubscriptionPeriod)
	assert.Equal(t, "monthly_sub", tp.Payload)
	assert.NotNil(t, tp.User)
	assert.Equal(t, int64(123456), tp.User.ID)
}

func TestTransactionPartnerWithGift(t *testing.T) {
	jsonData := `{
		"type": "user",
		"user": {
			"id": 789012,
			"first_name": "Jane",
			"is_bot": false
		},
		"gift": {
			"id": "gift_123",
			"sticker": {
				"file_id": "sticker_456",
				"file_unique_id": "unique_789",
				"type": "regular",
				"width": 512,
				"height": 512,
				"is_animated": false,
				"is_video": false
			},
			"star_count": 100,
			"upgrade_star_count": 500
		}
	}`

	var tp TransactionPartner
	err := json.Unmarshal([]byte(jsonData), &tp)
	assert.NoError(t, err)
	assert.Equal(t, TransactionTypeUser, tp.Type)
	assert.NotNil(t, tp.Gift)
	assert.Equal(t, "gift_123", tp.Gift.ID)
	assert.Equal(t, 100, tp.Gift.StarCount)
	assert.Equal(t, 500, tp.Gift.UpgradeStarCount)
	assert.NotNil(t, tp.Gift.Sticker)
	assert.Equal(t, "sticker_456", tp.Gift.Sticker.FileID)
}

func TestTransactionPartnerOmitEmpty(t *testing.T) {
	tp := TransactionPartner{
		Type:    TransactionTypeFragment,
		Payload: "test_payload",
	}

	data, err := json.Marshal(tp)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.NotContains(t, result, "subscription_period")
	assert.NotContains(t, result, "gift")
	assert.NotContains(t, result, "user")
}

func TestStarTransactionTime(t *testing.T) {
	tx := StarTransaction{
		ID:       "tx_123",
		Amount:   1000,
		Unixtime: 1700000000,
	}

	time := tx.Time()
	assert.Equal(t, int64(1700000000), time.Unix())
}

func TestStarTransactionWithNanostarAmount(t *testing.T) {
	jsonData := `{
		"id": "tx_nanostar_123",
		"amount": 1000,
		"nanostar_amount": 50000000,
		"date": 1700000000,
		"source": {
			"type": "user",
			"user": {
				"id": 123456,
				"first_name": "John",
				"is_bot": false
			}
		}
	}`

	var tx StarTransaction
	err := json.Unmarshal([]byte(jsonData), &tx)
	assert.NoError(t, err)
	assert.Equal(t, "tx_nanostar_123", tx.ID)
	assert.Equal(t, 1000, tx.Amount)
	assert.Equal(t, 50000000, tx.NanostarAmount)
	assert.True(t, tx.HasNanostarAmount())
	assert.Equal(t, 0.5, tx.NanostarAmountAsTON())
}

func TestStarTransactionNoNanostarAmount(t *testing.T) {
	tx := StarTransaction{
		ID:             "tx_no_nano",
		Amount:         1000,
		NanostarAmount: 0,
		Unixtime:       1700000000,
	}

	assert.False(t, tx.HasNanostarAmount())
	assert.Equal(t, 0.0, tx.NanostarAmountAsTON())
}

func TestTransactionPartnerAffiliateProgram(t *testing.T) {
	jsonData := `{
		"type": "affiliate_program",
		"sponsor_user": {
			"id": 987654,
			"first_name": "Sponsor",
			"is_bot": false
		},
		"commission_per_mille": 50
	}`

	var ap TransactionPartnerAffiliateProgram
	err := json.Unmarshal([]byte(jsonData), &ap)
	assert.NoError(t, err)
	assert.Equal(t, TransactionTypeAffiliateProgram, ap.Type)
	assert.Equal(t, 50, ap.CommissionPerMille)
	assert.NotNil(t, ap.SponsorUser)
	assert.Equal(t, int64(987654), ap.SponsorUser.ID)
}

func TestAffiliateInfo(t *testing.T) {
	jsonData := `{
		"affiliate_program": {
			"type": "affiliate_program",
			"sponsor_user": {
				"id": 987654,
				"first_name": "Sponsor",
				"is_bot": false
			},
			"commission_per_mille": 50
		},
		"commission_per_mille": 30,
		"amount": 100,
		"nanostar_amount": 25000000
	}`

	var ai AffiliateInfo
	err := json.Unmarshal([]byte(jsonData), &ai)
	assert.NoError(t, err)
	assert.True(t, ai.IsAffiliateProgram())
	assert.Equal(t, 30, ai.CommissionPerMille)
	assert.Equal(t, 3.0, ai.CommissionPercentage())
	assert.Equal(t, 100, ai.Amount)
	assert.Equal(t, 25000000, ai.NanostarAmount)
	assert.True(t, ai.HasNanostarAmount())
	assert.Equal(t, 0.25, ai.NanostarAmountAsTON())
}

func TestAffiliateInfoNoProgram(t *testing.T) {
	ai := AffiliateInfo{
		CommissionPerMille: 25,
		Amount:             50,
	}

	assert.False(t, ai.IsAffiliateProgram())
	assert.Equal(t, 2.5, ai.CommissionPercentage())
	assert.Equal(t, 50, ai.Amount)
	assert.False(t, ai.HasNanostarAmount())
}

func TestTransactionPartnerWithAffiliate(t *testing.T) {
	jsonData := `{
		"type": "user",
		"user": {
			"id": 123456,
			"first_name": "John",
			"is_bot": false
		},
		"invoice_payload": "test_payload",
		"affiliate": {
			"affiliate_program": {
				"type": "affiliate_program",
				"commission_per_mille": 50
			},
			"commission_per_mille": 30,
			"amount": 100
		}
	}`

	var tp TransactionPartner
	err := json.Unmarshal([]byte(jsonData), &tp)
	assert.NoError(t, err)
	assert.Equal(t, TransactionTypeUser, tp.Type)
	assert.NotNil(t, tp.User)
	assert.Equal(t, int64(123456), tp.User.ID)
	assert.Equal(t, "test_payload", tp.Payload)
	assert.NotNil(t, tp.Affiliate)
	assert.True(t, tp.Affiliate.IsAffiliateProgram())
	assert.Equal(t, 30, tp.Affiliate.CommissionPerMille)
	assert.Equal(t, 100, tp.Affiliate.Amount)
}

func TestTransactionPartnerAffiliateTypeConstant(t *testing.T) {
	assert.Equal(t, "affiliate_program", TransactionTypeAffiliateProgram)
}

func TestTransactionPartnerFields(t *testing.T) {
	var chat TransactionPartner
	require.NoError(t, json.Unmarshal([]byte(`{"type":"chat","chat":{"id":-100,"type":"channel"},"gift":{"id":"g","star_count":5}}`), &chat))
	require.NotNil(t, chat.Chat)
	assert.Equal(t, int64(-100), chat.Chat.ID)
	assert.Equal(t, "g", chat.Gift.ID)

	data, err := json.Marshal(chat)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"chat":{"id":-100`)

	var user TransactionPartner
	require.NoError(t, json.Unmarshal([]byte(`{"type":"user","transaction_type":"paid_media_payment","user":{"id":7},
		"paid_media":[{"type":"photo","photo":[{"file_id":"p"}]}],"paid_media_payload":"pl",
		"affiliate":{"affiliate_chat":{"id":-5},"commission_per_mille":100,"amount":3}}`), &user))
	assert.Equal(t, "paid_media_payment", user.TransactionType)
	assert.Equal(t, "pl", user.PaidMediaPayload)
	assert.Len(t, user.PaidMedia, 1)
	assert.Equal(t, int64(-5), user.Affiliate.AffiliateChat.ID)

	var api TransactionPartner
	require.NoError(t, json.Unmarshal([]byte(`{"type":"telegram_api","request_count":12}`), &api))
	assert.Equal(t, 12, api.RequestCount)

	var premium TransactionPartner
	require.NoError(t, json.Unmarshal([]byte(`{"type":"user","transaction_type":"premium_purchase","user":{"id":7},"premium_subscription_duration":6}`), &premium))
	assert.Equal(t, 6, premium.PremiumSubscriptionDuration)
}
