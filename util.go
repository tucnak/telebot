package telebot

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

func wrapSystem(err error) error {
	return errors.Wrap(err, "system error")
}

func (b *Bot) debug(err error) {
	if b.Errors != nil {
		b.Errors <- errors.WithStack(err)
	}
}

func (b *Bot) sendText(to Recipient, text string, opt *SendOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": to.Recipient(),
		"text":    text,
	}
	embedSendOptions(params, opt)

	respJSON, err := b.sendCommand("sendMessage", params)
	if err != nil {
		return nil, err
	}

	return extractMsgResponse(respJSON)
}

func extractMsgResponse(respJSON []byte) (*Message, error) {
	var resp struct {
		Ok          bool
		Result      *Message
		Description string
	}

	err := json.Unmarshal(respJSON, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return nil, errors.Errorf("api error: %s", resp.Description)
	}

	return resp.Result, nil
}

func extractOkResponse(respJSON []byte) error {
	var resp struct {
		Ok          bool
		Description string
	}

	err := json.Unmarshal(respJSON, &resp)
	if err != nil {
		return errors.Wrap(err, "bad response json")
	}

	if !resp.Ok {
		return errors.Errorf("api error: %s", resp.Description)
	}

	return nil
}

func extractOptions(how []interface{}) *SendOptions {
	var opts *SendOptions

	for _, prop := range how {
		switch opt := prop.(type) {
		case *SendOptions:
			opts = opt

		case *ReplyMarkup:
			if opts == nil {
				opts = &SendOptions{}
			}
			opts.ReplyMarkup = opt

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
				panic("telebot: unsupported option")
			}

		case ParseMode:
			if opts == nil {
				opts = &SendOptions{}
			}
			opts.ParseMode = opt

		default:
			panic(fmt.Sprintf("telebot: %v is not a send-option", opt))
		}
	}

	return opts
}

func embedSendOptions(params map[string]string, opt *SendOptions) {
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
		params["parse_mode"] = string(opt.ParseMode)
	}

	if opt.ReplyMarkup != nil {
		replyMarkup, _ := json.Marshal(opt.ReplyMarkup)
		params["reply_markup"] = string(replyMarkup)
	}
}

func embedRights(p map[string]string, prv Rights) {
	jsonRepr, _ := json.Marshal(prv)
	json.Unmarshal(jsonRepr, p)
}
