package telebot

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		ID:           "tx_no_nano",
		Amount:       1000,
		NanostarAmount: 0,
		Unixtime:     1700000000,
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
		Amount:            50,
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

func TestVerificationLevelConstants(t *testing.T) {
	assert.Equal(t, VerificationLevel("none"), VerificationLevelNone)
	assert.Equal(t, VerificationLevel("basic"), VerificationLevelBasic)
	assert.Equal(t, VerificationLevel("detailed"), VerificationLevelDetailed)
	assert.Equal(t, VerificationLevel("full"), VerificationLevelFull)
}

func TestVerificationRequirementsMarshal(t *testing.T) {
	req := VerificationRequirements{
		RequirePhoneNumber: true,
		RequireRealName:    true,
		CustomRequirements: []string{"passport", "address_verification"},
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var unmarshaled VerificationRequirements
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.True(t, unmarshaled.RequirePhoneNumber)
	assert.True(t, unmarshaled.RequireRealName)
	assert.False(t, unmarshaled.RequireBirthDate)
	assert.Equal(t, []string{"passport", "address_verification"}, unmarshaled.CustomRequirements)
}

func TestVerificationInfoMethods(t *testing.T) {
	// Use actual current time for proper testing
	now := time.Now().Unix()
	future := now + 86400 // 24 hours later
	past := now - 86400    // 24 hours ago

	tests := []struct {
		name           string
		info           VerificationInfo
		isExpired      bool
		isValid        bool
		verificationTime int64
	}{
		{
			name: "Active verification without expiration",
			info: VerificationInfo{
				Level:            VerificationLevelBasic,
				VerificationDate: now,
				IsActive:         true,
			},
			isExpired:        false,
			isValid:          true,
			verificationTime: now,
		},
		{
			name: "Inactive verification",
			info: VerificationInfo{
				Level:            VerificationLevelBasic,
				VerificationDate: now,
				IsActive:         false,
			},
			isExpired:        false,
			isValid:          false,
			verificationTime: now,
		},
		{
			name: "Expired verification",
			info: VerificationInfo{
				Level:            VerificationLevelBasic,
				VerificationDate: now,
				ExpirationDate:   past,
				IsActive:         true,
			},
			isExpired:        true,
			isValid:          false,
			verificationTime: now,
		},
		{
			name: "Valid future verification",
			info: VerificationInfo{
				Level:            VerificationLevelBasic,
				VerificationDate: now,
				ExpirationDate:   future,
				IsActive:         true,
			},
			isExpired:        false,
			isValid:          true,
			verificationTime: now,
		},
		{
			name: "No expiration date",
			info: VerificationInfo{
				Level:            VerificationLevelBasic,
				VerificationDate: now,
				IsActive:         true,
			},
			isExpired:        false,
			isValid:          true,
			verificationTime: now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, int64(tt.info.VerificationTime().Unix()), tt.verificationTime)
			assert.Equal(t, tt.isExpired, tt.info.IsExpired())
			assert.Equal(t, tt.isValid, tt.info.IsValid())
		})
	}
}

func TestVerificationInfoJSON(t *testing.T) {
	jsonData := `{
		"level": "basic",
		"verification_date": 1700000000,
		"expiration_date": 1700086400,
		"verified_by": "official_bot",
		"details": "User identity verified",
		"is_active": true
	}`

	var info VerificationInfo
	err := json.Unmarshal([]byte(jsonData), &info)
	assert.NoError(t, err)
	assert.Equal(t, VerificationLevelBasic, info.Level)
	assert.Equal(t, int64(1700000000), info.VerificationDate)
	assert.Equal(t, int64(1700086400), info.ExpirationDate)
	assert.Equal(t, "official_bot", info.VerifiedBy)
	assert.Equal(t, "User identity verified", info.Details)
	assert.True(t, info.IsActive)
}
