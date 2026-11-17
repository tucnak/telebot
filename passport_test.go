package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPassport(t *testing.T) {
	b, srv := newFakeBot(t)

	var got *PassportData
	b.Handle(OnPassportData, func(c Context) error { got = c.Message().PassportData; return nil })

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":1,"date":1,"chat":{"id":7},
		"passport_data":{"data":[{"type":"passport","data":"ZGF0YQ==","front_side":{"file_id":"f","file_unique_id":"u","file_size":10,"file_date":5},"hash":"h"}],
		"credentials":{"data":"d","hash":"h","secret":"s"}}}}`), &u))
	b.ProcessUpdate(u)
	require.NotNil(t, got)
	assert.Equal(t, "passport", got.Data[0].Type)
	assert.Equal(t, int64(10), got.Data[0].FrontSide.FileSize)
	assert.Equal(t, "s", got.Credentials.Secret)

	require.NoError(t, b.SetPassportDataErrors(&User{ID: 7}, []PassportElementError{
		{Source: PassportErrorDataField, Type: "passport", FieldName: "document_no", DataHash: "dh", Message: "bad"},
		{Source: PassportErrorFiles, Type: "utility_bill", FileHashes: []string{"a", "b"}, Message: "blurry"},
	}))
	assert.Equal(t, "setPassportDataErrors", srv.Last().Method)
	assert.JSONEq(t, `[{"source":"data","type":"passport","field_name":"document_no","data_hash":"dh","message":"bad"},
		{"source":"files","type":"utility_bill","file_hashes":["a","b"],"message":"blurry"}]`, srv.Last().Params["errors"])
}
