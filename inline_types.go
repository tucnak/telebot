package telebot

import (
	"encoding/json"
)

// InlineQueryResultArticle represents a link to an article or web page.
// See also: https://core.telegram.org/bots/api#inlinequeryresultarticle
type InlineQueryResultArticle struct {
	// Unique identifier for this result, 1-64 Bytes.
	// If left unspecified, a 64-bit FNV-1 hash will be calculated
	// from the other fields and used automatically.
	ID string `json:"id",hash:"ignore"`

	// Title of the result.
	Title string `json:"title"`

	// Message text. Shortcut (and mutually exclusive to) specifying
	// InputMessageContent.
	Text string `json:"message_text,omitempty"`

	// Content of the message to be sent.
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`

	// Optional. Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`

	// Optional. URL of the result.
	URL string `json:"url,omitempty"`

	// Optional. Pass True, if you don't want the URL to be shown in the message.
	HideURL bool `json:"hide_url,omitempty"`

	// Optional. Short description of the result.
	Description string `json:"description,omitempty"`

	// Optional. Url of the thumbnail for the result.
	ThumbURL string `json:"thumb_url,omitempty"`

	// Optional. Thumbnail width.
	ThumbWidth int `json:"thumb_width,omitempty"`

	// Optional. Thumbnail height.
	ThumbHeight int `json:"thumb_height,omitempty"`
}

func (r *InlineQueryResultArticle) MarshalJSON() ([]byte, error) {
	// avoiding endless self-recursion
	type wrapper InlineQueryResultArticle

	id, err := r.id()
	if err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		wrapper
		Type string `json:"type"`
		ID   string `json:"id",hash:"ignore"`
	}{
		wrapper: wrapper(*r),
		ID:      id,
		Type:    "article",
	})
}

func (r *InlineQueryResultArticle) id() (string, error) {
	if r.ID == "" {
		return hashInlineQueryResult(r)
	}

	return r.ID, nil
}

// InlineQueryResultAudio represents a link to an mp3 audio file.
type InlineQueryResultAudio struct {
	// Unique identifier for this result, 1-64 Bytes.
	// If left unspecified, a 64-bit FNV-1 hash will be calculated
	// from the other fields and used automatically.
	ID string `json:"id",hash:"ignore"`

	// A valid URL for the audio file.
	AudioURL string `json:"audio_url"`

	// Title.
	Title string `json:"title"`

	// Optional. Performer.
	Performer string `json:"performer,omitempty"`

	// Optional. Audio duration in seconds.
	Duration int `json:"audio_duration,omitempty"`

	// Optional. Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`

	// Optional. Content of the message to be sent instead of the audio.
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (r *InlineQueryResultAudio) MarshalJSON() ([]byte, error) {
	// avoiding endless self-recursion
	type wrapper InlineQueryResultAudio

	id, err := r.id()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		wrapper
		Type string `json:"type"`
		ID   string `json:"id",hash:"ignore"`
	}{
		wrapper: wrapper(*r),
		ID:      id,
		Type:    "audio",
	})
}

func (r *InlineQueryResultAudio) id() (string, error) {
	if r.ID == "" {
		return hashInlineQueryResult(r)
	}

	return r.ID, nil
}

