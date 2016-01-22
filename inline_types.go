package telebot

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// ArticleResult represents a link to an article or web page.
type ArticleResult struct {
	// [Required!] Title of the result.
	Title string

	// [Required!] Text of the message to be sent, 1-512 characters.
	Text string

	// Short description of the result.
	Description string

	// Markdown, HTML?
	Mode ParseMode

	// Disables link previews for links in the sent message.
	DisableWebPagePreview bool

	// URL of the result
	URL string

	// If true, the URL won't be shown in the message.
	HideURL bool

	// Result's thumbnail URL.
	ThumbURL string
}

func (r ArticleResult) id() string {
	sum := md5.Sum([]byte(r.Title + r.Text))
	return string(hex.EncodeToString(sum[:]))
}

// MarshalJSON ...
func (r ArticleResult) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer

	bind := func(key, value string) string {
		return fmt.Sprintf("\"%s\": \"%s\",", key, value)
	}

	bindl := func(key, value string) string {
		return fmt.Sprintf("\"%s\": \"%s\"", key, value)
	}

	b.WriteRune('{')

	b.WriteString(bind("type", "article"))
	b.WriteString(bind("id", r.id()))
	b.WriteString(bind("title", r.Title))
	b.WriteString(bind("description", r.Description))
	b.WriteString(bind("message_text", r.Text))

	if r.URL != "" {
		b.WriteString(bind("url", r.URL))
	}

	if r.ThumbURL != "" {
		b.WriteString(bind("thumb_url", r.URL))
	}

	if r.HideURL {
		b.WriteString(bind("hide_url", "true"))
	}

	if r.DisableWebPagePreview {
		b.WriteString(bind("disable_web_page_preview", "true"))
	}

	if r.Mode != ModeDefault {
		b.WriteString(bindl("parse_mode", string(r.Mode)))
	}

	b.WriteRune('}')

	return b.Bytes(), nil
}
