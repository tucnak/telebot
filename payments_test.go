package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentSubscriptionFields(t *testing.T) {
	jsonData := `{
		"currency": "XTR",
		"total_amount": 1000,
		"invoice_payload": "test_payload",
		"telegram_payment_charge_id": "charge_123",
		"provider_payment_charge_id": "provider_456",
		"subscription_expiration_date": 1700000000,
		"is_recurring": true,
		"is_first_recurring": false
	}`

	var p Payment
	err := json.Unmarshal([]byte(jsonData), &p)
	assert.NoError(t, err)
	assert.Equal(t, "XTR", p.Currency)
	assert.Equal(t, 1000, p.Total)
	assert.Equal(t, int64(1700000000), p.SubscriptionExpirationDate)
	assert.True(t, p.IsRecurring)
	assert.False(t, p.IsFirstRecurring)
}

func TestPaymentSubscriptionFieldsOmitEmpty(t *testing.T) {
	p := Payment{
		Currency:         "XTR",
		Total:            500,
		TelegramChargeID: "charge_789",
	}

	data, err := json.Marshal(p)
	assert.NoError(t, err)

	// Verify subscription fields are omitted when zero value
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.NotContains(t, result, "subscription_expiration_date")
	assert.NotContains(t, result, "is_recurring")
	assert.NotContains(t, result, "is_first_recurring")
}

func TestInvoiceSubscriptionFields(t *testing.T) {
	invoice := Invoice{
		Title:                "Test Subscription",
		Description:          "Monthly subscription",
		Currency:             "XTR",
		SubscriptionPeriod:   2592000, // 30 days in seconds
		BusinessConnectionID: "conn_123",
	}

	params := invoice.params()
	assert.Equal(t, "Test Subscription", params["title"])
	assert.Equal(t, "Monthly subscription", params["description"])
	assert.Equal(t, "XTR", params["currency"])
	assert.Equal(t, "2592000", params["subscription_period"])
	assert.Equal(t, "conn_123", params["business_connection_id"])
}

func TestInvoiceParamsOmitEmptySubscription(t *testing.T) {
	invoice := Invoice{
		Title:       "Test Invoice",
		Description: "One-time payment",
		Currency:    "USD",
	}

	params := invoice.params()
	assert.NotContains(t, params, "subscription_period")
	assert.NotContains(t, params, "business_connection_id")
}