// InlineQueryResultContact represents a contact with a phone number.
// See also: https://core.telegram.org/bots/api#inlinequeryresultcontact
type InlineQueryResultContact struct {
	// Unique identifier for this result, 1-64 Bytes.
	// If left unspecified, a 64-bit FNV-1 hash will be calculated
	// from the other fields and used automatically.
	ID string `json:"id",hash:"ignore"`

	// Contact's phone number.
	PhoneNumber string `json:"phone_number"`

	// Contact's first name.
	FirstName string `json:"first_name"`

	// Optional. Contact's last name.
	LastName string `json:"last_name,omitempty"`

	// Optional. Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`

	// Optional. Content of the message to be sent instead of the audio.
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`

	// Optional. Url of the thumbnail for the result.
	ThumbURL string `json:"thumb_url,omitempty"`

	// Optional. Thumbnail width.
	ThumbWidth int `json:"thumb_width,omitempty"`

	// Optional. Thumbnail height.
	ThumbHeight int `json:"thumb_height,omitempty"`
}

func (r *InlineQueryResultContact) MarshalJSON() ([]byte, error) {
	// avoiding endless self-recursion
	type wrapper InlineQueryResultContact

	id, err := r.id()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		wrapper
		Type string `json:"type"`
		ID   string `json:"id",hash:"ignore"`
	}{
		wrapper: wrapper(*r),
		ID:      id,
		Type:    "contact",
	})
}

func (r *InlineQueryResultContact) id() (string, error) {
	if r.ID == "" {
		return hashInlineQueryResult(r)
	}

	return r.ID, nil
}

// InlineQueryResultDocument represents a link to a file.
// See also: https://core.telegram.org/bots/api#inlinequeryresultdocument
type InlineQueryResultDocument struct {
	// Unique identifier for this result, 1-64 Bytes.
	// If left unspecified, a 64-bit FNV-1 hash will be calculated
	// from the other fields and used automatically.
	ID string `json:"id",hash:"ignore"`

	// Title for the result.
	Title string `json:"title"`

	// A valid URL for the file
	DocumentURL string `json:"document_url"`

	// Mime type of the content of the file, either “application/pdf” or
	// “application/zip”.
	MimeType string `json:"mime_type"`

	// Optional. Caption of the document to be sent, 0-200 characters.
	Caption string `json:"caption,omitempty"`

	// Optional. Short description of the result.
	Description string `json:"description,omitempty"`

	// Optional. Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`

	// Optional. Content of the message to be sent instead of the audio.
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`

	// Optional. URL of the thumbnail (jpeg only) for the file.
	ThumbURL string `json:"thumb_url,omitempty"`

	// Optional. Thumbnail width.
	ThumbWidth int `json:"thumb_width,omitempty"`

	// Optional. Thumbnail height.
	ThumbHeight int `json:"thumb_height,omitempty"`
}

func (r *InlineQueryResultDocument) MarshalJSON() ([]byte, error) {
	// avoiding endless self-recursion
	type wrapper InlineQueryResultDocument

	id, err := r.id()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		wrapper
		Type string `json:"type"`
		ID   string `json:"id",hash:"ignore"`
	}{
		wrapper: wrapper(*r),
		ID:      id,
		Type:    "document",
	})
}

func (r *InlineQueryResultDocument) id() (string, error) {
	if r.ID == "" {
		return hashInlineQueryResult(r)
	}

	return r.ID, nil
}

// InlineQueryResultGif represents a link to an animated GIF file.
// See also: https://core.telegram.org/bots/api#inlinequeryresultgif
type InlineQueryResultGif struct {
	// Unique identifier for this result, 1-64 Bytes.
	// If left unspecified, a 64-bit FNV-1 hash will be calculated
	// from the other fields and used automatically.
	ID string `json:"id",hash:"ignore"`

	// A valid URL for the GIF file. File size must not exceed 1MB.
	GifURL string `json:"gif_url"`

	// URL of the static thumbnail for the result (jpeg or gif).
	ThumbURL string `json:"thumb_url"`

	// Optional. Width of the GIF.
	GifWidth int `json:"gif_width,omitempty"`

	// Optional. Height of the GIF.
	GifHeight int `json:"gif_height,omitempty"`

	// Optional. Title for the result.
	Title string `json:"title,omitempty"`

	// Optional. Caption of the GIF file to be sent, 0-200 characters.
	Caption string `json:"caption,omitempty"`

	// Optional. Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`

	// Optional. Content of the message to be sent instead of the audio.
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (r *InlineQueryResultGif) MarshalJSON() ([]byte, error) {
	// avoiding endless self-recursion
	type wrapper InlineQueryResultGif

	id, err := r.id()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		wrapper
		Type string `json:"type"`
		ID   string `json:"id",hash:"ignore"`
	}{
		wrapper: wrapper(*r),
		ID:      id,
		Type:    "gif",
	})
}

func (r *InlineQueryResultGif) id() (string, error) {
	if r.ID == "" {
		return hashInlineQueryResult(r)
	}

	return r.ID, nil
}

