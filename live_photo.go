package telebot

import (
	"errors"
)

// LivePhoto is a photo with a short video (Bot API 10.0). Its File is the
// video; Photo holds the sizes of the static photo.
//
// To send one, set File to the video, no longer than 10 seconds and 10 MB,
// and Still to the static photo; a received live photo is sent again as it
// is, with its largest static photo.
//
// Bot API: LivePhoto.
type LivePhoto struct {
	File

	// (Optional) Sizes of the static photo, as received.
	Photo []Photo `json:"photo,omitempty"`

	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Duration int    `json:"duration"`
	MIME     string `json:"mime_type,omitempty"`

	// The static photo to send; if nil, the largest of Photo is used.
	Still *Photo `json:"-"`

	// (Optional) Caption and its placement, and a spoiler, when sending.
	Caption      string `json:"caption,omitempty"`
	CaptionAbove bool   `json:"show_caption_above_media,omitempty"`
	HasSpoiler   bool   `json:"has_spoiler,omitempty"`
}

// MediaType returns "live_photo".
func (p *LivePhoto) MediaType() string {
	return "live_photo"
}

// MediaFile returns the video of the live photo.
func (p *LivePhoto) MediaFile() *File {
	return &p.File
}

// InputMedia returns the live photo as an InputMediaLivePhoto, without its
// files.
func (p *LivePhoto) InputMedia() InputMedia {
	return InputMedia{
		Type:         p.MediaType(),
		Caption:      p.Caption,
		CaptionAbove: p.CaptionAbove,
		HasSpoiler:   p.HasSpoiler,
	}
}

// Paid reports that a live photo can be sent as paid media.
func (p *LivePhoto) Paid() bool {
	return true
}

// still returns the static photo to send.
func (p *LivePhoto) still() *File {
	if p.Still != nil {
		return &p.Still.File
	}
	if n := len(p.Photo); n > 0 {
		return &p.Photo[n-1].File
	}
	return nil
}

// embedStill sets the static photo of a live photo in im, adding it to files
// as name when it is uploaded. Other media are left as they are.
func embedStill(m Inputtable, im *InputMedia, name string, files map[string]File) error {
	lp, ok := m.(*LivePhoto)
	if !ok {
		return nil
	}
	still := lp.still()
	if still == nil {
		return errors.New("telebot: live photo has no static photo")
	}
	if im.Photo = still.process(name, files); im.Photo == "" {
		return errors.New("telebot: static photo of the live photo does not exist")
	}
	return nil
}

// Send delivers the live photo through bot b to recipient.
//
// Bot API: sendLivePhoto.
func (p *LivePhoto) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"caption": p.Caption,
	}
	b.embedSendOptions(params, opt)
	embedMediaFlags(params, p.HasSpoiler, p.CaptionAbove)

	still := p.still()
	if still == nil {
		return nil, errors.New("telebot: live photo has no static photo")
	}
	files := map[string]File{"live_photo": p.File, "photo": *still}

	data, err := b.sendFiles("sendLivePhoto", files, params)
	if err != nil {
		return nil, err
	}
	msg, err := extractMessage(data)
	if err != nil {
		return nil, err
	}

	if lp := msg.LivePhoto; lp != nil {
		lp.File.stealRef(&p.File)
		caption, above, spoiler := p.Caption, p.CaptionAbove, p.HasSpoiler
		*p = *lp
		p.Caption, p.CaptionAbove, p.HasSpoiler = caption, above, spoiler
	}
	return msg, nil
}
