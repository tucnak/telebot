package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTonWalletCard(t *testing.T) {
	b, srv := newFakeBot(t)
	srv.Reply("getChat", json.RawMessage(`{"id":7,"type":"private","ton_wallet_card":{
		"address":"EQA","amount":"1500000000","background_nft_address":"EQN",
		"stickers":[{"position":{"x_percentage":10.5,"y_percentage":20,"width_percentage":30,"height_percentage":40,"rotation_angle":45,"corner_radius_percentage":5},
			"nft":{"address":"EQN","collection_address":"EQC","index":"12","name":"N","description":"D",
				"attributes":[{"trait_type":"Color","value":"Red"}],"gift":`+uniqueGiftJSON+`}}]}}`))

	chat, err := b.ChatByID(7)
	require.NoError(t, err)
	card := chat.TonWalletCard
	require.NotNil(t, card)
	assert.Equal(t, "1500000000", card.Amount)
	assert.Equal(t, "EQN", card.BackgroundNFTAddress)
	require.Len(t, card.Stickers, 1)
	st := card.Stickers[0]
	assert.Equal(t, AreaPosition{X: 10.5, Y: 20, Width: 30, Height: 40, Rotation: 45, CornerRadius: 5}, st.Position)
	assert.Equal(t, []TonNftAttribute{{TraitType: "Color", Value: "Red"}}, st.NFT.Attributes)
	assert.Equal(t, "Cake-7", st.NFT.Gift.Name)

	info, err := b.ChatFullInfo(ChatID(7))
	require.NoError(t, err)
	assert.Equal(t, "EQA", info.TonWalletCard.Address)
}

func TestTonWalletTransfer(t *testing.T) {
	b, _ := newFakeBot(t)

	var got *TonWalletTransfer
	b.Handle(OnTonWalletTransfer, func(c Context) error { got = c.Message().TonWalletTransfer; return nil })

	var u Update
	require.NoError(t, json.Unmarshal([]byte(`{"update_id":1,"message":{"message_id":1,"date":1,"chat":{"id":7,"type":"private"},
		"ton_wallet_transfer":{"transaction_id":"tx","peer_address":"EQP","amount":"-100","comment":"thanks"}}}`), &u))
	b.ProcessUpdate(u)
	assert.Equal(t, &TonWalletTransfer{TransactionID: "tx", PeerAddress: "EQP", Amount: "-100", Comment: "thanks"}, got)
}

func TestRichTonAddress(t *testing.T) {
	in := `{"type":"ton_address","text":"EQA…","address":"EQA"}`
	var r RichText
	require.NoError(t, json.Unmarshal([]byte(in), &r))
	assert.Equal(t, RichTonAddress, r.Type)
	assert.Equal(t, "EQA", r.Address)
	out, err := json.Marshal(r)
	require.NoError(t, err)
	assert.JSONEq(t, in, string(out))
}
