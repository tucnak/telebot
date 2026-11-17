package telebot

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
)

// Recipient is any possible endpoint you can send
// messages to: either user, group or a channel.
type Recipient interface {
	Recipient() string // must return legit Telegram chat_id or username
}

// Sendable is any object that can send itself.
//
// This is pretty cool, since it lets bots implement
// custom Sendables for complex kind of media or
// chat objects spanning across multiple messages.
type Sendable interface {
	Send(*Bot, Recipient, *SendOptions) (*Message, error)
}

// embedMediaFlags writes the spoiler and caption-above flags a photo, video
// or animation carries itself, beside any SendOptions.HasSpoiler.
func embedMediaFlags(params map[string]string, spoiler, captionAbove bool) {
	if spoiler {
		params["has_spoiler"] = "true"
	}
	if captionAbove {
		params["show_caption_above_media"] = "true"
	}
}

// Send delivers media through bot b to recipient.
func (p *Photo) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"caption": p.Caption,
	}
	b.embedSendOptions(params, opt)
	embedMediaFlags(params, p.HasSpoiler, p.CaptionAbove)

	msg, err := b.sendMedia(p, params, nil)
	if err != nil {
		return nil, err
	}

	if msg.Photo != nil {
		msg.Photo.File.stealRef(&p.File)
		spoiler, above := p.HasSpoiler, p.CaptionAbove
		*p = *msg.Photo
		p.Caption, p.HasSpoiler, p.CaptionAbove = msg.Caption, spoiler, above
	}

	return msg, nil
}

// Send delivers media through bot b to recipient.
func (a *Audio) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":   to.Recipient(),
		"caption":   a.Caption,
		"performer": a.Performer,
		"title":     a.Title,
		"file_name": a.FileName,
	}
	b.embedSendOptions(params, opt)

	if a.Duration != 0 {
		params["duration"] = strconv.Itoa(a.Duration)
	}

	msg, err := b.sendMedia(a, params, thumbnailToFilemap(a.Thumbnail))
	if err != nil {
		return nil, err
	}

	if msg.Audio != nil {
		msg.Audio.File.stealRef(&a.File)
		*a = *msg.Audio
		a.Caption = msg.Caption
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
		"chat_id":   to.Recipient(),
		"caption":   d.Caption,
		"file_name": d.FileName,
	}
	b.embedSendOptions(params, opt)

	if d.FileSize != 0 {
		params["file_size"] = strconv.FormatInt(d.FileSize, 10)
	}
	if d.DisableTypeDetection {
		params["disable_content_type_detection"] = "true"
	}

	msg, err := b.sendMedia(d, params, thumbnailToFilemap(d.Thumbnail))
	if err != nil {
		return nil, err
	}

	if doc := msg.Document; doc != nil {
		doc.File.stealRef(&d.File)
		*d = *doc
		d.Caption = msg.Caption
	} else if vid := msg.Video; vid != nil {
		vid.File.stealRef(&d.File)
		d.Caption = vid.Caption
		d.MIME = vid.MIME
		d.Thumbnail = vid.Thumbnail
	}

	return msg, nil
}

// Send delivers media through bot b to recipient.
func (s *Sticker) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"emoji":   s.Emoji,
	}
	b.embedSendOptions(params, opt)

	msg, err := b.sendMedia(s, params, nil)
	if err != nil {
		return nil, err
	}

	if msg.Sticker != nil {
		msg.Sticker.File.stealRef(&s.File)
		*s = *msg.Sticker
	}

	return msg, nil
}

// Send delivers media through bot b to recipient.
func (v *Video) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":   to.Recipient(),
		"caption":   v.Caption,
		"file_name": v.FileName,
	}
	b.embedSendOptions(params, opt)
	embedMediaFlags(params, v.HasSpoiler, v.CaptionAbove)

	if v.Duration != 0 {
		params["duration"] = strconv.Itoa(v.Duration)
	}
	if v.Width != 0 {
		params["width"] = strconv.Itoa(v.Width)
	}
	if v.Height != 0 {
		params["height"] = strconv.Itoa(v.Height)
	}
	if v.Streaming {
		params["supports_streaming"] = "true"
	}
	if v.StartTimestamp != 0 {
		params["start_timestamp"] = strconv.FormatInt(v.StartTimestamp, 10)
	}

	// Prepare files map with thumbnail and cover
	files := thumbnailToFilemap(v.Thumbnail)
	if v.Cover != nil {
		files["cover"] = v.Cover.File
	}

	msg, err := b.sendMedia(v, params, files)
	if err != nil {
		return nil, err
	}

	if vid := msg.Video; vid != nil {
		vid.File.stealRef(&v.File)
		spoiler, above := v.HasSpoiler, v.CaptionAbove
		*v = *vid
		v.Caption, v.HasSpoiler, v.CaptionAbove = msg.Caption, spoiler, above
	} else if doc := msg.Document; doc != nil {
		// If video has no sound, Telegram can turn it into Document (GIF)
		doc.File.stealRef(&v.File)

		v.Caption = doc.Caption
		v.MIME = doc.MIME
		v.Thumbnail = doc.Thumbnail
	}

	return msg, nil
}

