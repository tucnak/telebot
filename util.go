package telebot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

func (b *Bot) debug(err error) {
	err = errors.WithStack(err)
	if b.reporter != nil {
		b.reporter(err)
	} else {
		log.Printf("%+v\n", err)
	}
}

func (b *Bot) deferDebug() {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			b.debug(err)
		} else if str, ok := r.(string); ok {
			b.debug(errors.Errorf("%s", str))
		}
	}
}

func (b *Bot) runHandler(handler func()) {
	f := func() {
		defer b.deferDebug()
		handler()
	}
	if b.synchronous {
		f()
	} else {
		go f()
	}
}

// wrapError returns new wrapped telebot-related error.
func wrapError(err error) error {
	return errors.Wrap(err, "telebot")
}

// extractOk checks given result for error. If result is ok returns nil.
// In other cases it extracts API error. If error is not presented
// in errors.go, it will be prefixed with `unknown` keyword.
func extractOk(data []byte) error {
	match := errorRx.FindStringSubmatch(string(data))
	if match == nil || len(match) < 3 {
		return nil
	}

	desc := match[2]
	err := ErrByDescription(desc)

	if err == nil {
		code, _ := strconv.Atoi(match[1])

		switch code {
		case http.StatusTooManyRequests:
			retry, _ := strconv.Atoi(match[3])
			err = FloodError{
				APIError:   NewAPIError(429, desc),
				RetryAfter: retry,
			}
		default:
			err = fmt.Errorf("telegram unknown: %s (%d)", desc, code)
		}
	}

	return err
}

// extractMessage extracts common Message result from given data.
// Should be called after extractOk or b.Raw() to handle possible errors.
func extractMessage(data []byte) (*Message, error) {
	var resp struct {
		Result *Message
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		var resp struct {
			Result bool
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, wrapError(err)
		}
		if resp.Result {
			return nil, ErrTrueResult
		}
		return nil, wrapError(err)
	}
	return resp.Result, nil
}

func extractOptions(how []interface{}) *SendOptions {
	var opts *SendOptions

	for _, prop := range how {
		switch opt := prop.(type) {
		case *SendOptions:
			opts = opt.copy()
		case *ReplyMarkup:
			if opts == nil {
				opts = &SendOptions{}
			}
			opts.ReplyMarkup = opt.copy()
		case Option:
			if opts == nil {
				opts = &SendOptions{}
			}

			switch opt {
			case NoPreview:
				opts.DisableWebPagePreview = true
			case Silent:
				opts.DisableNotification = true
			case ForceReply:
				if opts.ReplyMarkup == nil {
					opts.ReplyMarkup = &ReplyMarkup{}
				}
				opts.ReplyMarkup.ForceReply = true
			case OneTimeKeyboard:
				if opts.ReplyMarkup == nil {
					opts.ReplyMarkup = &ReplyMarkup{}
				}
				opts.ReplyMarkup.OneTimeKeyboard = true
			default:
				panic("telebot: unsupported flag-option")
			}
		case ParseMode:
			if opts == nil {
				opts = &SendOptions{}
			}
			opts.ParseMode = opt
		default:
			panic("telebot: unsupported send-option")
		}
	}

	return opts
}

func (b *Bot) embedSendOptions(params map[string]string, opt *SendOptions) {
	if b.parseMode != ModeDefault {
		params["parse_mode"] = b.parseMode
	}

	if opt == nil {
		return
	}

	if opt.ReplyTo != nil && opt.ReplyTo.ID != 0 {
		params["reply_to_message_id"] = strconv.Itoa(opt.ReplyTo.ID)
	}

	if opt.DisableWebPagePreview {
		params["disable_web_page_preview"] = "true"
	}

	if opt.DisableNotification {
		params["disable_notification"] = "true"
	}

	if opt.ParseMode != ModeDefault {
		params["parse_mode"] = opt.ParseMode
	}

	if opt.ReplyMarkup != nil {
		processButtons(opt.ReplyMarkup.InlineKeyboard)
		replyMarkup, _ := json.Marshal(opt.ReplyMarkup)
		params["reply_markup"] = string(replyMarkup)
	}
}

func processButtons(keys [][]InlineButton) {
	if keys == nil || len(keys) < 1 || len(keys[0]) < 1 {
		return
	}

	for i := range keys {
		for j := range keys[i] {
			key := &keys[i][j]
			if key.Unique != "" {
				// Format: "\f<callback_name>|<data>"
				data := key.Data
				if data == "" {
					key.Data = "\f" + key.Unique
				} else {
					key.Data = "\f" + key.Unique + "|" + data
				}
			}
		}
	}
}

func embedRights(p map[string]interface{}, rights Rights) {
	data, _ := json.Marshal(rights)
	_ = json.Unmarshal(data, &p)
}

func thumbnailToFilemap(thumb *Photo) map[string]File {
	if thumb != nil {
		return map[string]File{"thumb": thumb.File}
	}
	return nil
}

func isUserInList(user *User, list []User) bool {
	for _, user2 := range list {
		if user.ID == user2.ID {
			return true
		}
	}
	return false
}
