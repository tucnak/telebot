package telebot

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Recipient is any possible endpoint you can send
// messages to: either user, group or a channel.
type Recipient interface {
	// Must return legit Telegram chat_id or username
	Recipient() string
}

// Sendable is any object that can send itself.
//
// This is pretty cool, since it lets bots implement
// custom Sendables for complex kind of media or
// chat objects spanning across multiple messages.
type Sendable interface {
	Send(*Bot, Recipient, *SendOptions) (*Message, error)
}

// Send delivers media through bot b to recipient.
func (p *Photo) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"caption": p.Caption,
	}

	embedSendOptions(params, opt)

	msg, err := b.sendObject(&p.File, "photo", params, nil)
	if err != nil {
		return nil, err
	}

	msg.Photo.File.stealRef(&p.File)
	*p = *msg.Photo

	return msg, nil
}

// Send delivers media through bot b to recipient.
func (a *Audio) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":   to.Recipient(),
		"caption":   a.Caption,
		"performer": a.Performer,
		"title":     a.Title,
	}

	if a.Duration != 0 {
		params["duration"] = strconv.Itoa(a.Duration)
	}

	embedSendOptions(params, opt)

	msg, err := b.sendObject(&a.File, "audio", params, nil)
	if err != nil {
		return nil, err
	}

	if msg.Audio != nil {
		msg.Audio.File.stealRef(&a.File)
		*a = *msg.Audio
	}

	if msg.Document != nil {
		msg.Document.File.stealRef(&a.File)
		a.File = msg.Document.File
	}

	return msg, nil
}

// Send delivers media through bot b to recipient.
func (d *Document) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"caption": d.Caption,
	}

	if d.FileSize != 0 {
		params["file_size"] = strconv.Itoa(d.FileSize)
	}

	embedSendOptions(params, opt)

	msg, err := b.sendObject(&d.File, "document", params, thumbnailToFilemap(d.Thumbnail))
	if err != nil {
		return nil, err
	}

	msg.Document.File.stealRef(&d.File)
	*d = *msg.Document

	return msg, nil
}

// Send delivers media through bot b to recipient.
func (s *Sticker) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
	}
	embedSendOptions(params, opt)

	msg, err := b.sendObject(&s.File, "sticker", params, nil)
	if err != nil {
		return nil, err
	}

	msg.Sticker.File.stealRef(&s.File)
	*s = *msg.Sticker

	return msg, nil
}

// Send delivers media through bot b to recipient.
func (v *Video) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"caption": v.Caption,
	}

	if v.Duration != 0 {
		params["duration"] = strconv.Itoa(v.Duration)
	}
	if v.Width != 0 {
		params["width"] = strconv.Itoa(v.Width)
	}
	if v.Height != 0 {
		params["height"] = strconv.Itoa(v.Height)
	}
	if v.SupportsStreaming {
		params["supports_streaming"] = "true"
	}

	embedSendOptions(params, opt)

	msg, err := b.sendObject(&v.File, "video", params, thumbnailToFilemap(v.Thumbnail))
	if err != nil {
		return nil, err
	}

	if vid := msg.Video; vid != nil {
		vid.File.stealRef(&v.File)
		*v = *vid
	} else if doc := msg.Document; doc != nil {
		// If video has no sound, Telegram can turn it into Document (GIF)
		doc.File.stealRef(&v.File)

		v.Caption = doc.Caption
		v.MIME = doc.MIME
		v.Thumbnail = doc.Thumbnail
	}

	return msg, nil
}

// Send delivers media through bot b to recipient.
func (v *Voice) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
	}

	if v.Duration != 0 {
		params["duration"] = strconv.Itoa(v.Duration)
	}

	embedSendOptions(params, opt)

	msg, err := b.sendObject(&v.File, "voice", params, nil)
	if err != nil {
		return nil, err
	}

	msg.Voice.File.stealRef(&v.File)
	*v = *msg.Voice

	return msg, nil
}

// Send delivers media through bot b to recipient.
func (v *VideoNote) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
	}

	if v.Duration != 0 {
		params["duration"] = strconv.Itoa(v.Duration)
	}
	if v.Length != 0 {
		params["length"] = strconv.Itoa(v.Length)
	}

	embedSendOptions(params, opt)

	msg, err := b.sendObject(&v.File, "videoNote", params, thumbnailToFilemap(v.Thumbnail))
	if err != nil {
		return nil, err
	}

	msg.VideoNote.File.stealRef(&v.File)
	*v = *msg.VideoNote

	return msg, nil
}

// Send delivers media through bot b to recipient.
func (x *Location) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":     to.Recipient(),
		"latitude":    fmt.Sprintf("%f", x.Lat),
		"longitude":   fmt.Sprintf("%f", x.Lng),
		"live_period": strconv.Itoa(x.LivePeriod),
	}
	embedSendOptions(params, opt)

	respJSON, err := b.Raw("sendLocation", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}

// Send delivers media through bot b to recipient.
func (v *Venue) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":         to.Recipient(),
		"latitude":        fmt.Sprintf("%f", v.Location.Lat),
		"longitude":       fmt.Sprintf("%f", v.Location.Lng),
		"title":           v.Title,
		"address":         v.Address,
		"foursquare_id":   v.FoursquareID,
		"foursquare_type": v.FoursquareType,
	}
	embedSendOptions(params, opt)

	respJSON, err := b.Raw("sendVenue", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}

// Send delivers media through bot b to recipient.
func (i *Invoice) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	prices, _ := json.Marshal(i.Prices)

	params := map[string]string{
		"chat_id":         to.Recipient(),
		"title":           i.Title,
		"description":     i.Description,
		"start_parameter": i.Start,
		"payload":         i.Payload,
		"provider_token":  i.Token,
		"currency":        i.Currency,
		"prices":          string(prices),
	}
	embedSendOptions(params, opt)

	respJSON, err := b.Raw("sendInvoice", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}