// Send delivers animation through bot b to recipient.
func (a *Animation) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":   to.Recipient(),
		"caption":   a.Caption,
		"file_name": a.FileName,
	}
	b.embedSendOptions(params, opt)
	embedMediaFlags(params, a.HasSpoiler, a.CaptionAbove)

	if a.Duration != 0 {
		params["duration"] = strconv.Itoa(a.Duration)
	}
	if a.Width != 0 {
		params["width"] = strconv.Itoa(a.Width)
	}
	if a.Height != 0 {
		params["height"] = strconv.Itoa(a.Height)
	}

	// file_name is required, without it animation sends as a document
	if params["file_name"] == "" && a.File.OnDisk() {
		params["file_name"] = filepath.Base(a.File.FileLocal)
	}

	msg, err := b.sendMedia(a, params, thumbnailToFilemap(a.Thumbnail))
	if err != nil {
		return nil, err
	}

	spoiler, above := a.HasSpoiler, a.CaptionAbove
	if anim := msg.Animation; anim != nil {
		anim.File.stealRef(&a.File)
		*a = *msg.Animation
	} else if doc := msg.Document; doc != nil {
		*a = Animation{
			File:      doc.File,
			Thumbnail: doc.Thumbnail,
			MIME:      doc.MIME,
			FileName:  doc.FileName,
		}
	}

	a.Caption, a.HasSpoiler, a.CaptionAbove = msg.Caption, spoiler, above
	return msg, nil
}

// Send delivers media through bot b to recipient.
func (v *Voice) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"caption": v.Caption,
	}
	b.embedSendOptions(params, opt)

	if v.Duration != 0 {
		params["duration"] = strconv.Itoa(v.Duration)
	}

	msg, err := b.sendMedia(v, params, nil)
	if err != nil {
		return nil, err
	}

	// Telegram sends a voice note that is not OGG, MP3 or M4A as an audio
	// or a document.
	if msg.Voice != nil {
		msg.Voice.File.stealRef(&v.File)
		*v = *msg.Voice
	}

	return msg, nil
}

// Send delivers media through bot b to recipient.
func (v *VideoNote) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
	}
	b.embedSendOptions(params, opt)

	if v.Duration != 0 {
		params["duration"] = strconv.Itoa(v.Duration)
	}
	if v.Length != 0 {
		params["length"] = strconv.Itoa(v.Length)
	}

	msg, err := b.sendMedia(v, params, thumbnailToFilemap(v.Thumbnail))
	if err != nil {
		return nil, err
	}

	if msg.VideoNote != nil {
		msg.VideoNote.File.stealRef(&v.File)
		*v = *msg.VideoNote
	}

	return msg, nil
}

// Send delivers the contact through bot b to recipient.
//
// Bot API: sendContact.
func (c *Contact) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":      to.Recipient(),
		"phone_number": c.PhoneNumber,
		"first_name":   c.FirstName,
	}
	if c.LastName != "" {
		params["last_name"] = c.LastName
	}
	if c.VCard != "" {
		params["vcard"] = c.VCard
	}
	b.embedSendOptions(params, opt)

	data, err := b.Raw("sendContact", params)
	if err != nil {
		return nil, err
	}
	return extractMessage(data)
}

// Send delivers media through bot b to recipient.
func (x *Location) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":     to.Recipient(),
		"latitude":    fmt.Sprintf("%f", x.Lat),
		"longitude":   fmt.Sprintf("%f", x.Lng),
		"live_period": strconv.Itoa(x.LivePeriod),
	}
	if x.HorizontalAccuracy != nil {
		params["horizontal_accuracy"] = fmt.Sprintf("%f", *x.HorizontalAccuracy)
	}
	if x.Heading != 0 {
		params["heading"] = strconv.Itoa(x.Heading)
	}
	if x.AlertRadius != 0 {
		params["proximity_alert_radius"] = strconv.Itoa(x.AlertRadius)
	}
	b.embedSendOptions(params, opt)

	data, err := b.Raw("sendLocation", params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
}

