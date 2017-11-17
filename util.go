package telebot

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

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
	var options *SendOptions

	for _, object := range how {
		switch option := object.(type) {
		case *SendOptions:
			options = option
			break

		case *ReplyMarkup:
			if options == nil {
				options = &SendOptions{}
			}
			options.ReplyMarkup = option
			break

		default:
			panic(fmt.Sprintf("telebot: %v is not a send-option", option))
		}
	}

	return options
}

func embedSendOptions(params map[string]string, opt *SendOptions) {
	if opt == nil {
		return
	}

	if opt.ReplyTo.ID != 0 {
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
		force := opt.ReplyMarkup.ForceReply
		reply := opt.ReplyMarkup.ReplyKeyboard != nil
		inline := opt.ReplyMarkup.InlineKeyboard != nil
		hidden := opt.ReplyMarkup.HideReplyKeyboard
		resize := opt.ReplyMarkup.ResizeReplyKeyboard
		if force || reply || inline || hidden || resize {
			replyMarkup, _ := json.Marshal(opt.ReplyMarkup)
			params["reply_markup"] = string(replyMarkup)
		}
	}
}
