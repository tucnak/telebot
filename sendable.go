package telebot

import "fmt"

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

func (p *Photo) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"caption": p.Caption,
	}

	embedSendOptions(params, opt)

	msg, err := b.sendObject(&p.File, "photo", params)
	if err != nil {
		return nil, err
	}

	thumbnails := msg.Photo
	fname := p.filename
	p.File = thumbnails[len(thumbnails)-1].File
	p.filename = fname

	return msg, nil
}

func (a *Audio) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"caption": a.Caption,
	}
	embedSendOptions(params, opt)

	msg, err := b.sendObject(&a.File, "audio", params)
	if err != nil {
		return nil, err
	}

	fname := a.filename
	*a = *msg.Audio
	a.filename = fname

	return msg, nil
}

func (d *Document) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"caption": d.Caption,
	}
	embedSendOptions(params, opt)

	msg, err := b.sendObject(&d.File, "document", params)
	if err != nil {
		return nil, err
	}

	fname := d.filename
	*d = *msg.Document
	d.filename = fname

	return msg, nil
}

func (s *Sticker) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
	}
	embedSendOptions(params, opt)

	msg, err := b.sendObject(&s.File, "sticker", params)
	if err != nil {
		return nil, err
	}

	fname := s.filename
	*s = *msg.Sticker
	s.filename = fname

	return msg, nil
}

func (v *Video) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"caption": v.Caption,
	}
	embedSendOptions(params, opt)

	msg, err := b.sendObject(&v.File, "video", params)
	if err != nil {
		return nil, err
	}

	fname := v.filename
	*v = *msg.Video
	v.filename = fname

	return msg, nil
}

func (v *Voice) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
	}
	embedSendOptions(params, opt)

	msg, err := b.sendObject(&v.File, "voice", params)
	if err != nil {
		return nil, err
	}

	fname := v.filename
	*v = *msg.Voice
	v.filename = fname

	return msg, nil
}

func (v *VideoNote) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
	}
	embedSendOptions(params, opt)

	msg, err := b.sendObject(&v.File, "videoNote", params)
	if err != nil {
		return nil, err
	}

	fname := v.filename
	*v = *msg.VideoNote
	v.filename = fname

	return msg, nil
}

func (x *Location) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":   to.Recipient(),
		"latitude":  fmt.Sprintf("%f", x.Lat),
		"longitude": fmt.Sprintf("%f", x.Lng),
	}
	embedSendOptions(params, opt)

	respJSON, err := b.sendCommand("sendLocation", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}

func (v *Venue) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":       to.Recipient(),
		"latitude":      fmt.Sprintf("%f", v.Location.Lat),
		"longitude":     fmt.Sprintf("%f", v.Location.Lng),
		"title":         v.Title,
		"address":       v.Address,
		"foursquare_id": v.FoursquareID,
	}
	embedSendOptions(params, opt)

	respJSON, err := b.sendCommand("sendVenue", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}
