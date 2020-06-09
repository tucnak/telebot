package telebot

import (
	"encoding/json"
	"strconv"
)

// Sticker object represents a WebP image, so-called sticker.
type Sticker struct {
	File
	Width        int           `json:"width"`
	Height       int           `json:"height"`
	Animated     bool          `json:"is_animated"`
	Thumbnail    *Photo        `json:"thumb"`
	Emoji        string        `json:"emoji"`
	SetName      string        `json:"set_name"`
	MaskPosition *MaskPosition `json:"mask_position"`
}

// StickerSet represents a sticker set.
type StickerSet struct {
	Name          string        `json:"name"`
	Title         string        `json:"title"`
	Animated      bool          `json:"is_animated"`
	Stickers      []Sticker     `json:"stickers"`
	Thumbnail     *Photo        `json:"thumb"`
	PNG           *File         `json:"png_sticker"`
	TGS           *File         `json:"tgs_sticker"`
	Emojis        string        `json:"emojis"`
	ContainsMasks bool          `json:"contains_masks"`
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

// UploadStickerFile uploads a .PNG file with a sticker for later use.
func (b *Bot) UploadStickerFile(to Recipient, png *File) (*File, error) {
	files := map[string]File{
		"png_sticker": *png,
	}
	params := map[string]string{
		"user_id": to.Recipient(),
	}

	data, err := b.sendFiles("uploadStickerFile", files, params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result File
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return &resp.Result, nil
}

// GetStickerSet returns a StickerSet on success.
func (b *Bot) GetStickerSet(name string) (*StickerSet, error) {
	data, err := b.Raw("getStickerSet", map[string]string{"name": name})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *StickerSet
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// CreateNewStickerSet creates a new sticker set.
func (b *Bot) CreateNewStickerSet(to Recipient, s StickerSet) error {
	files := make(map[string]File)
	if s.PNG != nil {
		files["png_sticker"] = *s.PNG
	}
	if s.TGS != nil {
		files["tgs_sticker"] = *s.TGS
	}

	params := map[string]string{
		"user_id":        to.Recipient(),
		"name":           s.Name,
		"title":          s.Title,
		"emojis":         s.Emojis,
		"contains_masks": strconv.FormatBool(s.ContainsMasks),
	}

	if s.MaskPosition != nil {
		data, err := json.Marshal(&s.MaskPosition)
		if err != nil {
			return err
		}
		params["mask_position"] = string(data)
	}

	_, err := b.sendFiles("createNewStickerSet", files, params)
	return err
}

// AddStickerToSet adds new sticker to existing sticker set.
func (b *Bot) AddStickerToSet(to Recipient, s StickerSet) error {
	files := make(map[string]File)
	if s.PNG != nil {
		files["png_sticker"] = *s.PNG
	} else if s.TGS != nil {
		files["tgs_sticker"] = *s.TGS
	}

	params := map[string]string{
		"user_id": to.Recipient(),
		"name":    s.Name,
		"emojis":  s.Emojis,
	}

	if s.MaskPosition != nil {
		data, err := json.Marshal(&s.MaskPosition)
		if err != nil {
			return err
		}
		params["mask_position"] = string(data)
	}

	_, err := b.sendFiles("addStickerToSet", files, params)
	return err
}

// SetStickerPositionInSet moves a sticker in set to a specific position.
func (b *Bot) SetStickerPositionInSet(sticker string, position int) error {
	params := map[string]string{
		"sticker":  sticker,
		"position": strconv.Itoa(position),
	}

	_, err := b.Raw("setStickerPositionInSet", params)
	return err
}

// DeleteStickerFromSet deletes sticker from set created by the bot.
func (b *Bot) DeleteStickerFromSet(sticker string) error {
	_, err := b.Raw("deleteStickerFromSet", map[string]string{"sticker": sticker})
	return err

}

// SetStickerSetThumb sets the thumbnail of a sticker set.
// Animated thumbnails can be set for animated sticker sets only.
//
// Thumbnail must be a PNG image, up to 128 kilobytes in size
// and have width and height exactly 100px, or a TGS animation
// up to 32 kilobytes in size.
//
// Animated sticker set thumbnail can't be uploaded via HTTP URL.
//
func (b *Bot) SetStickerSetThumb(to Recipient, s StickerSet) error {
	files := map[string]File{}
	if s.PNG != nil {
		files["thumb"] = *s.PNG
	} else if s.TGS != nil {
		files["thumb"] = *s.TGS
	}

	params := map[string]string{
		"name":    s.Name,
		"user_id": to.Recipient(),
	}

	_, err := b.sendFiles("setStickerSetThumb", files, params)
	return err
}
