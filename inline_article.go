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

	// Sends the message silently. iOS users will not receive a notification, Android users will receive a notification with no sound.
	DisableNotification bool

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

// MarshalJSON is a serializer.
func (r ArticleResult) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer

	props := map[string]string{}

	props["type"] = "article"
	props["id"] = r.id()
	props["title"] = r.Title
	props["description"] = r.Description
	props["message_text"] = r.Text

	if r.URL != "" {
		props["url"] = r.URL
	}

	if r.ThumbURL != "" {
		props["thumb_url"] = r.ThumbURL
	}

	if r.HideURL {
		props["hide_url"] = "true"
	}

	if r.DisableWebPagePreview {
		props["disable_web_page_preview"] = "true"
	}

	if r.DisableNotification {
		props["disable_notification"] = "true"
	}

	if r.Mode != ModeDefault {
		props["parse_mode"] = string(r.Mode)
	}

	b.WriteRune('{')

	if len(props) > 0 {
		const tpl = `"%s":"%s",`

		for key, value := range props {
			b.WriteString(fmt.Sprintf(tpl, key, value))
		}

		// the last
		b.WriteString(`"":""`)
	}

	b.WriteRune('}')

	return b.Bytes(), nil
}
