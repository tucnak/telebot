package telebot

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type StickerSetType = string

const (
	StickerRegular     = "regular"
	StickerMask        = "mask"
	StickerCustomEmoji = "custom_emoji"
)

type StickerSetFormat = string

const (
	StickerStatic   = "static"
	StickerAnimated = "animated"
	StickerVideo    = "video"
)

// StickerSet represents a sticker set.
type StickerSet struct {
	Type          StickerSetType   `json:"sticker_type"`
	Format        StickerSetFormat `json:"sticker_format"`
	Name          string           `json:"name"`
	Title         string           `json:"title"`
	Animated      bool             `json:"is_animated"`
	Video         bool             `json:"is_video"`
	Stickers      []Sticker        `json:"stickers"`
	Sticker       Sticker          `json:"sticker"`
	Thumbnail     *Photo           `json:"thumbnail"`
	Emojis        string           `json:"emojis"`
	ContainsMasks bool             `json:"contains_masks"` // FIXME: can be removed
	MaskPosition  *MaskPosition    `json:"mask_position"`
	Repaint       bool             `json:"needs_repainting"`
}

// MaskPosition describes the position on faces where
// a mask should be placed by default.
type MaskPosition struct {
	Feature MaskFeature `json:"point"`
	XShift  float32     `json:"x_shift"`
	YShift  float32     `json:"y_shift"`
	Scale   float32     `json:"scale"`
}

// MaskFeature defines sticker mask position.
type MaskFeature string

const (
	FeatureForehead MaskFeature = "forehead"
	FeatureEyes     MaskFeature = "eyes"
	FeatureMouth    MaskFeature = "mouth"
	FeatureChin     MaskFeature = "chin"
)