// InlineQueryResultLocation represents a location on a map.
// See also: https://core.telegram.org/bots/api#inlinequeryresultlocation
type InlineQueryResultLocation struct {

	// Unique identifier for this result, 1-64 Bytes.
	// If left unspecified, a 64-bit FNV-1 hash will be calculated
	// from the other fields and used automatically.
	ID string `json:"id",hash:"ignore"`

	// Latitude of the location in degrees.
	Latitude float32 `json:"latitude"`

	// Longitude of the location in degrees.
	Longitude float32 `json:"longitude"`

	// Location title.
	Title string `json:"title"`

	// Optional. Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`

	// Optional. Content of the message to be sent instead of the audio.
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`

	// Optional. Url of the thumbnail for the result.
	ThumbURL string `json:"thumb_url,omitempty"`

	// Optional. Thumbnail width.
	ThumbWidth int `json:"thumb_width,omitempty"`

	// Optional. Thumbnail height.
	ThumbHeight int `json:"thumb_height,omitempty"`
}

func (r *InlineQueryResultLocation) MarshalJSON() ([]byte, error) {
	// avoiding endless self-recursion
	type wrapper InlineQueryResultLocation

	id, err := r.id()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		wrapper
		Type string `json:"type"`
		ID   string `json:"id",hash:"ignore"`
	}{
		wrapper: wrapper(*r),
		ID:      id,
		Type:    "location",
	})
}

func (r *InlineQueryResultLocation) id() (string, error) {
	if r.ID == "" {
		return hashInlineQueryResult(r)
	}

	return r.ID, nil
}

// InlineQueryResultMpeg4Gif represents a link to a video animation
// (H.264/MPEG-4 AVC video without sound).
// See also: https://core.telegram.org/bots/api#inlinequeryresultmpeg4gif
type InlineQueryResultMpeg4Gif struct {
	// Unique identifier for this result, 1-64 Bytes.
	// If left unspecified, a 64-bit FNV-1 hash will be calculated
	// from the other fields and used automatically.
	ID string `json:"id",hash:"ignore"`

	// A valid URL for the MP4 file.
	URL string `json:"mpeg4_url"`

	// Optional. Video width.
	Width int `json:"mpeg4_width,omitempty"`

	// Optional. Video height.
	Height int `json:"mpeg4_height,omitempty"`

	// URL of the static thumbnail (jpeg or gif) for the result.
	ThumbURL string `json:"thumb_url,omitempty"`

	// Optional. Title for the result.
	Title string `json:"title,omitempty"`

	// Optional. Caption of the MPEG-4 file to be sent, 0-200 characters.
	Caption string `json:"caption,omitempty"`

	// Optional. Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`

	// Optional. Content of the message to be sent instead of the audio.
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (r *InlineQueryResultMpeg4Gif) MarshalJSON() ([]byte, error) {
	// avoiding endless self-recursion
	type wrapper InlineQueryResultMpeg4Gif

	id, err := r.id()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		wrapper
		Type string `json:"type"`
		ID   string `json:"id",hash:"ignore"`
	}{
		wrapper: wrapper(*r),
		ID:      id,
		Type:    "mpeg4_gif",
	})
}

func (r *InlineQueryResultMpeg4Gif) id() (string, error) {
	if r.ID == "" {
		return hashInlineQueryResult(r)
	}

	return r.ID, nil
}

// InlineQueryResultPhoto represents a link to a photo.
// See also: https://core.telegram.org/bots/api#inlinequeryresultphoto
type InlineQueryResultPhoto struct {
	// Unique identifier for this result, 1-64 Bytes.
	// If left unspecified, a 64-bit FNV-1 hash will be calculated
	// from the other fields and used automatically.
	ID string `json:"id",hash:"ignore"`

	// A valid URL of the photo. Photo must be in jpeg format.
	// Photo size must not exceed 5MB.
	PhotoURL string `json:"photo_url"`

	// URL of the thumbnail for the photo.
	ThumbURL string `json:"thumb_url"`

	// Optional. Width of the photo.
	PhotoWidth int `json:"photo_width,omitempty"`

	// Optional. Height of the photo.
	PhotoHeight int `json:"photo_height,omitempty"`

	// Optional. Title for the result.
	Title string `json:"title,omitempty"`

	// Optional. Short description of the result.
	Description string `json:"description,omitempty"`

	// Optional. Caption of the photo to be sent, 0-200 characters.
	Caption string `json:"caption,omitempty"`

	// Optional. Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`

	// Optional. Content of the message to be sent instead of the audio.
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (r *InlineQueryResultPhoto) MarshalJSON() ([]byte, error) {
	// avoiding endless self-recursion
	type wrapper InlineQueryResultPhoto

	id, err := r.id()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		wrapper
		Type string `json:"type"`
		ID   string `json:"id",hash:"ignore"`
	}{
		wrapper: wrapper(*r),
		ID:      id,
		Type:    "photo",
	})
}

func (r *InlineQueryResultPhoto) id() (string, error) {
	if r.ID == "" {
		return hashInlineQueryResult(r)
	}

	return r.ID, nil
}

