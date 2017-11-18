package telebot

// Sticker object represents a WebP image, so-called sticker.
type Sticker struct {
	File

	Width  int `json:"width"`
	Height int `json:"height"`

	Thumbnail    *Photo        `json:"thumb,omitempty"`
	Emoji        string        `json:"emoji,omitempty"`
	SetName      string        `json:"set_name,omitempty"`
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
}

type MaskPosition struct {
}