// Send delivers media through bot b to recipient.
func (v *Venue) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":           to.Recipient(),
		"latitude":          fmt.Sprintf("%f", v.Location.Lat),
		"longitude":         fmt.Sprintf("%f", v.Location.Lng),
		"title":             v.Title,
		"address":           v.Address,
		"foursquare_id":     v.FoursquareID,
		"foursquare_type":   v.FoursquareType,
		"google_place_id":   v.GooglePlaceID,
		"google_place_type": v.GooglePlaceType,
	}
	b.embedSendOptions(params, opt)

	data, err := b.Raw("sendVenue", params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
}

// Send delivers invoice through bot b to recipient.
func (i *Invoice) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := i.params()
	params["chat_id"] = to.Recipient()
	b.embedSendOptions(params, opt)

	data, err := b.Raw("sendInvoice", params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
}

// Send delivers poll through bot b to recipient.
func (p *Poll) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":                 to.Recipient(),
		"question":                p.Question,
		"type":                    string(p.Type),
		"is_closed":               strconv.FormatBool(p.Closed),
		"is_anonymous":            strconv.FormatBool(p.Anonymous),
		"allows_multiple_answers": strconv.FormatBool(p.MultipleAnswers),
	}
	if p.Type == PollQuiz || len(p.CorrectOptions) > 0 {
		correct := p.CorrectOptions
		if len(correct) == 0 {
			correct = []int{p.CorrectOption}
		}
		data, _ := json.Marshal(correct)
		params["correct_option_ids"] = string(data)
	}
	if p.QuestionParseMode != "" {
		params["question_parse_mode"] = p.QuestionParseMode
	}
	if len(p.QuestionEntities) > 0 {
		data, _ := json.Marshal(p.QuestionEntities)
		params["question_entities"] = string(data)
	}
	if p.Explanation != "" {
		params["explanation"] = p.Explanation
		params["explanation_parse_mode"] = p.ParseMode
		if len(p.Entities) > 0 {
			data, _ := json.Marshal(p.Entities)
			params["explanation_entities"] = string(data)
			delete(params, "explanation_parse_mode")
		}
	}
	if p.Description != "" {
		params["description"] = p.Description
		if len(p.DescriptionEntities) > 0 {
			data, _ := json.Marshal(p.DescriptionEntities)
			params["description_entities"] = string(data)
		} else if p.DescriptionParseMode != "" {
			params["description_parse_mode"] = p.DescriptionParseMode
		}
	}
	if p.OpenPeriod != 0 {
		params["open_period"] = strconv.Itoa(p.OpenPeriod)
	} else if p.CloseUnixdate != 0 {
		params["close_date"] = strconv.FormatInt(p.CloseUnixdate, 10)
	}
	for key, on := range map[string]bool{
		"allows_revoting":           p.AllowsRevoting,
		"shuffle_options":           p.ShuffleOptions,
		"allow_adding_options":      p.AllowAddingOptions,
		"hide_results_until_closes": p.HideResultsUntilCloses,
		"members_only":              p.MembersOnly,
	} {
		if on {
			params[key] = "true"
		}
	}
	if p.DisallowRevoting {
		params["allows_revoting"] = "false"
	}
	if len(p.CountryCodes) > 0 {
		data, _ := json.Marshal(p.CountryCodes)
		params["country_codes"] = string(data)
	}
	b.embedSendOptions(params, opt)

	files := make(map[string]File)
	for key, m := range map[string]InputPollMedia{"media": p.InputMedia, "explanation_media": p.InputExplanationMedia} {
		if m == nil {
			continue
		}
		media, err := encodePollMedia(m, key, files)
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(media)
		params[key] = string(data)
	}

	options, err := encodePollOptions(p.Options, files)
	if err != nil {
		return nil, err
	}
	params["options"] = options

	data, err := b.sendFiles("sendPoll", files, params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
}

// Send delivers dice through bot b to recipient.
func (d *Dice) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"emoji":   string(d.Type),
	}
	b.embedSendOptions(params, opt)

	data, err := b.Raw("sendDice", params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
}

// Send delivers game through bot b to recipient.
func (g *Game) Send(b *Bot, to Recipient, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":         to.Recipient(),
		"game_short_name": g.Name,
	}
	b.embedSendOptions(params, opt)

	data, err := b.Raw("sendGame", params)
	if err != nil {
		return nil, err
	}

	return extractMessage(data)
}

func thumbnailToFilemap(thumb *Photo) map[string]File {
	files := make(map[string]File)
	if thumb != nil {
		files["thumbnail"] = thumb.File
	}
	return files
}