// InlineQueryResultVenue represents a venue.
// See also: https://core.telegram.org/bots/api#inlinequeryresultvenue
type InlineQueryResultVenue struct {
	// Unique identifier for this result, 1-64 Bytes.
	// If left unspecified, a 64-bit FNV-1 hash will be calculated
	// from the other fields and used automatically.
	ID string `json:"id",hash:"ignore"`

	// Latitude of the venue location in degrees.
	Latitude float32 `json:"latitude"`

	// Longitude of the venue location in degrees.
	Longitude float32 `json:"longitude"`

	// Title of the venue.
	Title string `json:"title"`

	// Address of the venue.
	Address string `json:"address"`

	// Optional. Foursquare identifier of the venue if known.
	FoursquareID string `json:"foursquare_id,omitempty"`

	// Optional. Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`

	// Optional. Content of the message to be sent instead of the audio.
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`

	// Optional. Url of the thumbnail for the result.
	ThumbURL string `json:"thumb_url,omitempty"`

	// Optional. Thumbnail width.
	ThumbWidth int `json:"thumb_width,omitempty"`

	// Optional. Thumbnail height.
	ThumbHeight int `json:"thumb_height,omitempty"`
}

func (r *InlineQueryResultVenue) MarshalJSON() ([]byte, error) {
	// avoiding endless self-recursion
	type wrapper InlineQueryResultVenue

	id, err := r.id()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		wrapper
		Type string `json:"type"`
		ID   string `json:"id",hash:"ignore"`
	}{
		wrapper: wrapper(*r),
		ID:      id,
		Type:    "venue",
	})
}

func (r *InlineQueryResultVenue) id() (string, error) {
	if r.ID == "" {
		return hashInlineQueryResult(r)
	}

	return r.ID, nil
}

// InlineQueryResultVideo represents a link to a page containing an embedded
// video player or a video file.
// See also: https://core.telegram.org/bots/api#inlinequeryresultvideo
type InlineQueryResultVideo struct {
	// Unique identifier for this result, 1-64 Bytes.
	// If left unspecified, a 64-bit FNV-1 hash will be calculated
	// from the other fields and used automatically.
	ID string `json:"id",hash:"ignore"`

	// A valid URL for the embedded video player or video file.
	VideoURL string `json:"video_url"`

	// Mime type of the content of video url, “text/html” or “video/mp4”.
	MimeType string `json:"mime_type"`

	// URL of the thumbnail (jpeg only) for the video.
	ThumbURL string `json:"thumb_url"`

	// Title for the result.
	Title string `json:"title"`

	// Optional. Caption of the video to be sent, 0-200 characters.
	Caption string `json:"caption,omitempty"`

	// Optional. Video width.
	VideoWidth int `json:"video_width,omitempty"`

	// Optional. Video height.
	VideoHeight int `json:"video_height,omitempty"`

	// Optional. Video duration in seconds.
	VideoDuration int `json:"video_duration,omitempty"`

	// Optional. Short description of the result.
	Description string `json:"description,omitempty"`

	// Optional. Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`

	// Optional. Content of the message to be sent instead of the audio.
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (r *InlineQueryResultVideo) MarshalJSON() ([]byte, error) {
	// avoiding endless self-recursion
	type wrapper InlineQueryResultVideo

	id, err := r.id()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		wrapper
		Type string `json:"type"`
		ID   string `json:"id",hash:"ignore"`
	}{
		wrapper: wrapper(*r),
		ID:      id,
		Type:    "video",
	})
}

func (r *InlineQueryResultVideo) id() (string, error) {
	if r.ID == "" {
		return hashInlineQueryResult(r)
	}

	return r.ID, nil
}

// InlineQueryResultVoice represents a link to a voice recording in a
// .ogg container encoded with OPUS.
// See also: https://core.telegram.org/bots/api#inlinequeryresultvoice
type InlineQueryResultVoice struct {

	// Unique identifier for this result, 1-64 Bytes.
	// If left unspecified, a 64-bit FNV-1 hash will be calculated
	// from the other fields and used automatically.
	ID string `json:"id",hash:"ignore"`

	// A valid URL for the voice recording.
	VoiceURL string `json:"voice_url"`

	// Recording title.
	Title string `json:"title"`

	// Optional. Recording duration in seconds.
	VoiceDuration int `json:"voice_duration"`

	// Optional. Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`

	// Optional. Content of the message to be sent instead of the audio.
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
}

func (r *InlineQueryResultVoice) MarshalJSON() ([]byte, error) {
	// avoiding endless self-recursion
	type wrapper InlineQueryResultVoice

	id, err := r.id()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		wrapper
		Type string `json:"type"`
		ID   string `json:"id",hash:"ignore"`
	}{
		wrapper: wrapper(*r),
		ID:      id,
		Type:    "voice",
	})
}

func (r *InlineQueryResultVoice) id() (string, error) {
	if r.ID == "" {
		return hashInlineQueryResult(r)
	}

	return r.ID, nil
}
