package telebot

// AreaPosition is the position of a rectangle on a story's media or on a TON
// wallet card, in percentages of the media or card size (Bot API 10.4; named
// StoryAreaPosition before).
//
// Bot API: AreaPosition.
type AreaPosition struct {
	X            float64 `json:"x_percentage"`
	Y            float64 `json:"y_percentage"`
	Width        float64 `json:"width_percentage"`
	Height       float64 `json:"height_percentage"`
	Rotation     float64 `json:"rotation_angle"`
	CornerRadius float64 `json:"corner_radius_percentage"`
}

// TonWalletCard is the TON wallet card shown to everyone on a user's profile
// (Bot API 10.4).
//
// Bot API: TonWalletCard.
type TonWalletCard struct {
	// Address of the wallet.
	Address string `json:"address"`

	// Balance in the smallest units of the cryptocurrency, as a decimal
	// number.
	Amount string `json:"amount"`

	// Stickers placed on the card.
	Stickers []TonWalletCardSticker `json:"stickers"`

	// (Optional) Address of the sticker NFT whose colors fill the card
	// background.
	BackgroundNFTAddress string `json:"background_nft_address,omitempty"`
}

// TonWalletCardSticker is a sticker on a TON wallet card.
//
// Bot API: TonWalletCardSticker.
type TonWalletCardSticker struct {
	Position AreaPosition `json:"position"`
	NFT      TonNft       `json:"nft"`
}

// TonNft is an NFT in the TON blockchain.
//
// Bot API: TonNft.
type TonNft struct {
	Address string `json:"address"`

	// (Optional) Address of the NFT's collection.
	CollectionAddress string `json:"collection_address,omitempty"`

	// Index of the NFT in its collection, as a decimal number.
	Index string `json:"index"`

	Name        string            `json:"name"`
	Description string            `json:"description"`
	Attributes  []TonNftAttribute `json:"attributes,omitempty"`

	// (Optional) The unique gift corresponding to the NFT.
	Gift *UniqueGift `json:"gift,omitempty"`
}

// TonNftAttribute is an attribute of a TON NFT.
//
// Bot API: TonNftAttribute.
type TonNftAttribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

// TonWalletTransfer is a service message about a transfer with a TON wallet
// (Bot API 10.4).
//
// Bot API: TonWalletTransfer.
type TonWalletTransfer struct {
	// Identifier of the transaction; empty if it hasn't been sent yet.
	TransactionID string `json:"transaction_id"`

	// Address of the peer in the TON blockchain.
	PeerAddress string `json:"peer_address"`

	// Amount in the smallest units of the cryptocurrency, as a decimal
	// number; negative for outgoing transfers.
	Amount string `json:"amount"`

	// (Optional) Comment of the transfer; omitted if it is encrypted.
	Comment string `json:"comment,omitempty"`
}
