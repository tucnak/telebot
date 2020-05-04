package telebot

// Sticker object represents a WebP image, so-called sticker.
type Sticker struct {
	File
	Width        int           `json:"width"`
	Height       int           `json:"height"`
	Animated     bool          `json:"is_animated"`
	Thumbnail    *Photo        `json:"thumb"`
	Emoji        string        `json:"emoji"`
	Name         string        `json:"name"`
	SetName      string        `json:"set_name"`
	PNG          *File         `json:"png_sticker"`
	TGS          *File         `json:"tgs_file"`
	Emojis       string        `json:"emojis"`
	MaskPosition *MaskPosition `json:"mask_position"`
}

// StickerSet represents a sticker set
type StickerSet struct {
	Name          string        `json:"name"`
	Title         string        `json:"title"`
	Animated      bool          `json:"is_animated"`
	ContainsMasks bool          `json:"contains_masks"`
	Stickers      []Sticker     `json:"stickers"`
	Thumbnail     *Photo        `json:"thumb"`
	PNG           *File         `json:"png_sticker"`
	TGS           *File         `json:"tgs_file"`
	Emojis        string        `json:"emojis"`
	MaskPosition  *MaskPosition `json:"mask_position"`
}

// MaskPosition describes the position on faces where
// a mask should be placed by default.
type MaskPosition struct {
	Feature MaskFeature `json:"point"`
	XShift  float32     `json:"x_shift"`
	YShift  float32     `json:"y_shift"`
	Scale   float32     `json:"scale"`
}

// StickerSetParams describes the payload in creating new sticker set api-method.
type StickerSetParams struct {
	UserID     int
	Name       string
	Title      string
	PngSticker *File
	Emojis     string
}
