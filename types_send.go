package telebot

import "fmt"

func (b *Bot) sendText(to Recipient, text string, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Destination(),
		"text":    text,
	}
	embedSendOptions(params, opt)

	respJSON, err := b.sendCommand("sendMessage", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}

func (p *Photo) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Destination(),
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
		"chat_id": to.Destination(),
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
		"chat_id": to.Destination(),
	}
	embedSendOptions(params, opt)

	msg, err := b.sendObject(&d.File, "audio", params)
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
		"chat_id": to.Destination(),
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
		"chat_id": to.Destination(),
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

func (x *Location) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":   to.Destination(),
		"latitude":  fmt.Sprintf("%f", x.Latitude),
		"longitude": fmt.Sprintf("%f", x.Longitude),
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
		"chat_id":       to.Destination(),
		"latitude":      fmt.Sprintf("%f", v.Location.Latitude),
		"longitude":     fmt.Sprintf("%f", v.Location.Longitude),
		"title":         v.Title,
		"address":       v.Address,
		"foursquare_id": v.FoursquareID,
	}
	embedSendOptions(params, opt)

	respJSON, err := b.sendCommand("sendLocation", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}