// UploadSticker uploads a PNG file with a sticker for later use.
func (b *Bot) UploadSticker(to Recipient, s StickerSet) (*File, error) {
	files := map[string]File{
		"sticker": s.Sticker.File,
	}

	params := map[string]string{
		"user_id":        to.Recipient(),
		"sticker_format": s.Format,
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

// StickerSet returns a sticker set on success.
func (b *Bot) StickerSet(name string) (*StickerSet, error) {
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

// CreateStickerSet creates a new sticker set.
func (b *Bot) CreateStickerSet(to Recipient, s StickerSet) error {
	files := make(map[string]File)
	for i, sticker := range s.Stickers {
		key := fmt.Sprint("sticker", i)
		files[key] = sticker.File
	}

	data, err := json.Marshal(s.Stickers)
	if err != nil {
		return err
	}

	params := map[string]string{
		"user_id":          to.Recipient(),
		"name":             s.Name,
		"title":            s.Title,
		"sticker_type":     s.Type,
		"sticker_format":   s.Format,
		"stickers":         string(data),
		"needs_repainting": strconv.FormatBool(s.Repaint),
	}

	_, err = b.sendFiles("createNewStickerSet", files, params)
	return err
}

// AddStickerToSet adds a new sticker to the existing sticker set.
func (b *Bot) AddStickerToSet(to Recipient, s StickerSet) error {
	var (
		files   = make(map[string]File)
		sticker = s.Sticker
	)
	files["sticker"] = sticker.File

	params := map[string]string{
		"user_id": to.Recipient(),
		"name":    s.Name,
	}

	if sticker.Emojis != nil {
		data, _ := json.Marshal(s.Emojis)
		params["emoji_list"] = string(data)
	}
	if s.MaskPosition != nil {
		data, _ := json.Marshal(s.MaskPosition)
		params["mask_position"] = string(data)
	}
	if sticker.Keywords != nil {
		data, _ := json.Marshal(sticker.Keywords)
		params["keywords"] = string(data)
	}

	_, err := b.sendFiles("addStickerToSet", files, params)
	return err
}

// SetStickerPosition moves a sticker in set to a specific position.
func (b *Bot) SetStickerPosition(sticker string, position int) error {
	params := map[string]string{
		"sticker":  sticker,
		"position": strconv.Itoa(position),
	}

	_, err := b.Raw("setStickerPositionInSet", params)
	return err
}

// DeleteSticker deletes a sticker from a set created by the bot.
func (b *Bot) DeleteSticker(sticker string) error {
	_, err := b.Raw("deleteStickerFromSet", map[string]string{"sticker": sticker})
	return err

}

// SetStickerSetThumb sets a thumbnail of the sticker set.
// Animated thumbnails can be set for animated sticker sets only.
//
// Thumbnail must be a PNG image, up to 128 kilobytes in size
// and have width and height exactly 100px, or a TGS animation
// up to 32 kilobytes in size.
//
// Animated sticker set thumbnail can't be uploaded via HTTP URL.
func (b *Bot) SetStickerSetThumb(to Recipient, s StickerSet) error {
	var (
		sticker = s.Sticker
		files   = make(map[string]File)
	)
	files["thumbnail"] = sticker.File

	data, err := json.Marshal(sticker.File)
	if err != nil {
		return err
	}

	params := map[string]string{
		"name":      s.Name,
		"user_id":   to.Recipient(),
		"thumbnail": string(data),
	}

	_, err = b.sendFiles("setStickerSetThumbnail", files, params)
	return err
}

// SetStickerSetTitle sets the title of a created sticker set.
func (b *Bot) SetStickerSetTitle(s StickerSet) error {
	params := map[string]string{
		"name":  s.Name,
		"title": s.Title,
	}

	_, err := b.Raw("setStickerSetTitle", params)
	return err
}

// DeleteStickerSet deletes a sticker set that was created by the bot.
func (b *Bot) DeleteStickerSet(name string) error {
	params := map[string]string{"name": name}

	_, err := b.Raw("deleteStickerSet", params)
	return err
}

// SetStickerEmojiList changes the list of emoji assigned to a regular or custom emoji sticker.
func (b *Bot) SetStickerEmojiList(sticker string, emojis []string) error {
	data, err := json.Marshal(emojis)
	if err != nil {
		return err
	}

	params := map[string]string{
		"sticker":    sticker,
		"emoji_list": string(data),
	}

	_, err = b.Raw("setStickerEmojiList", params)
	return err
}

// SetStickerKeywords changes search keywords assigned to a regular or custom emoji sticker.
func (b *Bot) SetStickerKeywords(sticker string, keywords []string) error {
	mk, err := json.Marshal(keywords)
	if err != nil {
		return err
	}

	params := map[string]string{
		"sticker":  sticker,
		"keywords": string(mk),
	}

	_, err = b.Raw("setStickerKeywords", params)
	return err
}

// SetStickerMaskPosition changes the mask position of a mask sticker.
func (b *Bot) SetStickerMaskPosition(sticker string, mask MaskPosition) error {
	data, err := json.Marshal(mask)
	if err != nil {
		return err
	}

	params := map[string]string{
		"sticker":       sticker,
		"mask_position": string(data),
	}

	_, err = b.Raw("setStickerMaskPosition", params)
	return err
}

// CustomEmojiStickers returns the information about custom emoji stickers by their ids.
func (b *Bot) CustomEmojiStickers(ids []string) ([]Sticker, error) {
	data, _ := json.Marshal(ids)

	params := map[string]string{
		"custom_emoji_ids": string(data),
	}

	data, err := b.Raw("getCustomEmojiStickers", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result []Sticker
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

// SetCustomEmojiStickerSetThumb sets the thumbnail of a custom emoji sticker set.
func (b *Bot) SetCustomEmojiStickerSetThumb(name, id string) error {
	params := map[string]string{
		"name":            name,
		"custom_emoji_id": id,
	}

	_, err := b.Raw("setCustomEmojiStickerSetThumbnail", params)
	return err
}
